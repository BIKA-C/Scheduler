package router

import "scheduler/util"

func concat(s ...string) string {
	size := 0
	for i := 0; i < len(s); i++ {
		size += len(s[i])
	}

	buf := make([]byte, 0, size)

	for i := 0; i < len(s); i++ {
		buf = append(buf, util.S2B(s[i])...)
	}

	return util.B2S(buf)
}
