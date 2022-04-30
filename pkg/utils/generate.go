package utils

import (
	"fmt"
	"hash/crc32"
)

func GenerateHash(url string) string {
	return fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(url)))
}
