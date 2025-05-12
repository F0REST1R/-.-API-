package examples

import (
	"log"
	"os"
	"prog/internal/api"
	"prog/internal/models"
)

func ExamplMaps() {
    client := api.NewYaMapsClient()
    
    // Тестируем геокодирование
    lat, lon, err := client.Geocode("Москва, Красная площадь")
    if err != nil {
        log.Fatalf("Geocode error: %v", err)
    }
    log.Printf("Координаты: %.6f, %.6f", lat, lon)
    
    // Тестируем статическую карту
    img, err := client.GenerateStaticMap(
        models.GeoPoint{Lat: lat, Lon: lon},
        15,
        models.MapSize{Width: 650, Height: 450},
        []models.Marker{
            {
                GeoPoint: models.GeoPoint{Lat: lat, Lon: lon},
                Color:    "0xFF0000",
                Label:    "Красная площадь",
            },
        },
    )
    if err != nil {
        log.Fatalf("Map generation error: %v", err)
    }
    
    if err := os.WriteFile("map.png", img, 0644); err != nil {
        log.Fatalf("File save error: %v", err)
    }
    log.Println("Карта сохранена как map.png")
}