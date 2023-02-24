// Package util 安全的协程和函数调用,其中协程数量不支持配置热更新.
package utils

// 注意：recover不能在defer中的其他函数内调用，不能跳出当前执行栈，否则recover不住.

// IsRelease 是否为release模式运行.
var IsRelease string

// GOErrHandler panic error handler.
func GOErrHandler(err any) {
	// var buf [1024 * 4]byte
	// n := runtime.Stack(buf[:], false)
	// log.Fatal("<panic> msg=%v type=%T traceback:\n%s.", err, err, string(buf[:n])) // nolint:logcheck

	// // 刷新日志
	// if e := log.Sync(); e != nil {
	// 	log.Error("log Sync err:%s", e.Error())
	// }
	// // release模式才恢复
	// if !isReleaseMode() {
	// 	panic(err)
	// }
}

// func isReleaseMode() bool {
// 	if len(IsRelease) == 0 {
// 		return false
// 	}
// 	x, e := strconv.Atoi(IsRelease)
// 	if e != nil {
// 		log.Error("IsRelease=%s conv err:%s", IsRelease, e.Error())
// 		return false
// 	}
// 	return x != 0
// }

// SecureGO 带recover的go.
func SecureGO(f func()) {
	// 函数封装
	// nolint:securego
	go func() {
		defer func() {
			if err := recover(); err != nil {
				GOErrHandler(err)
			}
		}()
		// 真正的调用
		f()
	}()
}
