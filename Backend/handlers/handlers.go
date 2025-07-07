package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Club struct {
	ClubName     string  `json:"club_name"`
	CityName     string  `json:"city_name"`
	Titles       int     `json:"titles"`
	AvgPlayerAge float64 `json:"avg_player_age"`
}

func SearchClubs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		field := c.Query("field")
		value := c.Query("value")

		if field == "" || value == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Field and value parameters are required"})
			return
		}

		query := `
            SELECT c.club_name, ci.city_name, c.titles, c.avg_player_age
            FROM clubs c
            JOIN cities ci ON c.city_id = ci.city_id
            WHERE `
		var params []interface{}

		switch field {
		case "club_name":
			query += "c.club_name ILIKE $1"
			params = append(params, "%"+value+"%")
		case "city_name":
			query += "ci.city_name ILIKE $1"
			params = append(params, "%"+value+"%")
		case "titles":
			titles, err := strconv.Atoi(value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Titles must be an integer"})
				return
			}
			query += "c.titles = $1"
			params = append(params, titles)
		case "avg_player_age":
			age, err := strconv.ParseFloat(value, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Avg_player_age must be a number"})
				return
			}
			query += "c.avg_player_age BETWEEN $1 - 0.1 AND $1 + 0.1"
			params = append(params, age)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field"})
			return
		}

		rows, err := db.Query(query, params...)
		if err != nil {
			log.Printf("Database query error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		}
		defer rows.Close()

		var clubs []Club
		for rows.Next() {
			var club Club
			if err := rows.Scan(&club.ClubName, &club.CityName, &club.Titles, &club.AvgPlayerAge); err != nil {
				log.Printf("Row scan error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
				return
			}
			clubs = append(clubs, club)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Rows iteration error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating rows"})
			return
		}

		c.JSON(http.StatusOK, clubs)
	}
}
