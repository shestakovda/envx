package envx_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/shestakovda/envx"
	"github.com/stretchr/testify/suite"
)

const (
	wtf  = `$#&?!`
	name = "test_value"
)

// TestDrivers - индивидуальные тесты драйверов
func TestDrivers(t *testing.T) {
}

// TestProviders - индивидуальные тесты провайдеров
func TestProviders(t *testing.T) {
}

// TestArgs - интеграционные тесты пакета на драйвере-эмуляторе
func TestArgs(t *testing.T) {
	suite.Run(t, new(ArgsSuite))
}

type ArgsSuite struct {
	suite.Suite

	prv envx.Provider
}

func (s *ArgsSuite) SetupTest() {
	s.prv = envx.NewProvider(envx.NewMemDriver(16))
}

func (s *ArgsSuite) TestString() {
	const def = "Жареная рыба"

	s.prv.Del(name)

	s.Equal(def, s.prv.String(name, def))

	s.prv.Set(name, wtf)
	s.Equal(wtf, s.prv.String(name, def))
}

func (s *ArgsSuite) TestBool() {
	s.prv.Del(name)

	s.True(s.prv.Bool(name, true))
	s.False(s.prv.Bool(name, false))

	no := []string{"0", "f", "false", "n", "no", "н", "нет", " нет "}
	yes := []string{"1", "t", "true", "y", "yes", "д", "да", " да "}

	for i := range yes {
		s.prv.Set(name, yes[i])
		s.True(s.prv.Bool(name, false))
		s.prv.Set(name, no[i])
		s.False(s.prv.Bool(name, true))
	}
}

func (s *ArgsSuite) TestURL() {
	const def = "https://example.com"

	s.prv.Del(name)

	v, err := s.prv.URL(name, def+"/")
	s.NoError(err)
	s.Equal(def, v)

	if _, err = s.prv.URL(name, ""); s.Error(err) {
		s.True(errors.Is(err, envx.ErrURLEmpty))
	}

	s.prv.Set(name, wtf)
	if _, err = s.prv.URL(name, def); s.Error(err) {
		s.True(errors.Is(err, envx.ErrURLInvalid))
	}

	str := def + "test"
	s.prv.Set(name, str+"/")
	v, err = s.prv.URL(name, def)
	s.NoError(err)
	s.Equal(str, v)
}

func (s *ArgsSuite) TestUUID() {
	const def = "123456781234123412341234123412AF"

	s.prv.Del(name)

	v, err := s.prv.UUID(name, def)
	s.NoError(err)
	s.Equal(strings.ToLower(def), v)

	if _, err = s.prv.UUID(name, ""); s.Error(err) {
		s.True(errors.Is(err, envx.ErrUUIDEmpty))
	}

	s.prv.Set(name, wtf)
	if _, err = s.prv.UUID(name, def); s.Error(err) {
		s.True(errors.Is(err, envx.ErrUUIDInvalid))
	}

	str := "4321432143af43AF43af432143214321"
	s.prv.Set(name, str)
	v, err = s.prv.UUID(name, def)
	s.NoError(err)
	s.Equal(strings.ToLower(str), v)
}

func (s *ArgsSuite) TestGUID() {
	const def = "12345678-1234-1234-12af-123412341234"

	s.prv.Del(name)

	v, err := s.prv.GUID(name, def)
	s.NoError(err)
	s.Equal(strings.ToUpper(def), v)

	if _, err = s.prv.GUID(name, ""); s.Error(err) {
		s.True(errors.Is(err, envx.ErrGUIDEmpty))
	}

	s.prv.Set(name, wtf)
	if _, err = s.prv.GUID(name, def); s.Error(err) {
		s.True(errors.Is(err, envx.ErrGUIDInvalid))
	}

	str := "43214321-43af-43af-43af-432143214321"
	s.prv.Set(name, str)
	v, err = s.prv.GUID(name, def)
	s.NoError(err)
	s.Equal(strings.ToUpper(str), v)
}

func (s *ArgsSuite) TestUint64() {
	const def = uint64(42)

	s.prv.Del(name)

	v, err := s.prv.Uint64(name, def)
	s.NoError(err)
	s.Equal(def, v)

	s.prv.Set(name, wtf)
	if _, err = s.prv.Uint64(name, def); s.Error(err) {
		s.True(errors.Is(err, envx.ErrUint64Invalid))
	}

	s.prv.Set(name, "-60")
	if _, err = s.prv.Uint64(name, def); s.Error(err) {
		s.True(errors.Is(err, envx.ErrUint64Invalid))
	}

	s.prv.Set(name, "23.4")
	if _, err = s.prv.Uint64(name, def); s.Error(err) {
		s.True(errors.Is(err, envx.ErrUint64Invalid))
	}

	s.prv.Set(name, "60")
	v, err = s.prv.Uint64(name, def)
	s.NoError(err)
	s.Equal(uint64(60), v)
}

func (s *ArgsSuite) TestTimezone() {
	const def = "Europe/Moscow"

	s.prv.Del(name)

	v, err := s.prv.Timezone(name, def)
	s.NoError(err)
	s.Equal(def, v.String())

	if _, err = s.prv.Timezone(name, ""); s.Error(err) {
		s.True(errors.Is(err, envx.ErrTimezoneEmpty))
	}

	s.prv.Set(name, wtf)
	if _, err = s.prv.Timezone(name, ""); s.Error(err) {
		s.True(errors.Is(err, envx.ErrTimezoneInvalid))
	}

	str := "America/New_York"
	s.prv.Set(name, str)
	v, err = s.prv.Timezone(name, def)
	s.NoError(err)
	s.Equal(str, v.String())
}

func (s *ArgsSuite) TestDuration() {
	const def = time.Second

	s.prv.Del(name)

	v, err := s.prv.Duration(name, def)
	s.NoError(err)
	s.Equal(def, v)

	s.prv.Set(name, wtf)
	if _, err = s.prv.Duration(name, def); s.Error(err) {
		s.True(errors.Is(err, envx.ErrDurationInvalid))
	}

	s.prv.Set(name, "60s")
	v, err = s.prv.Duration(name, def)
	s.NoError(err)
	s.Equal(time.Minute, v)
}

func (s *ArgsSuite) TestRFC3339() {
	var def = time.Now()

	s.prv.Del(name)

	v, err := s.prv.TimeRFC3339(name, def)
	s.NoError(err)
	s.Equal(def, v)

	s.prv.Set(name, wtf)
	if _, err = s.prv.TimeRFC3339(name, def); s.Error(err) {
		s.True(errors.Is(err, envx.ErrRFC3339Invalid))
	}

	now := time.Now()
	s.prv.Set(name, now.Format(time.RFC3339))
	v, err = s.prv.TimeRFC3339(name, def)
	s.NoError(err)
	s.Equal(now.Unix(), v.Unix())
}

func (s *ArgsSuite) TestStringArray() {
	var def = []string{"item"}

	s.prv.Del(name)

	v, err := s.prv.StringArray(name, def)
	s.NoError(err)
	s.Equal(def, v)

	s.prv.Set(name, wtf)
	v, err = s.prv.StringArray(name, def)
	s.NoError(err)
	s.Equal([]string{wtf}, v)
}

func (s *ArgsSuite) TestJSON() {
	var def = `["test", "wtf", "ololo"]`
	var good = []string{"test", "wtf", "ololo"}

	s.prv.Del(name)

	var list []string
	s.NoError(s.prv.JSON(name, def, &list))
	s.Equal(good, list)

	type testType struct {
		Test string `json:"ololo"`
	}

	item := new(testType)

	s.prv.Set(name, wtf)
	if err := s.prv.JSON(name, def, item); s.Error(err) {
		s.True(errors.Is(err, envx.ErrJSONInvalid))
	}

	s.prv.Set(name, `{"ololo": "purpur"}`)
	s.NoError(s.prv.JSON(name, def, item))
	s.Equal("purpur", item.Test)
}
