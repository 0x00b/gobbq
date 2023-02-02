package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/0x00b/gobbq/engine/entity"
	"github.com/sirupsen/logrus"
)

var (
	// defaultHook for getCaller
	defaultHook = &callerHook{}

	AddedHooks []logrus.Hook
)

type Formatter interface {
	FormatHook(*Hook) ([]byte, error)
}

// TagGetter
type TagGetter interface {
	GetTags(c *entity.Context) map[string]string
}

// goroutine unsafe, 里面的变量只能在初始化时赋值
// Hook hook for log
type Hook struct {
	ReportCaller bool

	Formatter Formatter

	Out io.Writer

	// Contains all the fields set by the user.
	Data logrus.Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level logrus.Level

	// Calling method, with package name
	Caller *runtime.Frame

	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// Contains the context set by the user. Useful for hook processing etc.
	Context *entity.Context

	// When formatter is called in entry.log(), a Buffer may be set to entry
	Buffer *bytes.Buffer

	mu sync.Mutex
}

func (h *Hook) HasCaller() bool {
	return h.Caller != nil && h.ReportCaller
}

func (h *Hook) Dup() *Hook {
	return &Hook{
		ReportCaller: h.ReportCaller,
		Formatter:    h.Formatter,
	}
}

func copyEntry(entry *logrus.Entry, hook *Hook) {
	hook.Data = entry.Data
	hook.Time = entry.Time
	hook.Level = entry.Level
	hook.Message = entry.Message
	// hook.Context = entry.Context
}

func (h *Hook) Fire(entry *logrus.Entry) error {

	newHook := h.Dup()

	if newHook.ReportCaller {
		newHook.Caller = getCaller()
	}

	copyEntry(entry, newHook)

	buffer, _ := bufPool.Get().(*bytes.Buffer)
	defer func() {
		newHook.Buffer = nil
		bufPool.Put(buffer)
	}()
	buffer.Reset()
	newHook.Buffer = buffer

	serialized, err := newHook.Formatter.FormatHook(newHook)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return nil
	}
	str := string(serialized)
	_ = str
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, err := h.Out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}

	return nil
}

// Levels giving the level you care, you can rewrite it
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type callerHook struct {
}

func (h *callerHook) Fire(entry *logrus.Entry) error {
	if entry.Logger.ReportCaller {
		entry.Caller = getCaller()
	}
	return nil
}

// Levels giving the level you care, you can rewrite it
func (h *callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var (
	// qualified package name, cached at first use
	logPackage []string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once

	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(4, pcs)

		packageNum := 0

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logPackage = append(logPackage, getPackageName(funcName))
				packageNum++
			} else if strings.Contains(funcName, "logrus") {
				logPackage = append(logPackage, getPackageName(funcName))
				packageNum++
			}
			if packageNum == 2 {
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if !contains(logPackage, pkg) {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
