package common

import (
	"errors"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ErrorMessageMap var PrivilegeChangeMark sync.Map
var ErrorMessageMap map[string]Error
var (
	ErrBind                 = "ErrBind"
	ErrValidation           = "ErrValidation"
	ErrEncrypt              = "ErrEncrypt"
	ErrDatabase             = "ErrDatabase"
	ErrRecordNotFound       = "ErrRecordNotFound"
	ErrTokenInvalid         = "ErrTokenInvalid"
	ErrIamForbidden         = "ErrIamForbidden"
	ErrUserIncorrect        = "ErrUserIncorrect"
	ErrTokenParse           = "ErrTokenParse"
	ErrTokenSign            = "ErrTokenSign"
	ErrMissingAuthorization = "ErrMissingAuthorization"
	ErrRateLimit            = "ErrRateLimit"
	ErrDBProxy              = "ErrDBProxy"
	ErrMonitorProxy         = "ErrMonitorProxy"
	ErrInstanceNotFound     = "ErrInstanceNotFound"
	ErrAnthenaProxy         = "ErrAnthenaProxy"
	ErrZeusProxy            = "ErrZeusProxy"
	ErrCreateLogMonitor     = "ErrCreateLogMonitor"
	ErrLogMonitorNotFound   = "ErrLogMonitorNotFound"
	ErrSSHProxy             = "ErrSSHProxy"
	ErrMissingHeader        = errors.New("The length of the `Authorization` header is zero.")
	ErrDecode               = "ErrDecode"
	ErrDuplicate            = "ErrDuplicate"
	ErrConnectFailed        = "ErrConnectFailed"
	ErrConnectError         = "ErrConnectError"
	ErrDatabaseType         = "ErrDatabaseType"
	ErrResponse             = "ErrResponse"
	ErrSession              = "ErrSession"
	ErrConnectorPause       = "ErrConnectorPause"
	ErrConnectorStart       = "ErrConnectorStart"
	ErrPasswordRule         = "ErrPasswordRule"
	ErrPasswordExpired      = "ErrPasswordExpired"
	ErrCreateInstance       = "ErrCreateInstance"
	ErrUpdateInstance       = "ErrUpdateInstance"
	ErrDeleteInstance       = "ErrDeleteInstance"
	ErrCreateInstanceGroup  = "ErrCreateInstanceGroup"
	ErrQueryInstanceGroup   = "ErrQueryInstanceGroup"
	ErrUpdateInstanceGroup  = "ErrUpdateInstanceGroup"
	ErrDeleteInstanceGroup  = "ErrDeleteInstanceGroup"
	ErrReadFile             = "ErrReadFile"
	ErrGoogleVerify         = "ErrGoogleVerify"
	ErrWindowsAdError       = "ErrWindowsAdError"
	ErrWindowsADFailed      = "ErrWindowsADFailed"
	ErrCreateDraftBox       = "ErrCreateDraftBox"
	ErrDeleteDraftBox       = "ErrDeleteDraftBox"
	ErrDeleteTemplate       = "ErrDeleteTemplate"
	ErrUpdateTemplate       = "ErrUpdateTemplate"

	ErrCreateWebsocket    = "ErrCreateWebsocket"
	ErrValidate           = "ErrValidate"
	ErrQueryTask          = "ErrQueryTask"
	ErrBuildParams        = "ErrBuildParams"
	ErrExecSql            = "ErrExecSql"
	ErrQueryOverSql       = "ErrQueryOverSql"
	ErrTaskIsStillRunning = "ErrTaskIsStillRunning"
	ErrCreateTask         = "ErrCreateTask"
	ErrAnalyzeSql         = "ErrAnalyzeSql"
	ErrUpdateTask         = "ErrUpdateTask"
	ErrDeleteTask         = "ErrDeleteTask"
	ErrUserStatus         = "ErrUserStatus"
	ErrNodeFailed         = "ErrNodeFailed"
	ErrNodeError          = "ErrNodeError"

	ErrSupportDataBase  = "ErrSupportDataBase"
	ErrGetCurrentSchema = "ErrGetCurrentSchema"
	ErrSetCurrentSchema = "ErrSetCurrentSchema"

	ErrCreateColumnReflect = "ErrCreateColumnReflect"
	ErrUpdateColumnReflect = "ErrUpdateColumnReflect"

	ErrCancelTask          = "ErrCancelTask"
	ErrOutOfLimit          = "ErrOutOfLimit"
	ErrSupportSqlClassType = "ErrSupportSqlClassType"
	ErrUpdateDraftBox      = "ErrUpdateDraftBox"
	ErrCreateUUID          = "ErrCreateUUID"
	ErrCreateExportTask    = "ErrCreateExportTask"
	ErrUpdateExportTask    = "ErrUpdateExportTask"
	ErrImportData          = "ErrImportData"
	ErrGetExportTask       = "ErrGetExportTask"
	ErrPermission          = "ErrPermission"
	ErrConnectToDatabase   = "ErrConnectToDatabase"
	ErrRemoveFile          = "ErrRemoveFile"
	ErrFileIsNotExist      = "ErrFileIsNotExist"
	ErrReloadExportTask    = "ErrReloadExportTask"
	ErrCloseConnection     = "ErrCloseConnection"
	ErrGetExportFile       = "ErrGetExportFile"

	ErrCreateOrganization = "ErrCreateOrganization"
	ErrUpdateOrganization = "ErrUpdateOrganization"

	ErrBindUserToInstance = "ErrBindUserToInstance"
	ErrUserAccess         = "ErrUserAccess"
	ErrGetAccess          = "ErrGetAccess"

	ErrDownloadOperationLog  = "ErrDownloadOperationLog"
	ErrSendNotify            = "ErrSendNotify"
	ErrSendTestNotify        = "ErrSendTestNotify"
	ErrCreateNotify          = "ErrCreateNotify"
	ErrCreatePubSubConn      = "ErrCreatePubSubConn"
	ErrSubscribeNotify       = "ErrSubscribeNotify"
	ErrReceiveNotify         = "ErrReceiveNotify"
	ErrBindRestrictWithUser  = "ErrBindRestrictWithUser"
	ErrDeleteRestrictAccess  = "ErrDeleteRestrictAccess"
	ErrRestrictAccessInvalid = "ErrRestrictAccessInvalid"

	ErrInterceptRule        = "ErrInterceptRule"
	ErrInterceptRuleCheck   = "ErrInterceptRuleCheck"
	ErrInterceptRuleExecute = "ErrInterceptRuleExecute"

	ErrCreateDataMasking = "ErrCreateDataMasking"

	ErrLicenseExpired       = "ErrLicenseExpired"
	ErrLicenseAuthorization = "ErrLicenseAuthorization"
	ErrActivation           = "ErrActivation"
	ErrLicenseParse         = "ErrLicenseParse"
	ErrKillSession          = "ErrKillSession"
	ErrGenerateSuggestion   = "ErrGenerateSuggestion"

	ErrIgnoreErrRestrictAccess = "ErrIgnoreErrRestrictAccess"
	ErrIgnoreErrDataMasking    = "ErrIgnoreErrDataMasking"

	ErrIpAccess            = "ErrIpAccess"
	ErrSaveIpAccessSetting = "ErrSaveIpAccessSetting"
	ErrLoginIpAccess       = "ErrLoginIpAccess"
)

type Req struct {
	ErrorCode int         `json:"error_code" example:"0"`
	Msg       Err         `json:"msg"`
	Data      interface{} `json:"data"`
	Lock      bool        `json:"lock"`
}

//PageData 分页数据
/*
	对于接口中需要返回查询数量或者分页的结果，需要使用该结构体传给SuccessReturn
	SuccessReturn(PageData={result,total})
*/
type PageData struct {
	Total  int64       `json:"total" default:"0"` // 查询总数
	Result interface{} `json:"result"`            // 查询结果分页数据
}

type Error struct {
	ID        int    `json:"id"`
	Service   string `json:"service"`
	ErrType   string `json:"type"`
	ErrorCode int    `json:"error_code"`
	MessageEn string `json:"message_en"`
	MessageZh string `json:"message_zh"`
}

// Err represents an error, `Code`, `File`, `Line`, `Func` will be automatically filled.
type Err struct {
	ErrType   string `json:"code" example:"ErrNone"`
	MessageEn string `json:"message_en" example:"this is an error"`
	MessageZh string `json:"message_zh" example:"这是一个错误"`
	Detail    string `json:"detail" example:"xxx"`
	File      string `json:"file" example:"test.go"`
	Line      int    `json:"line" example:"111"`
	Func      string `json:"func" example:"/golang/test.go/test():123"`
}

// Error returns the error message.
func (e *Err) Error() string {
	return e.MessageEn
}

// Fill the error struct with the detail error information.
func fill(e *Err) *Err {
	// Fill the error occurred path, line, code.
	pc, fn, line, _ := runtime.Caller(2)
	e.File = strings.Replace(fn, os.Getenv("GOPATH")+"/src/", "", -1)
	e.Line = line
	f := runtime.FuncForPC(pc).Name()
	e.Func = strings.Split(f, ".")[1] + "()"
	return e
}

// Abort the current request with the specified error code.
func Abort(code string, err error, c *gin.Context) {
	isLock := c.Value("lock")

	if isLock == nil {
		isLock = false
	}

	var req Req
	d := ErrorMessageMap[code]
	req.ErrorCode = d.ErrorCode

	// Get error StatusCode, Code, Message from errno.ERROR_MESSAGE
	req.Msg.ErrType = d.ErrType
	req.Msg.MessageEn = d.MessageEn
	req.Msg.MessageZh = d.MessageZh
	req.Lock = isLock.(bool)
	if err == nil {
		err = errors.New(req.Msg.MessageEn)
		req.Msg.Detail = "No more detail, see `Message`."
	} else {
		req.Msg.Detail = err.Error()
	}

	_ = c.Error(err)
	_ = c.Error(fill(&req.Msg))
	c.JSON(http.StatusInternalServerError, req)
	c.Abort()
}

func SuccessReturn(result interface{}, c *gin.Context) {
	isLock := c.Value("lock")

	if isLock == nil {
		isLock = false
	}

	var req Req
	req.Data = result
	req.ErrorCode = 0
	req.Lock = isLock.(bool)

	c.JSON(http.StatusOK, req)
}

// GetErrnoMessages getErrnoMessages get the errno message.
func GetErrnoMessages(DB *gorm.DB) {
	var errorList []Error

	if err := DB.Raw("select * from miku_errors").Find(&errorList).Error; err != nil {
		LogFatalf("Get errno messags failed.", logrus.Fields{
			"error": err.Error(),
		})
	}
	data := make(map[string]Error)
	for _, v := range errorList {
		data[v.ErrType] = v
	}

	ErrorMessageMap = data
}
