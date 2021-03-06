package win32

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var user32 = syscall.NewLazyDLL("user32.dll")

var (
	procEnumWindows          = user32.NewProc("EnumWindows")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
)

// EnumWindows 枚举窗口
func EnumWindows(cb func(hwnd windows.HWND, args uintptr) bool, args uintptr) (err error) {
	fn := func(hwnd windows.HWND, args uintptr) int {
		return fromBool(cb(hwnd, args))
	}

	if _, _, err = procEnumWindows.Call(syscall.NewCallback(fn), args); err != nil {
		err = fromError(err)
	}

	return
}

// GetWindowTextLengthW 获取窗口文本长度
func GetWindowTextLengthW(hwnd windows.HWND) (length uint, err error) {
	ret, _, err := procGetWindowTextLengthW.Call(uintptr(hwnd))
	return uint(ret), fromError(err)
}

// GetWindowTextW 获取窗口文本，指定最大长度
func GetWindowTextW(hwnd windows.HWND, maxCount uint) (text string, err error) {
	var buf = make([]uint16, maxCount)
	if _, _, err = procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(maxCount)); err != nil {
		err = fromError(err)
	}

	return syscall.UTF16ToString(buf), err
}

// GetWindowText 获取窗口文本
func GetWindowText(hwnd windows.HWND) (text string, err error) {
	length, err := GetWindowTextLengthW(hwnd)

	return GetWindowTextW(hwnd, length+1)
}

// GetProcessIDByName 根据进程名称获取进程ID
func GetProcessIDByName(name string) (pid []uint32, err error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))

	for err = windows.Process32First(snapshot, &pe); err == nil; err = windows.Process32Next(snapshot, &pe) {
		if windows.UTF16ToString(pe.ExeFile[:]) == name {
			pid = append(pid, pe.ProcessID)
		}
	}

	if errno, ok := err.(windows.Errno); ok && errno == windows.ERROR_NO_MORE_FILES {
		err = nil
	}

	return
}

// GetWindowThreadProcessID 根据窗口句柄获取线程ID和进程ID
func GetWindowThreadProcessID(hwnd windows.HWND) (tid, pid uint32, err error) {
	if tid, err = windows.GetWindowThreadProcessId(hwnd, &pid); err != nil {
		return
	}

	return
}

// GetWindowHwndByPID 获取进程窗口句柄
func GetWindowHwndByPID(pid uint32) (hwnds []windows.HWND) {
	_ = EnumWindows(func(hwnd windows.HWND, args uintptr) bool {
		var (
			err error
			pid uint32
		)

		if _, pid, err = GetWindowThreadProcessID(hwnd); err != nil {
			return true
		}

		if pid == uint32(args) {
			hwnds = append(hwnds, hwnd)
		}

		return true
	}, uintptr(pid))

	return
}
