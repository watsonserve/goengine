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

func (sess *Session) Save(resp http.ResponseWriter, maxAge int) error {
	return sess.sm.Save(resp, sess.SessionInfo, maxAge)
}
