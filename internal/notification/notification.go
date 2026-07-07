package notification

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Notify(title, message string) error {
	return nil
}
