package mhook

import (
	"strconv"
	"syscall"
)

type Result struct {
	Value    uintptr
	Continue bool
}

type HookFunc func(Hook, []uintptr) Result

type Hook interface {
	Name() string
	Set(HookFunc) error
	Restore() error
	Call([]uintptr) (r1, r2 uintptr, err error)
}

func NewHook(name string) (Hook, error) {
	if fn, err := dynLookup(name); err != nil {
		return nil, err
	} else {
		return fn, nil
	}
}

func call(p uintptr, a ...uintptr) (r1, r2 uintptr, lastErr error) {
	switch len(a) {
	case 0:
		return syscall.Syscall(p, uintptr(len(a)), 0, 0, 0)
	case 1:
		return syscall.Syscall(p, uintptr(len(a)), a[0], 0, 0)
	case 2:
		return syscall.Syscall(p, uintptr(len(a)), a[0], a[1], 0)
	case 3:
		return syscall.Syscall(p, uintptr(len(a)), a[0], a[1], a[2])
	case 4:
		return syscall.Syscall6(p, uintptr(len(a)), a[0], a[1], a[2], a[3], 0, 0)
	case 5:
		return syscall.Syscall6(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], 0)
	case 6:
		return syscall.Syscall6(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5])
	case 7:
		return syscall.Syscall9(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], 0, 0)
	case 8:
		return syscall.Syscall9(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], 0)
	case 9:
		return syscall.Syscall9(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8])
	case 10:
		return syscall.Syscall12(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], 0, 0)
	case 11:
		return syscall.Syscall12(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], 0)
	case 12:
		return syscall.Syscall12(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11])
	case 13:
		return syscall.Syscall15(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], 0, 0)
	case 14:
		return syscall.Syscall15(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], 0)
	case 15:
		return syscall.Syscall15(p, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14])
	default:
		panic("Call is impossible, too many arguments: " + strconv.Itoa(len(a)) + ".")
	}
}
