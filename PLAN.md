# Project: Terminal P2P Chat (Go)

## High-Level Goals

* Go-only implementation
* No servers
* No backend
* No account registration
* No central authority
* Terminal UI (TUI)
* Linux/macOS support
* Windows support where platform-independent
* End-to-end encrypted
* Peer-to-peer networking
* Feature-rich collaboration platform
* Single binary
* Minimal configuration
* Portable

---

# Core Principles

* Every client is an equal peer.
* Identity is local.
* Data ownership is local.
* Encryption everywhere.
* Offline-first local storage.
* Automatic peer discovery when possible.
* Manual bootstrap always available.
* No cloud dependency.
* No telemetry.
* No analytics.

---

# Suggested Stack

## Language

* Go 1.24+

## TUI

* Bubble Tea
* Lip Gloss
* Bubbles

## Networking

* libp2p

Features:

* QUIC
* TCP
* NAT traversal
* Hole punching
* mDNS
* DHT
* GossipSub
* Relay v2

## Cryptography

* X25519
* Ed25519
* ChaCha20-Poly1305
* HKDF
* SHA-256
* Argon2id

## Storage

SQLite

Tables:

* identities
* peers
* organizations
* channels
* memberships
* messages
* attachments
* invites
* reactions
* settings
* sessions

---

# Project Structure

```
cmd/
    chat/

internal/

    app/
    tui/
    network/
    crypto/
    storage/
    discovery/
    protocol/
    organization/
    channel/
    peer/
    sync/
    file/
    notification/
    config/
    commands/

pkg/

assets/

configs/
```

---

# Identity

Generated automatically on first launch.

Contains:

* display name
* unique peer id
* Ed25519 keypair
* avatar color
* creation timestamp

Stored locally.

Rename anytime.

---

# Startup Flow

```
launch

↓

load config

↓

load identity

↓

start networking

↓

discover peers

↓

sync

↓

enter TUI
```

---

# Peer Discovery

Support all:

## Local

* mDNS

## Internet

* DHT

## Bootstrap peers

User-configurable.

## Manual

```
/connect multiaddr
```

## QR code

Generate connect code.

## Invite code

Single string.

---

# Connection

Peer connects using:

```
libp2p multiaddress
```

Automatic:

* reconnect
* retry
* session restore

---

# Authentication

Mutual identity verification.

Trust-first-use.

Fingerprint comparison.

Commands

```
/fingerprint

/verify

/untrust
```

---

# Encryption

Every transport encrypted.

Every message encrypted.

Forward secrecy.

Unique session keys.

Key rotation.

---

# Organizations

Unlimited organizations.

Fields

* id
* name
* description
* icon
* owner
* created
* updated

Permissions

* owner
* admin
* moderator
* member
* guest

Operations

Create

Join

Leave

Delete

Rename

Archive

Restore

Clone

Export

Import

---

# Organization Settings

Name

Description

Visibility

Private

Invite only

Public

Max members

History retention

Attachment limit

Pinned channels

Roles

Invites

---

# Channels

Unlimited.

Types

Text

Announcements

Voice placeholder

Private

Public

Read-only

Categories.

Nested categories.

Mute.

Archive.

Favorite.

Pin.

---

# Channel Settings

Topic

Description

Permissions

Retention

Slow mode

Read-only

History sync

Attachment limit

Emoji reactions

---

# Direct Messages

Peer-to-peer.

Encrypted.

Typing indicator.

Read receipts optional.

Pinned messages.

Search.

Export.

---

# Group Chats

Independent from organizations.

Temporary.

Permanent.

Invite-only.

---

# Invitations

Invite peer.

Invite organization.

Invite channel.

Invite group.

Expiration.

One-time.

Unlimited.

Password protected.

---

# Membership

Invite

Accept

Decline

Kick

Ban

Unban

Promote

Demote

Transfer ownership

---

# Roles

Owner

Admin

Moderator

Member

Guest

Custom roles.

Permission matrix.

---

# Permissions

Read

Write

Delete

Edit

Pin

Invite

Manage users

Manage channels

Manage organization

Manage roles

Manage settings

Upload

Download

---

# Messages

Support

Text

Markdown

Code blocks

Quotes

Lists

Tables

Inline formatting

Links

Mentions

Replies

Threads

Edits

Delete

Undo delete

Pin

Copy

Forward

Search

Bookmarks

Favorites

Jump to message

Timestamp

Sender fingerprint

Delivery state

---

# Reactions

Emoji

Custom emoji

Remove reaction

Reaction count

---

# Attachments

Images

Videos

Documents

Archives

Source code

Audio

Drag-drop

Clipboard

Progress

Resume

Integrity verification

Hash verification

Deduplication

---

# File Transfer

P2P.

Encrypted.

Chunked.

Resume.

Pause.

Cancel.

Integrity verification.

Large file support.

---

# Message Sync

Automatic.

Peer history reconciliation.

Missing message recovery.

Conflict resolution.

Duplicate detection.

Incremental sync.

---

# Offline Support

Queue outgoing messages.

Replay later.

Automatic sync.

---

# Search

Global.

Organization.

Channel.

DM.

Message.

Attachment.

User.

Date.

Regex optional.

---

# Notifications

Desktop.

Terminal bell.

Badge.

Mention.

DM.

Invite.

Custom rules.

---

# Presence

Online

Away

Busy

Offline

Invisible

Typing

Idle

Last seen optional

---

# Profile

Display name

Bio

Avatar color

Status

Custom status

Timezone

Pronouns optional

---

# Contacts

Favorite

Blocked

Muted

Trusted

Notes

Tags

---

# Blocking

Ignore messages.

Ignore invites.

Ignore files.

Ignore presence.

---

# TUI Layout

```
+------------------------------------------------------+

Organizations

Channels

Members

Chat

Input

Status

Command palette

Notifications

Logs

```

Panels resizable.

Keyboard driven.

Mouse optional.

---

# Keyboard Shortcuts

Navigation

Search

Switch organization

Switch channel

Open DM

Reply

Edit

Delete

Upload

Download

Toggle sidebar

Jump

Command palette

---

# Commands

```
/help

/connect

/disconnect

/peer

/peers

/org

/org create

/org join

/org leave

/channel

/channel create

/channel delete

/channel archive

/dm

/msg

/invite

/export

/import

/profile

/status

/search

/block

/unblock

/pin

/unpin

/file

/settings

/quit
```

---

# Synchronization

Organization metadata.

Channels.

Permissions.

Messages.

Files.

Pins.

Reactions.

Profiles.

Incremental state updates.

---

# Storage

Local SQLite.

Encrypted database optional.

Automatic migration.

Indexes.

Garbage collection.

Compaction.

---

# Backup

Identity export.

Organization export.

Full export.

Encrypted backup.

Restore.

---

# Import

Identity.

Organizations.

Settings.

History.

---

# Configuration

YAML.

TOML.

JSON.

CLI overrides.

Environment variables.

---

# Logging

Levels

Debug

Info

Warn

Error

Trace

Rotating logs.

---

# Security

Encrypted transport.

Encrypted payload.

Authenticated peers.

Identity fingerprints.

Replay protection.

Nonce verification.

Key rotation.

Rate limiting.

Spam protection.

Message signing.

Secure random generation.

Encrypted local database option.

Memory zeroization where possible.

Secure file deletion option.

---

# Networking

QUIC preferred.

TCP fallback.

Relay fallback.

Hole punching.

NAT traversal.

IPv4.

IPv6.

LAN discovery.

Internet discovery.

Bandwidth adaptation.

Compression.

Connection pooling.

---

# Reliability

Automatic reconnect.

Heartbeat.

Health checks.

Timeout recovery.

Retry logic.

Peer scoring.

Duplicate suppression.

---

# CLI

```
chat

chat serve

chat connect

chat export

chat import

chat identity

chat config

chat reset

chat doctor

chat version
```

---

# Configuration Example

```
identity

network

bootstrap peers

relay peers

database

downloads

uploads

theme

notifications

logging

security

key rotation

history retention
```

---

# Themes

Dark

Light

Nord

Gruvbox

Catppuccin

Solarized

Custom themes.

---

# Plugin System (Future)

Commands.

Themes.

Bots.

Integrations.

Custom panels.

Hooks.

---

# APIs (Future)

Local IPC.

Unix socket.

Named pipe.

JSON-RPC.

Plugin API.

---

# Testing

Unit tests.

Integration tests.

Network simulation.

NAT simulation.

Message sync tests.

Large file transfer tests.

Cross-platform tests.

Load testing.

Fuzz testing.

---

# Cross-Platform Targets

Linux amd64

Linux arm64

macOS Intel

macOS Apple Silicon

Windows amd64

Windows arm64

---

# Release Pipeline

```
lint

↓

unit tests

↓

integration tests

↓

cross compilation

↓

binary signing

↓

release archives

↓

checksums

↓

GitHub release
```

---

# Milestone 1 — Foundation

* Project structure
* Configuration
* Identity generation
* SQLite storage
* Logging
* Basic TUI
* libp2p initialization
* Peer discovery
* Manual connections
* Basic encrypted messaging

---

# Milestone 2 — Chat

* DMs
* Chat history
* Read receipts
* Typing indicator
* Search
* Notifications
* Message editing
* Reactions
* File transfer

---

# Milestone 3 — Organizations

* Organization creation
* Membership
* Roles
* Permissions
* Channel management
* Invites
* Moderation

---

# Milestone 4 — Synchronization

* Multi-peer state synchronization
* Offline message queue
* Conflict resolution
* Incremental synchronization
* History recovery
* Attachment synchronization

---

# Milestone 5 — Advanced Collaboration

* Threads
* Bookmarks
* Pins
* Favorites
* Categories
* Custom roles
* Rich Markdown
* Command palette
* Themes
* Encrypted backups

---

# Milestone 6 — Polish

* Performance optimization
* Memory optimization
* UI refinement
* Cross-platform validation
* Security audit
* Stress testing
* Documentation
* Packaging
* Automated releases

