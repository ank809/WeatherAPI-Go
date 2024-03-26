package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	router.GET("/details/:city", GetweatherDetails)
	router.Run(":8081")

}

type LongLat struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

func GetweatherDetails(c *gin.Context) {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err)
	}
	apiKey := os.Getenv("API_KEY")

	city := c.Param("city")
	fmt.Println(city)
	apiUrl := "http://api.openweathermap.org/geo/1.0/direct?q=" + city + "&limit=5&appid=" + apiKey
	fmt.Println(apiUrl)
	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	var longlat []LongLat
	var longitude string
	var latitude string
	if response.StatusCode == 200 {
		if err := json.NewDecoder(response.Body).Decode(&longlat); err != nil {
			fmt.Println(err)
			return
		}
		a := longlat[0].Longitude
		b := longlat[0].Latitude
		longitude = strconv.FormatFloat(a, 'f', -1, 64)
		latitude = strconv.FormatFloat(b, 'f', -1, 64)
		fmt.Println(longlat[0].Latitude, longlat[0].Longitude)
	}
	webapiURL := "https://api.openweathermap.org/data/2.5/weather?lat=" + latitude + "&lon=" + longitude + "&appid=" + apiKey

	fmt.Println(webapiURL)
	res, err := http.Get(webapiURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	var details WeatherData
	if res.StatusCode == 200 {
		json.NewDecoder(res.Body).Decode(&details)
		c.JSON(http.StatusOK, gin.H{"message": details.Main.Temp,
			"rain": details.Rain.OneHour})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": res.Body})
	}

}
