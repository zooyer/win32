package win32

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	SE_BACKUP_NAME   = windows.StringToUTF16Ptr("SeBackupPrivilege")
	SE_RESTORE_NAME  = windows.StringToUTF16Ptr("SeRestorePrivilege")
	SE_SHUTDOWN_NAME = windows.StringToUTF16Ptr("SeShutdownPrivilege")
	SE_DEBUG_NAME    = windows.StringToUTF16Ptr("SeDebugPrivilege")
)

// readProcessMemory 读取内存数据到缓冲区
func readProcessMemory(handle windows.Handle, baseAddress uintptr, buf unsafe.Pointer, size int) (n uint, err error) {
	if err = windows.ReadProcessMemory(handle, baseAddress, (*byte)(buf), uintptr(size), (*uintptr)(unsafe.Pointer(&n))); err != nil {
		return
	}

	return
}

// ReadProcessMemory 读取内存数据到缓冲区
func ReadProcessMemory(handle windows.Handle, baseAddress uintptr, buf []byte) (n uint, err error) {
	return readProcessMemory(handle, baseAddress, unsafe.Pointer(&buf[0]), len(buf))
}

// ReadProcessMemoryData 读取内存数据到变量
func ReadProcessMemoryData[T any](handle windows.Handle, baseAddress uintptr) (val T, err error) {
	if _, err = readProcessMemory(handle, baseAddress, unsafe.Pointer(&val), int(unsafe.Sizeof(val))); err != nil {
		return
	}

	return
}

// ReadProcessMemoryCStringA 读取内存数据为ASCII字符串
func ReadProcessMemoryCStringA(handle windows.Handle, baseAddress uintptr) (str string, err error) {
	var utf8 = make([]byte, 0, 128)
	ptr, err := ReadProcessMemoryData[uintptr](handle, baseAddress)
	if err != nil {
		return
	}

	var char byte
	for {
		if char, err = ReadProcessMemoryData[byte](handle, ptr); err != nil {
			return
		}

		if char == 0x00 || char == 0xFF {
			break
		}

		utf8 = append(utf8, char)

		ptr++
	}

	return string(utf8), nil
}

// ReadProcessMemoryCStringW 读取内存数据为Unicode字符串
func ReadProcessMemoryCStringW(handle windows.Handle, baseAddress uintptr) (str string, err error) {
	var utf16 = make([]uint16, 0, 128)
	ptr, err := ReadProcessMemoryData[uintptr](handle, baseAddress)
	if err != nil {
		return
	}

	var wchar uint16
	for {
		if wchar, err = ReadProcessMemoryData[uint16](handle, ptr); err != nil {
			return
		}

		utf16 = append(utf16, wchar)

		if wchar == 0x00 || wchar == 0xFF {
			break
		}

		ptr += 2
	}

	return windows.UTF16ToString(utf16), nil
}

// EnablePrivilege 进程开启特权模式
func EnablePrivilege() (err error) {
	var (
		tkp   windows.Tokenprivileges
		token windows.Token
	)

	if err = windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY, &token); err != nil {
		return
	}

	defer token.Close()

	if err = windows.LookupPrivilegeValue(nil, SE_DEBUG_NAME, &tkp.Privileges[0].Luid); err != nil {
		return
	}

	tkp.PrivilegeCount = 1
	tkp.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED

	if err = windows.AdjustTokenPrivileges(token, false, &tkp, uint32(unsafe.Sizeof(tkp)), nil, nil); err != nil {
		return
	}

	return
}
