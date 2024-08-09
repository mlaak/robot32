//go:build !ttd

package ttd

func TTD(rq int64, message string, pars ...interface{}) {

}

func __TTD(rq int64, message string, pars ...interface{}) {

}

func TTDLEV(c int64, lev int64) int64 {
	return c
}

func TTX[A any](c int64, v A) A {
	return v
}

func TTX2[A any, B any](c int64, v1 A, v2 B) (A, B) {
	return v1, v2
}

func TTX3[A any, B any, C any](c int64, v1 A, v2 B, v3 C) (A, B, C) {
	return v1, v2, v3
}
