package tui

import (
	"fmt"
	"strings"

	"go-chat/internal/app"
	"go-chat/internal/crypto"
	"go-chat/internal/storage"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MessageItem struct {
	Sender    string
	Content   string
	Timestamp string
}

type focusPanel int

const (
	focusChannels focusPanel = iota
	focusOrgs
)

type Model struct {
	app          *app.App
	ready        bool
	width        int
	height       int
	firstLaunch  bool

	orgList      []*storage.Organization
	selectedOrg  int

	channelList  []*storage.Channel
	selectedChan int

	chatView     viewport.Model
	messages     []MessageItem

	input        textinput.Model
	inputMode    bool
	focus        focusPanel

	statusText   string
	statusLog    []string

	peerList     []*storage.Peer

	showHelp     bool
	showPeers    bool

	err          error
}

func NewModel(a *app.App) *Model {
	ti := textinput.New()
	ti.Placeholder = "Type a message or /help..."
	ti.Focus()
	ti.CharLimit = 2000
	ti.Width = 60

	return &Model{
		app:         a,
		input:       ti,
		inputMode:   true,
		focus:       focusChannels,
		firstLaunch: true,
		statusText:  fmt.Sprintf("PeerID: %s | /myaddr to see shareable address", a.PeerID()),
	}
}

func (m *Model) Init() tea.Cmd {
	m.loadOrgs()
	if len(m.orgList) > 0 {
		m.loadChannels()
	}
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		leftPanelWidth := 50
		inputHeight := 3
		statusHeight := 1
		chatHeight := m.height - inputHeight - statusHeight - 4

		if m.chatView.Height == 0 {
			m.chatView = viewport.New(m.width-leftPanelWidth-6, chatHeight)
		} else {
			m.chatView.Width = m.width - leftPanelWidth - 6
			m.chatView.Height = chatHeight
		}

		m.input.Width = m.width - 56

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		case "tab":
			if m.inputMode {
				m.inputMode = false
				m.focus = focusChannels
			} else {
				switch m.focus {
				case focusChannels:
					m.focus = focusOrgs
				case focusOrgs:
					m.focus = focusChannels
				}
			}

		case "up":
			if m.inputMode {
				m.chatView.LineUp(1)
			} else if m.selectedChan > 0 {
				m.selectedChan--
				m.loadMessages()
			}

		case "down":
			if m.inputMode {
				m.chatView.LineDown(1)
			} else if m.selectedChan < len(m.channelList)-1 {
				m.selectedChan++
				m.loadMessages()
			}

		case "left":
			if !m.inputMode && m.selectedOrg > 0 {
				m.selectedOrg--
				m.loadChannels()
			}

		case "right":
			if !m.inputMode && m.selectedOrg < len(m.orgList)-1 {
				m.selectedOrg++
				m.loadChannels()
			}

		case "enter":
			if m.inputMode {
				return m, m.handleInput()
			}
			m.inputMode = true

		case "?":
			m.showHelp = !m.showHelp

		case "P":
			m.showPeers = !m.showPeers
			if m.showPeers {
				m.loadPeers()
			}
		}
	}

	if m.inputMode {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.chatView, _ = m.chatView.Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleInput() tea.Cmd {
	text := strings.TrimSpace(m.input.Value())
	m.input.SetValue("")

	if text == "" {
		return nil
	}

	if strings.HasPrefix(text, "/") {
		return m.handleCommand(text)
	}

	if len(m.channelList) == 0 {
		m.addStatus("No channel selected. Use Tab to navigate and select a channel.")
		return nil
	}

	channelID := m.channelList[m.selectedChan].ChannelID
	if err := m.app.SendMessage(channelID, text, "text"); err != nil {
		m.addStatus(fmt.Sprintf("Error: %v", err))
		return nil
	}

	msg := MessageItem{
		Sender:    m.app.Identity().DisplayName,
		Content:   text,
		Timestamp: "now",
	}
	m.messages = append(m.messages, msg)
	m.chatView.SetContent(m.renderMessages())

	m.chatView.GotoBottom()

	return nil
}

func (m *Model) handleCommand(text string) tea.Cmd {
	parts := strings.Fields(text)
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "/help":
		m.showHelp = !m.showHelp

	case "/connect":
		if len(parts) < 2 {
			m.addStatus("Usage: /connect <multiaddr>")
			return nil
		}
		if err := m.app.Connect(parts[1]); err != nil {
			m.addStatus(fmt.Sprintf("Connect error: %v", err))
			return nil
		}
		m.addStatus(fmt.Sprintf("Connected to %s", parts[1]))

	case "/disconnect":
		m.app.DisconnectAll()
		m.addStatus("Disconnected")

	case "/peers":
		m.loadPeers()
		m.showPeers = !m.showPeers

	case "/org":
		if len(parts) < 2 {
			m.addStatus("Usage: /org create <name> or /org list")
			return nil
		}
		switch parts[1] {
		case "create":
			if len(parts) < 3 {
				m.addStatus("Usage: /org create <name>")
				return nil
			}
			name := strings.Join(parts[2:], " ")
			org, err := m.app.CreateOrg(name, "")
			if err != nil {
				m.addStatus(fmt.Sprintf("Error: %v", err))
				return nil
			}
			m.orgList = append(m.orgList, org)
			m.selectedOrg = len(m.orgList) - 1
			m.loadChannels()
			m.addStatus(fmt.Sprintf("Organization '%s' created", name))
		case "list":
			m.loadOrgs()
		default:
			m.addStatus("Unknown org command")
		}

	case "/channel":
		if len(parts) < 2 {
			m.addStatus("Usage: /channel create <name> or /channel list")
			return nil
		}
		switch parts[1] {
		case "create":
			if len(parts) < 3 {
				m.addStatus("Usage: /channel create <name>")
				return nil
			}
			name := strings.Join(parts[2:], " ")
			orgID := ""
			if len(m.orgList) > 0 {
				orgID = m.orgList[m.selectedOrg].OrgID
			}
			ch, err := m.app.CreateChannel(orgID, name, "text", "")
			if err != nil {
				m.addStatus(fmt.Sprintf("Error: %v", err))
				return nil
			}
			m.channelList = append(m.channelList, ch)
			m.selectedChan = len(m.channelList) - 1
			m.loadMessages()
			m.addStatus(fmt.Sprintf("Channel '%s' created", name))
		case "list":
			m.loadChannels()
		default:
			m.addStatus("Unknown channel command")
		}

	case "/dm":
		if len(parts) < 2 {
			m.addStatus("Usage: /dm <peer_id>")
			return nil
		}
		peerID := parts[1]
		if err := m.app.OpenDM(peerID); err != nil {
			m.addStatus(fmt.Sprintf("Error: %v", err))
			return nil
		}
		m.addStatus(fmt.Sprintf("DM with %s", peerID))

	case "/myaddr":
		peerID := m.app.PeerID()
		allAddrs := m.app.AllAddrs()
		m.addStatus(fmt.Sprintf("=== Peer ID: %s ===", peerID))
		for _, addr := range allAddrs {
			m.addStatus(fmt.Sprintf("  %s", addr))
		}
		m.addStatus("---")
		if len(allAddrs) > 0 {
			m.addStatus(fmt.Sprintf("Share: /connect %s", allAddrs[0]))
		}
		m.addStatus("Peer can find you via mDNS (LAN), DHT (internet), or a relay")

	case "/relay":
		if len(parts) < 2 {
			m.addStatus("Usage: /relay <multiaddr>  or  /relay connect <addr>")
			m.addStatus("Set relay_peers in config.yaml to auto-connect on startup")
			return nil
		}
		if parts[1] == "connect" && len(parts) >= 3 {
			if err := m.app.Connect(parts[2]); err != nil {
				m.addStatus(fmt.Sprintf("Relay connect error: %v", err))
				return nil
			}
			m.addStatus(fmt.Sprintf("Connected to relay: %s", parts[2]))
		} else {
			if err := m.app.Connect(parts[1]); err != nil {
				m.addStatus(fmt.Sprintf("Relay connect error: %v", err))
				return nil
			}
			m.addStatus(fmt.Sprintf("Connected to relay: %s", parts[1]))
		}

	case "/profile":
		id := m.app.Identity()
		fp := crypto.Fingerprint(id.PublicKey)
		m.addStatus(fmt.Sprintf("Profile: %s | PeerID: %s | Fingerprint: %s",
			id.DisplayName, id.PeerID, fp))

	case "/quit":
		return tea.Quit

	default:
		m.addStatus(fmt.Sprintf("Unknown command: %s. Type /help for help.", cmd))
	}

	return nil
}

func (m *Model) addStatus(msg string) {
	m.statusText = msg
	m.statusLog = append(m.statusLog, msg)
	if len(m.statusLog) > 100 {
		m.statusLog = m.statusLog[len(m.statusLog)-100:]
	}
}

func (m *Model) loadOrgs() {
	orgs, err := m.app.ListOrgs()
	if err != nil {
		m.addStatus(fmt.Sprintf("Error loading orgs: %v", err))
		return
	}
	m.orgList = orgs
}

func (m *Model) loadChannels() {
	if len(m.orgList) == 0 {
		m.channelList = nil
		return
	}
	orgID := m.orgList[m.selectedOrg].OrgID
	channels, err := m.app.ListChannels(orgID)
	if err != nil {
		m.addStatus(fmt.Sprintf("Error loading channels: %v", err))
		return
	}
	m.channelList = channels
	m.selectedChan = 0
	m.loadMessages()
}

func (m *Model) loadMessages() {
	if len(m.channelList) == 0 {
		m.messages = nil
		m.chatView.SetContent("")
		return
	}
	channelID := m.channelList[m.selectedChan].ChannelID
	msgs, err := m.app.ListMessages(channelID, 100, 0)
	if err != nil {
		m.addStatus(fmt.Sprintf("Error loading messages: %v", err))
		return
	}

	m.messages = nil
	for i := len(msgs) - 1; i >= 0; i-- {
		msg := msgs[i]
		sender := msg.SenderPeerID
		if len(sender) > 12 {
			sender = sender[:12]
		}
		m.messages = append(m.messages, MessageItem{
			Sender:    sender,
			Content:   msg.Content,
			Timestamp: msg.CreatedAt.Format("15:04"),
		})
	}
	m.chatView.SetContent(m.renderMessages())
	m.chatView.GotoBottom()
}

func (m *Model) loadPeers() {
	peers, err := m.app.ListPeers()
	if err != nil {
		m.addStatus(fmt.Sprintf("Error loading peers: %v", err))
		return
	}
	m.peerList = peers
}

func (m *Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	orgPanel := m.renderOrgPanel()
	channelPanel := m.renderChannelPanel()
	chatPanel := m.renderChatPanel()
	statusBar := m.renderStatusBar()

	leftSide := lipgloss.JoinHorizontal(lipgloss.Top, orgPanel, channelPanel)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, chatPanel)

	input := InputStyle.Render(m.input.View())
	if !m.inputMode {
		input = DimmedInputStyle.Render(m.input.View())
	}

	if m.firstLaunch {
		m.firstLaunch = false
		m.addStatus("Press Tab to navigate orgs/channels | ? for help")
	}

	body := lipgloss.JoinVertical(lipgloss.Left, topRow, input, statusBar)

	return AppStyle.Render(body)
}

func (m *Model) renderOrgPanel() string {
	style := OrgPanelStyle
	if !m.inputMode && m.focus == focusOrgs {
		style = OrgPanelFocusedStyle
	}

	if len(m.orgList) == 0 {
		return style.Render(DimmedStyle.Render("No orgs\n\n/org create <name>"))
	}

	var items []string
	for i, org := range m.orgList {
		name := org.Name
		if len(name) > 18 {
			name = name[:18]
		}
		if i == m.selectedOrg {
			items = append(items, SelectedOrgStyle.Render("  "+name+"  "))
		} else {
			items = append(items, OrgItemStyle.Render("  "+name))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return style.Render(TitleStyle.Render("Orgs") + "\n" + content)
}

func (m *Model) renderChannelPanel() string {
	style := ChannelPanelStyle
	if !m.inputMode && m.focus == focusChannels {
		style = ChannelPanelFocusedStyle
	}

	if len(m.channelList) == 0 {
		return style.Render(DimmedStyle.Render("No channels\n\n/channel create <name>"))
	}

	var items []string
	for i, ch := range m.channelList {
		name := ch.Name
		if len(name) > 22 {
			name = name[:22]
		}
		prefix := "# "
		if ch.ChannelType == "announcement" {
			prefix = "! "
		}
		if ch.ReadOnly {
			prefix = "🔒 "
		}
		if i == m.selectedChan {
			items = append(items, SelectedChannelStyle.Render("  "+prefix+name+"  "))
		} else {
			items = append(items, ChannelItemStyle.Render("  "+prefix+name))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	return style.Render(TitleStyle.Render("Channels") + "\n" + content)
}

func (m *Model) renderChatPanel() string {
	messages := m.renderMessages()

	header := ""
	if len(m.channelList) > 0 {
		ch := m.channelList[m.selectedChan]
		header = ChannelHeaderStyle.Render(" # "+ch.Name+" ") + "\n"
	}

	chatContent := header + "\n" + messages
	m.chatView.SetContent(chatContent)

	return ChatPanelStyle.Render(m.chatView.View())
}

func (m *Model) renderMessages() string {
	if m.showHelp {
		return m.helpView()
	}
	if m.showPeers {
		return m.peersView()
	}
	if len(m.messages) == 0 {
		return DimmedStyle.Render("  No messages yet. Start typing!")
	}

	var items []string
	for _, msg := range m.messages {
		line := fmt.Sprintf("%s %s %s",
			TimeStyle.Render(msg.Timestamp),
			SenderStyle.Render(msg.Sender),
			msg.Content,
		)
		items = append(items, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func (m *Model) renderStatusBar() string {
	modeBadge := ModeBadgeInput.Render(" INPUT ")
	if !m.inputMode {
		modeBadge = ModeBadgeNav.Render(" NAV ")
	}

	ctx := ""
	if len(m.orgList) > 0 {
		ctx = m.orgList[m.selectedOrg].Name
	}
	if len(m.channelList) > 0 {
		ctx = ctx + " / #" + m.channelList[m.selectedChan].Name
	}

	statusText := m.statusText
	if ctx != "" {
		statusText = ctx + " │ " + m.statusText
	}

	left := StatusStyle.Render(statusText)
	right := StatusStyle.Render(fmt.Sprintf("%s  Peers: %d", modeBadge, len(m.peerList)))

	barWidth := m.width - 4
	leftLen := lipgloss.Width(left)
	rightLen := lipgloss.Width(right)
	gap := barWidth - leftLen - rightLen
	if gap < 1 {
		gap = 1
	}

	spacer := strings.Repeat(" ", gap)
	return lipgloss.JoinHorizontal(lipgloss.Left, left, spacer, right)
}

func (m Model) helpView() string {
	return HelpStyle.Render(`Commands:
  /help             Show this help
  /myaddr           Show your shareable multiaddress
  /connect <addr>   Connect to a peer directly
  /relay <addr>     Connect via a relay peer
  /disconnect       Disconnect all peers
  /peers            List known peers
  /org create       Create an organization
  /channel create   Create a channel
  /dm <peer>        Open direct message
  /profile          Show your profile
  /quit             Quit

Keys:
  Tab        Cycle: input ─ channels ─ orgs
  Arrows     Navigate orgs/channels (nav mode)
  Enter      Send message / confirm
  ?          Toggle help
  P          Toggle peers
  Ctrl+C     Quit`)
}

func (m Model) peersView() string {
	if len(m.peerList) == 0 {
		return DimmedStyle.Render("  No peers connected.")
	}
	var items []string
	for _, p := range m.peerList {
		items = append(items, fmt.Sprintf("  %s (%s)", p.DisplayName, p.Status))
	}
	return strings.Join(items, "\n")
}
