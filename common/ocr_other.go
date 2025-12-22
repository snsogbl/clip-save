//go:build !darwin
// +build !darwin

package common

// RecognizeTextInImage 非 macOS 平台暂不支持 OCR
func RecognizeTextInImage(imageData []byte) string {
	return ""
}
