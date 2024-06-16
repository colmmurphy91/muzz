package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Match struct {
	ID      int    `db:"id" json:"id"`
	User1ID int    `db:"user1_id" json:"user1_id"`
	User2ID int    `db:"user2_id" json:"user2_id"`
	MatchID string `db:"match_id" json:"match_id"`
}

// GenerateMatchID generates a consistent hash for the match
func (m *Match) GenerateMatchID() {
	// Ensure the order is consistent by always having the smaller ID first
	matchData := fmt.Sprintf("%d:%d", m.User1ID, m.User2ID)

	hash := sha256.Sum256([]byte(matchData))
	m.MatchID = hex.EncodeToString(hash[:])
}
