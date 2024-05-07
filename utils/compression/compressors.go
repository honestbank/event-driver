package compression

import (
	"bytes"
	"compress/gzip"
	"io"
)

type Compressor interface {
	Compress(input []byte) ([]byte, error)
	Decompress(input []byte) ([]byte, error)
}

type gzipCompressor struct {
	level int
}

func Gzip(level int) Compressor {
	return gzipCompressor{
		level: level,
	}
}

func (g gzipCompressor) Compress(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return input, nil
	}
	var buffer bytes.Buffer
	gzipWriter, _ := gzip.NewWriterLevel(&buffer, g.level)

	_, err := gzipWriter.Write(input)
	if err != nil {
		return nil, err
	}
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), err
}

func (g gzipCompressor) Decompress(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return input, nil
	}
	byteReader := bytes.NewReader(input)
	gzipReader, err := gzip.NewReader(byteReader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	return io.ReadAll(gzipReader)
}

type noop struct {
}

func Noop() Compressor {
	return noop{}
}

func (n noop) Compress(input []byte) ([]byte, error) {
	return input, nil
}

func (n noop) Decompress(input []byte) ([]byte, error) {
	return input, nil
}
