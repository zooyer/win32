package win32

import (
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestReadProcessMemory(t *testing.T) {
	var (
		data = []byte("hello,world")
		buf  = make([]byte, len(data))
	)

	n, err := ReadProcessMemory(windows.CurrentProcess(), uintptr(unsafe.Pointer(&data[0])), buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(n)
	t.Log(int(n) == len(data))
	t.Log(string(buf))
}

func TestReadProcessMemoryData(t *testing.T) {
	var values = []any{
		true,
		123456,
		"hello,world",
	}

	for _, value := range values {
		switch value := value.(type) {
		case int:
			t.Log(ReadProcessMemoryData[int](windows.CurrentProcess(), uintptr(unsafe.Pointer(&value))))
		case bool:
			t.Log(ReadProcessMemoryData[bool](windows.CurrentProcess(), uintptr(unsafe.Pointer(&value))))
		case string:
			t.Log(ReadProcessMemoryData[string](windows.CurrentProcess(), uintptr(unsafe.Pointer(&value))))
		}
	}
}

func TestReadProcessMemoryCStringA(t *testing.T) {
	var (
		str = append([]byte("hello,world!"), 0x00)
		ptr = &str[0]
	)

	data, err := ReadProcessMemoryCStringA(windows.CurrentProcess(), uintptr(unsafe.Pointer(&ptr)))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data)
}

func TestReadProcessMemoryCStringW(t *testing.T) {
	var str = windows.StringToUTF16Ptr("hello,world!!!")

	data, err := ReadProcessMemoryCStringW(windows.CurrentProcess(), uintptr(unsafe.Pointer(&str)))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data)
}

func TestEnablePrivilege(t *testing.T) {
	if err := EnablePrivilege(); err != nil {
		t.Fatal(err)
	}
}
