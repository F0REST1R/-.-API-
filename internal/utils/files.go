package utils

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func SaveImageFromURL(url, filename string) error {
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("ошибка загрузки: %v", err)
    }
    defer resp.Body.Close()

    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("ошибка создания файла: %v", err)
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    return err
}