package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveProductImage(file multipart.File, header *multipart.FileHeader) (string, error){
	defer file.Close()
	ext := filepath.Ext(header.Filename)
	filename := filepath.Join("uploads/products",
time.Now().Format("20060102150405")+ext)
out, err := os.Create(filename)
if err !=nil{
	return  "", err
}
defer out.Close()

_, err = io.Copy(out, file)
if err != nil{
	return  "", err
}
return filename, nil
}