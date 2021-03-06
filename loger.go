package loger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	WStdout = 1 << iota
	WFile
)

type LogFile struct {
	Root     string
	FileName string
}

func (l *LogFile) path() string {
	now := time.Now()

	l.FileName = strconv.Itoa(now.Year()) + "_" +
		strconv.Itoa(int(now.Month())) + "_" +
		strconv.Itoa(now.Day()) +
		".log"

	return filepath.Join(l.Root, l.FileName)
}

type Loger struct {
	LogFile
	WType  int
	File   *os.File
	lock   *sync.RWMutex
	logger *log.Logger
}

func NewLoger(root string) *Loger {
	l := &Loger{
		LogFile: LogFile{
			Root: root,
		},
		lock:   new(sync.RWMutex),
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
		WType:  WStdout | WFile,
	}
	l.initDir()
	return l
}

func (l Loger) initDir() (err error) {
	_, err = os.Stat(l.Root)

	if err != nil {
		if os.IsNotExist(err) {
			// use 0777 for test. unsafe
			os.Mkdir(l.Root, 0755)
		} else {
			return
		}
	}
	return
}

// open and set io writer
func (l *Loger) open() (err error) {
	if l.WType&WFile != 0 {
		l.File, err = os.OpenFile(l.path(),
			os.O_APPEND|os.O_CREATE|os.O_RDWR,
			0755)

		if err != nil {
			log.Printf("open log file error %s \r\n", err)
			return
		}
		var w io.Writer
		if err != nil {
			w = os.Stdout
			return err
		} else {
			if l.WType&WStdout != 0 {
				w = io.MultiWriter(l.File, os.Stdout)
			} else {
				w = io.MultiWriter(l.File)
			}
		}
		l.setOutput(w)
	} else {
		// just use std out for default
		l.setOutput(os.Stdout)
		return
	}
	return
}

func (l *Loger) setOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// short for close file
func (l *Loger) close() (err error) {
	err = l.File.Close()
	return
}

func (l *Loger) SetFlag(flag int) {
	l.logger.SetFlags(flag)
}

func (l *Loger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *Loger) SetWriteType(t int) {
	l.WType = t
}

func (l *Loger) Info(v ...interface{}) {
	defer l.close()
	defer l.lock.Unlock()
	l.open()
	l.lock.Lock()
	l.logger.SetPrefix("[INFO]: ")
	l.logger.Output(2, fmt.Sprintln(v...))
}
func (l *Loger) Warn(v ...interface{}) {
	defer l.close()
	defer l.lock.Unlock()
	l.open()
	l.lock.Lock()
	l.logger.SetPrefix("[WARN]: ")
	l.logger.Output(2, fmt.Sprintln(v...))
}

func (l *Loger) Error(v ...interface{}) {
	defer l.close()
	defer l.lock.Unlock()
	l.open()
	l.lock.Lock()
	l.logger.SetPrefix("[ERROR]:")
	l.logger.Output(2, fmt.Sprintln(v...))
}

func (l *Loger) Panic(v ...interface{}) {
	defer l.close()
	defer l.lock.Unlock()
	l.open()
	l.lock.Lock()
	l.logger.SetPrefix("[PANIC]:")
	l.logger.Panicln(v...)
}

func (l *Loger) Println(v ...interface{}) {
	defer l.close()
	defer l.lock.Unlock()
	l.open()
	l.lock.Lock()
	l.logger.SetPrefix("")
	l.logger.Output(2, fmt.Sprintln(v...))
}

func (l *Loger) Printf(format string, v ...interface{}) {
	defer l.close()
	defer l.lock.Unlock()
	l.open()
	l.lock.Lock()
	l.logger.SetPrefix("")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}
