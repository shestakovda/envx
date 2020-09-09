package envx_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/shestakovda/envx"
	"github.com/stretchr/testify/assert"
)

const k, v, tv, e = "Key", " Value ", "Value", ""

func testDriver(t *testing.T, d envx.Driver) {
	assert.Equal(t, e, d.Get(k))
	assert.Nil(t, d.GetArray(k))

	d.Set(k, v)
	assert.Equal(t, tv, d.Get(k))
	assert.Equal(t, []string{tv}, d.GetArray(k))

	d.Del(k)
	assert.Equal(t, e, d.Get(k))
	assert.Nil(t, d.GetArray(k))
}

func TestMemDriver(t *testing.T) {
	testDriver(t, envx.NewMemDriver(16))
}

func TestEnvDriver(t *testing.T) {
	testDriver(t, envx.NewEnvDriver("test"))
}

func TestHTTPDriver(t *testing.T) {
	req, err := http.NewRequest("POST", "", nil)
	assert.NoError(t, err)

	drv, err := envx.NewHTTPDriver(req)
	assert.True(t, errors.Is(err, envx.ErrHTTPInvalid))
	assert.Nil(t, drv)

	req, err = http.NewRequest("GET", "", nil)
	assert.NoError(t, err)

	drv, err = envx.NewHTTPDriver(req)
	assert.NoError(t, err)
	assert.NotNil(t, drv)

	testDriver(t, drv)
}

func TestDriverJSON(t *testing.T) {
	js := []byte(`{
		"test": "ololo",
		"meow": [
			"purpur",
			"furfur"
		],
		"тачки": [
			{
				"модель": "vaz",
				"год": 1995,
				"владельцы": [
					"Иванов",
					"Петров В."
				]
			},
			{
				"модель": "gaz",
				"год": 1986,
				"владельцы": [
					"Сидоров-Пражский",
					"Жужелица А.В.",
					"П. Лут"
				]
			}
		],
		"Владельцы": {
			"Иванов": {
				"Город": "Москва",
				"Ник в PUBG": "wado",
			},
			"Сидоров-Пражский": {
				"Город": "Калуга",
				"Ник в PUBG": "sidor",
			},
			"П. Лут": {
				"Город": "Усть-Каменогорск",
				"Ник в PUBG": "pluto",
			}
		}
	}`)

	drv := envx.NewDriverJSON(js)

	assert.Equal(t, "ololo", drv.Get("test"))
	assert.Equal(t, "Петров В.", drv.Get("тачки.0.владельцы.1"))
	assert.Equal(t, []string{"purpur", "furfur"}, drv.GetArray("meow"))
	assert.Equal(t, []string{"vaz", "gaz"}, drv.GetArray("тачки.#.модель"))
	assert.Equal(t, []string{"1995", "1986"}, drv.GetArray("тачки.#.год"))
	assert.Equal(t, []string{"2", "3"}, drv.GetArray("тачки.#.владельцы.#"))
	assert.Equal(t, []string{
		"Иванов", "Петров В.", "Сидоров-Пражский", "Жужелица А.В.", "П. Лут",
	}, drv.GetArray("тачки.#.владельцы"))
	assert.Equal(t, []string{
		"Москва", "Калуга", "Усть-Каменогорск",
	}, drv.GetArray("Владельцы.$.Город"))

	drv = envx.NewDriverJSON([]byte(jsBenchEvent))

	want := "Документы.#.Подписи.#.ИдФайла"
	good := []string{
		"09e3968080fd447aa56954b34872f9f1/39d19ae2ddec4b98af57e60487c69a9c",
		"09e3968080fd447aa56954b34872f9f1/024a426f619745f3a3ddb1a3ba9e2012",
		"09e3968080fd447aa56954b34872f9f1/4c3577b9e4b64a7eb42d9162bf98ca4d",
	}

	assert.Equal(t, good, drv.GetArray(want))
}

//nolint:lll
const jsBenchEvent = `{
  "Направление": "ФНС",
  "Отправитель": "9999",
  "Получатель": "1AEBC3A424B-A27C-4F6D-9BA3-973D5B6F3C69",
  "ИдАбонента": "bc3a424b-a27c-4f6d-9ba3-973d5b6f3c69",
  "ИдПакетаДокументов": "09e3968080fd447aa56954b34872f9f1",
  "Документооборот": "1ead8656-1bac-4d21-a183-0827f747bcbd",
  "ТипДокументооборота": "Документ",
  "КодДокументооборота": "10",
  "ТипТранзакции": "ДокументНО",
  "КодТранзакции": "01",
  "ИдФайлаАрхива": "09e3968080fd447aa56954b34872f9f1/088d6a3504174e5589c973ffab39596c",
  "Документы": [
    {
      "ТипДокумента": "документ",
      "КодДокумента": "01",
      "ИдФайла": "09e3968080fd447aa56954b34872f9f1/f13e7bb93f654790bb4129c18ad15346",
      "Наименование": "ON_DOCNPNO_9653460472999901001_9653460472999901001_9999_20200305_d28de959882147918797c7459bcfef09.xml",
      "Зашифрован": false,
      "Подписи": [
        {
          "ИдФайла": "09e3968080fd447aa56954b34872f9f1/39d19ae2ddec4b98af57e60487c69a9c",
          "Присоединенная": false
        }
      ]
    },
    {
      "ТипДокумента": "приложение",
      "КодДокумента": "03",
      "ИдФайла": "09e3968080fd447aa56954b34872f9f1/45e7897b885a47248a792eb0e14845a1",
      "Наименование": "1165050_9999_9653460472999901001_78a4bc8a-5880-43c5-9f8a-817e43db58fd_20200305_fb290b5a9f28425393ce00e013afc54c.xml",
      "Зашифрован": false,
      "Подписи": [
        {
          "ИдФайла": "09e3968080fd447aa56954b34872f9f1/024a426f619745f3a3ddb1a3ba9e2012",
          "Присоединенная": false
        }
      ]
    },
    {
      "ТипДокумента": "описание",
      "КодДокумента": "02",
      "ИдФайла": "09e3968080fd447aa56954b34872f9f1/1df72c432b654ac980f9e4a91fc1b662",
      "Наименование": "TR_INFSOOB.xml",
      "Зашифрован": false,
      "Подписи": [
        {
          "ИдФайла": "09e3968080fd447aa56954b34872f9f1/4c3577b9e4b64a7eb42d9162bf98ca4d",
          "Присоединенная": false
        }
      ]
    }
  ]
}`

func BenchmarkJsGetArray(b *testing.B) {
	drv := envx.NewDriverJSON([]byte(jsBenchEvent))

	want := "Документы.#.Подписи.#.ИдФайла"
	good := []string{
		"09e3968080fd447aa56954b34872f9f1/39d19ae2ddec4b98af57e60487c69a9c",
		"09e3968080fd447aa56954b34872f9f1/024a426f619745f3a3ddb1a3ba9e2012",
		"09e3968080fd447aa56954b34872f9f1/4c3577b9e4b64a7eb42d9162bf98ca4d",
	}

	for i := 0; i < b.N; i++ {
		assert.Equal(b, good, drv.GetArray(want))
	}
}
