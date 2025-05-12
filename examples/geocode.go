package examples

import (
    "fmt"
    "log"
    "prog/internal/api"
)

func RunGeocodingExample() {
    client := api.NewGeocodeClient()

    // Тестовые координаты (Кремль, Эрмитаж, случайная точка)
    testPoints := []struct {
        name      string
        lat, lon float64
    }{
        {"Московский Кремль", 55.752023, 37.617499},
        {"Эрмитаж", 59.934280, 30.335098},
        {"Рандомная точка в Москве (точку вводим самостоятельно)", 55.733842, 37.588144}, //Рандомную точку вводим самостоятельно
    }

    for _, point := range testPoints {
        fmt.Printf("\n🔍 %s\n", point.name)
        fmt.Printf("Координаты: %.6f, %.6f\n", point.lat, point.lon)

        address, err := client.ReverseGeocode(point.lat, point.lon)
        if err != nil {
            log.Printf("Ошибка: %v", err)
            continue
        }

        fmt.Println("\n📌 Результат геокодирования:")
        fmt.Printf("Краткий адрес: %s\n", address.Formatted)
        
        if address.Description != "" {
            fmt.Printf("Объект: %s\n", address.Description)
        }
        
        if address.Street != "" {
            streetInfo := address.Street
            if address.HouseNumber != "" {
                streetInfo += ", " + address.HouseNumber
            }
            fmt.Printf("Улица: %s\n", streetInfo)
        }
        
        if address.City != "" {
            fmt.Printf("Город: %s\n", address.City)
        }
        
        if address.PostalCode != "" {
            fmt.Printf("Почтовый индекс: %s\n", address.PostalCode)
        }
        
        fmt.Printf("Полный адрес: %s\n", address.FullAddress)
    }
}