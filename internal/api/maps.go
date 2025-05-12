package api

import (
	"fmt"
	"prog/config"
	"prog/internal/models"
	"strconv"
	"github.com/go-resty/resty/v2"
)


func (c *YaMapsClient) GenerateStaticMap(center models.GeoPoint, zoom int, size models.MapSize, markers []models.Marker) ([]byte, error) {
    client := resty.New()
    
    // Проверка входных параметров
    if zoom < 1 || zoom > 17 {
        return nil, fmt.Errorf("недопустимый уровень масштабирования (должен быть 1-17)")
    }
    if size.Width < 1 || size.Height < 1 || size.Width > 650 || size.Height > 450 {
        return nil, fmt.Errorf("недопустимый размер карты (макс. 650x450 пикселей)")
    }

    // Базовые параметры запроса
    params := map[string]string{
        "ll":    fmt.Sprintf("%f,%f", center.Lon, center.Lat),
        "z":     strconv.Itoa(zoom),
        "l":     "map",
        "size":  fmt.Sprintf("%d,%d", size.Width, size.Height),
        "apikey": config.GetAPIMap(),
    }

    // Добавляем маркеры (исправленный формат)
    for i, m := range markers {
        if i >= 50 { // Лимит Яндекс.Карт
            break
        }
        params[fmt.Sprintf("pt%d", i+1)] = fmt.Sprintf("%f,%f,pm2%dl%d", m.Lon, m.Lat, m.Color, i+1)
    }

    // Отправляем запрос
    resp, err := client.R().
        SetQueryParams(params).
        Get("https://static-maps.yandex.ru/1.x/")

    if err != nil {
        return nil, fmt.Errorf("ошибка запроса: %v", err)
    }

    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("ошибка API (статус %d): %s", resp.StatusCode(), string(resp.Body()))
    }

    return resp.Body(), nil
}

