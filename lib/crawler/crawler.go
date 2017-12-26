package crawler

import (
	"log"
	"github.com/u35s/shell/lib"
	"github.com/u35s/shell/lib/download"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func DownloadMedia(url, path string) {
	fullpath := url + path
	bar, err := download.Download(lib.RemoveHttpPrefix(fullpath), fullpath)
	if err == nil {
		log.Printf("[crawler],开始下载:%v", fullpath)
		timeOutNum := 0
		for {
			lastLength := bar.Progress()
			if !bar.Complete() && bar.Error() == nil {
				time.Sleep(time.Millisecond * 500)
				prg, length := bar.Progress(), bar.Length()
				if timeOutNum > 5 {
					bar.Break()
					log.Printf("[下载],%v,无速率打断,次数:%v", bar.Name(), timeOutNum)
				}
				if lastLength == prg {
					log.Printf("[下载],%v,下载无速率,上一次大小:%v,进度:%v/%v,次数:%v", bar.Name(), lastLength, prg, length, timeOutNum)
					timeOutNum++
				} else {
					timeOutNum = 0
					log.Printf("[下载],%v,下载进度:%v/%v,百分比:%v/100", bar.Name(), prg, length, prg*100/length)
				}
			} else if bar.Complete() {
				log.Printf("[下载],%v,下载完成", bar.Name())
				break
			} else {
				log.Printf("[下载],%v,下载错误,%v", bar.Name(), bar.Error())
				break
			}
		}
	} else {
		log.Printf("[crawler],下载错误:%v,%v", fullpath, err)
	}
}

func AssetsCrawling(url, path string, check map[string]bool) {
	if _, ok := check[url+path]; ok {
		return
	} else {
		check[url+path] = true
	}
	doc, err := goquery.NewDocument(url + path)
	if err != nil {
		log.Printf("[crawler],new document err:%v", err)
		return
	}
	log.Printf("[crawler],查找中,%v,%v", url, path)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			media := lib.IsMedia(href)
			dir := lib.IsDir(href)
			switch {
			case media && lib.AbsUrl(href):
				DownloadMedia(href, "")
			case media && lib.AbsPath(href):
				DownloadMedia(url, href)
			case media && !lib.AbsPath(href):
				DownloadMedia(url, path+href)
			case dir && lib.AbsUrl(href):
				AssetsCrawling(url, "", check)
				log.Printf("[crawler],递归,URL:%v,绝对url:%v", url, href)
			case dir && lib.AbsPath(href) && len(href) > len(path):
				AssetsCrawling(url, href, check)
				log.Printf("[crawler],递归,URL:%v,绝对路径:%v", url, href)
			case dir && !lib.AbsPath(href):
				AssetsCrawling(url, path+href, check)
				log.Printf("[crawler],递归,URL:%v,相对路径:%v,%v", url, path, href)
			default:
				log.Printf("[crawler],URL:%v,PATH:%v,丢弃href:%v", url, path, href)
			}
		}
	})
}
