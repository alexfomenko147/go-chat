package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"go-chat/internal/storage"

	"github.com/libp2p/go-libp2p/core/network"
)

type StreamHandler struct {
	node *Node
}

func NewStreamHandler(node *Node) *StreamHandler {
	return &StreamHandler{node: node}
}

func (h *StreamHandler) Handle(s network.Stream) {
	defer s.Close()

	peerID := s.Conn().RemotePeer().String()
	r := bufio.NewReader(s)

	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				h.node.Logger.Debug("stream read error from %s: %v", peerID, err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			h.node.Logger.Debug("invalid message from %s: %v", peerID, err)
			continue
		}

		h.handleMessage(&msg)
	}
}

func (h *StreamHandler) handleMessage(msg *Message) {
	h.node.Logger.Info("received %s from %s", msg.Type, msg.SenderID)

	switch msg.Type {
	case "message":
		h.handleSyncMessage(msg)
	case "sync_org":
		h.handleSyncOrg(msg)
	case "sync_channel":
		h.handleSyncChannel(msg)
	default:
		h.node.Logger.Debug("unknown message type: %s", msg.Type)
	}
}

func (h *StreamHandler) handleSyncMessage(msg *Message) {
	if h.node.Store == nil {
		return
	}
	storeMsg := &storage.Message{
		MessageID:    msg.MessageID,
		ChannelID:    msg.ChannelID,
		SenderPeerID: msg.SenderID,
		Content:      msg.Content,
		ContentType:  msg.ContentType,
		DeliveryState: "received",
		CreatedAt:    time.UnixMilli(msg.Timestamp).UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := h.node.Store.SaveMessage(storeMsg); err != nil {
		h.node.Logger.Warn("save synced message: %v", err)
	}
}

func (h *StreamHandler) handleSyncOrg(msg *Message) {
	if h.node.Store == nil {
		return
	}
	existing, _ := h.node.Store.GetOrganization(msg.OrgID)
	if existing != nil {
		return
	}
	org := &storage.Organization{
		OrgID:     msg.OrgID,
		Name:      msg.Content,
		OwnerPeerID: msg.SenderID,
		CreatedAt: time.UnixMilli(msg.Timestamp).UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := h.node.Store.SaveOrganization(org); err != nil {
		h.node.Logger.Warn("save synced org: %v", err)
	}
}

func (h *StreamHandler) handleSyncChannel(msg *Message) {
	if h.node.Store == nil {
		return
	}
	existing, _ := h.node.Store.GetChannel(msg.ChannelID)
	if existing != nil {
		return
	}
	ch := &storage.Channel{
		ChannelID:   msg.ChannelID,
		OrgID:       msg.OrgID,
		Name:        msg.Content,
		ChannelType: "text",
		CreatedAt:   time.UnixMilli(msg.Timestamp).UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if err := h.node.Store.SaveChannel(ch); err != nil {
		h.node.Logger.Warn("save synced channel: %v", err)
	}
}

func (h *StreamHandler) SendMessage(s network.Stream, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	data = append(data, '\n')

	if _, err := s.Write(data); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	return nil
}
