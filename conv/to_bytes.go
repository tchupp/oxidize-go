package conv

import "encoding/binary"

func U32ToBytes(num uint32) []byte {
	enc := make([]byte, 4)
	binary.BigEndian.PutUint32(enc, num)
	return enc
}

func U64ToBytes(num uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, num)
	return enc
}
