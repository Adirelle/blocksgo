package ipc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	assert := assert.New(t)

	buf := bytes.Buffer{}
	val := []int{4, 5}
	err := writeJSON(&buf, &val)

	assert.Nil(err)
	assert.JSONEq("[4,5]\n", buf.String())
}

func TestBlockIdentifierMatches(t *testing.T) {
	assert := assert.New(t)

	a := BlockIdentifier{"", ""}
	b := BlockIdentifier{"x", ""}
	c := BlockIdentifier{"", "y"}
	d := BlockIdentifier{"x", "y"}
	e := BlockIdentifier{"x", "z"}

	assert.False(a.Matches(a))
	assert.False(a.Matches(b))
	assert.False(a.Matches(c))
	assert.False(b.Matches(d))
	assert.False(b.Matches(e))
	assert.False(c.Matches(d))
	assert.False(c.Matches(e))
	assert.True(d.Matches(d))
	assert.False(d.Matches(e))
	assert.True(e.Matches(e))
}
