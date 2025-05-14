package api

import (
    "encoding/json"
    "fmt"
    "log"
    "prog/config"
    "prog/internal/models"
    "strings"
    "time"

    "github.com/go-resty/resty/v2"
)

type GeocodeClient struct {
    client *resty.Client
    apiKey string
}

func NewGeocodeClient() *GeocodeClient {
    apiKey := config.GetAPIGeo()
    if apiKey == "" {
        log.Fatal("API ключ для геокодирования не установлен в .env (YA_MAPS_API_KEY_GEO)")
    }

    return &GeocodeClient{
        client: resty.New().
            SetBaseURL("https://geocode-maps.yandex.ru").
            SetTimeout(10 * time.Second),
        apiKey: apiKey,
    }
}

// Геообратное кодирование
func (c *GeocodeClient) ReverseGeocode(lat, lon float64) (*models.Address, error) {
    resp, err := c.client.R().
        SetQueryParams(map[string]string{
            "apikey":   c.apiKey,
            "format":   "json",
            "geocode":  fmt.Sprintf("%f,%f", lon, lat), // Яндекс использует lon,lat
            "lang":     "ru_RU",
            "results":  "1",
            "kind":     "house", // Точность до дома
        }).
        Get("/1.x/")

    if err != nil {
        return nil, fmt.Errorf("ошибка запроса: %v", err)
    }

    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("API вернул ошибку %d: %s", resp.StatusCode(), resp.Body())
    }

    var response struct {
        Response struct {
            GeoObjectCollection struct {
                FeatureMember []struct {
                    GeoObject struct {
                        Name string `json:"name"`
                        MetaDataProperty struct {
                            GeocoderMetaData struct {
                                Text    string `json:"text"`
                                Address struct {
                                    CountryCode string `json:"country_code"`
                                    PostalCode  string `json:"postal_code"`
                                    Province   string `json:"province"`   // Область/край
                                    Locality   string `json:"locality"`   // Город
                                    District   string `json:"district"`   // Район города
                                    Thoroughfare string `json:"thoroughfare"` // Улица
                                    Premise     string `json:"premise"`       // Дом
                                    Building    string `json:"building"`     // Корпус
                                } `json:"address"`
                            } `json:"GeocoderMetaData"`
                        } `json:"metaDataProperty"`
                    } `json:"GeoObject"`
                } `json:"featureMember"`
            } `json:"GeoObjectCollection"`
        } `json:"response"`
    }

    if err := json.Unmarshal(resp.Body(), &response); err != nil {
        return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
    }

    if len(response.Response.GeoObjectCollection.FeatureMember) == 0 {
        return nil, fmt.Errorf("адрес не найден для координат %.6f,%.6f", lat, lon)
    }

    geo := response.Response.GeoObjectCollection.FeatureMember[0].GeoObject
    meta := geo.MetaDataProperty.GeocoderMetaData
    addr := meta.Address

    // Формируем краткий адрес
    var parts []string
    if addr.Thoroughfare != "" {
        parts = append(parts, addr.Thoroughfare)
    }
    if addr.Premise != "" {
        parts = append(parts, addr.Premise)
    }
    shortAddress := strings.Join(parts, ", ")

    // Если нет улицы/дома, используем название объекта
    if shortAddress == "" && geo.Name != "" {
        shortAddress = geo.Name
    }

    return &models.Address{
        Formatted:    shortAddress,
        FullAddress:  meta.Text,
        Country:      addr.CountryCode,
        Region:       addr.Province,
        City:         addr.Locality,
        District:     addr.District,
        Street:       addr.Thoroughfare,
        HouseNumber:  addr.Premise,
        Building:     addr.Building,
        PostalCode:   addr.PostalCode,
        Description:  geo.Name,
    }, nil
}