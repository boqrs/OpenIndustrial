package codec

// Endian represents the byte order for multi-byte values.
type Endian string

const (
	BigEndian    Endian = "big_endian"
	LittleEndian Endian = "little_endian"
)

// WordOrder represents the word order for 32-bit and 64-bit values composed of multiple 16-bit words.
type WordOrder string

const (
	HighWordFirst WordOrder = "high_word_first" // e.g., [HW LW]
	LowWordFirst  WordOrder = "low_word_first"  // e.g., [LW HW]
)