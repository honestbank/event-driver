package compression_test

import (
	"compress/gzip"
	"testing"

	"github.com/lukecold/event-driver/utils/compression"
	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	input := []byte("test string")
	compressor := compression.Gzip(gzip.BestSpeed)

	compressed, err := compressor.Compress(input)
	assert.NoError(t, err)
	assert.NotEqual(t, input, compressed)

	roundTripped, err := compressor.Decompress(compressed)
	assert.NoError(t, err)
	assert.Equal(t, input, roundTripped)
}

func TestNoop(t *testing.T) {
	input := []byte("test string")
	compressor := compression.Noop()

	compressed, err := compressor.Compress(input)
	assert.NoError(t, err)
	assert.Equal(t, input, compressed)

	roundTripped, err := compressor.Decompress(compressed)
	assert.NoError(t, err)
	assert.Equal(t, input, roundTripped)
}
