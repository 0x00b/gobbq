package log

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"

	//"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	FieldKeyMsg            = logrus.FieldKeyMsg
	FieldKeyLevel          = logrus.FieldKeyLevel
	FieldKeyTime           = logrus.FieldKeyTime
	FieldKeyLogrusError    = logrus.FieldKeyLogrusError
	FieldKeyFunc           = logrus.FieldKeyFunc
	FieldKeyFile           = logrus.FieldKeyFile
	FieldKeyLine           = "line"
	TagBL                  = "["
	TagBR                  = "]"
	TagColon               = ":"
	TagVBar                = "|"
	defaultTimestampFormat = time.RFC3339
)

var (
	// defaultFormat = fmt.Sprintf("[%%%v%%] %%%v%% %%%v%%:%%line%% - %%%v%%",
	// FieldKeyTime, FieldKeyLevel, FieldKeyFunc, FieldKeyMsg)
	defaultFormatArray = []string{TagBL, FieldKeyTime, TagBR, FieldKeyLevel, FieldKeyMsg}

	FunctionNameLength = 25
	FileNameLength     = 20
)

// TextFormatter formats logs into text
type TextFormatter struct {
	// TagGetter implement this to get all tag you want, like traceid , requestid etc...
	TagGetters []TagGetter

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool
	NoQuoteFields    bool

	//LogFormat
	//LogFormat string

	FormatFuncName HandlerFormatFunc
	FormatFileName HandlerFormatFile

	TagSource bool

	//hasTime  bool
	//hasLevel bool
	//hasMsg   bool
	//hasFunc  bool
	//hasFile  bool
	//hasLine  bool

	keyArray  []string
	FieldKeys []string

	TagSourceFormatter
}

type TagSourceFormatter interface {
	Format(file, fun string, line int) (tag, value string)
}

// HandlerFormatFunc format function name
type HandlerFormatFunc func(funcName string) string

// HandlerFormatFile format file name
type HandlerFormatFile func(fileName string) string

func defaultFormatFunc(funcName string) string {
	length := len(funcName)
	if length > FunctionNameLength {
		return "." + funcName[length-FunctionNameLength+1:length]
	}
	return strings.Repeat(" ", FunctionNameLength-length) + funcName
}
func defaultFormatFuncTag(funcName string) string {
	idx := strings.LastIndex(funcName, "/")
	if -1 == idx {
		return funcName
	}
	return funcName[idx:]
}
func defaultFormatFile(fileName string) string {
	r := []rune(fileName)
	length := len(r)
	if length > FileNameLength {
		return "." + fileName[length-FileNameLength+1:length]
	}
	return strings.Repeat(" ", FileNameLength-length) + fileName
}

func isTag(s string) bool {
	switch s {
	case TagBR:
		return true
	case TagBL:
		return true
	case TagColon:
		return true
	case TagVBar:
		return true
	}
	return false
}

func isBR(s string) bool {
	return s == TagBR
}

func needBlank(s string) bool {
	switch s {
	case TagBR:
		return false
	case TagColon:
		return false
	case TagVBar:
		return false
	}
	return true
}

//func (f *TextFormatter) setHasKey(k string) {
//	switch true {
//	case k == FieldKeyTime:
//		f.hasTime = true
//	case k == FieldKeyLevel:
//		f.hasLevel = true
//	case k == FieldKeyMsg:
//		f.hasMsg = true
//	case k == FieldKeyFunc:
//		f.hasFunc = true
//	case k == FieldKeyFile:
//		f.hasFile = true
//	case k == FieldKeyLine:
//		f.hasLine = true
//	}
//}

func (f *TextFormatter) SetFormat(msgFmts ...string) {
	f.keyArray = msgFmts
}
func (f *TextFormatter) RegisterFields(args ...string) {
	f.FieldKeys = args
}
func (f *TextFormatter) SetFormatAndTagSource(sourceFmt TagSourceFormatter, msgFmts ...string) {
	if sourceFmt == nil {
		f.TagSourceFormatter = defualtTagSourceFormatter{}
	}
	f.keyArray = msgFmts
	f.TagSource = true
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	hook := &Hook{}

	copyEntry(entry, hook)
	hook.Caller = entry.Caller
	hook.ReportCaller = entry.HasCaller()
	hook.Buffer = entry.Buffer

	return f.FormatHook(hook)
}

// nolint
func (f *TextFormatter) FormatHook(hook *Hook) ([]byte, error) {
	var buf *bytes.Buffer
	if hook.Buffer != nil {
		buf = hook.Buffer
	} else {
		buf = &bytes.Buffer{}
	}

	if len(f.keyArray) == 0 {
		f.keyArray = defaultFormatArray
	}
	endWithLn := false
	for idx, k := range f.keyArray {
		if isTag(k) {
			buf.WriteString(k)
			if isBR(k) {
				buf.WriteByte(' ')
			}
		} else {
			switch k {
			case FieldKeyMsg:
				buf.WriteString(f.quoteValue(hook.Message))
			case FieldKeyLevel:
				level := Level(hook.Level).String()
				buf.WriteString(level)
			case FieldKeyTime:
				timestampFormat := f.TimestampFormat
				if timestampFormat == "" {
					timestampFormat = defaultTimestampFormat
				}
				buf.WriteString(hook.Time.Format(timestampFormat))
			case FieldKeyFunc:
				if hook.HasCaller() && hook.Caller != nil {
					function := hook.Caller.File
					if f.FormatFuncName != nil {
						function = f.FormatFuncName(function)
					} else {
						function = defaultFormatFunc(function)
					}
					buf.WriteString(function)
				}
			case FieldKeyFile:
				if hook.HasCaller() && hook.Caller != nil {
					fileName := hook.Caller.File
					if f.FormatFileName != nil {
						fileName = f.FormatFileName(fileName)
					} else {
						fileName = defaultFormatFile(fileName)
					}
					buf.WriteString(fileName)
				}
			case FieldKeyLine:
				if hook.HasCaller() && hook.Caller != nil {
					line := fmt.Sprintf("%-4v", strconv.FormatInt(int64(hook.Caller.Line), 10))
					buf.WriteString(line)
				}
			default:
				if contains(f.FieldKeys, k) {
					if hook.Data[k] != nil {
						buf.WriteString(f.quoteValue(fmt.Sprintf("%v", hook.Data[k])))
					}
				} else {
					if idx < len(f.keyArray)-1 || k != "\n" {
						buf.WriteString(k)
						continue
					} else {
						endWithLn = true
					}
				}
			}
			if idx < len(f.keyArray)-1 && needBlank(f.keyArray[idx+1]) {
				buf.WriteByte(' ')
			}
		}
	}

	tagSource := f.TagSource && hook.HasCaller() && hook.Caller != nil
	hasTag := tagSource
	length := len(hook.Data)
	if tagSource {
		fileName := hook.Caller.File
		if f.FormatFileName != nil {
			fileName = f.FormatFileName(fileName)
		} else {
			fileName = defaultFormatFile(fileName)
		}
		function := hook.Caller.Function
		if f.FormatFuncName != nil {
			function = f.FormatFuncName(function)
		} else {
			function = defaultFormatFuncTag(function)
		}
		t, v := f.TagSourceFormatter.Format(fileName, function, hook.Caller.Line)
		buf.WriteString(fmt.Sprintf(" (%v=%v", t, v))
	}
	if length > 0 {
		bMoreField := false
		moreCnt := 0
		for k := range hook.Data {
			if !contains(f.FieldKeys, k) {
				bMoreField = true
				moreCnt++
			}
		}
		hasTag = hasTag || bMoreField

		if bMoreField {
			var idx int
			if !tagSource {
				buf.WriteString(" (")
			} else {
				buf.WriteByte(' ')
			}
			printCnt := 0
			for k, v := range hook.Data {
				if contains(f.FieldKeys, k) {
					continue
				}
				printCnt++
				if s, ok := v.(string); ok {
					buf.WriteString(fmt.Sprintf("%v=%q", k, s))
				} else {
					buf.WriteString(fmt.Sprintf("%v=%v", k, v))
				}
				//if idx < length-1 {
				if printCnt < moreCnt {
					buf.WriteByte(' ')
				}
				idx++
			}
		}
	}

	datas := make(map[string]string)

	for _, tags := range f.TagGetters {
		for k, v := range tags.GetTags(hook.Context) {
			datas[k] = v
		}
	}
	lenDatas := len(datas)
	if hasTag {
		if lenDatas > 0 {
			buf.WriteByte(' ')
		}
	} else if lenDatas > 0 {
		buf.WriteString(" (")
	}
	idx := 0
	for k, v := range datas {
		buf.WriteString(fmt.Sprintf("%v=%v", k, v))
		idx++
		if idx != lenDatas {
			buf.WriteByte(' ')
		}
	}
	hasTag = hasTag || (lenDatas > 0)

	if hasTag {
		buf.WriteByte(')')
	}

	buf.WriteByte('\n')
	if endWithLn {
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

func contains(a []string, s string) bool {
	for _, value := range a {
		if value == s {
			return true
		}
	}
	return false
}

type Level logrus.Level

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	}
	return "unknown"
}

func (level Level) MarshalText() ([]byte, error) {
	switch logrus.Level(level) {
	case logrus.TraceLevel:
		return []byte("TRAC"), nil
	case logrus.DebugLevel:
		return []byte("DEBG"), nil
	case logrus.InfoLevel:
		return []byte("INFO"), nil
	case logrus.WarnLevel:
		return []byte("WARN"), nil
	case logrus.ErrorLevel:
		return []byte("ERRO"), nil
	case logrus.FatalLevel:
		return []byte("FATA"), nil
	case logrus.PanicLevel:
		return []byte("PANC"), nil
	}

	return nil, fmt.Errorf("not a valid lorus level %q", level)
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if !f.NoQuoteFields {
		for _, ch := range text {
			if !((ch >= 'a' && ch <= 'z') ||
				(ch >= 'A' && ch <= 'Z') ||
				(ch >= '0' && ch <= '9') ||
				ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
				return true
			}
		}
	}
	return false
}

func (f *TextFormatter) quoteValue(value any) string {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		return stringVal
	}
	return fmt.Sprintf("%q", stringVal)
}

type defualtTagSourceFormatter struct{}

func (defualtTagSourceFormatter) Format(file, fun string, line int) (tag, value string) {
	return "F", fmt.Sprintf("%v:%v", file, line)
}
