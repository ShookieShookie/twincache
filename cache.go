package twincache

type Cache struct {
	capacity int
}

func New(cap int) *Cache {
	return &Cache{capacity: cap}
}

func (c *Cache) Add(interface{}, interface{}) error {
	return nil
}
func (c *Cache) get(interface{}) (interface{}, error) {
	return nil, nil
}
