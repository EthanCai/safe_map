package safemap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	assert := assert.New(t)

	sm := New()
	assert.NotNil(sm)

	trans := sm.BeginTransaction()
	trans.Set("k1", "v1")
	trans.Set("k2", "v2")
	trans.Set("k3", "v3")
	oldv := trans.Del("k3")
	assert.Equal("v3", oldv)
	trans.EndTransaction()

	oldv, ok := sm.Get("k1")
	if assert.True(ok) {
		assert.Equal("v1", oldv)
	}
}

func TestMultiTransactions(t *testing.T) {
	assert := assert.New(t)

	sm := New()
	assert.NotNil(sm)

	go func() {
		trans := sm.BeginTransaction()
		trans.Set("k1", "v1.1")
		time.Sleep(100 * time.Millisecond)
		trans.Set("k2", "v2.1")
		time.Sleep(100 * time.Millisecond)
		trans.Set("k3", "v3.1")
		time.Sleep(100 * time.Millisecond)
		trans.EndTransaction()
	}()

	go func() {
		trans := sm.BeginTransaction()
		time.Sleep(100 * time.Millisecond)
		trans.Set("k1", "v1.2")
		time.Sleep(100 * time.Millisecond)
		trans.Set("k2", "v2.2")
		time.Sleep(100 * time.Millisecond)
		trans.Set("k3", "v3.2")
		trans.EndTransaction()
	}()

	time.Sleep(1 * time.Second)
	v1, _ := sm.Get("k1")
	v2, _ := sm.Get("k2")
	v3, _ := sm.Get("k3")
	if v1 == "v1.1" {
		assert.Equal("v2.1", v2)
		assert.Equal("v3.1", v3)
	} else {
		assert.Equal("v2.2", v2)
		assert.Equal("v3.2", v3)
	}
}
