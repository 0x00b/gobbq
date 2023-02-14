package entity

// LegalHanlder 判断一个需要注册的业务接口是否符合要求(UnaryHandler)
// func LegalHanlder(fun any) error {
// 	handler := reflect.ValueOf(fun)
// 	if handler.Kind() != reflect.Func {
// 		return errors.New("handler wrong service")
// 	}
// 	if handler.Type().NumIn() != 2 {
// 		/* || handler.Type().NumOut() != 2 || !handler.Type().Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) */
// 		return errors.New("handler wrong parameter")
// 	}
// 	// p1 := handler.Type().In(0)
// 	// xlog.Println(handler.Type().In(0).Kind(), reflect.TypeOf((* entity.Context)(nil)).Elem().Kind())
// 	if !handler.Type().In(0).Implements(reflect.TypeOf((* entity.Context)(nil)).Elem()) {
// 		return errors.New("handler wrong parameter, first param must be entity.Context")
// 	}

// 	return nil
// }
