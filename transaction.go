package safemap

import ()

type Transaction struct {
	c chan request
}

func (t *Transaction) Get(key string) (interface{}, bool) {
	return get(t.c, key)
}

func (t *Transaction) Set(key string, value interface{}) (oldv interface{}) {
	return set(t.c, key, value)
}

func (t *Transaction) Del(key string) (oldv interface{}) {
	return del(t.c, key)
}

func (t *Transaction) EndTransaction() {
	endTransaction(t.c)
}
