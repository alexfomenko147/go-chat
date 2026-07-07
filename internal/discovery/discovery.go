package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-chat/internal/config"
	"go-chat/internal/logging"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Discoverer struct {
	host   host.Host
	config *config.NetworkConfig
	logger *logging.Logger
	mu     sync.RWMutex
	peers  map[string]peer.AddrInfo
}

func New(h host.Host, cfg *config.NetworkConfig, log *logging.Logger) *Discoverer {
	return &Discoverer{
		host:   h,
		config: cfg,
		logger: log,
		peers:  make(map[string]peer.AddrInfo),
	}
}

func (d *Discoverer) DiscoveredPeers() []peer.AddrInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]peer.AddrInfo, 0, len(d.peers))
	for _, pi := range d.peers {
		result = append(result, pi)
	}
	return result
}

func (d *Discoverer) Bootstrap(ctx context.Context) error {
	if len(d.config.BootstrapPeers) == 0 {
		d.logger.Debug("no bootstrap peers configured")
		return nil
	}

	for _, addr := range d.config.BootstrapPeers {
		addrInfo, err := peer.AddrInfoFromString(addr)
		if err != nil {
			d.logger.Warn("invalid bootstrap addr %s: %v", addr, err)
			continue
		}
		pi := *addrInfo

		d.mu.Lock()
		d.peers[pi.ID.String()] = pi
		d.mu.Unlock()

		go func(pi peer.AddrInfo) {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			if err := d.host.Connect(ctx, pi); err != nil {
				d.logger.Debug("bootstrap connect %s: %v", pi.ID.String(), err)
				return
			}
			d.logger.Info("bootstrapped with: %s", pi.ID.String())
		}(pi)
	}

	return nil
}

func (d *Discoverer) Discover(ctx context.Context) error {
	if d.config.EnableDHT {
		return d.discoverDHT(ctx)
	}
	return nil
}

func (d *Discoverer) discoverDHT(ctx context.Context) error {
	return fmt.Errorf("DHT discovery not yet implemented")
}
