package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Hotel struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	City              string  `json:"city"`
	Capacity          int     `json:"capacity"`
	StandardRoomPrice float64 `json:"standard_room_price"`
}

func main() {
	// Подключение к БД
	connStr := "user=hoteluser password=hotelpassword dbname=hoteldb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database!")

	// Инициализация Gin
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// Маршрут для поиска
	r.GET("/hotels", func(c *gin.Context) {
		field := c.Query("field")
		value := c.Query("value")

		query := `
            SELECT h.id, h.name, c.name as city, h.capacity, h.standard_room_price
            FROM hotels h
            JOIN hotel_cities hc ON h.id = hc.hotel_id
            JOIN cities c ON hc.city_id = c.id
            WHERE `

		switch field {
		case "name":
			query += "h.name ILIKE $1"
		case "city":
			query += "c.name ILIKE $1"
		case "capacity":
			query += "h.capacity::text ILIKE $1"
		case "price":
			query += "h.standard_room_price::text ILIKE $1"
		default:
			query += "h.name ILIKE $1"
		}

		rows, err := db.Query(query, "%"+value+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var hotels []Hotel
		for rows.Next() {
			var h Hotel
			if err := rows.Scan(&h.ID, &h.Name, &h.City, &h.Capacity, &h.StandardRoomPrice); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			hotels = append(hotels, h)
		}

		c.JSON(http.StatusOK, hotels)
	})

	// Запуск сервера
	r.Run(":8080")
}
