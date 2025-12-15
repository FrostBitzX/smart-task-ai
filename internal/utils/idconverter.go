package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

// ShortUUIDWithPrefix convert UUID to short UUID with prefix
func ShortUUIDWithPrefix(id uuid.UUID, prefix string) string {
	shortID := shortuuid.DefaultEncoder.Encode(id)
	return fmt.Sprintf("%s_%s", prefix, shortID)
}

// ParseShortUUID convert hort UUID with prefix to UUID
func ParseShortUUID(shortID string, prefix string) (uuid.UUID, error) {
	prefixLen := len(prefix) + 1
	if len(shortID) < prefixLen {
		return uuid.Nil, fmt.Errorf("invalid short UUID format")
	}
	
	actualShortID := shortID[prefixLen:]
	return shortuuid.DefaultEncoder.Decode(actualShortID)
}
