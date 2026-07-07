package sync

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Sync() error {
	return nil
}

func (m *Manager) Reconcile() error {
	return nil
}
