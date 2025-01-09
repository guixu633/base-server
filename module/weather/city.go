package weather

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Location 结构体定义
type Location struct {
	LocationID      string  `csv:"Location_ID"`
	LocationNameEN  string  `csv:"Location_Name_EN"`
	LocationNameZH  string  `csv:"Location_Name_ZH"`
	ISO31661        string  `csv:"ISO_3166_1"`
	CountryRegionEN string  `csv:"Country_Region_EN"`
	CountryRegionZH string  `csv:"Country_Region_ZH"`
	Adm1NameEN      string  `csv:"Adm1_Name_EN"`
	Adm1NameZH      string  `csv:"Adm1_Name_ZH"`
	Adm2NameEN      string  `csv:"Adm2_Name_EN"`
	Adm2NameZH      string  `csv:"Adm2_Name_ZH"`
	Timezone        string  `csv:"Timezone"`
	Latitude        float64 `csv:"Latitude"`
	Longitude       float64 `csv:"Longitude"`
	ADCode          string  `csv:"AD_code"`
}

// ParseCityList 解析CSV文件
func CityList(filename string) ([]Location, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 跳过版本信息行
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("读取版本信息失败: %v", err)
	}

	// 跳过标题行
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("读取标题行失败: %v", err)
	}

	var locations []Location

	// 逐行读取数据
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取数据行失败: %v", err)
		}

		// 将字符串转换为float64
		lat, err := strconv.ParseFloat(record[11], 64)
		if err != nil {
			return nil, fmt.Errorf("解析Latitude失败: %v", err)
		}
		lng, err := strconv.ParseFloat(record[12], 64)
		if err != nil {
			return nil, fmt.Errorf("解析Longitude失败: %v", err)
		}

		// 构建Location对象
		loc := Location{
			LocationID:      record[0],
			LocationNameEN:  record[1],
			LocationNameZH:  record[2],
			ISO31661:        record[3],
			CountryRegionEN: record[4],
			CountryRegionZH: record[5],
			Adm1NameEN:      record[6],
			Adm1NameZH:      record[7],
			Adm2NameEN:      record[8],
			Adm2NameZH:      record[9],
			Timezone:        record[10],
			Latitude:        lat,
			Longitude:       lng,
			ADCode:          record[13],
		}

		locations = append(locations, loc)
	}

	return locations, nil
}

func CityMap() map[string]Location {
	locations, err := CityList("China-City-List-latest.csv")
	if err != nil {
		fmt.Printf("解析CSV文件失败: %v", err)
	}
	cityMap := make(map[string]Location)
	for _, location := range locations {
		cityMap[location.LocationNameZH] = location
	}
	return cityMap
}
