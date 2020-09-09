package envx

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/shestakovda/errx"
)

// NewHTTPDriver - получение аргументов из HTTP-запроса (в т.ч. из POST-формы)
func NewHTTPDriver(req *http.Request) (_ Driver, err error) {
	if err = req.ParseForm(); err != nil {
		var detail string

		switch msg := err.Error(); {
		case strings.Contains(msg, "invalid URL escape"):
			detail = "Некорректное URL-кодирование"
		case strings.Contains(msg, "missing form body"):
			detail = "Отсутствует тело запроса"
		case strings.Contains(msg, "too large"):
			detail = "Слишком большое тело запроса"
		case strings.Contains(msg, "mime"):
			detail = "Некорректное содержимое заголовка `Content-Type`"
		default:
			detail = "Тело запроса повреждено или сформировано некорректно"
		}

		return nil, ErrHTTPInvalid.WithDetail(detail).WithDebug(errx.Debug{
			"Адрес":     req.URL.RawPath,
			"Заголовки": req.Header,
		})
	}

	return &httpDriver{
		values: req.Form,
	}, nil
}

type httpDriver struct {
	values url.Values
}

func (d *httpDriver) Set(name, value string) {
	d.values.Set(name, value)
}

func (d *httpDriver) Get(name string) string {
	return strings.TrimSpace(d.values.Get(name))
}

func (d *httpDriver) GetArray(name string) []string {
	if _, ok := d.values[name]; !ok {
		return nil
	}

	for i := range d.values[name] {
		d.values[name][i] = strings.TrimSpace(d.values[name][i])
	}

	return d.values[name]
}

func (d *httpDriver) Del(name string) {
	d.values.Del(name)
}
