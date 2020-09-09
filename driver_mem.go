package envx

import (
	"strings"
	"sync"
)

func NewMemDriver(size uint) Driver {
	return &memDriver{
		data: make(map[string]string, int(size)),
	}
}

type memDriver struct {
	sync.RWMutex
	data map[string]string
}

func (d *memDriver) Set(name, value string) {
	d.Lock()
	defer d.Unlock()
	d.data[name] = value
}

func (d *memDriver) Get(name string) string {
	d.RLock()
	defer d.RUnlock()

	return strings.TrimSpace(d.data[name])
}

func (d *memDriver) GetArray(name string) []string {
	d.RLock()
	defer d.RUnlock()

	val, ok := d.data[name]

	if ok {
		return []string{strings.TrimSpace(val)}
	}

	return nil
}

func (d *memDriver) Del(name string) {
	d.Lock()
	defer d.Unlock()
	delete(d.data, name)
}
