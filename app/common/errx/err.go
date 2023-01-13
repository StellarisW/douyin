package errx

import (
	"math"
	"reflect"
)

const (
	PkgName = "douyin/app/common/errx"
)

type (
	Error interface {
		Code() uint32
		Error() string
	}
)

type myErr struct {
	c uint32 // 错误编号, 即错误码最后一级
	s string
}

func New(code uint32, text string) Error {
	return &myErr{
		c: code,
		s: text,
	}
}

func (e *myErr) Code() uint32 {
	return e.c
}

func (e *myErr) Error() string {
	return e.s
}

func Convert(err error) Error {
	if e, ok := err.(Error); ok {
		return e
	}

	return New(0, err.Error())
}

func GetCode(err error) uint32 {
	if err == nil {
		return math.MaxUint32
	}

	rv := reflect.ValueOf(err)

	rvk := rv.Kind()
	for rvk == reflect.Ptr {
		rv = rv.Elem()
		rvk = rv.Kind()
	}

	rvt := rv.Type()

	if rvt == reflect.TypeOf(myErr{}) {
		nf := rv.NumField()
		for i := 0; i < nf; i++ {
			rvtf := rvt.Field(i)
			rvf := rv.Field(i)

			if rvtf.PkgPath == PkgName && rvtf.Name == "c" && rvf.Kind() == reflect.Uint32 {
				if rvf.CanUint() {
					return uint32(rvf.Uint())
				}
				return math.MaxUint32
			}
		}

	}

	return math.MaxUint32
}
