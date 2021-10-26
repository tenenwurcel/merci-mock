package mock

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
)

type Helper interface {
	NewMockHelper(i interface{}, failMethods []string) Helper
	GetForcedError() (err error)
}

type helper struct {
	failMap map[string]bool
}

func (h *helper) NewMockHelper(i interface{}, failMethods []string) Helper {
	h.failMap = map[string]bool{}
	e := reflect.TypeOf(i)

	for i := 0; i < e.NumMethod(); i++ {
		methodName := strings.ToLower(e.Method(i).Name)
		h.failMap[methodName] = false

		for _, failName := range failMethods {
			if methodName == failName {
				h.failMap[methodName] = true
			}
		}
	}

	return h
}

func (h *helper) getGrandParentCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	nameParts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	name := nameParts[len(nameParts) - 1]

	return strings.ToLower(name)
}

func (h *helper) GetForcedError() (err error) {
	if h.failMap[h.getGrandParentCallerName()] {
		err = errors.New("this is a mock error")
	}

	return
}