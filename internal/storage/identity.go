package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type Identity struct {
	ID          int64     `json:"id"`
	DisplayName string    `json:"display_name"`
	PeerID      string    `json:"peer_id"`
	PrivateKey  []byte    `json:"-"`
	PublicKey   []byte    `json:"public_key"`
	AvatarColor string    `json:"avatar_color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Store) SaveIdentity(id *Identity) error {
	_, err := s.db.Exec(`INSERT INTO identities (display_name, peer_id, private_key, public_key, avatar_color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(peer_id) DO UPDATE SET
			display_name=excluded.display_name,
			private_key=excluded.private_key,
			public_key=excluded.public_key,
			avatar_color=excluded.avatar_color,
			updated_at=excluded.updated_at`,
		id.DisplayName, id.PeerID, id.PrivateKey, id.PublicKey, id.AvatarColor, id.CreatedAt, id.UpdatedAt)
	if err != nil {
		return fmt.Errorf("save identity: %w", err)
	}
	return nil
}

func (s *Store) GetIdentity() (*Identity, error) {
	id := &Identity{}
	err := s.db.QueryRow(`SELECT id, display_name, peer_id, private_key, public_key, avatar_color, created_at, updated_at
		FROM identities ORDER BY id LIMIT 1`).Scan(
		&id.ID, &id.DisplayName, &id.PeerID, &id.PrivateKey, &id.PublicKey, &id.AvatarColor, &id.CreatedAt, &id.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get identity: %w", err)
	}
	return id, nil
}

func (s *Store) UpdateIdentityDisplayName(peerID, displayName string) error {
	_, err := s.db.Exec(`UPDATE identities SET display_name=?, updated_at=? WHERE peer_id=?`,
		displayName, time.Now().UTC(), peerID)
	return err
}
