package nopreempt

func GetGID() int64 {
	return getg().goid
}

func GetMID() int64 {
	return getg().m.id
}

func DisablePreempt() {
	getg().m.locks++
}

func EnablePreempt() {
	getg().m.locks--
}
