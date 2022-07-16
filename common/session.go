package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

// SessionHash 记录登陆session的哈希表
const SessionHash = "SessionLogin"

//Session 保存信息
type Session struct {
	SessionID string `json:"session_id"` // SessionID
	UserID    string `json:"user_id"`    // 用户ID
	UserName  string `json:"user_name"`  // 用户名
	Lock      bool   `json:"lock"`       // 锁定
	Appid     string `json:"appid"`
}

//CreateSessionID 创建sessionID
func (s *Session) CreateSessionID() (res string) {
	res = fmt.Sprintf("Miku_%s_%s_%d", s.Appid, s.UserID, time.Now().Unix())
	return
}

//GetSessionUserID 根据sessionID获取用户ID
func (s *Session) GetSessionUserID() {
	sessionSlice := strings.Split(s.SessionID, "_")
	s.Appid = sessionSlice[1]
	s.UserID = sessionSlice[2]
	return
}

//SessionCheck 检查
func SessionCheck(c *gin.Context) {
	path := strings.Split(c.Request.URL.Path, "/")
	if path[len(path)-1] == "health" || path[len(path)-1] == "tenant" {
		c.Next()
		return
	}

	// 无需session的接口
	switch path[len(path)-1] {
	case "login", "logout", "favicon.ico", "token", "licenses", "upgrade", "passwordmsg", "execute", "auth", "results", "receive", "outworker", "terminal":
		c.Next()
		return
	}

	Session_id := c.Request.Header.Get("Session")
	if len(Session_id) == 0 {
		Abort(ErrMissingAuthorization, nil, c)
		return
	}
	se := Session{
		SessionID: Session_id,
	}
	//根据session_id获取用户id
	se.GetSessionUserID()
	// 检查session是否存在且一致未过期
	status, err := se.CheckSession()
	if !status || err != nil {
		Abort(ErrSession, err, c)
		return
	} else {
		// 获取session具体信息
		if err = se.GetSession(); err != nil {
			Abort(ErrSession, nil, c)
			return
		}
	}

	c.Set("userid", se.UserID)
	c.Set("username", se.UserName)
	c.Set("lock", se.Lock)
}

//SessionRegister session 注册
// 设置session_id及其过期时间，并保存记录session_id 用于判断单用户登陆
func (s *Session) SessionRegister() (err error) {
	conn := Pool.Get()
	defer conn.Close()
	if conn == nil {
		return errors.New("redis connection is nil")
	}
	var value []byte
	value, err = json.Marshal(s)
	if err != nil {
		return
	}
	// 记录session_id
	_, err = conn.Do("SET", s.SessionID, value)
	if err != nil {
		LogErrorf("set session_id error", logrus.Fields{"err": err})
		return
	}
	// 设置过期时间
	if _, err = conn.Do("EXPIRE", s.SessionID, CONFIG.SessionExpireTime); err != nil {
		LogErrorf("EXPIRE session_id error", logrus.Fields{"err": err})
		return
	}
	// 保存当前session
	if _, err = conn.Do("HSET", SessionHash, fmt.Sprintf("%s_%s", s.Appid, s.UserID), s.SessionID); err != nil {
		LogErrorf("hset sessionHash error", logrus.Fields{"err": err})
		return
	}
	return
}

//CheckSession 检查session，用于每次请求判断
//判断是否是有符合的用户登陆且session_id一致, 否 则要退出重新登录
func (s *Session) CheckSession() (status bool, err error) {
	conn := Pool.Get()
	defer conn.Close()
	if conn == nil {
		return false, errors.New("redis connection is nil")
	}
	if status, err = redis.Bool(conn.Do("HEXISTS", SessionHash, fmt.Sprintf("%s_%s", s.Appid, s.UserID))); err != nil || !status {
		status = false
		LogErrorf("redis HEXISTS Session HashMap error,may be the param session is not equal to the save one", logrus.Fields{"err": err})
		return
	} else {
		//  判断是否和记录的session_id一致
		exisSessionid, _ := redis.String(conn.Do("HGET", SessionHash, fmt.Sprintf("%s_%s", s.Appid, s.UserID)))
		if exisSessionid != s.SessionID {
			return false, nil
		}
		// 如果存在记录session_id ,且还未过期
		if status, err = redis.Bool(conn.Do("EXISTS", s.SessionID)); err != nil || !status {
			status = false
			LogErrorf("redis EXISTS session_id error,or is already expired", logrus.Fields{"err": err})
			return
		} else {
			return
		}
	}
}

//SetSessionExpire 设置过期时间,用于每次请求不断重新设置过期时间
func (s *Session) SetSessionExpire() (err error) {
	conn := Pool.Get()
	if conn == nil {
		return errors.New("redis connection is nil")
	}
	defer conn.Close()
	_, err = conn.Do("EXPIRE", s.SessionID, CONFIG.SessionExpireTime)
	if err != nil {
		return
	}
	return
}

//SetSession 设置session
func (s *Session) SetSession() (err error) {
	conn := Pool.Get()
	if conn == nil {
		return errors.New("redis connection is nil")
	}
	defer conn.Close()
	var value []byte
	value, err = json.Marshal(s)
	if err != nil {
		return
	}
	_, err = conn.Do("SET", s.SessionID, value)
	if err != nil {
		return
	}
	return
}

//GetSession 获取session
func (s *Session) GetSession() (err error) {
	conn := Pool.Get()
	if conn == nil {
		return errors.New("redis connection is nil")
	}
	defer conn.Close()
	var unMarshal []byte
	unMarshal, err = redis.Bytes(conn.Do("GET", s.SessionID))
	err = json.Unmarshal(unMarshal, s)
	return
}

//DeleteSession 删除用户登陆信息
func (s *Session) DeleteSession() (err error) {
	conn := Pool.Get()
	if conn == nil {
		return errors.New("redis connection is nil")
	}
	defer conn.Close()
	_, err = conn.Do("HDEL", SessionHash, s.UserID)
	var status bool
	if status, err = redis.Bool(conn.Do("HEXISTS", SessionHash, fmt.Sprintf("%s_%s", s.Appid, s.UserID))); err != nil {
		LogErrorf("redis DeleteSession ,HEXISTS User Token Error ", logrus.Fields{"err": err, "status": status})
	} else if !status {
		LogInfo("redis DeleteSession Success User Is Not Login")
	} else if status {
		_, err = conn.Do("HDEL", SessionHash, fmt.Sprintf("%s_%s", s.Appid, s.UserID))
		if err != nil {
			LogErrorf("redis DeleteSession ,HDEL Error ", logrus.Fields{"err": err, "status": status})
		}
	}
	return
}
