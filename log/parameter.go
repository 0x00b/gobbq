package log

// func handleLogParameter(args []interface{}) (context.Context, []interface{}) {
// 	if len(args) < 1 {
// 		return context.Background(), args
// 	}
// 	switch args[0].(type) {
// 	case context.Context:
// 		return args[0].(context.Context), args[1:]
// 	default:
// 		return context.Background(), args
// 	}
// }

// func handleLogfParameter(args []interface{}) (context.Context, string, []interface{}) {
// 	if len(args) < 1 {
// 		return context.Background(), "", args
// 	}
// 	switch args[0].(type) {
// 	case context.Context:
// 		if len(args) < 2 {
// 			return args[0].(context.Context), "", args[1:]
// 		}
// 		switch args[1].(type) {
// 		case string:
// 			return args[0].(context.Context), args[1].(string), args[2:]
// 		default:
// 			return args[0].(context.Context), "", args[1:]
// 		}
// 	case string:
// 		return context.Background(), args[0].(string), args[1:]
// 	default:
// 		return context.Background(), "", args
// 	}
// }
