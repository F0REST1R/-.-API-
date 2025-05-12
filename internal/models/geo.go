package models

// GeoPoint представляет географические координаты
type GeoPoint struct {
    Lat float64 `json:"lat"` // Широта (-90 до 90)
    Lon float64 `json:"lon"` // Долгота (-180 до 180)
}