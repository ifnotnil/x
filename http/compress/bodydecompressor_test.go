package compress

import (
	"bytes"
	_ "embed"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/compressed/file_1.txt.gz
var file1GZIPCompressed []byte

//go:embed testdata/compressed/file_1.txt.zst
var file1ZSTDCompressed []byte

//go:embed testdata/compressed/file_1.txt.br
var file1BRCompressed []byte

//go:embed testdata/compressed/file_1.txt
var file1Original []byte

func TestBodyDecompressor(t *testing.T) {
	tests := map[string]struct {
		BodyDecompressor BodyDecoder
		Compressed       []byte
		Expected         []byte
		ExpectedError    assert.ErrorAssertionFunc
	}{
		"gzip": {
			BodyDecompressor: NewGZIPBodyDecompressor(),
			Compressed:       file1GZIPCompressed,
			ExpectedError:    assert.NoError,
		},
		"gzip pool": {
			BodyDecompressor: NewGZIPBodyDecompressorPool(),
			Compressed:       file1GZIPCompressed,
			ExpectedError:    assert.NoError,
		},
		"zstd": {
			BodyDecompressor: NewZSTDBodyDecompressor(),
			Compressed:       file1ZSTDCompressed,
			ExpectedError:    assert.NoError,
		},
		"zstd pool": {
			BodyDecompressor: NewZSTDBodyDecompressorPool(),
			Compressed:       file1ZSTDCompressed,
			ExpectedError:    assert.NoError,
		},
		"br": {
			BodyDecompressor: NewBRBodyDecompressor(),
			Compressed:       file1BRCompressed,
			ExpectedError:    assert.NoError,
		},
		"br pool": {
			BodyDecompressor: NewBRBodyDecompressorPool(),
			Compressed:       file1BRCompressed,
			ExpectedError:    assert.NoError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			wb := tc.BodyDecompressor.WrapBody(io.NopCloser(bytes.NewBuffer(tc.Compressed)))

			got, err := io.ReadAll(wb)
			tc.ExpectedError(t, err)

			err = wb.Close()
			require.NoError(t, err)

			assert.Equal(t, string(file1Original), string(got))
		})
	}
}
