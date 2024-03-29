//go:build go1.19

package nopreempt

import "unsafe"

func getg() *g

// Stack describes a Go execution stack.
// The bounds of the stack are exactly [lo, hi),
// with no implicit data structures on either side.
type stack struct {
	lo uintptr
	hi uintptr
}

type guintptr uintptr
type puintptr uintptr

type gobuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	//
	// ctxt is unusual with respect to GC: it may be a
	// heap-allocated funcval, so GC needs to track it, but it
	// needs to be set and cleared from assembly, where it's
	// difficult to have write barriers. However, ctxt is really a
	// saved, live register, and we only ever exchange it between
	// the real register and the gobuf. Hence, we treat it as a
	// root during stack scanning, which means assembly that saves
	// and restores it doesn't need write barriers. It's still
	// typed as a pointer so that any other writes from Go get
	// write barriers.
	sp   uintptr
	pc   uintptr
	g    guintptr
	ctxt uintptr // unsafe.Pointer
	ret  uintptr
	lr   uintptr
	bp   uintptr // for framepointer-enabled architectures
}

type g struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
	// stackguard1 is the stack pointer compared in the C stack growth prologue.
	// It is stack.lo+StackGuard on g0 and gsignal stacks.
	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	_panic    uintptr //*_panic // innermost panic - offset known to liblink
	_defer    uintptr // *_defer // innermost defer
	m         *m      // current m; offset known to arm liblink
	sched     gobuf
	syscallsp uintptr // if status==Gsyscall, syscallsp = sched.sp to use during gc
	syscallpc uintptr // if status==Gsyscall, syscallpc = sched.pc to use during gc
	stktopsp  uintptr // expected sp at top of stack, to check in traceback
	// param is a generic pointer parameter field used to pass
	// values in particular contexts where other storage for the
	// parameter would be difficult to find. It is currently used
	// in three ways:
	// 1. When a channel operation wakes up a blocked goroutine, it sets param to
	//    point to the sudog of the completed blocking operation.
	// 2. By gcAssistAlloc1 to signal back to its caller that the goroutine completed
	//    the GC cycle. It is unsafe to do so in any other way, because the goroutine's
	//    stack may have moved in the meantime.
	// 3. By debugCallWrap to pass parameters to a new goroutine because allocating a
	//    closure in the runtime is forbidden.
	param        unsafe.Pointer
	atomicstatus uint32
	stackLock    uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	goid         int64
}

const (
	// tlsSlots is the number of pointer-sized slots reserved for TLS on some platforms,
	// like Windows.
	tlsSlots = 6
)

type sigset struct{}
type throwType uint32

type m struct {
	g0      uintptr //*g     // goroutine with scheduling stack
	morebuf gobuf   // gobuf arg to morestack
	divmod  uint32  // div/mod denominator for arm - known to liblink
	_       uint32  // align next field to 8 bytes

	// Fields not known to debuggers.
	procid     uint64            // for debuggers, but offset not hard-coded
	gsignal    uintptr           // *g                // signal-handling g
	goSigStack gsignalStack      // Go-allocated signal handling stack
	sigmask    sigset            // storage for saved signal mask
	tls        [tlsSlots]uintptr // thread-local storage (for x86 extern register)
	mstartfn   func()
	curg       uintptr  //*g       // current running goroutine
	caughtsig  guintptr // goroutine running during fatal signal
	p          puintptr // attached p for executing go code (nil if not executing go code)
	nextp      puintptr
	oldp       puintptr // the p that was attached before executing a syscall
	id         int64
	mallocing  int32
	throwing   throwType
	preemptoff string // if != "", keep curg running on this m
	locks      int32
	dying      int32
	profilehz  int32
}
