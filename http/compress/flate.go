package compress

import (
	"compress/flate"
	"errors"
	"io"
	"sync"
)

type FlateBodyDecompressorPool struct {
	readerPool sync.Pool
}

func NewFlateBodyDecompressorPool() *FlateBodyDecompressorPool {
	return &FlateBodyDecompressorPool{
		readerPool: sync.Pool{New: func() any { return &flateWrapper{r: flate.NewReader(nil)} }},
	}
}

func (d *FlateBodyDecompressorPool) WrapBody(compressBody io.ReadCloser) io.ReadCloser {
	return &decompressorBodyWrapper[*flateWrapper]{
		CompressedBody: compressBody,
		GetDecoderFn: func(compressedBody io.ReadCloser) (*flateWrapper, error) {
			gz, _ := d.readerPool.Get().(*flateWrapper)
			return gz, gz.Reset(compressedBody)
		},
		ReturnDecoderFn: func(decoder *flateWrapper) error {
			d.readerPool.Put(decoder)
			return nil
		},
	}
}

type FlateBodyDecompressor struct{}

func NewFlateBodyDecompressor() *FlateBodyDecompressor {
	return &FlateBodyDecompressor{}
}

func (d *FlateBodyDecompressor) WrapBody(compressBody io.ReadCloser) io.ReadCloser {
	return &decompressorBodyWrapper[*flateWrapper]{
		CompressedBody: compressBody,
		GetDecoderFn: func(compressedBody io.ReadCloser) (*flateWrapper, error) {
			return &flateWrapper{r: flate.NewReader(compressedBody)}, nil
		},
		ReturnDecoderFn: func(_ *flateWrapper) error {
			return nil
		},
	}
}

type flateWrapper struct {
	r io.ReadCloser
}

func (f *flateWrapper) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}

func (f *flateWrapper) Close() error {
	return f.r.Close()
}

func (f *flateWrapper) Reset(compressedBody io.Reader) error {
	if rs, is := f.r.(interface {
		Reset(r io.Reader, dict []byte) error
	}); is {
		return rs.Reset(compressedBody, nil)
	}

	return ErrDeflateMissingReset
}

var ErrDeflateMissingReset = errors.New("the internal flate reader does not implement reset")
