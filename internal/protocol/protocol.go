package protocol

import (
	"encoding/json"
	"fmt"
	"time"

	"go-chat/internal/crypto"
)

type Envelope struct {
	Type          string `json:"type"`
	SenderID      string `json:"sender_id"`
	ChannelID     string `json:"channel_id,omitempty"`
	MessageID     string `json:"message_id,omitempty"`
	EncryptedData []byte `json:"encrypted_data,omitempty"`
	Nonce         []byte `json:"nonce,omitempty"`
	Timestamp     int64  `json:"timestamp"`
}

type PlaintextPayload struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	ReplyTo     string `json:"reply_to,omitempty"`
}

type EncryptedMessage struct {
	Envelope
}

func NewEncryptedMessage(senderID, channelID, messageID, content, contentType, replyTo string, cipher *crypto.Cipher) (*Envelope, error) {
	payload := PlaintextPayload{
		Content:     content,
		ContentType: contentType,
		ReplyTo:     replyTo,
	}

	plaintext, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	encrypted, err := cipher.Encrypt(plaintext)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	return &Envelope{
		Type:          "message",
		SenderID:      senderID,
		ChannelID:     channelID,
		MessageID:     messageID,
		EncryptedData: encrypted,
		Timestamp:     time.Now().UnixMilli(),
	}, nil
}

func DecryptEnvelope(env *Envelope, cipher *crypto.Cipher) (*PlaintextPayload, error) {
	plaintext, err := cipher.Decrypt(env.EncryptedData)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	var payload PlaintextPayload
	if err := json.Unmarshal(plaintext, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	return &payload, nil
}
