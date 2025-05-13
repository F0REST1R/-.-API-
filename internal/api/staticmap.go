package api

import (
    "fmt"
    "net/url"
    "prog/internal/models"
)

const (
    StaticMapURL = "https://static-maps.yandex.ru/1.x/"
    MapWidth     = 650
    MapHeight    = 450
)

func GenerateClusterMapURL(clusters []models.Cluster, center models.GeoPoint, zoom int) string {
    params := url.Values{}
    params.Add("l", "map")
    params.Add("size", fmt.Sprintf("%d,%d", MapWidth, MapHeight))
    params.Add("ll", fmt.Sprintf("%f,%f", center.Lon, center.Lat))
    params.Add("z", fmt.Sprintf("%d", zoom))

    for i, cluster := range clusters {
        params.Add("pt", fmt.Sprintf("%f,%f,pm2lbl%d", cluster.Center.Lon, cluster.Center.Lat, i+1))
    }

    return StaticMapURL + "?" + params.Encode()
}