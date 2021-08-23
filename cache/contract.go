package cache

type Contract interface {
	Remove(key string) error
	RemovePrefix(prefix string) error
	GetMarshal(key string, unMarshal interface{}) error
	SetMarshal(key string, canMarshalVal interface{}, seconds int) error
	Exist(key string) bool
	GetOrErr(key string) (string, error)
	Get(key string, fallback string) string
	Set(key string, value string, seconds int)
	SetIfNotExist(key string, value string, seconds int) bool
	Reset()
	Incr(key string) (result int64, err error)
	GetInt64(key string, fallback int64) int64
}
