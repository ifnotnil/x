package compress

import (
	"errors"
	"io"
	"sync"
)

type Decompressor interface {
	io.Reader
	Reset(compressedBody io.Reader) error
}

type decompressorBodyWrapper[D Decompressor] struct {
	CompressedBody  io.ReadCloser
	GetDecoderFn    func(compressedBody io.ReadCloser) (D, error)
	ReturnDecoderFn func(decoder D) error

	stickyError error
	decoder     D
	onceInit    sync.Once
	onceClose   sync.Once
}

func (d *decompressorBodyWrapper[D]) initDecoder() {
	d.onceInit.Do(func() {
		var err error
		d.decoder, err = d.GetDecoderFn(d.CompressedBody)
		if err != nil {
			d.stickyError = err
		}
	})
}

func (d *decompressorBodyWrapper[D]) closeDecoder() {
	d.onceClose.Do(func() {
		var dErr, rErr error

		if cl, is := Decompressor(d.decoder).(io.Closer); is {
			if err := cl.Close(); err != nil {
				dErr = err
			}
		}

		rErr = d.ReturnDecoderFn(d.decoder)

		d.stickyError = errors.Join(dErr, rErr)
	})
}

func (d *decompressorBodyWrapper[D]) Read(p []byte) (int, error) {
	d.initDecoder()

	if d.stickyError != nil {
		return 0, d.stickyError
	}

	return d.decoder.Read(p)
}

func (d *decompressorBodyWrapper[D]) Close() error {
	d.closeDecoder()
	return errors.Join(d.CompressedBody.Close(), d.stickyError)
}
