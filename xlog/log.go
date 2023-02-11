package xlog

import (
	"context"
)

func Tracef(args ...any) {
	ctx, format, args := handleLogfParameter(args)
	log(ctx).Tracef(format, args...)
}

func Debugf(args ...any) {
	ctx, format, args := handleLogfParameter(args)
	log(ctx).Debugf(format, args...)
}

func Infof(args ...any) {
	ctx, format, args := handleLogfParameter(args)
	log(ctx).Infof(format, args...)
}

func Printf(args ...any) {
	ctx, format, args := handleLogfParameter(args)
	log(ctx).Printf(format, args...)
}

func Warnf(args ...any) {
	ctx, format, args := handleLogfParameter(args)
	log(ctx).Warnf(format, args...)
}

func Warningf(args ...any) {
	Warnf(args...)
}

func Errorf(args ...any) {
	ctx, format, args := handleLogfParameter(args)
	log(ctx).Errorf(format, args...)
}

func Fatalf(args ...any) {
	ctx, format, args := handleLogfParameter(args)
	log(ctx).Fatalf(format, args...)
}

func Panicf(args ...any) {
	ctx, format, args := handleLogfParameter(args)
	log(ctx).Panicf(format, args...)
}

func Trace(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Trace(args...)
}

func Debug(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Debug(args...)
}

func Info(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Info(args...)
}

func Print(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Print(args...)
}

func Warn(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Warn(args...)
}

func Warning(args ...any) {
	Warn(args...)
}

func Error(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Error(args...)
}

func Fatal(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Fatal(args...)
}

func Panic(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Panic(args...)
}

func Traceln(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Traceln(args...)
}

func Debugln(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Debugln(args...)
}

func Infoln(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Infoln(args...)
}

func Println(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Println(args...)
}

func Warnln(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Warnln(args...)
}

func Warningln(args ...any) {
	Warnln(args...)
}

func Errorln(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Errorln(args...)
}

func Fatalln(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Fatalln(args...)
}

func Panicln(args ...any) {
	ctx, args := handleLogParameter(args)
	log(ctx).Panicln(args...)
}

func handleLogParameter(args []any) (ctx context.Context, logArgs []any) {
	if len(args) <= 0 {
		return nil, args
	}
	if ctx, ok := args[0].(context.Context); ok {
		return ctx, args[1:]
	}
	return nil, args
}

func handleLogfParameter(args []any) (ctx context.Context, fmtter string, logArgs []any) {
	if len(args) <= 0 {
		return nil, "", args
	}
	idx := 0
	var ok bool
	ctx, ok = args[idx].(context.Context)
	if ok {
		idx++
	}
	if idx < len(args) {
		fmtter, _ = args[idx].(string)
	}
	// fmt.Println("idx", idx)
	// fmt.Println(fmtter)
	return ctx, fmtter, args[idx+1:]
}
