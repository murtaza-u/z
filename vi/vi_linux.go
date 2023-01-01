package vi

import (
	_ "embed"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

const MEMFILE = "/zvi.bin"

//go:embed bin/vi
var payload []byte

func vi(args ...string) error {
	if runtime.GOARCH != "amd64" {
		return fmt.Errorf("architecture not supported")
	}

	fd, err := memfdCreate(MEMFILE)
	if err != nil {
		return err
	}

	err = copyToMem(fd, payload)
	if err != nil {
		return err
	}

	if len(args) != 0 {
		args = append([]string{MEMFILE}, args...)
	}

	err = execveAt(fd, args...)
	if err != nil {
		return err
	}

	return nil
}

func memfdCreate(path string) (r1 uintptr, err error) {
	s, err := syscall.BytePtrFromString(path)
	if err != nil {
		return 0, err
	}

	r1, _, errno := syscall.Syscall(319, uintptr(unsafe.Pointer(s)), 0, 0)

	if int(r1) == -1 {
		return r1, errno
	}

	return r1, nil
}

func copyToMem(fd uintptr, buf []byte) (err error) {
	_, err = syscall.Write(int(fd), buf)
	if err != nil {
		return err
	}

	return nil
}

func execveAt(fd uintptr, argv ...string) (err error) {
	s, err := syscall.BytePtrFromString("")
	if err != nil {
		return err
	}

	argvp, err := syscall.SlicePtrFromStrings(argv)
	if err != nil {
		return err
	}

	ret, _, errno := syscall.Syscall6(
		322, fd, uintptr(unsafe.Pointer(s)),
		uintptr(unsafe.Pointer(&argvp[0])), 0, 0x1000, 0,
	)
	if int(ret) == -1 {
		return errno
	}

	return err
}
