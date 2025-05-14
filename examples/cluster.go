package examples

import (
    "log"
    "prog/internal/api"
    "prog/internal/models"
    "prog/internal/utils"
)

func Cluster() {
    locator := api.NewLocatorClient()
    
    // Тестовые координаты (Красная площадь)
    points, err := locator.GetPOI(55.753544, 37.621202, 1000, "кафе")
    if err != nil {
        log.Fatalf("Ошибка получения POI: %v", err)
    }

    if len(points) == 0 {
        // Тестовые точки для проверки
        points = []models.GeoPoint{
            {Lat: 55.753544, Lon: 37.621202},
            {Lat: 55.752023, Lon: 37.617499},
        }
        log.Println("Используются тестовые точки")
    }

    clusters := utils.ClusterPoints(points, 0.005)
    log.Printf("Создано кластеров: %d", len(clusters))

    center := models.GeoPoint{Lat: 55.753544, Lon: 37.621202}
    mapURL := api.GenerateClusterMapURL(clusters, center, 15)
    log.Printf("URL карты: %s", mapURL)

    if err := utils.SaveImageFromURL(mapURL, "cluster_map.png"); err != nil {
        log.Fatalf("Ошибка сохранения: %v", err)
    }
    log.Println("Карта сохранена как cluster_map.png")
}