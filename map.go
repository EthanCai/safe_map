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

	go result.handleRequest()

	return result
}

func (s *SafeMap) handleRequest() {
	for {
		req := <-s.c
		switch req.reqType {
		case Get:
			resp := response{}
			resp.value, resp.ok = s.m[req.key]
			req.respc <- resp
		case Set:
			resp := response{}
			resp.value, _ = s.m[req.key]
			s.m[req.key] = req.value
			req.respc <- resp
		case Delete:
			resp := response{}
			resp.value, _ = s.m[req.key]
			delete(s.m, req.key)
			req.respc <- resp
		}
	}
}

// Get the value by key
func (s *SafeMap) Get(key string) (interface{}, bool) {
	c := make(chan response)
	s.c <- request{
		reqType: Get,
		key:     key,
		respc:   c,
	}
	resp := <-c
	return resp.value, resp.ok
}

// Set value by key, return old value if exists
func (s *SafeMap) Set(key string, value interface{}) (oldv interface{}) {
	c := make(chan response)
	s.c <- request{
		reqType: Set,
		key:     key,
		value:   value,
		respc:   c,
	}
	resp := <-c
	return resp.value
}

// Delete a entry by key, return old value if exists
func (s *SafeMap) Delete(key string) (oldv interface{}) {
	c := make(chan response)
	s.c <- request{
		reqType: Delete,
		key:     key,
		respc:   c,
	}
	resp := <-c
	return resp.value
}
