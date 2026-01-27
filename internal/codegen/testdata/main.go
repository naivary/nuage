package main

type PathParamRequest struct {
	Int    int  `path:"integer"`
	PtrInt *int `path:"ptr_integer"`
}
