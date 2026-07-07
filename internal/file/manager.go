package file

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Transfer() error {
	return nil
}

func (m *Manager) Verify() error {
	return nil
}
