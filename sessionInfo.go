/**
 * @author JamesWatson
 * @time   2017-06-10
 */
package goengine

import (
	"encoding/json"
)

type SessionInfo interface {
	GetSid() string
	Get(key string) interface{}
	Set(key string, value interface{})
	ToJSON() ([]byte, error)
}

type sessionInfo struct {
	sid   string
	store map[string]interface{}
}

func NewSessionInfo(sid string, store map[string]interface{}) SessionInfo {
	return &sessionInfo{sid: sid, store: store}
}

func (si *sessionInfo) GetSid() string {
	return si.sid
}

func (si *sessionInfo) Set(key string, value interface{}) {
	si.store[key] = value
}

func (si *sessionInfo) Get(key string) interface{} {
	value, found := si.store[key]
	if found {
		return value
	}
	return nil
}

func (si *sessionInfo) ToJSON() ([]byte, error) {
	return json.Marshal(si.store)
}
