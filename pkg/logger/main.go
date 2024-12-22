package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/vandi37/vanerrors/vanstack"
)

// The logger
type Logger struct {
	out        io.Writer
	datePrefix string
	prefixes   [4]string
}

var stdPrefixes = [4]string{
	"^^info>>",
	"._.!Warn>>",
	"@_@:Error>>",
	">_<Fatal>>",
}

var stdDate = "02.01.06 3:04:05"

// Creates a new logger
func New(out io.Writer) *Logger {
	return NewWithSettings(out, stdDate, stdPrefixes)
}

func NewWithSettings(out io.Writer, date string, prefixes [4]string) *Logger {
	return &Logger{
		out:        out,
		datePrefix: date,
		prefixes:   prefixes,
	}
}

func (l *Logger) writeln(lvl int, a []any) {
	fmt.Fprintln(l.out, append([]any{time.Now().Format(l.datePrefix), l.prefixes[lvl]}, a...)...)
}

func (l *Logger) writef(lvl int, format string, a []any) {
	format = "%s %s " + format + "\n"
	fmt.Fprintf(l.out, format, append([]any{time.Now().Format(l.datePrefix), l.prefixes[lvl]}, a...)...)
}

// Prints a line
func (l *Logger) Println(a ...any) {
	l.writeln(0, a)
}

// Prints a formatted line
func (l *Logger) Printf(format string, a ...any) {
	l.writef(0, format, a)
}

// Prints a warn line
func (l *Logger) Warnln(a ...any) {
	l.writeln(1, a)

}

// Prints a warn formatted line
func (l *Logger) Warnf(format string, a ...any) {
	l.writef(1, format, a)
}

// Prints a error line
func (l *Logger) Errorln(a ...any) {
	l.writeln(2, a)
}

// Prints a error formatted line
func (l *Logger) Errorf(format string, a ...any) {
	l.writef(2, format, a)
}

// Prints a fatal line and exit
func (l *Logger) Fatalln(a ...any) {
	l.writeln(3, a)
	stack := vanstack.NewStack()
	stack.Fill("", 20)
	fmt.Fprintln(os.Stderr, stack)
	os.Exit(http.StatusTeapot)
}

// Prints a fatal formatted line and exit
func (l *Logger) Fatalf(format string, a ...any) {
	l.writef(4, format, a)
	stack := vanstack.NewStack()
	stack.Fill("", 20)
	fmt.Fprintln(os.Stderr, stack)
	os.Exit(http.StatusTeapot)
}
