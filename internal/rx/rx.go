package rx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/gabriel-vasile/mimetype"
	"github.com/isuquo/templatemaker/internal/models"
)

type TestStruct struct {
	Name string
	ID   string
}

func Test(t *models.Template, files []*multipart.FileHeader) (string, error) {

	fmt.Printf("Template: %+v\n", t)
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return "", err
		}
		defer file.Close()

		mime, err := mimetype.DetectReader(file)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return "", err
		}

		fileType := mime.String()
		file.Seek(0, io.SeekStart)

		if fileType == "application/json" {
			var jsonData interface{}
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&jsonData); err != nil {
				fmt.Printf("Failed to decode JSON: %s\n", err.Error())
				return "", err
			}
			fmt.Printf("JSON Content: %+v\n", jsonData)
		}
	}

	return "", errors.New("AAAAAAAAAA")
}

func GetStructs(t *models.Template) ([]TestStruct, error) {
	var structs []TestStruct

	structs = append(structs, TestStruct{
		Name: t.Name,
		ID:   t.ID,
	})

	return structs, nil
}
