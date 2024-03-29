package goengine

import (
	"net/http"
)

type filters_t struct {
	filter         []FilterFunc
}

func (this *filters_t) Use(handle FilterFunc) {
	this.filter = append(this.filter, handle)
}

// @return stop
func (this *filters_t) Range(res http.ResponseWriter, req *http.Request) bool {
	for i := range this.filter {
		if !this.filter[i](res, req) {
			return true
		}
	}
	return false
}
