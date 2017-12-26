package download

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"github.com/u35s/shell/lib"
)

type ProgressBar struct {
	start, end uint
	name       string
	length     int
	progress   int
	err        error
	rep        *http.Response
}

func (this *ProgressBar) Name() string  { return this.name }
func (this *ProgressBar) Length() int   { return this.length }
func (this *ProgressBar) Progress() int { return this.progress }

func (this *ProgressBar) Error() error   { return this.err }
func (this *ProgressBar) Break()         { this.rep.Body.Close() }
func (this *ProgressBar) Complete() bool { return this.length <= this.progress }

func (this *ProgressBar) Write(p []byte) (int, error) {
	length := len(p)
	this.progress += length
	return length, nil
}

func Download(filename, url string) (*ProgressBar, error) {
	dir, name := filepath.Split(filename)
	if name == "" {
		return nil, errors.New("no file name")
	}
	dir = "downloadd/" + dir
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	filename = filepath.Join(dir, name)
	if ok, err := lib.PathExists(filename); err != nil {
		return nil, err
	} else if ok {
		return nil, errors.New("文件已存在")
	}
	rep, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	output, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		rep.Body.Close()
		return nil, err
	}
	bar := &ProgressBar{name: filename, length: lib.Atoi(rep.Header.Get("Content-Length")), rep: rep}
	writer := io.MultiWriter(output, bar)
	go func() {
		defer rep.Body.Close()
		defer output.Close()
		_, bar.err = io.Copy(writer, rep.Body)
		if bar.err != nil {
			os.Remove(filename)
		}
	}()
	return bar, nil
}
