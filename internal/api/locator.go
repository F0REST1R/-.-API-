package api

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "net/url"
    "prog/config"
    "prog/internal/models"
    "time"

    "github.com/go-resty/resty/v2"
)

//реализует взаимодействие с API
type LocatorClient struct {
    client *resty.Client
    apiKey string
}

func NewLocatorClient() *LocatorClient {
    apiKey := config.GetAPILocator()
    return &LocatorClient{
        client: resty.New().
            SetBaseURL("https://search-maps.yandex.ru").
            SetTimeout(10 * time.Second),
        apiKey: apiKey,
    }
}

// ErrorResponse для обработки XML ошибок
type ErrorResponse struct {
    XMLName xml.Name `xml:"error"`
    Message string   `xml:"message"`
}

//Получение POI
func (c *LocatorClient) GetPOI(lat, lon float64, radius int, category string) ([]models.GeoPoint, error) {
    resp, err := c.client.R().
        SetQueryParams(map[string]string{
            "apikey":  c.apiKey,
            "text":    url.QueryEscape(category),
            "ll":      fmt.Sprintf("%f,%f", lon, lat),
            "spn":     fmt.Sprintf("%f,%f", 1, 1),
            "rspn":    "1",
            "lang":    "ru_RU",
            "results": "100",
            "format":  "json", // Явно запрашиваем JSON
        }).
        Get("/v1/")

    if err != nil {
        return nil, fmt.Errorf("ошибка запроса: %v", err)
    }

    // Проверяем Content-Type
    contentType := resp.Header().Get("Content-Type")
    if contentType == "application/xml" || contentType == "text/xml" {
        var apiError ErrorResponse
        if err := xml.Unmarshal(resp.Body(), &apiError); err == nil {
            return nil, fmt.Errorf("API ошибка: %s", apiError.Message)
        }
        return nil, fmt.Errorf("неожиданный XML-ответ: %s", resp.Body())
    }

    var result struct {
        Features []struct {
            Geometry struct {
                Coordinates []float64 `json:"coordinates"`
            } `json:"geometry"`
        } `json:"features"`
    }

    if err := json.Unmarshal(resp.Body(), &result); err != nil {
        return nil, fmt.Errorf("ошибка парсинга JSON: %v (тело ответа: %s)", err, resp.Body())
    }

    var points []models.GeoPoint
    for _, feature := range result.Features {
        if len(feature.Geometry.Coordinates) >= 2 {
            points = append(points, models.GeoPoint{
                Lat: feature.Geometry.Coordinates[1],
                Lon: feature.Geometry.Coordinates[0],
            })
        }
    }

    return points, nil
}