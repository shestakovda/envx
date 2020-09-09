package envx

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/shestakovda/errx"
)

const (
	argName  = "Параметр"
	argValue = "Значение"
)

// NewProvider - конструктор поставщика настроек из окружения
func NewProvider(driver Driver) Provider {
	return &provider{
		Driver: driver,
		rxUUID: regexp.MustCompile(`^[0-9a-f]{32}$`),
		rxGUID: regexp.MustCompile(`^[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}$`),
	}
}

type provider struct {
	Driver
	rxUUID *regexp.Regexp
	rxGUID *regexp.Regexp
}

func (p *provider) String(name string, def string) string {
	s := p.Get(name)

	if s == "" {
		return def
	}

	return s
}

func (p *provider) Bool(name string, def bool) bool {
	const t1, t2, t3, t4, t5, t6, t7 = "1", "t", "true", "y", "yes", "д", "да"

	if s := strings.ToLower(p.Get(name)); s != "" {
		return s == t1 || s == t2 || s == t3 || s == t4 || s == t5 || s == t6 || s == t7
	}

	return def
}

func (p *provider) URL(name string, def string) (string, error) {
	s := p.Get(name)

	if s == "" {
		if def == "" {
			return "", ErrURLEmpty.WithDebug(errx.Debug{argName: name})
		}

		return strings.TrimSuffix(def, "/"), nil
	}

	if !govalidator.IsURL(s) {
		return "", ErrURLInvalid.WithDebug(errx.Debug{argValue: s, argName: name})
	}

	return strings.TrimSuffix(s, "/"), nil
}

func (p *provider) UUID(name string, def string) (string, error) {
	s := strings.ToLower(p.Get(name))

	if s == "" {
		if def == "" {
			return "", ErrUUIDEmpty.WithDebug(errx.Debug{argName: name})
		}

		return strings.ToLower(def), nil
	}

	if !p.rxUUID.MatchString(s) {
		return "", ErrUUIDInvalid.WithDebug(errx.Debug{argValue: s, argName: name})
	}

	return s, nil
}

func (p *provider) GUID(name string, def string) (string, error) {
	s := strings.ToUpper(p.Get(name))

	if s == "" {
		if def == "" {
			return "", ErrGUIDEmpty.WithDebug(errx.Debug{argName: name})
		}

		return strings.ToUpper(def), nil
	}

	if !p.rxGUID.MatchString(s) {
		return "", ErrGUIDInvalid.WithDebug(errx.Debug{argValue: s, argName: name})
	}

	return s, nil
}

func (p *provider) Uint64(name string, def uint64) (uint64, error) {
	var err error
	var num uint64

	s := p.Get(name)

	if s == "" {
		return def, nil
	}

	if num, err = strconv.ParseUint(s, 10, 64); err != nil {
		return 0, ErrUint64Invalid.WithReason(err).WithDebug(errx.Debug{argValue: s, argName: name})
	}

	return num, nil
}

func (p *provider) Timezone(name string, def string) (*time.Location, error) {
	var err error
	var loc *time.Location

	s := p.Get(name)

	if s == "" {
		if def == "" {
			return nil, ErrTimezoneEmpty.WithDebug(errx.Debug{argName: name})
		}
		s = def
	}

	if loc, err = time.LoadLocation(s); err != nil {
		return nil, ErrTimezoneInvalid.WithReason(err).WithDebug(errx.Debug{argValue: s, argName: name})
	}

	return loc, nil
}

func (p *provider) Duration(name string, def time.Duration) (time.Duration, error) {
	var err error
	var dur time.Duration

	s := strings.ToLower(p.Get(name))

	if s == "" {
		return def, nil
	}

	if dur, err = time.ParseDuration(s); err != nil {
		return 0, ErrDurationInvalid.WithReason(err).WithDebug(errx.Debug{argValue: s, argName: name})
	}

	return dur, nil
}

func (p *provider) TimeRFC3339(name string, def time.Time) (time.Time, error) {
	var err error
	var rfc time.Time

	s := strings.ToUpper(p.Get(name))

	if s == "" {
		return def, nil
	}

	if rfc, err = time.Parse(time.RFC3339, s); err != nil {
		return rfc, ErrRFC3339Invalid.WithReason(err).WithDebug(errx.Debug{argValue: s, argName: name})
	}

	return rfc, nil
}

func (p *provider) StringArray(name string, def []string) ([]string, error) {
	if s := p.GetArray(name); s != nil {
		return s, nil
	}

	return def, nil
}

func (p *provider) JSON(name, def string, item interface{}) error {
	var js []byte

	if s := p.Get(name); s != "" {
		js = []byte(s)
	} else {
		js = []byte(def)
	}

	if err := json.Unmarshal(js, item); err != nil {
		return ErrJSONInvalid.WithReason(err).WithDebug(errx.Debug{argValue: js, argName: name})
	}

	return nil
}
