//go:build windows
// +build windows

package common

import (
	"encoding/json"
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// Windows APIs/DLLs
var (
	modUser32   = syscall.NewLazyDLL("user32.dll")
	modKernel32 = syscall.NewLazyDLL("kernel32.dll")
	modPsapi    = syscall.NewLazyDLL("psapi.dll")
	modShell32  = syscall.NewLazyDLL("shell32.dll")
	modGdi32    = syscall.NewLazyDLL("gdi32.dll") // not strictly required, kept for completeness

	procGetForegroundWindow        = modUser32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId   = modUser32.NewProc("GetWindowThreadProcessId")
	procOpenProcess                = modKernel32.NewProc("OpenProcess")
	procCloseHandle                = modKernel32.NewProc("CloseHandle")
	procQueryFullProcessImageName  = modKernel32.NewProc("QueryFullProcessImageNameW")
	procGetClipboardSequenceNumber = modUser32.NewProc("GetClipboardSequenceNumber")
	procOpenClipboard              = modUser32.NewProc("OpenClipboard")
	procCloseClipboard             = modUser32.NewProc("CloseClipboard")
	procGetClipboardData           = modUser32.NewProc("GetClipboardData")
	procEmptyClipboard             = modUser32.NewProc("EmptyClipboard")
	procSetClipboardData           = modUser32.NewProc("SetClipboardData")
	procGlobalAlloc                = modKernel32.NewProc("GlobalAlloc")
	procGlobalLock                 = modKernel32.NewProc("GlobalLock")
	procGlobalUnlock               = modKernel32.NewProc("GlobalUnlock")
	procDragQueryFile              = modShell32.NewProc("DragQueryFileW")
)

const (
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	GMEM_MOVEABLE                     = 0x0002
	CF_HDROP                          = 15
)

// DROPFILES structure (winuser.h)
type dropfiles struct {
	pFiles uint32 // offset of file list (in bytes) from beginning of this struct
	x      int32
	y      int32
	fNC    int32
	fWide  int32 // nonzero if Unicode
}

// GetFrontmostAppName 返回当前前台窗口所属进程的可执行名（或完整路径的末段）
func GetFrontmostAppName() string {
	// HWND GetForegroundWindow(void);
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return "Unknown"
	}

	// DWORD GetWindowThreadProcessId(HWND hWnd, LPDWORD lpdwProcessId);
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return "Unknown"
	}

	// HANDLE OpenProcess(DWORD dwDesiredAccess, BOOL bInheritHandle, DWORD dwProcessId);
	hProc, _, _ := procOpenProcess.Call(PROCESS_QUERY_LIMITED_INFORMATION, 0, uintptr(pid))
	if hProc == 0 {
		return "Unknown"
	}
	defer procCloseHandle.Call(hProc)

	// BOOL QueryFullProcessImageNameW(HANDLE hProcess, DWORD dwFlags, LPWSTR lpExeName, PDWORD lpdwSize);
	buf := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(buf))
	ret, _, _ := procQueryFullProcessImageName.Call(hProc, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if ret == 0 || size == 0 {
		return "Unknown"
	}
	full := syscall.UTF16ToString(buf[:size])
	// 取文件名部分
	for i := len(full) - 1; i >= 0; i-- {
		if full[i] == '\\' || full[i] == '/' {
			if i+1 < len(full) {
				return full[i+1:]
			}
			break
		}
	}
	return full
}

// GetPasteboardChangeCount 使用 Windows 的剪贴板序列号
func GetPasteboardChangeCount() int {
	seq, _, _ := procGetClipboardSequenceNumber.Call()
	return int(seq)
}

// ReadPasteboardData Windows 下不需要（由 golang.design/x/clipboard 兜底图片读取），返回 nil
func ReadPasteboardData(typeName string) []byte {
	return nil
}

// ReadFileURLs 读取剪贴板中的文件路径（CF_HDROP）
func ReadFileURLs() (string, int) {
	// BOOL OpenClipboard(HWND hWndNewOwner);
	if ok, _, _ := procOpenClipboard.Call(0); ok == 0 {
		return "", 0
	}
	defer procCloseClipboard.Call()

	// HANDLE GetClipboardData(UINT uFormat);
	hDrop, _, _ := procGetClipboardData.Call(CF_HDROP)
	if hDrop == 0 {
		return "", 0
	}

	// UINT DragQueryFileW(HDROP hDrop, UINT iFile, LPWSTR lpszFile, UINT cch)
	// 先获取文件数量
	count, _, _ := procDragQueryFile.Call(hDrop, 0xFFFFFFFF, 0, 0)
	if count == 0 {
		return "", 0
	}

	paths := make([]string, 0, count)
	for i := uintptr(0); i < count; i++ {
		// 查询长度（包含终止符）
		nameLen, _, _ := procDragQueryFile.Call(hDrop, i, 0, 0)
		if nameLen == 0 {
			continue
		}
		// 为终止符多分配 1，并将该大小传入 cch，避免最后一个字符被截断
		buf := make([]uint16, int(nameLen)+1)
		procDragQueryFile.Call(hDrop, i, uintptr(unsafe.Pointer(&buf[0])), nameLen+1)
		paths = append(paths, syscall.UTF16ToString(buf))
	}

	if len(paths) == 0 {
		return "", 0
	}
	jsonBytes, err := json.Marshal(paths)
	if err != nil {
		return "", 0
	}
	return string(jsonBytes), len(paths)
}

// WriteFileURLs 将文件路径数组写入剪贴板（CF_HDROP）
func WriteFileURLs(filePaths string) error {
	var paths []string
	if err := json.Unmarshal([]byte(filePaths), &paths); err != nil {
		return fmt.Errorf("解析文件路径 JSON 失败: %w", err)
	}
	if len(paths) == 0 {
		return fmt.Errorf("空文件列表")
	}

	// 计算所需内存：DROPFILES + UTF-16 文件列表，文件之间以'\u0000'分隔，末尾以额外'\u0000'终止
	utf16Lists := make([][]uint16, 0, len(paths))
	var chars uint32
	for _, p := range paths {
		u := utf16.Encode([]rune(p))
		utf16Lists = append(utf16Lists, append(u, 0)) // 每个以 0 结尾
		chars += uint32(len(u) + 1)
	}
	chars++ // 额外的终止 0

	bytesNeeded := uint32(unsafe.Sizeof(dropfiles{})) + chars*2

	// BOOL OpenClipboard(HWND)
	if ok, _, _ := procOpenClipboard.Call(0); ok == 0 {
		return fmt.Errorf("打开剪贴板失败")
	}
	defer procCloseClipboard.Call()

	// 置空再写入
	procEmptyClipboard.Call()

	// HGLOBAL GlobalAlloc(UINT uFlags, SIZE_T dwBytes)
	hMem, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, uintptr(bytesNeeded))
	if hMem == 0 {
		return fmt.Errorf("内存分配失败")
	}

	// LPVOID GlobalLock(HGLOBAL hMem)
	pMem, _, _ := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return fmt.Errorf("内存锁定失败")
	}

	// 写入 DROPFILES 头
	head := (*dropfiles)(unsafe.Pointer(pMem))
	head.pFiles = uint32(unsafe.Sizeof(dropfiles{}))
	head.x = 0
	head.y = 0
	head.fNC = 0
	head.fWide = 1 // Unicode

	// 写入 UTF-16 路径块
	dataStart := uintptr(pMem) + uintptr(unsafe.Sizeof(dropfiles{}))
	cursor := dataStart
	for _, u16 := range utf16Lists {
		bytes := u16ToBytes(u16)
		memcpy(unsafe.Pointer(cursor), unsafe.Pointer(&bytes[0]), uintptr(len(bytes)))
		cursor += uintptr(len(bytes))
	}
	// 末尾再补一个 0
	zero := []byte{0, 0}
	memcpy(unsafe.Pointer(cursor), unsafe.Pointer(&zero[0]), uintptr(len(zero)))

	procGlobalUnlock.Call(hMem)

	// HANDLE SetClipboardData(UINT uFormat, HANDLE hMem)
	ret, _, _ := procSetClipboardData.Call(CF_HDROP, hMem)
	if ret == 0 {
		return fmt.Errorf("设置剪贴板数据失败")
	}
	// 注意：成功后内存归系统管理，不应再释放 hMem
	return nil
}

// Helpers
func u16ToBytes(u16 []uint16) []byte {
	if len(u16) == 0 {
		return nil
	}
	b := make([]byte, len(u16)*2)
	for i, v := range u16 {
		b[2*i] = byte(v)
		b[2*i+1] = byte(v >> 8)
	}
	return b
}

// minimal memcpy using Go unsafe
func memcpy(dst, src unsafe.Pointer, n uintptr) {
	// convert to byte slices backed by the same memory
	d := (*[1 << 30]byte)(dst)[:n:n]
	s := (*[1 << 30]byte)(src)[:n:n]
	copy(d, s)
}
