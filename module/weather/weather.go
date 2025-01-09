package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WeatherResponse 表示完整的天气响应结构
type WeatherResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Now        Now    `json:"now"`
	Refer      Refer  `json:"refer"`
}

// Now 表示当前天气数据
type Now struct {
	ObsTime   string `json:"obsTime"`
	Temp      string `json:"temp"`
	FeelsLike string `json:"feelsLike"`
	Icon      string `json:"icon"`
	Text      string `json:"text"`
	Wind360   string `json:"wind360"`
	WindDir   string `json:"windDir"`
	WindScale string `json:"windScale"`
	WindSpeed string `json:"windSpeed"`
	Humidity  string `json:"humidity"`
	Precip    string `json:"precip"`
	Pressure  string `json:"pressure"`
	Vis       string `json:"vis"`
	Cloud     string `json:"cloud"`
	Dew       string `json:"dew"`
}

// Refer 表示数据引用信息
type Refer struct {
	Sources []string `json:"sources"`
	License []string `json:"license"`
}

const (
	apiKey     = "1d080d2f5d424d9caea63890f7d5acd6"
	weatherURL = "https://devapi.qweather.com/v7/weather/now"
)

var cityMap = CityMap()

func GetWeather(cityName string) (*WeatherResponse, error) {
	cityID, ok := cityMap[cityName]
	if !ok {
		return nil, fmt.Errorf("城市名称不存在")
	}
	return GetWeatherByID(cityID.LocationID)
}

func GetWeatherString(cityName string) (string, error) {
	weather, err := GetWeather(cityName)
	if err != nil {
		return "", err
	}
	json, err := json.Marshal(weather.Now)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func GetWeatherSchema() string {
	return `{
		"type": "object",
		"properties": {
			"temp": {"type": "string"},
		}
	}`
}

// GetWeather 获取指定城市的天气信息
func GetWeatherByID(cityID string) (*WeatherResponse, error) {
	// 构建请求URL
	url := fmt.Sprintf("%s?location=%s&key=%s", weatherURL, cityID, apiKey)

	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON响应
	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return &weather, nil
}
