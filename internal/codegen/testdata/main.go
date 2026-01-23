package main

type Int int

type PathParamRequest struct {
	Str string `path:"str"`

	Named Int   `path:"named"`
	I     int   `path:"int"`
	I8    int8  `path:"i8"`
	I16   int16 `path:"i16"`
	I32   int32 `path:"i32"`
	I64   int64 `path:"i64"`

	U8  uint8  `path:"u8"`
	U16 uint16 `path:"u16"`
	U32 uint32 `path:"u32"`
	U64 uint64 `path:"u64"`

	F32 float32 `path:"f32"`
	F64 float64 `path:"f64"`
}

type HeaderParamRequest struct {
	Str string `header:"str"`

	Named Int   `header:"named"`
	I     int   `header:"int"`
	I8    int8  `header:"i8"`
	I16   int16 `header:"i16"`
	I32   int32 `header:"i32"`
	I64   int64 `header:"i64"`

	U8  uint8  `header:"u8"`
	U16 uint16 `header:"u16"`
	U32 uint32 `header:"u32"`
	U64 uint64 `header:"u64"`

	F32 float32 `header:"f32"`
	F64 float64 `header:"f64"`
}
