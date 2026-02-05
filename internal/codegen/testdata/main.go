package main

type (
	String    string
	Int32     int32
	PtrString *string
	PtrInt32  *int32
)

type PathParamRequest struct {
	Str                  string      `path:"str"`
	PtrStr               *string     `path:"ptr_str"`
	Int                  int         `path:"int"`
	Int32                int32       `path:"int_32"`
	Int64                int64       `path:"int_64"`
	PtrInt64             *int64      `path:"ptr_int_64"`
	NamedStr             String      `path:"named_str"`
	NamedInt32           Int32       `path:"named_int32"`
	PtrString            PtrString   `path:"ptr_string"`
	PtrI32               PtrInt32    `path:"ptr_int32"`
	Slice                []string    `path:"slice"`
	SliceNamedElem       []Int32     `path:"slice_named_elem"`
	SlicePtrElem         []*int      `path:"slice_ptr_elem"`
	PtrNamedPtrString    *PtrString  `path:"ptr_named_ptr_str"`
	PtrNamedPtrInt32     *PtrInt32   `path:"ptr_named_ptr_int32"`
	SlicePtrNamedPtrElem []*PtrInt32 `path:"slice_ptr_named_ptr_elem"`
	Uint                 uint        `path:"uint"`
}

type QueryParamRequest struct {
	Str              string             `query:"str"`
	Bool             bool               `query:"boolean"`
	Int32            int32              `query:"int_32"`
	Int64            int64              `query:"int_64"`
	PtrInt64         *int64             `query:"ptr_int_64"`
	Uint32           uint32             `query:"uint_32"`
	Uint64           uint64             `query:"uint_64"`
	PtrUint64        *uint64            `query:"ptr_uint_64"`
	SliceString      []string           `query:"slice_string"`
	SliceInt         []int              `query:"slice_int"`
	Map              map[string]string  `query:"object"`
	MapNamed         map[String]string  `query:"named_object"`
	MapPtrKey        map[*string]string `query:"ptr_key_object"`
	MapPtrKeyExplode map[*string]string `query:"ptr_key_explode,explode=false"`
	DeepObject       struct {
		S1      string
		Int     int
		Int64   int64
		Ptr     *string
		PtrInt  *int32
		Invalid []string
	} `query:"deepObject,style=deepObject"`
}
