package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

func (ut *Utils) GenerateUUID() string {
	s := uuid.NewString()
	return s
}

func (ut *Utils) GenerateShortUUID() string {
    b := make([]byte, 3)
    _, err := rand.Read(b)
    if err != nil {
        return "000000"

    }
    return hex.EncodeToString(b)
}