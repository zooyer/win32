package win32

import (
	"fmt"
	"testing"

	"golang.org/x/sys/windows"
)

func TestEnumWindows(t *testing.T) {
	if err := EnumWindows(func(hwnd windows.HWND, args uintptr) bool {
		text, err := GetWindowText(hwnd)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(hwnd, text)

		return true
	}, 0); err != nil {
		t.Fatal(err)
	}
}

func TestGetProcessIDByName(t *testing.T) {
	pid, err := GetProcessIDByName("explorer.exe")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pid)
}

func TestGetWindowThreadProcessID(t *testing.T) {
	if err := EnumWindows(func(hwnd windows.HWND, args uintptr) bool {
		tid, pid, err := GetWindowThreadProcessID(hwnd)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(hwnd, tid, pid)

		return true
	}, 0); err != nil {
		t.Fatal(err)
	}
}

func TestGetWindowHwndByPID(t *testing.T) {
	pid, err := GetProcessIDByName("explorer.exe")
	if err != nil {
		t.Fatal(err)
	}

	for _, pid := range pid {
		hwnd := GetWindowHwndByPID(pid)
		t.Log(hwnd)
	}
}
