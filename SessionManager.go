/**
 * @author JamesWatson
 * @time   2017-06-10
 */
package goengine

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SessionStore interface {
	Get(string) (*map[string]interface{}, error)
	Save(string, []byte, int) error
}

type SessionManager interface {
	MaxAge() int
	Secure() bool
	LoadSession(req *http.Request) SessionInfo
	Get(req *http.Request) *Session
	Save(session SessionInfo, maxAge int) (*http.Cookie, error)
}

func GenerateSid() string {
	md5Gen := md5.New()
	uuid_v1 := uuid.NewV1()
	uuid_v4 := uuid.NewV4()

	md5Gen.Write([]byte(uuid_v1.String() + "-" + uuid_v4.String()))
	cipherStr := md5Gen.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func GetExpirationTime(maxAge int) time.Time {
	return time.Now().UTC().Add(time.Duration(maxAge) * time.Second)
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

func (sm *sessionManager) MaxAge() int {
	return sm.maxAge
}

func (sm *sessionManager) Secure() bool {
	return sm.secure
}

func (sessMgr *sessionManager) LoadSession(req *http.Request) SessionInfo {
	for true {
		cookie, err := req.Cookie(sessMgr.sessionName)
		if nil != err {
			break
		}
		sid := cookie.Value[len(sessMgr.cookiePrefix):]
		if "" == sid {
			break
		}
		valMap, err := sessMgr.storer.Get(sessMgr.sessionPrefix + sid)
		if nil != err {
			break
		}
		return NewSessionInfo(sid, *valMap)
	}

	// 生成新的sid 和 session
	return NewSessionInfo(GenerateSid(), make(map[string]interface{}))
}

func (sessMgr *sessionManager) Get(req *http.Request) *Session {
	info := sessMgr.LoadSession(req)
	return &Session{
		sm:          sessMgr,
		SessionInfo: info,
	}
}

func (sm *sessionManager) Save(session SessionInfo, maxAge int) (*http.Cookie, error) {
	sid := session.GetSid()
	data, err := session.ToJSON()
	if nil != err {
		return nil, err
	}
	cookie := &http.Cookie{
		Name:     sm.sessionName,
		Value:    sm.cookiePrefix + sid,
		Path:     "/",
		Domain:   sm.domain,
		Secure:   sm.Secure(),
		HttpOnly: true,
		Expires:  GetExpirationTime(maxAge),
	}

	return cookie, sm.storer.Save(sm.sessionPrefix+sid, data, maxAge)
}
