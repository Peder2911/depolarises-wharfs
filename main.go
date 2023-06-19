package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
        "path"
)

//go:embed home.html
var home_html []byte 

//go:embed upload.html
var upload_html []byte 

type FileHandler struct {
   upload []byte
   filefolder string
}

func NewFileHandler (upload []byte, filefolder string) FileHandler {
   return FileHandler{
      upload: upload,
      filefolder: filefolder,
   }
}

func (f *FileHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
   r.ParseMultipartForm(1000 << 20)
   uploaded, header, err := r.FormFile("upload")
   if err != nil {
      log.Println(fmt.Sprintf("Error when uploading file: %s", err))
      w.WriteHeader(400)
      return
   }
   defer uploaded.Close()

   file, err := os.Create(path.Join(f.filefolder, header.Filename))
   if err != nil {
      w.WriteHeader(500)
      return
   }
   defer file.Close()

   io.Copy(file, uploaded)
   w.WriteHeader(201)
   w.Write(f.upload)
}

type HomeHandler struct {
   home []byte 
}

func NewHomeHandler(home []byte) HomeHandler{
   return HomeHandler{home}
}

func (f *HomeHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
   if r.URL.Path != "/"{
      w.WriteHeader(404)
      return
   }
   w.Write(f.home)
}

func main(){
   err := os.MkdirAll("files", 0755)
   if err != nil {
      panic(fmt.Sprintf("Failed to create files directory: %s", err))
   }

   file_handler := NewFileHandler(upload_html, "files")
   home_handler := NewHomeHandler(home_html)

   mux := http.NewServeMux()
   mux.Handle("/upload", &file_handler)
   mux.Handle("/get", &home_handler) 
   mux.Handle("/", &home_handler) 
   http.ListenAndServe(":8000", mux)
}
