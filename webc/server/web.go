package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/u35s/shell/lib"
	"github.com/u35s/shell/lib/crawler"

	"github.com/go-xorm/xorm"
	"github.com/u35s/gmod/lib/gtime"
	"github.com/u35s/gmod/mods/gconf"
)

type Web struct {
	db *xorm.Engine
}

func (this *Web) init() {
	conf := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
		gconf.String("mysql_user"), gconf.String("mysql_password"),
		gconf.String("mysql_server"), gconf.String("mysql_server_port"), gconf.String("mysql_database"))
	db, err := lib.NewEngine(conf)
	if err == nil {
		this.db = db
	}
	log.Printf("[数据库],连接,%v", err)
	// 设置访问的路由
	http.HandleFunc("/", this.index)
	http.HandleFunc("/img", this.img)
	http.HandleFunc("/download", this.download)
	http.HandleFunc("/crawler", this.crawler)
	http.HandleFunc("/ajax", this.ajax)
	http.HandleFunc("/upload", this.upload)
	http.HandleFunc("/statistics", this.statistics)

	//静态文件
	wd, err := os.Getwd()
	if err == nil {
		dd := wd + "/downloadd/"
		pd := wd + "/public/"
		http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(dd))))
		http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(pd))))
		log.Printf("[http],静态文件服务器地址:%v", wd)
	}

	go func() {
		addr := gconf.String("web_client_addr")
		log.Printf("[http],监听地址:%v", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Printf("[http],退出,%v", err)
		}
	}()
}

func (this *Web) upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	data := make(map[string]string)
	if r.Method == "POST" && r.MultipartForm != nil && r.MultipartForm.File != nil {
		if fhs, ok := r.MultipartForm.File["file"]; ok {
			for i := range fhs {
				file, err := fhs[i].Open()
				if lib.CheckError(err) {
					dir := "./downloadd/upload/"
					os.MkdirAll(dir, 0755)
					f, err := os.OpenFile(dir+lib.Utoa(gtime.Time())+"-"+fhs[i].Filename, os.O_WRONLY|os.O_CREATE, 0755)
					if lib.CheckError(err) {
						io.Copy(f, file)
						f.Close()
					}
					file.Close()
				}
			}
			log.Printf("[upload],上传文件%v个", len(fhs))
		}
	}
	t, err := template.ParseFiles("./template/upload.html")
	if err != nil {
		log.Printf("[http],%v,%v,模板不存在,%v,%v,", r.Method, r.URL.String(), err)
		return
	}
	t.Execute(w, data)
}

func (this *Web) download(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := make(map[string]string)
	url := r.Form.Get("url")
	name := r.Form.Get("name")
	if url != "" && name != "" {
		srv.dm.Download(name, url)
	}

	t, err := template.ParseFiles("./template/download.html")
	if err != nil {
		log.Printf("[http],%v,%v,模板不存在,%v,%v,", r.Method, r.URL.String(), err)
		return
	}
	t.Execute(w, data)
}

func (this *Web) crawler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.Form.Get("url")
	path := r.Form.Get("path")
	if url != "" {
		go crawler.AssetsCrawling(url, path, make(map[string]bool))
	}
	t, err := template.ParseFiles("./template/crawler.html")
	if err != nil {
		log.Printf("[http],%v,%v,模板不存在,%v", r.Method, r.URL.String(), err)
		return
	}
	t.Execute(w, nil)
}

func (this *Web) statistics(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.Form.Get("url")
	if _, err := this.db.Exec("update asset set count=count+1 where imgsrc=?", url); err != nil {
		log.Printf("[statistics],sql,%v", err)
	}
	url = string([]byte(url)[:strings.LastIndex(url, "/")])
	http.Redirect(w, r, "/files/"+url, http.StatusFound)
}

func (this *Web) ajax(w http.ResponseWriter, r *http.Request) {
	ret, err := this.db.Query("select imgsrc from asset where rand_count = (select min(rand_count) from asset) order by rand() limit 15")
	slc := make([]string, 0)
	buf := &bytes.Buffer{}
	for _, v := range ret {
		file := string(v["imgsrc"])
		if !lib.IsMediaWithExts(file, []string{"MP4"}) {
			slc = append(slc, file)
			buf.WriteByte('"')
			buf.WriteString(file)
			buf.WriteByte('"')
			buf.WriteByte(',')
		}
	}
	if buf.Len() > 0 {
		sql := fmt.Sprintf("update asset set rand_count=rand_count+1 where imgsrc in (%v)", string(buf.Next(buf.Len()-1)))
		if _, err := this.db.Exec(sql); err != nil {
			log.Printf("[statistics],sql,%v", err)
		}
	}
	bts, err := json.Marshal(&slc)
	if err == nil {
		w.Write(bts)
	}
}

func (this *Web) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func (this *Web) img(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template/img.html")
	if err != nil {
		log.Printf("[http],%v,%v,模板不存在,%v", r.Method, r.URL.String(), err)
		return
	}
	t.Execute(w, nil)
}

func (this *Web) writeIPAddr() {
	ip := string(lib.HttpGet("http://ip.cn"))
	_, err := this.db.Exec("insert into ip_addr (ip_addr,time) values (?,?)", ip, gtime.Time())
	log.Printf("[ip addr],sql,%v,%v", ip, err)
}
