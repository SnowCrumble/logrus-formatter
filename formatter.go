package formatter

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultTimestampFormat = time.RFC3339

//Formatter A custom formatter
type Formatter struct {
	buffer *bytes.Buffer

	// Output line number and function name and other else debug info
	Debug bool

	sync.Once
}

func (f *Formatter) init() {
	f.buffer = bytes.NewBuffer(make([]byte, 0, 4096))
}

func (f *Formatter) appendKeyValue(key string, value interface{}) {
	if f.buffer.Len() > 0 {
		f.buffer.WriteByte(' ')
	}
	f.buffer.WriteString(key)
	f.buffer.WriteByte('=')
	f.appendValue(value)
	return
}

func (f *Formatter) appendValue(value interface{}) (n int) {
	stringVal := fmt.Sprint(value)

	if !f.needsQuoting(stringVal) {
		n, _ = f.buffer.WriteString(stringVal)
	} else {
		n, _ = f.buffer.WriteString(fmt.Sprintf("%q", stringVal))
	}
	return
}

func (f *Formatter) needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

//Format Format entry to []byte
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.Do(func() {
		f.init()
	})
	f.buffer.Reset()
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if f.Debug {
		// Fetch debug info from call stack
		var callers = make([]uintptr, 3)
		runtime.Callers(6, callers)
		frames := runtime.CallersFrames(callers)
		for {
			frame, more := frames.Next()
			if !more {
				break
			}
			if !strings.Contains(frame.Function, "github.com/sirupsen/logrus") {
				debugInfo := fmt.Sprintf("[%s][%s][%d]",
					filepath.Base(frame.File),
					filepath.Base(frame.Func.Name()),
					frame.Line,
				)
				f.appendKeyValue("debug", debugInfo)
				break
			}
		}
	}
	f.appendKeyValue("time", entry.Time.Local().Format(defaultTimestampFormat))
	f.appendKeyValue("level", entry.Level.String())
	if entry.Message != "" {
		f.appendKeyValue("msg", entry.Message)
	}
	for _, key := range keys {
		f.appendKeyValue(key, entry.Data[key])
	}
	f.buffer.WriteByte('\n')
	return f.buffer.Bytes(), nil
}
