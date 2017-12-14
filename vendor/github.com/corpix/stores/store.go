package stores

type Store interface {
	Name() string
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Remove(key string) error
	Keys() ([]string, error)
	Values() ([]interface{}, error)
	Map() (map[string]interface{}, error)
	Iter(func(key string, value interface{}) bool) error
	Close() error
}
