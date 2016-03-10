package safemap

import ()

// SafeMap is a thread-safe map based on channel
type SafeMap struct {
	m map[string]interface{}
	c chan request // operation request
}

// New is used to create a SafeMap, you should always use this function to create a SafeMap
func New() *SafeMap {
	result := &SafeMap{
		m: make(map[string]interface{}),
		c: make(chan request),
	}

	go handleRequest(result.m, result.c)

	return result
}

func handleRequest(m map[string]interface{}, c chan request) {
	for {
		req := <-c
		switch req.reqType {
		case Get:
			resp := response{}
			resp.value, resp.ok = m[req.key]
			req.respc <- resp
		case Set:
			resp := response{}
			resp.value, _ = m[req.key]
			m[req.key] = req.value
			req.respc <- resp
		case Delete:
			resp := response{}
			resp.value, _ = m[req.key]
			delete(m, req.key)
			req.respc <- resp
		case BeginTransaction:
			handleRequest(m, req.tc)
		case EndTransaction:
			return
		}
	}
}

// Get the value by key
func (s *SafeMap) Get(key string) (interface{}, bool) {
	return get(s.c, key)
}

// Set value by key, return old value if exists
func (s *SafeMap) Set(key string, value interface{}) (oldv interface{}) {
	return set(s.c, key, value)
}

// Delete a entry by key, return old value if exists
func (s *SafeMap) Del(key string) (oldv interface{}) {
	return del(s.c, key)
}

// BeginTransaction opens a transaction
func (s *SafeMap) BeginTransaction() *Transaction {
	c := make(chan request)
	s.c <- request{
		reqType: BeginTransaction,
		tc:      c,
	}
	return &Transaction{
		c: c,
	}
}

func get(c chan request, key string) (interface{}, bool) {
	respc := make(chan response)
	c <- request{
		reqType: Get,
		key:     key,
		respc:   respc,
	}
	resp := <-respc
	return resp.value, resp.ok
}

func set(c chan request, key string, value interface{}) (oldv interface{}) {
	respc := make(chan response)
	c <- request{
		reqType: Set,
		key:     key,
		value:   value,
		respc:   respc,
	}
	resp := <-respc
	return resp.value
}

func del(c chan request, key string) (oldv interface{}) {
	respc := make(chan response)
	c <- request{
		reqType: Delete,
		key:     key,
		respc:   respc,
	}
	resp := <-respc
	return resp.value
}

func endTransaction(c chan request) {
	c <- request{
		reqType: EndTransaction,
	}
}
