package log

import "context"

func Tracef(ctx context.Context, format string, args ...interface{}) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Tracef(format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Debugf(format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Infof(format, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Printf(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Warnf(format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	Warnf(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Errorf(format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Fatalf(format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Panicf(format, args...)
}

func Trace(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Trace(args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Debug(args...)
}

func Info(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Info(args...)
}

func Print(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Print(args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Warn(args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	Warn(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Error(args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Fatal(args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Panic(args...)
}

func Traceln(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Traceln(args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Debugln(args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Infoln(args...)
}

func Println(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Println(args...)
}

func Warnln(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Warnln(args...)
}

func Warningln(ctx context.Context, args ...interface{}) {
	Warnln(ctx, args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Errorln(args...)
}

func Fatalln(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Fatalln(args...)
}

func Panicln(ctx context.Context, args ...interface{}) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Panicln(args...)
}
