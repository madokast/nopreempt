package nopreempt

import (
	"unsafe"
)

type MP struct {
	mp uintptr
}

//go:linkname acquirem runtime.acquirem
func acquirem() unsafe.Pointer

//go:linkname releasem runtime.releasem
func releasem(unsafe.Pointer)

func GetGID() int64 {
	return getg().goid
}

func GetMID() int64 {
	return getg().m.id
}

func AcquireM() MP {
	return MP{
		mp: uintptr(acquirem()),
	}
}

func (mp MP) MID() int64 {
	return (*m)(unsafe.Pointer(mp.mp)).id
}

func (mp MP) Release() {
	releasem(unsafe.Pointer(mp.mp))
}
