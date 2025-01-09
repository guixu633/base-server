package weather

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {
	weather, err := GetWeather("泰兴")
	assert.NoError(t, err)
	fmt.Println(weather)
}

func TestGetWeatherString(t *testing.T) {
	weather, err := GetWeatherString("泰兴")
	assert.NoError(t, err)
	fmt.Println(weather)
}

func TestCityList(t *testing.T) {
	start := time.Now()
	locations, err := CityList("China-City-List-latest.csv")
	if err != nil {
		t.Fatalf("解析CSV文件失败: %v", err)
	}
	for _, location := range locations[:20] {
		fmt.Println(location.LocationNameZH)
	}
	fmt.Println(len(locations))
	fmt.Println(time.Since(start))
}

func TestCityMap(t *testing.T) {
	cityMap := CityMap()
	fmt.Println(cityMap["泰兴"])
}
