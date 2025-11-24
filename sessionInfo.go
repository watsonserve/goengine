/**
 * @author JamesWatson
 * @time   2017-06-10
 */
package goengine

import (
	"encoding/json"
	"errors"
)

type SessionInfo interface {
	GetSid() string
	Get(key string) interface{}
	Set(key string, value interface{}) error
	Load(key string, v any) error
	ToJSON() ([]byte, error)
}

type sessionInfo struct {
	sid   string
	store map[string]string
}

func NewSessionInfo(sid string, store map[string]string) SessionInfo {
	return &sessionInfo{sid: sid, store: store}
}

func (si *sessionInfo) GetSid() string {
	return si.sid
}

func (si *sessionInfo) Set(key string, value interface{}) error {
	if nil == value {
		delete(si.store, key)
		return nil
	}
	jsonValue, err := json.Marshal(value)
	if nil == err {
		si.store[key] = string(jsonValue)
	}
	return err
}

func (si *sessionInfo) Get(key string) interface{} {
	value, found := si.store[key]
	if found {
		ret := make(map[string]interface{})
		if err := json.Unmarshal([]byte(value), &ret); nil == err {
			return ret
		}
	}
	return nil
}

func (si *sessionInfo) Load(key string, v any) error {
	value, found := si.store[key]
	if !found {
		return errors.New("not found")
	}
	return json.Unmarshal([]byte(value), &v)
}

func (si *sessionInfo) ToJSON() ([]byte, error) {
	return json.Marshal(si.store)
}
