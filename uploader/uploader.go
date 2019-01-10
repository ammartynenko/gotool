package uploader

import (
	"sync"
	"os"
	"io"
	"path/filepath"
	"mime/multipart"
	"log"
	"net/http"
)

type (
	Uploader struct {
		stock []FileInfo
		sync.WaitGroup
		sync.RWMutex
		Log   *log.Logger
	}
	FileInfo struct {
		Name string
		Path string
		Ext  string
		Size uint
	}
	UploadConfig struct {
		Uploadpath string
		Ajax       bool
		Req        *http.Request
		FormName   string
	}
	fileConfig struct {
		Fh         *multipart.FileHeader
		Uploadpath string
	}
)

const (
	prefix = "[gotool][uploader] "
)

//---------------------------------------------------------------------------
//  construct
//---------------------------------------------------------------------------
func NewUploader() *Uploader {
	s := &Uploader{
		Log: log.New(os.Stdout, prefix, log.Lshortfile|log.Ldate|log.Ltime),
	}
	return s
}

//---------------------------------------------------------------------------
//  организация пайпа для загрузки файлов
//  урправляющий загрузкой  + горутины на каждый файл
//  управляющий парсит форму на список файлов  и запускает горутины на обработку
//  файловых дескрипторов
//---------------------------------------------------------------------------
func (u *Uploader) UploadFiles(c *UploadConfig) error {

	//parse form for get file handlers
	err := c.Req.ParseMultipartForm(32 << 20)
	if err != nil {
		u.Log.Printf(err.Error())
		return err
	}
	formdata := c.Req.MultipartForm
	listfiles := formdata.File[c.FormName]

	//run gorutines for  upload files
	for _, INfile := range listfiles {
		u.Add(1)
		go u.goup(&fileConfig{
			Fh:         INfile,
			Uploadpath: c.Uploadpath,
		})
	}
	//ожидаю завершение закачки
	u.Wait()
	u.Log.Printf("Success upload all files")
	for i, x := range u.stock {
		u.Log.Printf("%d. %s\n", i, x.Name)
	}
	return nil
}
func (u *Uploader) goup(f *fileConfig) {
	defer u.Done()
	fin, err := f.Fh.Open()
	if err != nil {
		u.Log.Println(err)
		return
	}
	defer fin.Close()

	//создаю файл на локальной машине - приемный файл
	fout, err := os.Create(filepath.Join(f.Uploadpath, f.Fh.Filename))
	if err != nil {
		u.Log.Println(err)
		return
	}
	defer fout.Close()

	//копирую файл
	_, err = io.Copy(fout, fin)
	if err != nil {
		u.Log.Println(err)
		return
	}

	//получаю данные по файлу
	info, err := fout.Stat()
	if err != nil {
		u.Log.Println(err)
		return
	}

	fu := FileInfo{
		Name: f.Fh.Filename,
		Path: filepath.Join(f.Uploadpath, f.Fh.Filename),
		Ext:  filepath.Ext(info.Name()),
		Size: uint(info.Size()),
	}
	//lock shared slice
	u.Lock()
	u.stock = append(u.stock, fu)
	u.Unlock()

	//success upload file
	u.Log.Printf("Success upload file `%s`\n", f.Fh.Filename)
	return
}
