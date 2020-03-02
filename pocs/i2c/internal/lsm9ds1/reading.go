package lsm9ds1

import "sync"

type Reading interface {
	SetReading(x, y, z int16)
	FromList([]int16)

	GetReading() (x, y, z int16)
	ToList() []int16
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

func (m *MutexReading) ToList() []int16 {
	x, y, z := m.GetReading()
	return []int16{x, y, z}
}

func (m *MutexReading) FromList(l []int16) {
	m.SetReading(l[0], l[1], l[2])
}

func (m *MutexReading) GetReading() (x, y, z int16) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.x, m.y, m.z
}
