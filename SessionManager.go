package goengine

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

func GenerateSid() string {
	md5Gen := md5.New()
	uuid_v1, _ := uuid.NewV1()
	uuid_v4, _ := uuid.NewV4()

	md5Gen.Write([]byte(uuid_v1.String() + "-" + uuid_v4.String()))
	cipherStr := md5Gen.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

/**
 * @class  Session
 * @author JamesWatson
 * @time   2017-06-10
 */
type Session struct {
	sid   string
	sm    *SessionManager
	res   *http.ResponseWriter
	store map[string]interface{}
}

func (this *Session) Set(key string, value interface{}) {
	this.store[key] = value
}

func (this *Session) Get(key string) interface{} {
	value, found := this.store[key]
	if found {
		return value
	}
	return nil
}

func (this *Session) Save(maxAge int) {
	if 0 == maxAge {
		maxAge = this.sm.MaxAge
	}
	this.sm.Save(this, maxAge)
}

/**
 * @class  SessionManager
 * @author JamesWatson
 * @time   2017-06-10
 */
type SessionManager struct {
	conn          redis.Conn
	sessionName   string
	cookiePrefix  string
	sessionPrefix string
	domain        string
	MaxAge        int
	Secure        bool
}

func InitSessionManager(redisConn redis.Conn, sessName string, cookiePrefix string, sessionPrefix string, domain string) *SessionManager {

	ret := &SessionManager{
		conn:          redisConn,
		sessionName:   sessName,
		cookiePrefix:  cookiePrefix,
		sessionPrefix: sessionPrefix,
		domain:        domain,
		MaxAge:        3600 * 24,
		Secure:        true,
	}
	return ret
}

func (this *SessionManager) getExpirationTime(maxAge int) time.Time {
	return time.Now().UTC().Add(time.Duration(maxAge) * time.Second)
}

func (this *SessionManager) Get(res *http.ResponseWriter, req *http.Request) *Session {
	sessionInfo := &Session{
		sm:  this,
		res: res,
	}
	var sid string

	cookie, err := req.Cookie(this.sessionName)
	if nil == err {
		sid = cookie.Value[len(this.cookiePrefix):]
		sessionInfo.sid = sid

		jsonSession, err := redis.String(this.conn.Do("GET", this.sessionPrefix+sid))

		if nil == err {
			json.Unmarshal([]byte(jsonSession), &(sessionInfo.store))
			return sessionInfo
		}
	}

	// 生成新的sid 和 session
	sid = GenerateSid()
	sessionInfo.sid = sid
	sessionInfo.store = make(map[string]interface{})

	return sessionInfo
}

func (this *SessionManager) Save(session *Session, maxAge int) error {
	cookie := &http.Cookie{
		Name:     this.sessionName,
		Value:    this.cookiePrefix + session.sid,
		Path:     "/",
		Domain:   this.domain,
		Secure:   this.Secure,
		HttpOnly: true,
		Expires:  this.getExpirationTime(maxAge),
	}
	http.SetCookie(*session.res, cookie)
	buf, err := json.Marshal(session.store)
	if nil != err {
		return err
	}
	_, err = redis.String(this.conn.Do("SETEX", this.sessionPrefix+session.sid, maxAge, string(buf)))

	return err
}
