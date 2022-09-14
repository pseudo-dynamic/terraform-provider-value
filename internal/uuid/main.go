package uuid

import (
	"hash/crc64"
	"math/rand"

	"github.com/google/uuid"
)

func DeterministicUuidFromString(seed string) uuid.UUID {
	hash := crc64.New(crc64.MakeTable(crc64.ISO))
	hash.Write([]byte(seed))
	hashSumUnsigned := hash.Sum64()
	hashSumSigned := int64(hashSumUnsigned)
	randomGenerator := rand.New(rand.NewSource(hashSumSigned))
	deterministicUuid, _ := uuid.NewRandomFromReader(randomGenerator)
	return deterministicUuid
}
