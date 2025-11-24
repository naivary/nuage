package nuage

type DataSize int64

const (
	KiB DataSize = 1
	MiB          = 1024 * KiB
	GiB          = 1024 * MiB
	TiB          = 1024 * GiB
)

func (d DataSize) Int() int {
	return int(d)
}
