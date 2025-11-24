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
	Get(string) (*map[string]string, error)
	/**
	 * @param json session data in json format, if nil or empty, delete the key
	 * @param maxAge in seconds, if maxAge ==0, no expiration; if maxAge <0, delete the key
	 */
	Save(string, []byte, int) error
}

type SessionManager interface {
	MaxAge() int
	Secure() bool
	LoadSession(req *http.Request) SessionInfo
	Get(req *http.Request) *Session
	Save(resp http.ResponseWriter, session SessionInfo, maxAge int) error
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
	return NewSessionInfo(GenerateSid(), make(map[string]string))
}

func (sessMgr *sessionManager) Get(req *http.Request) *Session {
	info := sessMgr.LoadSession(req)
	return &Session{
		sm:          sessMgr,
		SessionInfo: info,
	}
}

/**
 * @param resp http response writer
 * @param session session info
 * @param maxAge in seconds, if maxAge ==0, delete the key; if maxAge <0, use default maxAge
 */
func (sm *sessionManager) Save(resp http.ResponseWriter, session SessionInfo, maxAge int) error {
	sid := session.GetSid()
	data, err := session.ToJSON()
	if nil != err {
		return err
	}

	switch {
	case maxAge < 0:
		maxAge = sm.MaxAge()
	case 0 == maxAge:
		maxAge = -1
	}

	err = sm.storer.Save(sm.sessionPrefix+sid, data, maxAge)
	if nil == err {
		cookie := &http.Cookie{
			Name:     sm.sessionName,
			Value:    sm.cookiePrefix + sid,
			Path:     "/",
			Domain:   sm.domain,
			Secure:   sm.Secure(),
			HttpOnly: true,
			Expires:  GetExpirationTime(maxAge),
		}
		http.SetCookie(resp, cookie)
	}
	return err
}
