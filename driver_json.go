package envx

import (
	"strings"

	"github.com/shestakovda/errx"
	"github.com/tidwall/gjson"
)

const (
	objDelim     = ".$."
	jsonReadOnly = "json args driver is read-only"
)

func NewDriverJSON(js []byte) Driver {
	return &jsonDriver{
		src: js,
	}
}

type jsonDriver struct {
	src []byte
}

func (d *jsonDriver) Set(name, value string) { panic(errx.New(jsonReadOnly)) }
func (d *jsonDriver) Del(name string)        { panic(errx.New(jsonReadOnly)) }

func (d *jsonDriver) Get(name string) string {
	return gjson.GetBytes(d.src, name).String()
}

func (d *jsonDriver) GetArray(name string) []string {
	var collect func(items gjson.Result)

	list := make([]string, 0, 16)

	collect = func(res gjson.Result) {
		if !res.Exists() {
			return
		}

		if !res.IsArray() {
			list = append(list, res.String())
			return
		}

		items := res.Array()
		for i := range items {
			collect(items[i])
		}
	}

	if parts := strings.Split(name, objDelim); len(parts) > 1 {
		gjson.GetBytes(d.src, parts[0]).ForEach(func(key, value gjson.Result) bool {
			collect(value.Get(parts[1]))
			return true
		})
	} else {
		collect(gjson.GetBytes(d.src, name))
	}

	return list
}
