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
}
