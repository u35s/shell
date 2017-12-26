package lib

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"net/http"
	"io/ioutil"
	"time"
)

type uint = uint64

func StartCmd(cmd string) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	err := exec.Command(head, parts...).Start()
	if err == nil {
		log.Printf("[执行命令],%v", cmd)
	} else {
		log.Printf("[执行命令错误],%v,%v", cmd, err)
	}
}

func RunCmd(cmd string) []byte {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err == nil {
		log.Printf("[执行命令],%v,%v", cmd, string(out))
		return out
	} else {
		log.Printf("[执行命令错误],%v,%v", cmd, err)
	}
	return nil
}

var MEDIA_ARR []string = []string{"JPG", "PNG", "MP4"}

func IsMedia(path string) bool {
	return IsMediaWithExts(path, MEDIA_ARR)
}

func IsMediaWithExts(path string, exts []string) bool {
	_, name := filepath.Split(path)
	arr := strings.Split(name, ".")
	if len(arr) != 2 {
		return false
	}
	ext := strings.ToUpper(arr[1])
	for i := range exts {
		if exts[i] == ext {
			return true
		}
	}
	return false
}

const whttp, whttps string = "http://", "https://"

func AbsUrl(url string) bool { return strings.HasPrefix(url, whttp) || strings.HasPrefix(url, whttps) }

func AbsPath(path string) bool { return strings.HasPrefix(path, "/") }

func IsDir(path string) bool { return strings.HasSuffix(path, "/") }

func RemoveHttpPrefix(url string) string {
	if strings.HasPrefix(url, whttp) {
		return url[len(whttp):]
	} else if strings.HasPrefix(url, whttps) {
		return url[len(whttps):]
	}
	return url
}

func HttpGet(url string) []byte {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	if req, err := http.NewRequest("GET", url, nil); err == nil {
		req.Header.Add("User-agent", "curl/7.54.0")
		res, err := client.Do(req)
		if err == nil {
			result, err := ioutil.ReadAll(res.Body)
			if err == nil {
				return result
			}
		}
	}
	return nil
}


