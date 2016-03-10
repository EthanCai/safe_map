package safemap

import ()

type SafeMap struct {
	m map[string]interface{}
	c chan request
}

// You should always use New() to create a SafeMap
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
