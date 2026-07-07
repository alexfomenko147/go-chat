# go-chat

**Terminal P2P Chat** — A serverless, encrypted, peer-to-peer messaging platform that runs entirely in your terminal.

```
                   _                    
                  | |                   
   ___   __ _  ___| |__   ___  ___  ___ 
  / __| / _` |/ __| '_ \ / _ \/ __|/ __|
 | (__ | (_| | (__| | | | (_) \__ \ (__ 
  \___| \__, |\___|_| |_|\___/|___/\___|
         __/ |                          
        |___/                           
```

## Features

- **No servers, no accounts** — Pure P2P. Every client is an equal peer.
- **End-to-end encrypted** — X25519 key exchange, ChaCha20-Poly1305 AEAD, HKDF session keys.
- **Local-first** — All data stored in SQLite. You own everything.
- **Terminal UI** — Built with Bubble Tea, Lip Gloss, and Bubbles.
- **libp2p networking** — TCP, QUIC, mDNS discovery, relay, NAT traversal, hole punching.
- **Organizations & Channels** — Create orgs, channels, manage members and permissions.
- **Direct Messages** — Encrypted peer-to-peer conversations.
- **Automatic peer discovery** — mDNS for LAN, DHT for internet, bootstrap peers.
- **Ed25519 identities** — Generated locally on first launch. Fingerprint-based verification.
- **Offline-first** — Messages queued locally, synced when peers reconnect.
- **Cross-platform** — Linux amd64/arm64, macOS Intel/Apple Silicon, Windows amd64/arm64.

## Quick Start

### Prerequisites

- Go 1.24+ (`/usr/local/go/bin/go`)

### Build

```bash
go build -o chat ./cmd/chat
```

### Run

```bash
./chat
```

On first launch, go-chat automatically:
1. Generates an Ed25519 identity keypair
2. Derives your libp2p Peer ID
3. Creates a local SQLite database
4. Starts listening on a random TCP port
5. Discovers peers via mDNS on the local network

## Configuration

Config is loaded from `~/.config/go-chat/config.yaml` (created automatically with defaults).

### Example config.yaml

```yaml
identity:
  display_name: ""             # defaults to user-<hostname>

network:
  port: 0                      # 0 = random port
  enable_relay: true
  enable_mdns: true
  enable_dht: true
  enable_quic: true
  enable_tcp: true
  bootstrap_peers: []          # multiaddrs for bootstrap
  relay_peers: []

database:
  path: chat.db
  encrypt: false

downloads:
  path: downloads
  max_size: 104857600          # 100 MB

uploads:
  max_size: 52428800           # 50 MB

theme: dark
notifications:
  desktop: false
  bell: true
  mentions: true

logging:
  level: info                  # trace, debug, info, warn, error
  file: ""
  rotate: true

security:
  key_rotation_days: 30
  encrypt_database: false
```

Override with CLI flags:

```bash
./chat --config /path/to/config.yaml
./chat --version
```

## Usage

### TUI Layout

```
+------------------+------------------+---------------------------+
| Organizations    | Channels         | Chat                      |
|                  |                  |                           |
|  > My Org        |  > # general     |  12:30 alice: hello       |
|                  |    # random      |  12:31 bob: hi            |
|                  |                  |                           |
+------------------+------------------+                           |
| Status Bar       |                  |                           |
+------------------+------------------+---------------------------+
| Input / Command Line                                              |
+------------------------------------------------------------------+
```

### Navigation

| Key | Action |
|-----|--------|
| `Tab` | Toggle between input and navigation mode |
| `Up` / `Down` | Navigate channels or scroll chat |
| `Left` / `Right` | Switch organizations |
| `Enter` | Send message (input mode) |
| `?` | Toggle help |
| `P` | Toggle peers list |
| `Ctrl+C` / `Ctrl+Q` | Quit |

### Commands

| Command | Description |
|---------|-------------|
| `/help` | Show help |
| `/connect <multiaddr>` | Connect to a peer manually |
| `/disconnect` | Disconnect all peers |
| `/peers` | List known peers |
| `/org create <name>` | Create an organization |
| `/channel create <name>` | Create a channel |
| `/dm <peer_id>` | Open a direct message |
| `/profile` | Show your identity info |
| `/quit` | Exit |

### Connecting to Peers

#### LAN (automatic)

go-chat discovers peers on the same local network via mDNS automatically.

#### Manual

```bash
# Get your address (shown on startup or in logs)
/ip4/192.168.1.42/tcp/43987/p2p/12D3KooW...

# Connect to a peer
/connect /ip4/192.168.1.100/tcp/43987/p2p/12D3KooW...
```

#### Bootstrap

Add bootstrap peer multiaddrs to `network.bootstrap_peers` in config.

## Identity

Your identity is automatically generated on first launch:

- **Display name** — Configurable via config or future /nick command
- **Peer ID** — Derived from your Ed25519 public key (e.g., `12D3KooW...`)
- **Keypair** — Ed25519 for identity, X25519 for key exchange
- **Fingerprint** — First 16 bytes of SHA-256 of your public key

### View Identity

```
/profile
```

Example output:
```
Profile: user-myhost | PeerID: 12D3KooWAbCdEf... | Fingerprint: a1b2c3d4e5f6...
```

## Architecture

```
cmd/chat/            CLI entrypoint
internal/
  app/               Application coordinator
  config/            Configuration (YAML/JSON)
  logging/           Logging (levels, rotation)
  crypto/            X25519, Ed25519, ChaCha20-Poly1305, HKDF
  storage/           SQLite database layer
  network/           libp2p host, streams, mDNS
  protocol/          Message envelopes, encryption
  discovery/         Peer discovery (bootstrap, DHT)
  peer/              Peer management
  organization/      Organization CRUD
  channel/           Channel CRUD
  tui/               Bubble Tea terminal UI
  sync/              State synchronization
  file/              File transfer
  notification/      Desktop/bell notifications
  commands/          Command registry
```

## Storage

Local SQLite database (`chat.db`) with the following tables:

- `identities` — Your Ed25519 keypair and profile
- `peers` — Known peers, keys, trust status
- `organizations` — Org metadata and settings
- `channels` — Channel config and metadata
- `memberships` — Peer-org role mappings
- `messages` — Chat history (encrypted content)
- `attachments` — File transfer metadata
- `invites` — Pending invitations
- `reactions` — Emoji reactions on messages
- `settings` — Key-value settings store
- `sessions` — Cryptographic session keys

## Security

- **Transport** — Encrypted via libp2p's Noise protocol
- **Messages** — Encrypted with ChaCha20-Poly1305 using session keys derived from X25519 + HKDF
- **Identity** — Ed25519 keypairs, verified via fingerprints
- **Forward secrecy** — Unique session keys with rotation
- **Replay protection** — Timestamps and nonce verification
- **No telemetry, no analytics, no cloud**

## Development

### Project Structure

```
cmd/chat/main.go          Entrypoint
internal/app/             Core application logic
internal/tui/             Bubble Tea model, styles, panels
internal/network/         libp2p host setup, stream handler
internal/crypto/          Cryptographic primitives
internal/storage/         SQLite schema, queries, migrations
internal/protocol/        Message encoding and encryption
internal/config/          Config parsing (YAML, JSON)
internal/logging/         Structured logging
internal/discovery/       Peer discovery (mDNS, DHT, bootstrap)
internal/organization/    Organization management
internal/channel/         Channel management
internal/peermgr/         Peer management
internal/sync/            State synchronization
internal/file/            File transfer
internal/notification/    Notifications
internal/commands/        Command registry
```

### Build

```bash
go build -o chat ./cmd/chat
```

### Cross-compile

```bash
GOOS=linux GOARCH=amd64 go build -o chat-linux-amd64 ./cmd/chat
GOOS=darwin GOARCH=arm64 go build -o chat-darwin-arm64 ./cmd/chat
GOOS=windows GOARCH=amd64 go build -o chat-windows-amd64.exe ./cmd/chat
```

## Milestones

| # | Area | Status |
|---|------|--------|
| 1 | Foundation — project structure, config, identity, SQLite, libp2p, TUI | Done |
| 2 | Chat — DMs, history, typing, reactions, file transfer, notifications | In Progress |
| 3 | Organizations — membership, roles, permissions, invites, moderation | Planned |
| 4 | Synchronization — multi-peer state, offline queue, conflict resolution | Planned |
| 5 | Advanced — threads, bookmarks, themes, command palette, encrypted backups | Planned |
| 6 | Polish — performance, security audit, cross-platform validation, releases | Planned |

## License

MIT
