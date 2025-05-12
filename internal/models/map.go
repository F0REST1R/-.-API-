package models

// MapSize определяет размер изображения карты
type MapSize struct {
    Width  int `json:"width"`  // Ширина в пикселях (макс. 650)
    Height int `json:"height"` // Высота в пикселях (макс. 450)
}

// Marker представляет маркер на карте
type Marker struct {
    GeoPoint           // Встраиваем структуру GeoPoint (наследует Lat и Lon)
    Color    string    `json:"color"` // Цвет в HEX (без #), например 0xFF0000 - красный
    Label    string `json:"label"` // Текст метки (необязательно)
}

// Пример цветов для маркеров
const (
    MarkerColorRed    = 0xFF0000
    MarkerColorGreen  = 0x00FF00
    MarkerColorBlue   = 0x0000FF
    MarkerColorYellow = 0xFFFF00
)