package compress

import (
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
)

type GZIPBodyDecompressorPool struct {
	readerPool sync.Pool
}

func NewGZIPBodyDecompressorPool() *GZIPBodyDecompressorPool {
	return &GZIPBodyDecompressorPool{
		readerPool: sync.Pool{New: func() any { return &gzip.Reader{} }},
	}
}

func (d *GZIPBodyDecompressorPool) WrapBody(compressBody io.ReadCloser) io.ReadCloser {
	return &decompressorBodyWrapper[*gzip.Reader]{
		CompressedBody: compressBody,
		GetDecoderFn: func(compressedBody io.ReadCloser) (*gzip.Reader, error) {
			gz, _ := d.readerPool.Get().(*gzip.Reader)
			return gz, gz.Reset(compressedBody)
		},
		ReturnDecoderFn: func(decoder *gzip.Reader) error {
			d.readerPool.Put(decoder)
			return nil
		},
	}
}

type GZIPBodyDecompressor struct{}

func NewGZIPBodyDecompressor() *GZIPBodyDecompressor {
	return &GZIPBodyDecompressor{}
}

func (d *GZIPBodyDecompressor) WrapBody(compressBody io.ReadCloser) io.ReadCloser {
	return &decompressorBodyWrapper[*gzip.Reader]{
		CompressedBody: compressBody,
		GetDecoderFn: func(compressedBody io.ReadCloser) (*gzip.Reader, error) {
			return gzip.NewReader(compressedBody)
		},
		ReturnDecoderFn: func(_ *gzip.Reader) error {
			return nil
		},
	}
}
