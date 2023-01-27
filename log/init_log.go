package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// 对日志初始化：
// 1、请求中使用统一的日志格式

var (
	_organizeLog bool

	LogFileNameLen = 40

	JsonFormatter *JSONFormatter
	//	SetCtxRequestID SetCtxRequestIDFunc
)

var (
	StdFormatter         *TextFormatter
	SimpleFormatter      *TextFormatter
	SingleLineFormatter  *TextFormatter
	organizeLogFormatter *TextFormatter
)

// 对日志初始化,参数如下：
// level		:设置日志级别, 方便从配置文件中直接获取，所以使用string
// organizeLog	:是否整理日志, 本地调试或者不需要收集时建议true, 配合OrganizeLogMiddleware 使用
// reportCaller	:是否打印函数调用栈
// out			:日志打印到何处
// tags			:设置日志中获取上下文的接口
func Init(level string, organizeLog, reportCaller bool, out io.Writer, conf string, tags ...TagGetter) {
	_organizeLog = organizeLog

	llevel, e := logrus.ParseLevel(level)
	if e != nil {
		panic(e)
	}

	defaultFmt(organizeLog, tags...)

	logrus.SetLevel(llevel)
	logrus.SetFormatter(SingleLineFormatter)
	// logrus.SetFormatter(JsonFormatter)
	// logrus.SetFormatter(SimpleFormatter)

	// default hook equal default log, use defaultHook replace logrus.StandardLogger
	// defaultHook.Formatter = JsonFormatter
	// defaultHook.Formatter = SingleLineFormatter
	// use hook
	AddHook(defaultHook)

	// TODO:
	// 读取配置，注册各种各样的hook。
	// read(conf)
	// conf to Hook
	// AddHook(*Hook)

	// dumb standard log
	logrus.SetOutput(out)
	logrus.SetReportCaller(reportCaller)
}

// goroutine unsafe
// AddHook adds a hook to the standard logger hooks.
func AddHook(hooks ...logrus.Hook) {
	for _, hook := range hooks {
		AddedHooks = append(AddedHooks, hook)
		logrus.AddHook(hook)
	}
}

func defaultFmt(organizeLog bool, tags ...TagGetter) {
	// f := &TextFormatter{TimestampFormat: "01-02 15:04:05.000000"}
	// //formatter.SetFormat(TagBL, FieldKeyTime, TagBR, FieldKeyLevel, FieldKeyFile,
	//	 TaGColon, FieldKeyFunc, TaGColon, FieldKeyLine, FieldKeyMsg)
	// f.SetFormatAndTagSource(TagBL, FieldKeyTime, TagBR, FieldKeyLevel, FieldKeyMsg)
	// //文件名取后两级"/"之后，如果长度超过40字节就取40字节，不足不用管
	// f.FormatFileName = formatFileName
	// f.NoQuoteFields = true
	// StdFormatter = f

	f := &TextFormatter{TimestampFormat: "01-02 15:04:05.000000"}
	f.SetFormatAndTagSource(nil, TagBL, FieldKeyTime, TagBR, FieldKeyLevel, FieldKeyMsg)
	f.NoQuoteFields = true
	f.FormatFileName = formatFileName
	if !organizeLog {
		f.TagGetters = tags
	}
	SingleLineFormatter = f

	JsonFormatter = &JSONFormatter{
		TagGetters:       tags,
		TimestampFormat:  "01-02 15:04:05.000000",
		DisableTimestamp: false,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return f.Function, formatFileName(f.File) + fmt.Sprintf(":%d", f.Line)
		},
	}

	f = &TextFormatter{TimestampFormat: "01-02 15:04:05.000000", TagGetters: tags}
	f.SetFormat(TagBL, FieldKeyTime, TagBR, FieldKeyLevel, FieldKeyMsg)
	f.NoQuoteFields = true
	SimpleFormatter = f

	f = &TextFormatter{TimestampFormat: "01-02 15:04:05.000000", TagGetters: tags}
	f.SetFormat(FieldKeyMsg, "\n")
	f.NoQuoteFields = true
	organizeLogFormatter = f
}

func formatFileName(name string) string {
	idx := strings.LastIndex(name, "/")
	if -1 != idx {
		idx = strings.LastIndex(name[:idx], "/")
		if -1 != idx {
			name = name[idx:]
		}
	}
	if len(name) > LogFileNameLen {
		return name[len(name)-LogFileNameLen:]
	}
	return name
}
