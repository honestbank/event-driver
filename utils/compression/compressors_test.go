package compression_test

import (
	"compress/gzip"
	"testing"

	"github.com/honestbank/event-driver/utils/compression"
	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	input := []byte("test string")
	compressor := compression.Gzip(gzip.BestSpeed)

	t.Run("round trip", func(t *testing.T) {
		compressed, err := compressor.Compress(input)
		assert.NoError(t, err)
		assert.NotEqual(t, input, compressed)

		roundTripped, err := compressor.Decompress(compressed)
		assert.NoError(t, err)
		assert.Equal(t, input, roundTripped)
	})

	t.Run("empty bytes", func(t *testing.T) {
		compressed, err := compressor.Compress([]byte{})
		assert.NoError(t, err)
		assert.Equal(t, []byte{}, compressed)

		decompressed, err := compressor.Decompress([]byte{})
		assert.NoError(t, err)
		assert.Equal(t, []byte{}, decompressed)
	})
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
