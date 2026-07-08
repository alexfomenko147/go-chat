package network

import (
	"context"
	"time"

	"go-chat/internal/storage"

	"github.com/libp2p/go-libp2p/core/peer"
)

type mdnsNotifee struct {
	node *Node
}

func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID == m.node.Host.ID() {
		return
	}
	m.node.Logger.Info("mDNS discovered peer: %s", pi.ID.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := m.node.Host.Connect(ctx, pi); err != nil {
		m.node.Logger.Debug("mDNS connect to %s: %v", pi.ID.String(), err)
		return
	}
	m.node.Logger.Info("mDNS connected to: %s", pi.ID.String())

	if m.node.Store != nil {
		existing, err := m.node.Store.GetPeer(pi.ID.String())
		if err != nil {
			m.node.Logger.Warn("mDNS get peer: %v", err)
		}
		if existing == nil {
			if err := m.node.Store.SavePeer(&storage.Peer{
				PeerID:      pi.ID.String(),
				DisplayName: pi.ID.String(),
				Status:      "online",
				LastSeen:    time.Now().UTC(),
			}); err != nil {
				m.node.Logger.Warn("mDNS save peer: %v", err)
			}
		} else {
			if err := m.node.Store.UpdatePeerStatus(pi.ID.String(), "online"); err != nil {
				m.node.Logger.Warn("mDNS update peer status: %v", err)
			}
		}
	}

	if err := m.node.SyncWithPeer(ctx, pi.ID); err != nil {
		m.node.Logger.Warn("mDNS sync with %s: %v", pi.ID.String(), err)
	}
}
