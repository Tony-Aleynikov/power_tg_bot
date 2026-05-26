package repository

import (
	"encoding/json"
	"fmt"
	"os"
)

type repository struct {
	filePath string
}

func NewRepository() repository {
	return repository{
		filePath: os.Getenv("FILE_LOCATION"),
	}
}

func (r repository) Save(data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка создания JSON: %w", err)
	}
	if r.filePath == "" {
		return fmt.Errorf("Не задан filePath")
	}

	_, err = os.Stat(r.filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Файл не существует")
	}

	err = os.WriteFile(r.filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}
	return nil
}
