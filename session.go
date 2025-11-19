/**
 * @author JamesWatson
 * @time   2017-06-10
 */
package goengine

import (
	"net/http"
)

type Session struct {
	SessionInfo
	sm SessionManager
}

func (sess *Session) Save(res http.ResponseWriter, maxAge int) error {
	if 0 == maxAge {
		maxAge = sess.sm.MaxAge()
	}
	cookie, err := sess.sm.Save(sess.SessionInfo, maxAge)
	if nil == err {
		http.SetCookie(res, cookie)
	}
	return err
}
