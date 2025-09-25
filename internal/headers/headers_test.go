package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestHeadersParse(t *testing.T) {

	// Test: Valid single header
	h := Headers{}
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid spacing header
	h = Headers{}
	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	h = Headers{}
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: End of headers
	h = Headers{}
	data = []byte("\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: Already existing header
	h = Headers{}
	h["Random-Header"] = "RandomValue"
	assert.Equal(t, "RandomValue", h["Random-Header"])
	data = []byte("User-Agent: curl/7.81.0\r\n\r\n")
	n, done, err = h.Parse(data)
	assert.Equal(t, "RandomValue", h["Random-Header"])
	assert.Equal(t, "curl/7.81.0", h["user-agent"])
	require.NoError(t, err)
	assert.Equal(t, 25, n)
	assert.False(t, done)

	//Test: Invalid header with non-ASCII character in key
	h = Headers{}
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: Add multiple headers
	h = Headers{}
	data = []byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, "curl/7.81.0", h["user-agent"])
	assert.Equal(t, "*/*", h["accept"])
	assert.Equal(t, 61, n)
	assert.False(t, done)

}


