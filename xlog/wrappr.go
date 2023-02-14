package xlog

// Wrapper like using _organizeLog for new go routine
// func Wrapper(fun func( entity.Context)) func( entity.Context) {
// 	return func(c entity.Context) {
// 		if _organizeLog {
// 			_, ok := (c).Value(_ctxLogKey).(*logrus.Entry)
// 			if !ok {
// 				clog, buf := newCtxLog(c)
// 				c = context.WithValue(c, _ctxLogKey, clog)
// 				defer func() {
// 					organizeLog(c, buf)
// 				}()
// 			}
// 		}

// 		defer func() {
// 			if err := recover(); err != nil {
// 				//log
// 				Errorln(c, err)
// 				Errorln(c, "panic:\n", string(utils.CallStack(3)))
// 				e, ok := err.(error)
// 				if !ok {
// 					e = errors.New("panic")
// 				}
// 				//roport
// 				metrics.Metrics.Counter(c, "panic", 1, e)
// 			}
// 		}()
// 		c, span := trace.StartChildFromContext(c, ago.Meta(c).Action)
// 		defer span.Finish()

// 		fun(c)
// 	}
// }
