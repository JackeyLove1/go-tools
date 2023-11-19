package hashs

import (
    "encoding/hex"

    "github.com/spaolacci/murmur3"
)

func MurmurHash(str string) string {
    hash := murmur3.New64()
    hash.Write([]byte(str))
    hashBytes := hash.Sum(nil)
    hashStringHex := hex.EncodeToString(hashBytes)
    return hashStringHex
}
