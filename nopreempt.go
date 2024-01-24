package nopreempt

func GetGID() int64 {
	return getg().goid
}

func GetMID() int64 {
	return getg().m.id
}

func DisablePreempt() {
	// getg().m.locks++

	var m, m2 *m
	g := getg()
	m = g.m
retry:
	m.locks++
	if m2 = g.m; m != m2 {
		// stolen. recover and retry
		m.locks--
		m = m2
		goto retry
	}
}

func EnablePreempt() {
	getg().m.locks--
}
