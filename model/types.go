package model

const DefaultColor = "#696969"

type WeatherTemplateData struct {
	City      string `json:"city"`
	Weak      string `json:"weak"`
	Wea       string `json:"wea"`
	Tem       string `json:"tem"`
	Tem1      string `json:"tem1"`
	Tem2      string `json:"tem2"`
	Humidity  string `json:"humidity"`
	AirLevel  string `json:"air_level"`
	AirTips   string `json:"air_tips"`
	Alarm     Alarm  `json:"alarm"`
	Waichu    string `json:"waichu"`
	Kaichuang string `json:"kaichuang"`
}

type Alarm struct {
	AlarmContent string `json:"alarm_content"`
	AlarmLevel   string `json:"alarm_level"`
	AlarmType    string `json:"alarm_type"`
}

type WeatherApiData struct {
	Air      string `json:"air"`
	AirLevel string `json:"air_level"`
	AirPm25  string `json:"air_pm25"`
	AirTips  string `json:"air_tips"`
	Alarm    struct {
		AlarmContent string `json:"alarm_content"`
		AlarmLevel   string `json:"alarm_level"`
		AlarmType    string `json:"alarm_type"`
	} `json:"alarm"`
	Aqi struct {
		Air        string `json:"air"`
		AirLevel   string `json:"air_level"`
		AirTips    string `json:"air_tips"`
		City       string `json:"city"`
		CityEn     string `json:"cityEn"`
		Cityid     string `json:"cityid"`
		Co         string `json:"co"`
		CoDesc     string `json:"co_desc"`
		Country    string `json:"country"`
		CountryEn  string `json:"countryEn"`
		Jinghuaqi  string `json:"jinghuaqi"`
		Kaichuang  string `json:"kaichuang"`
		Kouzhao    string `json:"kouzhao"`
		No2        string `json:"no2"`
		No2Desc    string `json:"no2_desc"`
		O3         string `json:"o3"`
		O3Desc     string `json:"o3_desc"`
		Pm10       string `json:"pm10"`
		Pm10Desc   string `json:"pm10_desc"`
		Pm25       string `json:"pm25"`
		Pm25Desc   string `json:"pm25_desc"`
		So2        string `json:"so2"`
		So2Desc    string `json:"so2_desc"`
		UpdateTime string `json:"update_time"`
		Waichu     string `json:"waichu"`
		Yundong    string `json:"yundong"`
	} `json:"aqi"`
	City          string `json:"city"`
	CityEn        string `json:"cityEn"`
	Cityid        string `json:"cityid"`
	Country       string `json:"country"`
	CountryEn     string `json:"countryEn"`
	Date          string `json:"date"`
	Humidity      string `json:"humidity"`
	Pressure      string `json:"pressure"`
	Tem           string `json:"tem"`
	Tem1          string `json:"tem1"`
	Tem2          string `json:"tem2"`
	UpdateTime    string `json:"update_time"`
	Visibility    string `json:"visibility"`
	Wea           string `json:"wea"`
	WeaImg        string `json:"wea_img"`
	Week          string `json:"week"`
	Win           string `json:"win"`
	WinMeter      string `json:"win_meter"`
	WinSpeed      string `json:"win_speed"`
	WinSpeedDay   string `json:"win_speed_day"`
	WinSpeedNight string `json:"win_speed_night"`
}

type WXTemplateData struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

func CreateWXTemplateData(value string) WXTemplateData {
	return WXTemplateData{
		Value: value,
		Color: DefaultColor,
	}
}

func CreateWXTemplateDataSetColor(value, color string) WXTemplateData {
	return WXTemplateData{
		Value: value,
		Color: color,
	}
}

type WXWeatherTemplateData struct {
	City         WXTemplateData `json:"city"`
	Weak         WXTemplateData `json:"weak"`
	Wea          WXTemplateData `json:"wea"`
	Tem          WXTemplateData `json:"tem"`
	Tem1         WXTemplateData `json:"tem1"`
	Tem2         WXTemplateData `json:"tem2"`
	Humidity     WXTemplateData `json:"humidity"`
	AirLevel     WXTemplateData `json:"air_level"`
	AirTips      WXTemplateData `json:"air_tips"`
	AlarmContent WXTemplateData `json:"alarm_content"`
	AlarmLevel   WXTemplateData `json:"alarm_level"`
	AlarmType    WXTemplateData `json:"alarm_type"`
	Waichu       WXTemplateData `json:"waichu"`
	Kaichuang    WXTemplateData `json:"kaichuang"`
}

type WXWeatherTemplateNoAlarmData struct {
	City      WXTemplateData `json:"city"`
	Weak      WXTemplateData `json:"weak"`
	Wea       WXTemplateData `json:"wea"`
	Tem       WXTemplateData `json:"tem"`
	Tem1      WXTemplateData `json:"tem1"`
	Tem2      WXTemplateData `json:"tem2"`
	Humidity  WXTemplateData `json:"humidity"`
	AirLevel  WXTemplateData `json:"air_level"`
	AirTips   WXTemplateData `json:"air_tips"`
	Waichu    WXTemplateData `json:"waichu"`
	Kaichuang WXTemplateData `json:"kaichuang"`
}
