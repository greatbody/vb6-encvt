package detector

import (
	"io"
	"os"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Encoding string

const (
	UTF8    Encoding = "UTF-8"
	GBK     Encoding = "GBK"
	Unknown Encoding = "Unknown"
)

// Detect determines the encoding of the provided byte slice.
func Detect(data []byte) Encoding {
	// 1. Check for basic UTF-8 validity
	isUtf8 := utf8.Valid(data)

	// 2. Check for GBK validity
	isGbk := canDecodeGBK(data)

	if isUtf8 && isGbk {
		// If it's valid in both, it's likely ASCII or standard characters.
		// We default to UTF-8 as it is the target.
		return UTF8
	}

	if isUtf8 {
		return UTF8
	}

	if isGbk {
		return GBK
	}

	return Unknown
}

// canDecodeGBK checks if the data can be cleanly decoded using GBK decoder.
// It relies on the strictness of the transformation.
func canDecodeGBK(data []byte) bool {
	// We use a small buffer logic or just transform everything?
	// transform.Bytes is convenient.
	// Note: GBK decoder might be permissive.
	// If the file is NOT UTF-8, and we are checking GBK, we hope it fails if it's random binary?
	// But we already filtered binary in Walker.
	// So if it's "Unknown", it's likely some other encoding or malformed.

	// Create a transformer that strictly checks?
	// The standard simplifiedchinese.GBK decoder is generally robust.
	_, _, err := transform.Bytes(simplifiedchinese.GBK.NewDecoder(), data)
	return err == nil
}

// DetectFile reads the file and detects its encoding.
func DetectFile(path string) (Encoding, error) {
	f, err := os.Open(path)
	if err != nil {
		return Unknown, err
	}
	defer f.Close()

	// Limit read to avoid OOM on unexpectedly large files, though we filter binary.
	// 50MB should be plenty for any source file.
	const maxRead = 50 * 1024 * 1024
	content, err := io.ReadAll(io.LimitReader(f, maxRead))
	if err != nil {
		return Unknown, err
	}

	if len(content) == 0 {
		return UTF8, nil // Empty file is valid UTF-8
	}

	return Detect(content), nil
}
