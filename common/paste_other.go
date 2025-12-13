//go:build !darwin && !windows
// +build !darwin,!windows

package common

// InitAppSwitchListener 其他平台空实现
func InitAppSwitchListener() {
}

// RecordPreviousAppPID 其他平台空实现
func RecordPreviousAppPID() {
}

// ActivateAppByPID 其他平台空实现
func ActivateAppByPID(pid int) bool {
	return false
}

// PasteCmdV 其他平台空实现
func PasteCmdV() {
}

// ActivatePreviousApp 其他平台空实现
func ActivatePreviousApp() {
}

// PasteCmdVToPreviousApp 其他平台空实现
func PasteCmdVToPreviousApp() {
}

// TestPasteFunction 其他平台空实现
func TestPasteFunction() {
}

// SimpleTestPaste 其他平台空实现
func SimpleTestPaste() {
}

// TestLogWriting 其他平台空实现
func TestLogWriting() {
}
