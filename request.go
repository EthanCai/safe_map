package safemap

const (
	// request type
	Get = iota
	Set
	Delete
	BeginTransaction
	EndTransaction
)

type request struct {
	reqType int
	key     string
	value   interface{}
	respc   chan response
	tc      chan request // channel for transaction
}

type response struct {
	value interface{}
	ok    bool
}
