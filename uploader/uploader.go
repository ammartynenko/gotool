package uploader

import (
	"sync"
	"fmt"
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
		Log   *log.Logger
	}
	FileInfo struct {
		Name string
		Path string
		Ext  string
		Size uint
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
//  загрузка одиночного файла с участием AJAX
//---------------------------------------------------------------------------
func (u *Uploader) UploadSingleAJAX(formName, fileUploadPath string, r *http.Request) *FileInfo {
	var f FileInfo

	//parse form
	err := r.ParseForm()
	if err != nil {
		u.Log.Println(err.Error)
		return nil
	}
	//get element from form
	_, fh, errOpen := r.FormFile(formName)
	if errOpen != nil {
		u.Log.Println(err.Error)
		return nil
	}
	//open file handler getting from form for read and `upload`
	fin, errInopen := fh.Open()
	if errInopen != nil {
		u.Log.Println(errInopen.Error)
		return nil
	}
	defer fin.Close()

	//create local file (handler) for write byte pipe
	fout, errFout := os.Create(filepath.Join(fileUploadPath, fh.Filename))
	if errFout != nil {
		u.Log.Println(errInopen.Error)
		return nil
	}
	defer fout.Close()

	//copy file to file
	_, errRead := io.Copy(fout, fin)
	if errRead != nil {
		u.Log.Println(errInopen.Error)
		return nil
	}

	//получаю данные по файлу
	info, errFi := fout.Stat()
	if errFi != nil {
		u.Log.Println(errInopen.Error)
		return nil
	}
	f = FileInfo{
		Name: fh.Filename,
		Path: filepath.Join(fileUploadPath + fh.Filename),
		Ext:  filepath.Ext(info.Name()),
		Size: uint(info.Size()),
	}

	//success upload file
	u.Log.Printf("Success upload file `%s`\n", fh.Filename)
	return &f
}

//---------------------------------------------------------------------------
// загрузка одиночного файла, функцию оптимально использовать как горутину при обработке `multiple`
//---------------------------------------------------------------------------
func (u *Uploader) goUploadSingle(fh *multipart.FileHeader, ajax bool, r *http.Request) *FileInfo {
	var f FileInfo

	//не аякс, значит горутина
	if !ajax {
		defer func() {
			u.Stock = append(u.Stock, f)
			u.Done()
		}()
	}

	fin, err_inopen := fh.Open()
	if err_inopen != nil {
		sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_inopen.Error()))
		return nil
	} else {
		defer fin.Close()
		//создаю файл на локальной машине - приемный файл
		fout, err_fout := os.Create(sr.Config().UPLOADFilesPath + fh.Filename)
		if err_fout != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_fout.Error()))
			return nil
		} else {
			defer fout.Close()
			//копирую файл
			_, err_read := io.Copy(fout, fin)
			if err_read != nil {
				sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_read.Error()))
				return nil
			} else {
				//получаю данные по файлу
				info, err_fi := fout.Stat()
				if err_fi != nil {
					sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_fi.Error()))
					return nil
				} else {
					f = FileInfo{Name: fh.Filename, Path: sr.Config().UPLOADFilesPath + fh.Filename, Ext: filepath.Ext(info.Name()), Size: uint(info.Size())}
				}
				//success upload file
				sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOADOK, fh.Filename))
				return &f
			}
		}

	}
}

//---------------------------------------------------------------------------
//  загрузка single/multiple формы без участия сторонних асинхронных методов типа ajax
//---------------------------------------------------------------------------
func (u *SpoukUploader) Upload(nameForm string, ajax bool, sr *SpoukCarry) (error) {
	//получаю список файлов мультиформ
	if !ajax {
		//в цикле открывая дескрипторы файлов для загрузки и файлы для принятия
		//запуск на каждый дескриптор горутину с синхронизацией загрузки
		//ajax == false
		err := sr.request.ParseMultipartForm(32 << 20)
		if err != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err.Error()))
			return err
		}
		formdata := sr.Request().MultipartForm
		listfiles := formdata.File[nameForm]
		fmt.Printf("LISTFILES===> %v\n", listfiles)

		w.Add(len(listfiles))
		for _, INfile := range listfiles {
			go u.goUploadSingle(INfile, false, sr)
		}
		//ожидаю завершение закачки
		w.Wait()
	} else {
		//ajax, запуск без синхрона, т.к. каждый выхов аякса дергает жту функцию, которая сама
		//запускается как горутина
		//sr.request.ParseMultipartForm(32 << 20)
		err := sr.request.ParseForm()
		if err != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err.Error()))
			return err
		}
		_, handler, err_form := sr.request.FormFile("uploadfile")
		if err_form != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_form.Error()))
			return err_form
		}
		f := u.goUploadSingle(handler, true, sr)
		if f != nil {
			u.Stock = append(u.Stock, *f)
		}
	}
	fmt.Println("[sync] All uploading")
	fmt.Printf("StockFIles: %v\n", u.Stock)
	return nil
}
