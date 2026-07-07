package peermgr

import (
	"fmt"
	"time"

	"go-chat/internal/crypto"
	"go-chat/internal/storage"
)

type Manager struct {
	store *storage.Store
}

func NewManager(store *storage.Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) GetOrCreatePeer(peerID string, publicKey []byte) (*storage.Peer, error) {
	p, err := m.store.GetPeer(peerID)
	if err != nil {
		return nil, err
	}
	if p != nil {
		if publicKey != nil && len(p.PublicKey) == 0 {
			p.PublicKey = publicKey
			p.Fingerprint = crypto.Fingerprint(publicKey)
			if err := m.store.SavePeer(p); err != nil {
				return nil, err
			}
		}
		return p, nil
	}

	fingerprint := ""
	if len(publicKey) > 0 {
		fingerprint = crypto.Fingerprint(publicKey)
	}

	now := time.Now().UTC()
	p = &storage.Peer{
		PeerID:      peerID,
		DisplayName: peerID,
		PublicKey:   publicKey,
		Fingerprint: fingerprint,
		AvatarColor: "#5865F2",
		Status:      "online",
		FirstSeen:   now,
		LastSeen:    now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := m.store.SavePeer(p); err != nil {
		return nil, fmt.Errorf("save peer: %w", err)
	}

	return p, nil
}

func (m *Manager) GetPeer(peerID string) (*storage.Peer, error) {
	return m.store.GetPeer(peerID)
}

func (m *Manager) ListPeers() ([]*storage.Peer, error) {
	return m.store.ListPeers()
}

func (m *Manager) UpdateStatus(peerID, status string) error {
	return m.store.UpdatePeerStatus(peerID, status)
}
