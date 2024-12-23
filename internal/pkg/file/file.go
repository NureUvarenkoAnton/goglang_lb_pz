package file

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

func CreateFile(data any, ext string) string {
	file, _ := os.Create("tmp/" + uuid.New().String() + "." + ext)
	jsonData, _ := json.Marshal(data)
	file.Write(jsonData)
	return file.Name()
}

func DeleteFile(name string) error {
	err := os.Remove(name)
	return err
}
