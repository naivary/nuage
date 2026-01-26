package main

import "time"

type (
	Named    string
	PtrNamed *string
)

type UserID = int64

type QueryDeepObject struct {
	F1 string `query:"f1"`
	F2 int    `query:"f2"`
}

type QueryParamRequest struct {
	T             time.Time        `query:"t"`
	A             UserID           `query:"a"`
	N             Named            `query:"n"`
	PtrN          PtrNamed         `query:"ptr_n"`
	Deep          *QueryDeepObject `query:"deep"`
	DeepAnonymous struct{}         `query:"deep_anon"`
}
