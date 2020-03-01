package lsm9ds1

import "sync"

type Reading interface {
	SetReading(x, y, z int16)
	GetReading() (x, y, z int16)
}

type MutexReading struct {
	mu      sync.Mutex
	x, y, z int16
}

func (m *MutexReading) SetReading(x, y, z int16) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.x = x
	m.y = y
	m.z = z
}

func (m *MutexReading) GetReading() (x, y, z int16) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.x, m.y, m.z
}
