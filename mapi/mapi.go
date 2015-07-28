// Package mapi provides core api methods
package mapi

import (
	"errors"
	"fmt"
	"github.com/fuxiaohei/purine/log"
	"reflect"
	"runtime"
	"strings"
)

// api result
type Res struct {
	Status bool                   `json:"status"`
	Error  string                 `json:"error,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

// api func
type Func func(interface{}) *Res

// return success result
func Success(data map[string]interface{}) *Res {
	return &Res{
		Status: true,
		Data:   data,
	}
}

// return fail result
func Fail(err error) *Res {
	return &Res{
		Status: false,
		Data:   nil,
		Error:  err.Error(),
	}
}

// get function name
func funcName(fn Func) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	nameData := strings.Split(name, "/")
	if len(nameData) > 2 {
		nameData = nameData[len(nameData)-2:]
	}
	if runtime.GOOS == "windows" {
		return strings.TrimSuffix(strings.Join(nameData, "."), "-fm")
	}
	return strings.TrimSuffix(strings.Join(nameData, "."), "·fm")
}

// call api function with param
//
// usage:
//      mapi.Call(mapi.Article.Write,*ArticleForm)
//
func Call(fn Func, param interface{}) *Res {
	name := funcName(fn)
	log.Debug("Action | %-8s | %s", "Call", name)
	return fn(param)
}

// action param type error
func paramTypeError(v interface{}) error {
	return errors.New(fmt.Sprintf("need %s", reflect.TypeOf(v)))
}
