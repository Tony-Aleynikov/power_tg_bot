package repository

import (
	"encoding/json"
	"fmt"
	"os"
)

type repository struct {
	filePath string
}

func NewRepository() (repository, error) {
	if os.Getenv("FILE_LOCATION") == "" {
		return repository{}, fmt.Errorf("Не задан FILE_LOCATION")
	}

	_, err := os.Stat(os.Getenv("FILE_LOCATION"))
	if os.IsNotExist(err) {
		return repository{}, fmt.Errorf("Файл не существует")
	}

	return repository{
		filePath: os.Getenv("FILE_LOCATION"),
	}, nil
}

func (r repository) Save(data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка создания JSON: %w", err)
	}

	err = os.WriteFile(r.filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}
	return nil
}
