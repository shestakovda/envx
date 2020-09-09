package envx

import (
	"os"
	"strings"
)

func NewEnvDriver(pfx string) Driver {
	return &envDriver{
		pfx: strings.ToUpper(pfx) + "_",
	}
}

type envDriver struct {
	pfx string
}

func (d *envDriver) Set(name, value string) {
	os.Setenv(d.pfx+strings.ToUpper(name), value)
}

func (d *envDriver) Get(name string) string {
	return strings.TrimSpace(os.Getenv(d.pfx + strings.ToUpper(name)))
}

func (d *envDriver) Del(name string) {
	os.Unsetenv(d.pfx + strings.ToUpper(name))
}

func (d *envDriver) GetArray(name string) []string {
	val, ok := os.LookupEnv(d.pfx + strings.ToUpper(name))

	if ok {
		return []string{strings.TrimSpace(val)}
	}

	return nil
}
