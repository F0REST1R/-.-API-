package examples

import (
    "log"
    "prog/internal/api"
    "prog/internal/models"
    "prog/internal/utils"
)

func Cluster() {
    // 1. Получаем точки через API Локатора
    locator := api.NewLocatorClient()
    points, err := locator.GetPOI(55.751574, 37.573856, 1000, "кафе")
    if err != nil {
        log.Fatalf("Ошибка получения POI: %v", err)
    }

    // 2. Кластеризуем точки (радиус 0.01° ~ 1 км)
    clusters := utils.ClusterPoints(points, 0.01)

    // 3. Генерируем URL для статической карты
    center := models.GeoPoint{Lat: 55.751574, Lon: 37.573856}
    mapURL := api.GenerateClusterMapURL(clusters, center, 13)

    // 4. Сохраняем карту как изображение
    if err := utils.SaveImageFromURL(mapURL, "cluster_map.png"); err != nil {
        log.Fatalf("Ошибка сохранения карты: %v", err)
    }

    log.Println("Карта сохранена как cluster_map.png")
}