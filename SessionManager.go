/**
 * @author JamesWatson
 * @time   2017-06-10
 */
package goengine

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"slices"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SessionStore interface {
	Get(string, interface{}) error
	/**
	 * @param json session data in json format, if nil or empty, delete the key
	 * @param maxAge in seconds, if maxAge ==0, no expiration; if maxAge <0, delete the key
	 */
	Save(string, []byte, int) error
}

type SessionManager interface {
	MaxAge() int
	Secure() bool
	LoadSession(resp http.ResponseWriter, req *http.Request) SessionInfo
	UpData(session SessionInfo, maxAge int) error
	UpMaxAge(resp http.ResponseWriter, maxAge int) error
	Save(resp http.ResponseWriter, session SessionInfo, maxAge int) error
	Delete(resp http.ResponseWriter, req *http.Request) error
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

func InitSessionManager(storer SessionStore, sessName string, cookiePrefix string, sessionPrefix string, domain string) SessionManager {
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

func (sm *sessionManager) loadSession(req *http.Request) SessionInfo {
	valMap := make(map[string]string)
	for true {
		cookie, err := req.Cookie(sm.sessionName)
		if nil != err {
			break
		}
		sid := cookie.Value[len(sm.cookiePrefix):]
		if "" == sid {
			break
		}
		err = sm.storer.Get(sm.sessionPrefix+sid, &valMap)
		if nil != err {
			break
		}
		return NewSessionInfo(sid, valMap)
	}

	// 生成新的sid 和 session
	return NewSessionInfo(GenerateSid(), valMap)
}

func (sm *sessionManager) LoadSession(resp http.ResponseWriter, req *http.Request) SessionInfo {
	session := sm.loadSession(req)
	sid := session.GetSid()

	cookie := &http.Cookie{
		Name:     sm.sessionName,
		Value:    sm.cookiePrefix + sid,
		Path:     "/",
		Domain:   req.Host,
		Secure:   sm.Secure(),
		HttpOnly: true,
		Expires:  GetExpirationTime(sm.MaxAge()),
	}
	http.SetCookie(resp, cookie)
	return session
}

/**
 * @param resp http response writer
 * @param session session info
 * @param maxAge in seconds, if maxAge < 1 use default maxAge
 */
func (sm *sessionManager) UpData(session SessionInfo, maxAge int) error {
	sid := session.GetSid()
	data, err := session.ToJSON()
	if maxAge < 1 {
		maxAge = sm.maxAge
	}
	if nil == err {
		err = sm.storer.Save(sm.sessionPrefix+sid, data, maxAge)
	}
	return err
}

func (sm *sessionManager) UpMaxAge(resp http.ResponseWriter, maxAge int) error {
	header := resp.Header()
	dirtyCookies := header["Set-Cookie"]

	str := "" // cookie or domain
	for idx, item := range dirtyCookies {
		if strings.HasPrefix(item, sm.sessionName+"=") {
			str = item
			header["Set-Cookie"] = slices.Delete(dirtyCookies, idx, idx+1)
			break
		}
	}
	cookie, err := http.ParseSetCookie(str)
	if nil != err {
		return err
	}
	if maxAge < 1 {
		maxAge = sm.maxAge
	}
	cookie.Expires = GetExpirationTime(maxAge)

	http.SetCookie(resp, cookie)
	return nil
}

func (sm *sessionManager) Save(resp http.ResponseWriter, session SessionInfo, maxAge int) error {
	err := sm.UpData(session, maxAge)
	if nil == err {
		err = sm.UpMaxAge(resp, maxAge)
	}
	return err
}

func (sm *sessionManager) Delete(resp http.ResponseWriter, req *http.Request) error {
	cookie, err := req.Cookie(sm.sessionName)

	switch err {
	case nil:
	case http.ErrNoCookie:
		return nil
	default:
		return err
	}

	sid := cookie.Value[len(sm.cookiePrefix):]
	if "" == sid {
		return nil
	}

	sm.storer.Save(sm.sessionPrefix+sid, nil, -1)
	cookie = &http.Cookie{
		Name:     sm.sessionName,
		Value:    sm.cookiePrefix + sid,
		Path:     "/",
		Domain:   req.Host,
		Secure:   sm.Secure(),
		HttpOnly: true,
		Expires:  GetExpirationTime(0),
	}
	http.SetCookie(resp, cookie)
	return nil
}
