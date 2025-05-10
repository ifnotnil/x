package compress

import (
	"io"
	"sync"

	"github.com/klauspost/compress/zstd"
)

const defaultZSTDDecoderMaxWindow = 128 << 20

func defaultZSTDDecoder(r io.Reader) (*zstd.Decoder, error) {
	return zstd.NewReader(
		r,
		zstd.WithDecoderLowmem(true),
		zstd.WithDecoderMaxWindow(defaultZSTDDecoderMaxWindow),
		zstd.WithDecoderConcurrency(1),
	)
}

type zstdDecoderWrapper struct {
	decoder   *zstd.Decoder
	initError error
}

type ZSTDBodyDecompressorPool struct {
	readerPool sync.Pool
}

func NewZSTDBodyDecompressorPool() *ZSTDBodyDecompressorPool {
	return &ZSTDBodyDecompressorPool{
		readerPool: sync.Pool{New: func() any {
			w := &zstdDecoderWrapper{}
			w.decoder, w.initError = defaultZSTDDecoder(nil)
			return w
		}},
	}
}

func (d *ZSTDBodyDecompressorPool) WrapBody(compressBody io.ReadCloser) io.ReadCloser {
	return &decompressorBodyWrapper[*zstd.Decoder]{
		CompressedBody: compressBody,
		GetDecoderFn: func(compressedBody io.ReadCloser) (*zstd.Decoder, error) {
			w, _ := d.readerPool.Get().(*zstdDecoderWrapper)
			if w.initError != nil {
				return nil, w.initError
			}
			return w.decoder, w.decoder.Reset(compressedBody)
		},
		ReturnDecoderFn: func(decoder *zstd.Decoder) error {
			d.readerPool.Put(&zstdDecoderWrapper{decoder: decoder, initError: nil})
			return nil
		},
	}
}

type ZSTDBodyDecompressor struct{}

func NewZSTDBodyDecompressor() *ZSTDBodyDecompressor {
	return &ZSTDBodyDecompressor{}
}

func (d *ZSTDBodyDecompressor) WrapBody(compressBody io.ReadCloser) io.ReadCloser {
	return &decompressorBodyWrapper[*zstd.Decoder]{
		CompressedBody: compressBody,
		GetDecoderFn: func(compressedBody io.ReadCloser) (*zstd.Decoder, error) {
			return defaultZSTDDecoder(compressedBody)
		},
		ReturnDecoderFn: func(_ *zstd.Decoder) error { return nil },
	}
}
