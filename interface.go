package envx

import (
	"time"

	"github.com/shestakovda/errx"
)

// Provider - поставщик параметров для приложения
type Provider interface {
	/*
		Driver - ссылка на привязанный драйвер.

		* Позволяет работать с "сырыми" значениями
		* Прежде всего, очень удобно для тестов
	*/
	Driver

	/*
		Группа методов для получения значений в разных типах. Общая логика

		* Если значения нет, подставляется по-умолчанию
		* Если есть - валидируется по типу/схеме, если не подходит - возвращается ошибка
		* Методы могут как-то преобразовывать сырые значения, если это необходимо
	*/
	Bool(name string, def bool) bool
	String(name string, def string) string
	URL(name string, def string) (string, error)
	UUID(name string, def string) (string, error)
	GUID(name string, def string) (string, error)
	JSON(name, def string, item interface{}) error
	Uint64(name string, def uint64) (uint64, error)
	Timezone(name string, def string) (*time.Location, error)
	Duration(name string, def time.Duration) (time.Duration, error)
	StringArray(name string, def []string) ([]string, error)
	TimeRFC3339(name string, def time.Time) (time.Time, error)
}

// Driver - реализация конкретного поставщика параметров
type Driver interface {
	/*
		Get - получение значения по ключу.

		* Должен производить нормализацию ключа, если это необходимо
		* Должен возвращать первое значение, если ключу соответствует несколько
		* Должен быть потокобезопасным
	*/
	Get(name string) string

	/*
		Set - установка значения по ключу.

		* Должен производить нормализацию ключа и значения, если это необходимо
		* Должен добавлять значение к списку, если ключу соответствует несколько
		* Должен быть потокобезопасным
	*/
	Set(name, value string)

	/*
		Del - удаление значения по ключу.

		* Должен производить нормализацию ключа, если это необходимо
		* Должен удалять все значения, если ключу соответствует несколько
		* Должен быть потокобезопасным
	*/
	Del(name string)

	/*
		GetArray - получение списка значений по ключу.

		* Должен производить нормализацию ключа, если это необходимо
		* Должен возвращать список из одного элемента, если ключу соответствует только одно значение
		* Должен быть потокобезопасным
	*/
	GetArray(name string) []string
}

// Ошибки модуля
var (
	ErrURLEmpty        = errx.New("Пустой URL")
	ErrURLInvalid      = errx.New("Некорректный URL")
	ErrUUIDEmpty       = errx.New("Пустой UUID")
	ErrUUIDInvalid     = errx.New("Некорректный UUID")
	ErrGUIDEmpty       = errx.New("Пустой GUID")
	ErrGUIDInvalid     = errx.New("Некорректный GUID")
	ErrUint64Invalid   = errx.New("Некорректное целое")
	ErrTimezoneEmpty   = errx.New("Пустой часовой пояс")
	ErrTimezoneInvalid = errx.New("Некорректный часовой пояс")
	ErrDurationInvalid = errx.New("Некорректный промежуток времени")
	ErrRFC3339Invalid  = errx.New("Некорректная дата в формате RFC 3339")
	ErrJSONInvalid     = errx.New("Некорректный JSON")
	ErrHTTPInvalid     = errx.New("Некорректный HTTP-запрос")
)
