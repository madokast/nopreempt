# nopreempt

Prevent Golang goroutine preemption. Ref. https://github.com/petermattis/goid 

阻止 Go 协程抢占。参考 https://github.com/petermattis/goid 


Usage:

使用方法：

```golang
func fun() {
	DisablePreempt()
	defer EnablePreempt()

	// code will not be preempted
}
```

Func GetGID() and GetMID() are also available.

另外还提供 GetGID() 和 GetMID() 函数。
