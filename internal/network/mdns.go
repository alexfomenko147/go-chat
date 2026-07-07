package network

import (
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

	ctx := m.node.Ctx
	if err := m.node.Host.Connect(ctx, pi); err != nil {
		m.node.Logger.Debug("mDNS connect to %s: %v", pi.ID.String(), err)
		return
	}
	m.node.Logger.Info("mDNS connected to: %s", pi.ID.String())
}
