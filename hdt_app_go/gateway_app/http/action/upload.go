package action

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

import (
	"net/http"
)

func Upload(filename string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	return nil
}

func Upload2(q *http.Request, filename string) error {
	file, _, err := q.FormFile("face")
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer file.Close()
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}
