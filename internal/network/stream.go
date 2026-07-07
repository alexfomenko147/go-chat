package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

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

		h.node.Logger.Debug("received message type %s from %s", msg.Type, peerID)
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
