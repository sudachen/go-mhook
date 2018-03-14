//+build windows,amd64

package mhook

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type module struct {
	name string
	dll  *syscall.DLL
}

func (m *module) load() error {
	DLL, err := syscall.LoadDLL(m.name)
	if err != nil {
		return err
	}
	m.dll = DLL
	return nil
}

func (m *module) lookup(fn string) uintptr {
	proc, err := m.dll.FindProc(fn)
	if err != nil {
		return 0
	}
	return proc.Addr()
}

var modules = make(map[string]*module)
var functions = make(map[uintptr]*function)
var mu = sync.Mutex{}

func splitHookName(name string) (fnName, modName string, aLen int) {
	s := strings.Split(name, "!")
	modName = s[0]
	s = strings.Split(s[1], "#")
	fnName = s[0]
	if len(s) > 1 {
		aLen, _ = strconv.Atoi(s[1])
	}
	return
}

func dynLookup(name string) (*function, error) {
	mu.Lock()
	defer mu.Unlock()

	fnName, modName, aLen := splitHookName(name)
	var m *module
	m, ok := modules[modName]
	if !ok {
		m = &module{name: modName}
		if err := m.load(); err != nil {
			return nil, err
		}
		modules[modName] = m
	}

	p := m.lookup(fnName)
	if p == 0 {
		p = m.lookup(fmt.Sprintf("%s@%d", fnName, aLen*8))
	}
	if p == 0 {
		p = m.lookup(fmt.Sprintf("_%s", fnName))
	}
	if p == 0 {
		return nil, fmt.Errorf("module %s does not contain function %s", modName, fnName)
	}

	f, ok := functions[p]
	if !ok {
		f = &function{name: name, mod: m, fn: p, aLen: aLen}
	} else {
		if f.aLen != aLen {
			return nil, errors.New("arguments size conflict occured")
		}
	}

	return f, nil
}

type function struct {
	name string
	mod  *module
	fn   uintptr
	aLen int

	cb   uintptr
	orig uintptr
}

var freeGates []uintptr

const (
	allocBanch int = 256
	gateSize   int = 24
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	virtualAlloc          = kernel32.NewProc("VirtualAlloc")
	virtualFree           = kernel32.NewProc("VirtualFree")
	virtualProtect        = kernel32.NewProc("VirtualProtect")
	flushInstructionCache = kernel32.NewProc("FlushInstructionCache")
)

const (
	memCommit  = 0x00001000
	memRelease = 0x8000
)

//go:noescape
func setupWin64Gate(mem uintptr, cb uintptr)

func allocateGates(count int) error {

	addr, _, _ := virtualAlloc.Call(0, uintptr(count*gateSize), memCommit, syscall.PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		return errors.New("failed to allocate gates block")
	}
	for i := 0; i < count; i++ {
		freeGates = append(freeGates, addr+uintptr(i*gateSize))
	}

	return nil
}

func acquireGate(cb uintptr) (uintptr, error) {
	mu.Lock()
	defer mu.Unlock()

	if len(freeGates) < 1 {
		if err := allocateGates(allocBanch); err != nil {
			return 0, err
		}
	}

	L := len(freeGates) - 1
	g := freeGates[L]
	freeGates = freeGates[:L]

	setupWin64Gate(g, cb)
	return g, nil
}

func releaseGate(g uintptr) {
	mu.Lock()
	defer mu.Unlock()

	freeGates = append(freeGates, g)
}

func (f *function) Name() string {
	return f.name
}

func (f *function) readstack(stk uintptr) []uintptr {
	a := make([]uintptr, f.aLen)
	for i := 0; i < f.aLen; i++ {
		a[i] = *(*uintptr)(unsafe.Pointer(stk + uintptr(i*8)))
	}
	return a
}

func (f *function) Set(hf HookFunc) error {
	if f.cb == 0 {
		var err error
		f.cb, err = acquireGate(syscall.NewCallback(func(orig, stk uintptr) uintptr {
			a := f.readstack(stk)
			r := hf(f, a)
			if r.Continue {
				*(*uintptr)(unsafe.Pointer(orig)) = f.orig
			}
			return r.Value
		}))
		if err != nil {
			return err
		}
		//
	} else {
		return fmt.Errorf("%s is already binded to go calback", f.name)
	}
	return nil
}

func (f *function) Restore() error {
	//
	releaseGate(f.cb)
	return nil
}

func (f *function) Call(a []uintptr) (r1, r2 uintptr, err error) {
	return call(f.orig, a...)
}
