package nuage

const (
	// Bit represents the smallest unit of digital information.
	//
	// This constant is primarily provided for completeness and for
	// defining larger units in terms of bits and bytes.
	Bit int = 1

	// Byte represents a group of 8 bits.
	Byte = 8 * Bit

	// KiB represents a kibibyte (2¹⁰ bytes), as defined by IEC 80000-13.
	KiB = 1024 * Byte

	// MiB represents a mebibyte (2²⁰ bytes), as defined by IEC 80000-13.
	MiB = 1024 * KiB

	// GiB represents a gibibyte (2³⁰ bytes), as defined by IEC 80000-13.
	GiB = 1024 * MiB

	// TiB represents a tebibyte (2⁴⁰ bytes), as defined by IEC 80000-13.
	TiB = 1024 * GiB
)
