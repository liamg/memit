package memit

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/sys/unix" 
)

func Command(reader io.Reader, args ...string) (*exec.Cmd, *os.File, error) {
	memoryDescriptor, err := unix.MemfdCreate("", unix.MFD_CLOEXEC)
	if err != nil {
		return nil, nil, err
	}
	file := os.NewFile(
		uintptr(memoryDescriptor),
		fmt.Sprintf("/proc/self/fd/%d", memoryDescriptor),
	)
	if _, err := io.Copy(file, reader); err != nil {
		_ = file.Close()
		return nil, nil, err
	}
	return exec.Command(file.Name(), args...), file, nil
}
