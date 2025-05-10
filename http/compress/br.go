package compress

import (
	"io"
	"sync"

	"github.com/andybalholm/brotli"
)

type BRBodyDecompressorPool struct {
	readerPool sync.Pool
}

func NewBRBodyDecompressorPool() *BRBodyDecompressorPool {
	return &BRBodyDecompressorPool{
		readerPool: sync.Pool{New: func() any { return &brotli.Reader{} }},
	}
}

func (d *BRBodyDecompressorPool) WrapBody(compressBody io.ReadCloser) io.ReadCloser {
	return &decompressorBodyWrapper[*brotli.Reader]{
		CompressedBody: compressBody,
		GetDecoderFn: func(compressedBody io.ReadCloser) (*brotli.Reader, error) {
			br, _ := d.readerPool.Get().(*brotli.Reader)
			return br, br.Reset(compressedBody)
		},
		ReturnDecoderFn: func(decoder *brotli.Reader) error {
			d.readerPool.Put(decoder)
			return nil
		},
	}
}

type BRBodyDecompressor struct{}

func NewBRBodyDecompressor() *BRBodyDecompressor {
	return &BRBodyDecompressor{}
}

func (d *BRBodyDecompressor) WrapBody(compressBody io.ReadCloser) io.ReadCloser {
	return &decompressorBodyWrapper[*brotli.Reader]{
		CompressedBody: compressBody,
		GetDecoderFn: func(compressedBody io.ReadCloser) (*brotli.Reader, error) {
			return brotli.NewReader(compressedBody), nil
		},
		ReturnDecoderFn: func(_ *brotli.Reader) error {
			return nil
		},
	}
}
