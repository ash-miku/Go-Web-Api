package common

import (
	"bytes"
	"errors"
	"fmt"
	"go-api/utils"
	"io/ioutil"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/willf/pad"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// configKey is the key name of the config context in the Gin context.
const middleConfigKey = "config"

//PublicDB 公共表的连接DB
var PublicDB *gorm.DB

//ProjectName 对应组件名称
var ProjectName string

// MiddlewareConfig Config is a middleware function that initializes the config and attaches to
// the context of every request context.
func MiddlewareConfig(cli *cli.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleConfigKey, cli)
		c.Set("service-name", cli.String("service"))
		c.Set("jwt-secret", cli.String("jwt-secret"))
	}
}

// MiddlewareConfigContext ConfigContext returns the CLI context associated with this context.
func MiddlewareConfigContext(c *gin.Context) (ctx *cli.Context) {
	conf, _ := c.Get(middleConfigKey)
	ctx = conf.(*cli.Context)
	return
}

func open(driver, config string, serviceName string, values ...interface{}) (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open(config), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   strings.ToLower(serviceName) + "_",
			SingularTable: true,
		},
	})
	//db.LogMode(true)
	if err != nil {
		LogFatalf("Database connection failed.", logrus.Fields{
			"err": err,
		})
	}

	return
}

func DatabaseConnect(project string, values ...interface{}) gin.HandlerFunc {
	// 不同组件使用不同的表名称
	ProjectName = project

	Password, err := utils.Base64Decrypt(CONFIG.DB.Password)
	if err != nil {
		LogFatalf("Database Password Decrypt failed.", logrus.Fields{
			"err": err,
		})
	}

	DB := open(
		CONFIG.DB.Driver,
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
			CONFIG.DB.User,
			Password,
			CONFIG.DB.Host,
			CONFIG.DB.Port,
			CONFIG.DB.Name,
			CONFIG.DB.Charset,
			CONFIG.DB.ParseTime,
			CONFIG.DB.Loc,
		),
		project,
		values...,
	)

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(50)
	// 公共错误的信息获取
	GetErrnoMessages(DB)

	return func(c *gin.Context) {
		c.Set("DB", DB.WithContext(c))
		c.Next()
	}
}

// MiddleLogging Logging is a middleware function that logs the each request.
func MiddleLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		method := c.Request.Method
		// 解析请求体中的参数输出至日志
		body, _ := ioutil.ReadAll(c.Request.Body)
		// gin 的请求体是一次性的，如果后续还有步骤需要将数据重新放回请求
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		// Continue.
		c.Next()
		// 跳过某些轮询API，防止打印过多相同日志
		if path == "/v1/message/count" || (path == "/v1/tasks" && method == "GET") || (path == "/v1/notify/receive" && method == "GET") {
			return
		}
		// Calculates the latency.
		end := time.Now().UTC()
		latency := end.Sub(start)

		// The basic informations.
		status := c.Writer.Status()
		ip := c.GetHeader("X-real-ip")
		if ip == "" {
			ip = c.ClientIP()
		}
		userAgent := c.Request.UserAgent()

		// Create the symbols for each status.
		statusString := ""
		switch {
		case status >= 500:
			statusString = fmt.Sprintf("▲ %d", status)
		case status >= 400:
			statusString = fmt.Sprintf("▲ %d", status)
		case status >= 300:
			statusString = fmt.Sprintf("■ %d", status)
		case status >= 100:
			statusString = fmt.Sprintf("● %d", status)
		}

		// 请求体仅打印前 200 个字节
		// 注意，使用 byte 切分字符串可能会出现最后一个中文乱码的情况
		// 可以使用 rune 切片避免乱码，但是为了减少转换的次数，此处仍旧使用 byte 切分
		if len(body) > 200 {
			body = body[:200]
		}
		// Data fields that will be recorded into the log files.
		fields := logrus.Fields{
			"user_agent": userAgent,
			"body":       string(body),
		}
		// Append the error to the fields so we can record it.
		if len(c.Errors) != 0 {
			for k, v := range c.Errors {
				// Skip if it's the Gin internal error.
				if !v.IsType(gin.ErrorTypePrivate) {
					continue
				}
				// The field name with the `error_INDEX` format.
				errorKey := fmt.Sprintf("error_%d", k)

				switch v.Err.(type) {
				case *Err:
					e := v.Err.(*Err)
					fields[errorKey] = fmt.Sprintf("%s[%s:%d:%s]", e.ErrType, e.File, e.Line, e.Func)
				default:
					fields[errorKey] = fmt.Sprintf("%s", v.Err)
				}
			}
		}

		msg := fmt.Sprintf(" %s | %13s | %12s | %s %s", statusString, latency, ip, pad.Right(method, 5, " "), path)
		if len(c.Errors) == 0 {
			// Example: ● 200 |  102.268592ms |    127.0.0.1 | POST  /user (user_agent=xxx)
			LogInfof(msg, fields)
		} else {
			// Example: ▲ 403 |  102.268592ms |    127.0.0.1 | POST  /user (user_agent=xxx error_0=xxx)
			LogErrorf(msg, fields)
		}
	}
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		lmt := tollbooth.NewLimiter(0.01, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
		lmt.SetIPLookups([]string{"RemoteAddr"})
		httpErr := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)

		if httpErr != nil {
			Abort(ErrRateLimit, errors.New(httpErr.Message), c)
		}
	}
}
