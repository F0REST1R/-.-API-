package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"prog/config"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type YaMapsClient struct {
	geocoderClient  *resty.Client
	staticMapClient *resty.Client
}

func NewYaMapsClient() *YaMapsClient {
	return &YaMapsClient{
		geocoderClient: resty.New().SetBaseURL("https://geocode-maps.yandex.ru").
			SetTimeout(10 * time.Second),

		staticMapClient: resty.New().SetBaseURL("https://static-maps.yandex.ru").
			SetTimeout(10 * time.Second),
	}
}

func (c *YaMapsClient) Geocode(address string) (float64, float64, error) {
    resp, err := c.geocoderClient.R().
        SetQueryParams(map[string]string{
            "apikey":  config.GetAPIGeo(),
            "geocode": address,
            "format":  "json",
            "lang":    "ru_RU", // Добавляем язык
        }).
        Get("/1.x/")

    if err != nil {
        return 0, 0, fmt.Errorf("ошибка запроса: %v", err)
    }

    // Раскомментируйте для отладки
    //fmt.Printf("Raw response: %s\n", resp.Body())

    var result struct {
        Response struct {
            GeoObjectCollection struct {
                MetaDataProperty struct {
                    GeocoderResponseMetaData struct {
                        Request string `json:"request"`
                        Found   string `json:"found"`
                    } `json:"GeocoderResponseMetaData"`
                } `json:"metaDataProperty"`
                FeatureMember []struct {
                    GeoObject struct {
                        MetaDataProperty struct {
                            GeocoderMetaData struct {
                                Text string `json:"text"`
                                Kind string `json:"kind"`
                            } `json:"GeocoderMetaData"`
                        } `json:"metaDataProperty"`
                        Point struct {
                            Pos string `json:"pos"`
                        } `json:"Point"`
                    } `json:"GeoObject"`
                } `json:"featureMember"`
            } `json:"GeoObjectCollection"`
        } `json:"response"`
    }

    if err := json.Unmarshal(resp.Body(), &result); err != nil {
        return 0, 0, fmt.Errorf("ошибка парсинга JSON: %v", err)
    }

    // Проверяем наличие результатов
    if len(result.Response.GeoObjectCollection.FeatureMember) == 0 {
        return 0, 0, fmt.Errorf("адрес не найден. Ответ API: %s", resp.Body())
    }

    // Берем первый результат
    pos := strings.Split(
        result.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos, 
        " ",
    )
    if len(pos) != 2 {
        return 0, 0, errors.New("неверный формат координат")
    }

    lon, err := strconv.ParseFloat(pos[0], 64)
    if err != nil {
        return 0, 0, fmt.Errorf("ошибка преобразования долготы: %v", err)
    }

    lat, err := strconv.ParseFloat(pos[1], 64)
    if err != nil {
        return 0, 0, fmt.Errorf("ошибка преобразования широты: %v", err)
    }

    return lat, lon, nil
}
