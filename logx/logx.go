package logx

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Handler interface {
	Dofile(fpath string) error
}

type defaultHandler struct{}

func (def *defaultHandler) Dofile(fpath string) error {
	lr, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer func() {
		if lr != nil {
			lr.Close()
		}
	}()

	fw, err := os.Create(fpath + ".zip")
	if err != nil {
		return err
	}
	defer fw.Close()

	w := zip.NewWriter(fw)
	zf, err := w.Create(filepath.Base(fpath))
	if err != nil {
		return err
	}

	_, err = io.Copy(zf, lr)
	if err != nil {
		w.Close()
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	lr.Close()
	lr = nil
	return os.Remove(fpath)
}

var defhandler *defaultHandler = &defaultHandler{}

const (
	// Bits or'ed together to control what's printed. There is no control over the
	// order they appear (the order listed here) or the format they present (as
	// described in the comments).  A colon appears after these items:
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota              // the date: 2009/01/23
	Ltime                                  // the time: 01:23:23
	Lmicroseconds                          // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                              // full file name and line number: /a/b/c/d.go:23
	Lshortfile                             // final file name element and line number: d.go:23. overrides Llongfile
	Lmodule                                // module name
	Llevel                                 // level: 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
	LstdFlags     = Ldate | Ltime | Llevel // initial values for the standard logger
	Ldefault      = Llevel | LstdFlags
)
const (
	Ldebug = iota
	Linfo
	Lwarn
	Lerror
	Lpanic
	Lfatal
	Lmax
)

const (
	defmaxfilesize  = 1024 * 1024 * 10
	defmaxfilecount = 10
)

var levels = []string{
	"[DEB]",
	"[INF]",
	"[WAR]",
	"[ERR]",
	"[PAN]",
	"[FAT]",
}

// A Logger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu       sync.Mutex // ensures atomic writes; protects the following fields
	prefix   string     // prefix to write at beginning of each line
	flag     int        // properties
	out      io.Writer  // destination for output
	buf      []byte     // for accumulating text to write
	level    int
	bmulti   bool     // 是否是环保模式
	fcount   int      // 文件个数
	fmaxsize int      // 最大文件大小
	file     []string // 文件列表
	folder   string
	name     string
	hander   Handler // 文件处理
}

func NewLog(folder string, name string, lvl int) *Logger {
	l := &Logger{
		out:      nil,
		flag:     Ldefault,
		level:    lvl,
		bmulti:   true,
		fcount:   defmaxfilecount,
		fmaxsize: defmaxfilesize,
		folder:   folder,
		name:     name,
		hander:   defhandler,
	}
	wrter := l.createIo()
	if wrter == nil {
		return nil
	}
	l.out = wrter
	return l
}

func (l *Logger) Log(lvl int, format string, v ...interface{}) {
	l.Output(lvl, fmt.Sprintf(format, v...))
}
func Log(lvl int, format string, v ...interface{}) {
	std.Output(lvl, fmt.Sprintf(format, v...))
}
func getProcAbsDir() (string, error) {
	abs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", nil
	}
	return filepath.Dir(abs), nil
}
func isPathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
func getTimeStr() string {
	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	str := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", year, month, day, hour, minute, second)
	return str
}
func (l *Logger) createIo() io.Writer {
	procfolder, err := getProcAbsDir()
	if err != nil {
		return nil
	}
	//	procfolder = "c:/"

	timestr := getTimeStr()
	logfolder := fmt.Sprintf("%s/%s/%s", procfolder, l.folder, timestr)
	if len(l.file) > 0 {
		logfolder = filepath.Dir(l.file[0])
	}

	if isPathExist(logfolder) == false {
		if os.MkdirAll(logfolder, os.ModeDir) != nil {
			return nil
		}
	}

	fileext := ""
	filename := l.name
	extdotindex := strings.LastIndex(l.name, ".")
	if extdotindex != -1 {
		fileext = l.name[extdotindex:]
		filename = l.name[:extdotindex]
	}

	logfilename := fmt.Sprintf("%s/%s-%s%s", logfolder, filename, timestr, fileext)

	file, err := os.Create(logfilename)
	if err != nil {
		return nil
	}

	l.file = append(l.file, logfilename)
	if len(l.file) > l.fcount {
		oldestfile := l.file[0]
		if l.hander != nil {
			go func(path string) {
				l.hander.Dofile(path)
			}(oldestfile)
		}

		l.file = l.file[1:]
	}

	return file
}
func (l *Logger) SetFile(count, size int) {
	l.fcount = count
	l.fmaxsize = size
}
func (l *Logger) SetHandler(hander Handler) {
	l.hander = hander
}

// New creates a new Logger.   The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{out: out, prefix: prefix + " ", flag: flag, level: Linfo}
}

var std = New(os.Stdout, "", LstdFlags)

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(buf *[]byte, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		*buf = append(*buf, '0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}
	*buf = append(*buf, b[bp:]...)
}

func moduleOf(file string) string {
	pos := strings.LastIndex(file, "/")
	if pos != -1 {
		pos1 := strings.LastIndex(file[:pos], "/src/")
		if pos1 != -1 {
			return file[pos1+5 : pos]
		}
	}
	return "UNKNOWN"
}

func (l *Logger) formatHeader(buf *[]byte, t time.Time, lvl int, file string, line int) {
	*buf = append(*buf, l.prefix...)

	if l.flag&Lmodule != 0 {
		*buf = append(*buf, moduleOf(file)...)
		*buf = append(*buf, ' ')
	}

	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ' ')
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '-')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '-')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&Llevel != 0 {
		*buf = append(*buf, levels[lvl%Lmax]...)
		*buf = append(*buf, ' ')
	}
}

// Output writes the output for a logging event.  The string s contains
// the text to print after the prefix specified by the flags of the
// Logger.  A newline is appended if the last character of s is not
// already a newline.  Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(lvl int, s string) error {
	if lvl < l.level {
		return nil
	}
	calldepth := 2
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile|Lmodule) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, lvl, file, line)
	l.buf = append(l.buf, s...)
	l.buf = append(l.buf, '\n')
	_, err := l.out.Write(l.buf)
	if err != nil {
		return err
	}
	if l.bmulti {

		if fd, ok := l.out.(*os.File); ok {
			if fs, err := fd.Stat(); err == nil {
				if fs.Size() > int64(l.fmaxsize) {
					if newwrter := l.createIo(); newwrter != nil {
						fd.Close()
						l.out = newwrter
						return nil
					}

				}
			} else {
				return err
			}

		}
	}

	return err
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(Linfo, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) { l.Output(Linfo, fmt.Sprint(v...)) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) { l.Output(Linfo, fmt.Sprintln(v...)) }

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(Lfatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(Lfatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	l.Output(Lfatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(Lpanic, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(Lpanic, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(Lpanic, s)
	panic(s)
}

// Flags returns the output flags for the logger.
func (l *Logger) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

// SetFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

// Prefix returns the output prefix for the logger.
func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) Level() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}
func (l *Logger) SetLevel(lvl int) (old int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	old = l.level
	l.level = lvl
	return
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = w
}

// Flags returns the output flags for the standard logger.
func Flags() int {
	return std.Flags()
}

// SetFlags sets the output flags for the standard logger.
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

// These functions write to the standard logger.

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.

func Print(v ...interface{}) {
	std.Output(Linfo, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	std.Output(Linfo, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	std.Output(Linfo, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Output(Lfatal, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	std.Output(Lfatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	std.Output(Lfatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(Lpanic, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Output(Lpanic, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Output(Lpanic, s)
	panic(s)
}
func Level() int {
	return std.Level()
}
func SetLevel(lvl int) (old int) {
	old = std.SetLevel(lvl)
	return
}
