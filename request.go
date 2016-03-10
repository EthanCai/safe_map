package safemap

const (
	// request type
	Get = iota
	Set
	Delete
)

type request struct {
	reqType int
	key     string
	value   interface{}
	respc   chan response
}

type response struct {
	value interface{}
	ok    bool
}
