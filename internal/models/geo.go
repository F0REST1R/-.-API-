package models

// GeoPoint представляет географические координаты
type GeoPoint struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

// Address содержит полную информацию об адресе
type Address struct {
    Formatted    string `json:"formatted"`     // Краткий форматированный адрес
    FullAddress  string `json:"full_address"`  // Полный адрес от Яндекса
    Country      string `json:"country"`
    Region       string `json:"region"`        // Область/край
    City         string `json:"city"`
    District     string `json:"district"`      // Район города
    Street       string `json:"street"`
    HouseNumber  string `json:"house_number"`
    Building     string `json:"building"`     // Корпус/строение
    PostalCode   string `json:"postal_code"`
    Description  string `json:"description"`  // Название объекта (для Кремля и др.)
}