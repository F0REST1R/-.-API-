package utils

import (
    "math"
    "prog/internal/models"
)

//реализует алгоритм пространственной кластеризации географических точек
func ClusterPoints(points []models.GeoPoint, radius float64) []models.Cluster {
    var clusters []models.Cluster
    visited := make(map[int]bool)

    for i, point := range points {
        if visited[i] {
            continue
        }

        cluster := models.Cluster{Center: point, Points: []models.GeoPoint{point}}
        visited[i] = true

        for j := i + 1; j < len(points); j++ {
            if !visited[j] && distance(point, points[j]) <= radius {
                cluster.Points = append(cluster.Points, points[j])
                visited[j] = true
            }
        }

        if len(cluster.Points) > 1 {
            cluster.Center = calculateCentroid(cluster.Points)
        }

        clusters = append(clusters, cluster)
    }

    return clusters
}

//вычисляет евклидово расстояние между точками
func distance(p1, p2 models.GeoPoint) float64 {
    return math.Sqrt(math.Pow(p1.Lat-p2.Lat, 2) + math.Pow(p1.Lon-p2.Lon, 2))
}

//вычисляет географический центр масс группы точек
func calculateCentroid(points []models.GeoPoint) models.GeoPoint {
    var sumLat, sumLon float64
    for _, p := range points {
        sumLat += p.Lat
        sumLon += p.Lon
    }
    return models.GeoPoint{
        Lat: sumLat / float64(len(points)),
        Lon: sumLon / float64(len(points)),
    }
}