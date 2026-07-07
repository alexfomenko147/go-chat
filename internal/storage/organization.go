package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type Organization struct {
	ID               int64     `json:"id"`
	OrgID            string    `json:"org_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Icon             string    `json:"icon"`
	OwnerPeerID      string    `json:"owner_peer_id"`
	Visibility       string    `json:"visibility"`
	MaxMembers       int       `json:"max_members"`
	HistoryRetention int       `json:"history_retention"`
	AttachmentLimit  int64     `json:"attachment_limit"`
	Archived         bool      `json:"archived"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Membership struct {
	ID       int64     `json:"id"`
	PeerID   string    `json:"peer_id"`
	OrgID    string    `json:"org_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

func (s *Store) SaveOrganization(org *Organization) error {
	_, err := s.db.Exec(`INSERT INTO organizations (org_id, name, description, icon, owner_peer_id, visibility, max_members, history_retention, attachment_limit, archived, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(org_id) DO UPDATE SET
			name=excluded.name, description=excluded.description, icon=excluded.icon,
			visibility=excluded.visibility, max_members=excluded.max_members,
			history_retention=excluded.history_retention, attachment_limit=excluded.attachment_limit,
			archived=excluded.archived, updated_at=excluded.updated_at`,
		org.OrgID, org.Name, org.Description, org.Icon, org.OwnerPeerID, org.Visibility,
		org.MaxMembers, org.HistoryRetention, org.AttachmentLimit, org.Archived, org.CreatedAt, org.UpdatedAt)
	if err != nil {
		return fmt.Errorf("save org: %w", err)
	}
	return nil
}

func (s *Store) GetOrganization(orgID string) (*Organization, error) {
	org := &Organization{}
	err := s.db.QueryRow(`SELECT id, org_id, name, description, icon, owner_peer_id, visibility, max_members, history_retention, attachment_limit, archived, created_at, updated_at
		FROM organizations WHERE org_id=?`, orgID).Scan(
		&org.ID, &org.OrgID, &org.Name, &org.Description, &org.Icon, &org.OwnerPeerID,
		&org.Visibility, &org.MaxMembers, &org.HistoryRetention, &org.AttachmentLimit,
		&org.Archived, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get org: %w", err)
	}
	return org, nil
}

func (s *Store) ListOrganizations() ([]*Organization, error) {
	rows, err := s.db.Query(`SELECT id, org_id, name, description, icon, owner_peer_id, visibility, max_members, history_retention, attachment_limit, archived, created_at, updated_at FROM organizations ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list orgs: %w", err)
	}
	defer rows.Close()

	var orgs []*Organization
	for rows.Next() {
		org := &Organization{}
		if err := rows.Scan(&org.ID, &org.OrgID, &org.Name, &org.Description, &org.Icon, &org.OwnerPeerID,
			&org.Visibility, &org.MaxMembers, &org.HistoryRetention, &org.AttachmentLimit,
			&org.Archived, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan org: %w", err)
		}
		orgs = append(orgs, org)
	}
	return orgs, rows.Err()
}

func (s *Store) DeleteOrganization(orgID string) error {
	_, err := s.db.Exec(`DELETE FROM organizations WHERE org_id=?`, orgID)
	return err
}

func (s *Store) SaveMembership(m *Membership) error {
	_, err := s.db.Exec(`INSERT INTO memberships (peer_id, org_id, role, joined_at) VALUES (?, ?, ?, ?)
		ON CONFLICT(peer_id, org_id) DO UPDATE SET role=excluded.role`,
		m.PeerID, m.OrgID, m.Role, m.JoinedAt)
	if err != nil {
		return fmt.Errorf("save membership: %w", err)
	}
	return nil
}

func (s *Store) GetMembership(peerID, orgID string) (*Membership, error) {
	m := &Membership{}
	err := s.db.QueryRow(`SELECT id, peer_id, org_id, role, joined_at FROM memberships WHERE peer_id=? AND org_id=?`,
		peerID, orgID).Scan(&m.ID, &m.PeerID, &m.OrgID, &m.Role, &m.JoinedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get membership: %w", err)
	}
	return m, nil
}

func (s *Store) ListMembers(orgID string) ([]*Membership, error) {
	rows, err := s.db.Query(`SELECT id, peer_id, org_id, role, joined_at FROM memberships WHERE org_id=?`, orgID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	var members []*Membership
	for rows.Next() {
		m := &Membership{}
		if err := rows.Scan(&m.ID, &m.PeerID, &m.OrgID, &m.Role, &m.JoinedAt); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		members = append(members, m)
	}
	return members, rows.Err()
}
