// Test: Valid single header
package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Content-Type: application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers["Content-Type"])
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:       localhost:42069     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Content-Type:     application/json    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers["Content-Type"])
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Content-Type          : application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// TODO: Work in the Code Till it Pass This Test
	// // Test: Valid 2 headers with existing headers
	// headers = NewHeaders()
	// data = slices.Concat([]byte("Content-Type: application/json\r\n"), []byte("Host: localhost:42069\r\n\r\n"))
	// n, done, err = headers.Parse(data)
	// require.NoError(t, err)
	// require.NotNil(t, headers)
	// assert.Equal(t, "application/json", headers["Content-Type"])
	// assert.Equal(t, "localhost:42069", headers["Host"])
	// assert.Equal(t, len(data)-2, n)
	// assert.False(t, done)
}
