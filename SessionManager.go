package goengine

import (
	"net/http"
	"time"
)

type SessionStore interface {
	Get(string) (*map[string]interface{}, error)
	Save(string, *map[string]interface{}, int) error
}

type SessionManager interface {
	MaxAge() int
	Secure() bool
	Get(req *http.Request) *Session
	Save(session *Session, maxAge int) (*http.Cookie, error)
}

/**
 * @class  Session
 * @author JamesWatson
 * @time   2017-06-10
 */
type Session struct {
	sid   string
	sm    SessionManager
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

func (this *Session) Save(res http.ResponseWriter, maxAge int) error {
	if 0 == maxAge {
		maxAge = this.sm.MaxAge()
	}
	cookie, err := this.sm.Save(this, maxAge)
	if nil == err {
		http.SetCookie(res, cookie)
	}
	return err
}

/**
 * @class  sessionManager
 * @author JamesWatson
 * @time   2017-06-10
 */
type sessionManager struct {
	storer        SessionStore
	sessionName   string
	cookiePrefix  string
	sessionPrefix string
	domain        string
	maxAge        int
	secure        bool
}

func InitSessionManager(storer SessionStore, sessName string, cookiePrefix string, sessionPrefix string, domain string) *sessionManager {
	return &sessionManager{
		storer:        storer,
		sessionName:   sessName,
		cookiePrefix:  cookiePrefix,
		sessionPrefix: sessionPrefix,
		domain:        domain,
		maxAge:        3600 * 24,
		secure:        true,
	}
}

func (this *sessionManager) MaxAge() int {
	return this.maxAge
}

func (this *sessionManager) Secure() bool {
	return this.secure
}

func (this *sessionManager) getExpirationTime(maxAge int) time.Time {
	return time.Now().UTC().Add(time.Duration(maxAge) * time.Second)
}

func (this *sessionManager) Get(req *http.Request) *Session {
	sessionInfo := &Session{
		sm: this,
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

func (this *sessionManager) Save(session *Session, maxAge int) (*http.Cookie, error) {
	cookie := &http.Cookie{
		Name:     this.sessionName,
		Value:    this.cookiePrefix + session.sid,
		Path:     "/",
		Domain:   this.domain,
		Secure:   this.Secure(),
		HttpOnly: true,
		Expires:  this.getExpirationTime(maxAge),
	}
	return cookie, this.storer.Save(this.sessionPrefix+session.sid, &(session.store), maxAge)
}
