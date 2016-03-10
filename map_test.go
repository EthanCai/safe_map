package safemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	sm := New()
	assert.NotNil(sm)
	assert.NotNil(sm.m)
	assert.NotNil(sm.c)
}

func TestGetAndSet(t *testing.T) {
	assert := assert.New(t)

	sm := New()

	v, ok := sm.Get("k1")
	if assert.False(ok) {
		assert.Nil(v)
	}

	oldv := sm.Set("k1", "v1")
	assert.Nil(oldv)
	assert.Equal("v1", sm.m["k1"])

	v, ok = sm.Get("k1")
	if assert.True(ok) {
		assert.Equal("v1", v)
	}

	oldv = sm.Set("k1", "v1.1")
	assert.Equal("v1", oldv)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	sm := New()

	oldv := sm.Set("k1", "v1")
	assert.Nil(oldv)
	assert.Equal("v1", sm.m["k1"])

	oldv = sm.Delete("k1")
	assert.NotNil(oldv)
	assert.Equal("v1", oldv)
	assert.Equal(0, len(sm.m))
}
