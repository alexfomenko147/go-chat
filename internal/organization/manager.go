package organization

import (
	"fmt"
	"time"

	"go-chat/internal/storage"
)

type Manager struct {
	store *storage.Store
}

func NewManager(store *storage.Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) CreateOrganization(name, description, ownerPeerID string) (*storage.Organization, error) {
	orgID := fmt.Sprintf("org_%d", time.Now().UnixNano())
	now := time.Now().UTC()

	org := &storage.Organization{
		OrgID:       orgID,
		Name:        name,
		Description: description,
		OwnerPeerID: ownerPeerID,
		Visibility:  "private",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := m.store.SaveOrganization(org); err != nil {
		return nil, fmt.Errorf("save org: %w", err)
	}

	membership := &storage.Membership{
		PeerID:   ownerPeerID,
		OrgID:    orgID,
		Role:     "owner",
		JoinedAt: now,
	}

	if err := m.store.SaveMembership(membership); err != nil {
		return nil, fmt.Errorf("save membership: %w", err)
	}

	return org, nil
}

func (m *Manager) GetOrganization(orgID string) (*storage.Organization, error) {
	return m.store.GetOrganization(orgID)
}

func (m *Manager) ListOrganizations() ([]*storage.Organization, error) {
	return m.store.ListOrganizations()
}

func (m *Manager) DeleteOrganization(orgID string) error {
	return m.store.DeleteOrganization(orgID)
}
