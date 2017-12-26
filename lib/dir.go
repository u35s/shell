package lib

import (
	"io/ioutil"
	"os"
)

func Dir(path string) ([]string, error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, v := range dir {
		if v.IsDir() {
			files = append(files, v.Name()+"/")
		} else {
			files = append(files, v.Name())
		}
	}
	return files, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
