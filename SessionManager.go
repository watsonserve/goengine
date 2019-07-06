package goengine

import (
	"net/http"
	"time"
)

type SessionStore interface {
	Get(string) (*map[string]interface{}, error)
	Save(string, *map[string]interface{}, int) error
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

func (this *Session) Save(maxAge int) error {
	if 0 == maxAge {
		maxAge = this.sm.MaxAge
	}
	return this.sm.Save(this, maxAge)
}

/**
 * @class  SessionManager
 * @author JamesWatson
 * @time   2017-06-10
 */
type SessionManager struct {
	storer        SessionStore
	sessionName   string
	cookiePrefix  string
	sessionPrefix string
	domain        string
	MaxAge        int
	Secure        bool
}

func InitSessionManager(storer SessionStore, sessName string, cookiePrefix string, sessionPrefix string, domain string) *SessionManager {
	return &SessionManager{
		storer:        storer,
		sessionName:   sessName,
		cookiePrefix:  cookiePrefix,
		sessionPrefix: sessionPrefix,
		domain:        domain,
		MaxAge:        3600 * 24,
		Secure:        true,
	}
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

		valMap, err := this.storer.Get(this.sessionPrefix + sid)

		if nil == err {
			sessionInfo.store = *valMap
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
	return this.storer.Save(this.sessionPrefix + session.sid, &(session.store), maxAge)
}
