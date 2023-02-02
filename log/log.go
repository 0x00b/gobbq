package log

import "github.com/0x00b/gobbq/engine/entity"

func Tracef(ctx *entity.Context, format string, args ...any) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Tracef(format, args...)
}

func Debugf(ctx *entity.Context, format string, args ...any) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Debugf(format, args...)
}

func Infof(ctx *entity.Context, format string, args ...any) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Infof(format, args...)
}

func Printf(ctx *entity.Context, format string, args ...any) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Printf(format, args...)
}

func Warnf(ctx *entity.Context, format string, args ...any) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Warnf(format, args...)
}

func Warningf(ctx *entity.Context, format string, args ...any) {
	Warnf(ctx, format, args...)
}

func Errorf(ctx *entity.Context, format string, args ...any) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Errorf(format, args...)
}

func Fatalf(ctx *entity.Context, format string, args ...any) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Fatalf(format, args...)
}

func Panicf(ctx *entity.Context, format string, args ...any) {
	//ctx, format, args := handleLogfParameter(args)
	log(ctx).Panicf(format, args...)
}

func Trace(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Trace(args...)
}

func Debug(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Debug(args...)
}

func Info(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Info(args...)
}

func Print(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Print(args...)
}

func Warn(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Warn(args...)
}

func Warning(ctx *entity.Context, args ...any) {
	Warn(ctx, args...)
}

func Error(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Error(args...)
}

func Fatal(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Fatal(args...)
}

func Panic(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Panic(args...)
}

func Traceln(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Traceln(args...)
}

func Debugln(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Debugln(args...)
}

func Infoln(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Infoln(args...)
}

func Println(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Println(args...)
}

func Warnln(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Warnln(args...)
}

func Warningln(ctx *entity.Context, args ...any) {
	Warnln(ctx, args...)
}

func Errorln(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Errorln(args...)
}

func Fatalln(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Fatalln(args...)
}

func Panicln(ctx *entity.Context, args ...any) {
	//ctx, args := handleLogParameter(args)
	log(ctx).Panicln(args...)
}
