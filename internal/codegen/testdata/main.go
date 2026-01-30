package main

type (
	NamedInt int
	PtrNamed *int
	String   string
)

type PathParamRequest struct {
	Int      int       `path:"integer"`
	PtrInt   *int      `path:"ptr_integer"`
	Named    NamedInt  `path:"named"`
	PtrNamed *PtrNamed `path:"ptr_named"`
	Str      string    `path:"str"`
	StrNamed String    `path:"str_named"`
}
