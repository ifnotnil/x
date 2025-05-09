package compress

import (
	"io"
	"net/http"
	"strings"
)

func NewRoundTripper(next http.RoundTripper, opts ...RoundTripperOption) *RoundTripper {
	rt := &RoundTripper{
		next:                 next,
		omitCondition:        nil,
		acceptEncodingHeader: "",
		contentDecoders:      map[string]BodyDecoder{},
	}

	for _, o := range opts {
		o(rt)
	}

	if len(rt.contentDecoders) == 0 {
		rt.defaultInit()
	}

	return rt
}

type BodyDecoder interface {
	WrapBody(body io.ReadCloser) io.ReadCloser
}

const (
	contentEncodingGZIP    = "gzip"
	contentEncodingZSTD    = "zstd"
	contentEncodingBR      = "br"
	contentEncodingDeflate = "deflate"
)

const AcceptEncoding = "Accept-Encoding"

type OmitCondition func(req *http.Request) bool

type RoundTripperOption func(c *RoundTripper)

func WithOmitCondition(ec OmitCondition) RoundTripperOption {
	return func(c *RoundTripper) {
		c.omitCondition = ec
	}
}

func WithCompressionType(contentEncoding string, decompressor BodyDecoder) RoundTripperOption {
	return func(c *RoundTripper) {
		c.addDecoder(contentEncoding, decompressor)
	}
}

func WithCompressionTypeDeflate(useReaderPool bool) RoundTripperOption {
	return func(c *RoundTripper) {
		var d BodyDecoder
		if useReaderPool {
			d = NewFlateBodyDecompressorPool()
		} else {
			d = NewFlateBodyDecompressor()
		}

		c.addDecoder(contentEncodingDeflate, d)
	}
}

func WithCompressionTypeGZIP(useReaderPool bool) RoundTripperOption {
	return func(c *RoundTripper) {
		var d BodyDecoder
		if useReaderPool {
			d = NewGZIPBodyDecompressorPool()
		} else {
			d = NewGZIPBodyDecompressor()
		}

		c.addDecoder(contentEncodingGZIP, d)
	}
}

func WithCompressionTypeZSTD(useReaderPool bool) RoundTripperOption {
	return func(c *RoundTripper) {
		var d BodyDecoder
		if useReaderPool {
			d = NewZSTDBodyDecompressorPool()
		} else {
			d = NewZSTDBodyDecompressor()
		}

		c.addDecoder(contentEncodingZSTD, d)
	}
}

func WithCompressionTypeBR(useReaderPool bool) RoundTripperOption {
	return func(c *RoundTripper) {
		var d BodyDecoder
		if useReaderPool {
			d = NewBRBodyDecompressorPool()
		} else {
			d = NewBRBodyDecompressor()
		}

		c.addDecoder(contentEncodingBR, d)
	}
}

func WithAcceptEncoding(ae string) RoundTripperOption {
	return func(c *RoundTripper) {
		c.acceptEncodingHeader = ae
	}
}

func KeepContentHeaders() RoundTripperOption {
	return func(c *RoundTripper) {
		c.keepHeaders = true
	}
}

type RoundTripper struct {
	next                 http.RoundTripper
	omitCondition        OmitCondition
	contentDecoders      map[string]BodyDecoder
	acceptEncodingHeader string
	contentEncodings     []string
	keepHeaders          bool
}

func (rt *RoundTripper) addDecoder(contentEncoding string, decoder BodyDecoder) {
	contentEncoding = strings.ToLower(contentEncoding)
	_, exists := rt.contentDecoders[contentEncoding]
	rt.contentDecoders[contentEncoding] = decoder
	if !exists {
		rt.contentEncodings = append(rt.contentEncodings, contentEncoding)
		rt.initAcceptEncodingHeader()
	}
}

func (rt *RoundTripper) initAcceptEncodingHeader() {
	rt.acceptEncodingHeader = strings.Join(rt.contentEncodings, ", ")
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// check exclude conditions.
	if rt.omit(req) {
		return rt.next.RoundTrip(req)
	}

	// if no content encoding is selected, omit.
	if len(rt.contentDecoders) == 0 {
		return rt.next.RoundTrip(req)
	}
	req.Header.Set(AcceptEncoding, rt.acceptEncodingHeader)

	resp, err := rt.next.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	contentEncoding := strings.ToLower(resp.Header.Get("Content-Encoding"))
	decompressor, exists := rt.contentDecoders[contentEncoding]
	if !exists {
		return resp, nil
	}

	resp.Body = decompressor.WrapBody(resp.Body)

	if !rt.keepHeaders {
		resp.Header.Del("Content-Encoding")
		resp.Header.Del("Content-Length")
	}

	resp.ContentLength = -1
	resp.Uncompressed = true

	return resp, nil
}

func (rt *RoundTripper) omit(req *http.Request) bool {
	return req.Header.Get(AcceptEncoding) != "" ||
		req.Header.Get("Range") != "" ||
		req.Method == http.MethodHead ||
		(rt.omitCondition != nil && rt.omitCondition(req))
}

func (rt *RoundTripper) defaultInit() {
	WithCompressionTypeGZIP(true)(rt)
	WithCompressionTypeZSTD(true)(rt)
	WithCompressionTypeBR(true)(rt)
}
