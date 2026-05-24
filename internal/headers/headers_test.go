package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:8182\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:8182", headers["host"])
	assert.Equal(t, 22, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:8182\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host: localhost:8182 \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test Valid 2 headers with existing header
	//headers = NewHeaders()
	//data = []byte("Host: localhost:8182\r\n\r\n")
	//n, done, err = headers.Parse(data)
	//n, done, err = headers.Parse(data)
	//require.NoError(t, err)
	//assert.Equal(t, 23, n)
	//assert.False(t, done)

	// Captial letters + checking if key is lower
	headers = NewHeaders()
	data = []byte("HOST: localhost:8182\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 22, n)
	assert.False(t, done)
	_, ok := headers["host"]
	assert.True(t, ok)

	// Invalid character
	headers = NewHeaders()
	data = []byte("H©st: localhost:8182\r\n\r\n")
	n, done, err = headers.Parse(data)
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Duplicate valid headers
	headers = NewHeaders()
	data = []byte("batates: san\r\n\r\n")
	n, done, err = headers.Parse(data)
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 14, n)
	assert.Equal(t, "san, san", headers["batates"])
	assert.False(t, done)

	// malformed headers
	headers = NewHeaders()
	data = []byte("Host localhost:8182\r\n\r\n")
	n, done, err = headers.Parse(data)
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
