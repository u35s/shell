package download

import (
	"log"
)

type DownloadManager struct {
	Downloading map[*ProgressBar]uint
}

func (this *DownloadManager) Init() {
	this.Downloading = make(map[*ProgressBar]uint)
}

func (this *DownloadManager) Download(filename, url string) {
	bar, err := Download(filename, url)
	if err == nil {
		log.Printf("[下载],%v,开始下载,总大小:%v,", bar.Name(), bar.Length())
		this.Downloading[bar] = 1
		//this.Downloading = append(this.Downloading, bar)
	} else {
		log.Printf("[下载],%v下载错误", err)
	}
}

func (this *DownloadManager) Timer() {
	for bar := range this.Downloading {
		if bar.Complete() {
			log.Printf("[下载],%v,下载完成,总大小:%v,", bar.Name(), bar.Length())
			delete(this.Downloading, bar)
		} else if bar.Error() != nil {
			log.Printf("[下载],%v下载错误", bar.Error())
			delete(this.Downloading, bar)
		} else {
			prg, length := bar.Progress(), bar.Length()
			log.Printf("[下载],%v,下载进度,%v/%v,%v/100", bar.Name(), prg, length, prg*100/length)
		}
	}
}
