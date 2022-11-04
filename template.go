package main

import (
	"encoding/json"
	"example.com/m/v2/model"
	"fmt"
	"github.com/robfig/cron"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	APPID                    = "wxe102b21a40cd4603"
	WeatherAPPID             = "46965941"
	APPSecret                = "5709b03f4d8b2c0f853c235ba52ec6d5"
	WeatherSecret            = "BH8efN1b"
	WeatherTemplateID        = "WUyeUtq5AHNqaHbxUiQ2ydFR9VgdSL1GStfRyBqtGfs" //天气模板ID，替换成自己的
	WeatherTemplateIDNoAlarm = "1tIQnLLNiALJXVJsBZBsRpp9XjwZdots3bHxws9osX0" //天气模板ID，替换成自己的
	WeatherVersion           = "v61"
)

type token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {
	weather()
	select {}
	spec1 := "0 0 8 * * *" // 每天早晨8:00
	c := cron.New()
	err := c.AddFunc(spec1, weather)
	if err != nil {
		return
	}
	c.Start()
	fmt.Println("开启定时任务")
	select {}
	//weather()
}

//发送天气预报
func weather() {
	accessToken := getAccessToken()
	if accessToken == "" {
		return
	}

	fList := getFList(accessToken)
	if fList == nil {
		return
	}

	var cityId, city string
	for _, v := range fList {
		switch v.Str {
		case "oeZ6P5kyGsLKn3sIGRVfpb8oT4mg":
			city = "青岛"
		case "oeZ6P5jvFNh2y_h_2UcaoTXBaC2o":
			city = "西安"
		case "oQwrq5xAq1dWMJ5vg55MqL7Q9hj0":
			//cityId = "101270101"
			cityId = "101270106"
			//city = "成都"
			city = "双流"
			fallthrough
		default:
			go sendWeather(accessToken, cityId, city, v.Str)
		}
	}
	fmt.Println("weather is ok")
}

//获取微信accesstoken
func getAccessToken() string {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", APPID, APPSecret)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取微信token失败", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("微信token读取失败", err)
		return ""
	}

	token := token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("微信token解析json失败", err)
		return ""
	}

	return token.AccessToken
}

//获取关注者列表
func getFList(accessToken string) []gjson.Result {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + accessToken + "&next_openid="
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取关注列表失败", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return nil
	}
	fList := gjson.Get(string(body), "data.openid").Array()
	return fList
}

//发送模板消息
func templatePost(accessToken string, reqData string, fxUrl string, templateId string, openid string) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + accessToken

	reqBody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateId + "\", \"url\":\"" + fxUrl + "\", \"data\": " + reqData + "}"

	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

//获取天气
func getWeather(cityId, city string) (data model.WeatherTemplateData) {
	url := fmt.Sprintf("https://v0.yiketianqi.com/api?unescape=1&version=%s&appid=%s&appsecret=%s&cityid=%s&city=%s", WeatherVersion, WeatherAPPID, WeatherSecret, cityId, city)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取天气失败", err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return
	}

	weatherData := model.WeatherApiData{}
	err = json.Unmarshal(body, &weatherData)

	data = model.WeatherTemplateData{
		City:      weatherData.City,
		Weak:      weatherData.Week,
		Wea:       weatherData.Wea,
		Tem:       weatherData.Tem,
		Tem1:      weatherData.Tem1,
		Tem2:      weatherData.Tem2,
		Humidity:  weatherData.Humidity,
		AirLevel:  weatherData.AirLevel,
		AirTips:   weatherData.AirTips,
		Alarm:     weatherData.Alarm,
		Waichu:    weatherData.Aqi.Waichu,
		Kaichuang: weatherData.Aqi.Kaichuang,
	}

	return data
}

//发送天气
func sendWeather(accessToken, cityId, city, openid string) {
	weatherData := getWeather(cityId, city)
	if len(weatherData.City) == 0 {
		return
	}
	weatherTemplateId := WeatherTemplateIDNoAlarm
	//,"color":"#696969"

	tempData := model.WXWeatherTemplateNoAlarmData{
		City:      model.CreateWXTemplateData("城市：" + weatherData.City),
		Weak:      model.CreateWXTemplateData("星期：" + weatherData.Weak),
		Wea:       model.CreateWXTemplateData("天气情况：" + weatherData.Wea),
		Tem:       model.CreateWXTemplateDataSetColor("当前温度："+weatherData.Tem+"℃", getTemperatureColor(weatherData.Tem)),
		Tem1:      model.CreateWXTemplateDataSetColor("最高温度："+weatherData.Tem1+"℃", getTemperatureColor(weatherData.Tem1)),
		Tem2:      model.CreateWXTemplateDataSetColor("最低温度："+weatherData.Tem2+"℃", getTemperatureColor(weatherData.Tem2)),
		Humidity:  model.CreateWXTemplateData("湿度：" + weatherData.Humidity),
		AirLevel:  model.CreateWXTemplateData("空气质量等级：" + weatherData.AirLevel),
		AirTips:   model.CreateWXTemplateData("今日空气提示：" + weatherData.AirTips),
		Waichu:    model.CreateWXTemplateData("外出适宜：" + weatherData.Waichu),
		Kaichuang: model.CreateWXTemplateData("开窗适宜：" + weatherData.Kaichuang),
	}
	reqData, _ := json.Marshal(tempData)

	if len(weatherData.Alarm.AlarmType) != 0 {
		color := getAlarmColor(weatherData.Alarm.AlarmContent)
		tempHaveAlarmData := model.WXWeatherTemplateData{
			City:         model.CreateWXTemplateData("城市：" + weatherData.City),
			Weak:         model.CreateWXTemplateData("星期：" + weatherData.Weak),
			Wea:          model.CreateWXTemplateData("天气情况：" + weatherData.Wea),
			Tem:          model.CreateWXTemplateDataSetColor("当前温度："+weatherData.Tem+"℃", getTemperatureColor(weatherData.Tem)),
			Tem1:         model.CreateWXTemplateDataSetColor("最高温度："+weatherData.Tem1+"℃", getTemperatureColor(weatherData.Tem1)),
			Tem2:         model.CreateWXTemplateDataSetColor("最低温度："+weatherData.Tem2+"℃", getTemperatureColor(weatherData.Tem2)),
			Humidity:     model.CreateWXTemplateData("湿度：" + weatherData.Humidity),
			AirLevel:     model.CreateWXTemplateData("空气质量等级：" + weatherData.AirLevel),
			AirTips:      model.CreateWXTemplateData("今日空气提示：" + weatherData.AirTips),
			AlarmContent: model.CreateWXTemplateDataSetColor("预警类型："+weatherData.Alarm.AlarmContent, color),
			AlarmLevel:   model.CreateWXTemplateDataSetColor("预警等级："+weatherData.Alarm.AlarmLevel, color),
			AlarmType:    model.CreateWXTemplateDataSetColor("预警描述："+weatherData.Alarm.AlarmType, color),
			Waichu:       model.CreateWXTemplateData("外出适宜：" + weatherData.Waichu),
			Kaichuang:    model.CreateWXTemplateData("开窗适宜：" + weatherData.Kaichuang),
		}
		reqData, _ = json.Marshal(tempHaveAlarmData)
		weatherTemplateId = WeatherTemplateID
	}
	templatePost(accessToken, string(reqData), "https://www.qweather.com/", weatherTemplateId, openid)
}

func getTemperatureColor(temperatureStr string) (color string) {
	temperature, _ := strconv.Atoi(temperatureStr)

	//内陆上酷热标准是气温等于或超过38度，炎热时气温在35-37度，闷热是气温在28-35度 温暖是气温在10-27度，凉爽是气温在0-9度， 寒冷是气温在-9到-1度，严寒是气温等于或低于-10度。
	if temperature > 38 {
		color = "#fe0104"
	} else if 35 <= temperature && temperature <= 37 {
		color = "#ff6205"
	} else if 28 <= temperature && temperature < 35 {
		color = "#fd9907"
	} else if 10 <= temperature && temperature <= 27 {
		color = "#8AC76F"
	} else if 0 <= temperature && temperature <= 9 {
		color = "#00ab65"
	} else if -9 <= temperature && temperature <= -1 {
		color = "#05929c"
	} else if temperature <= -10 {
		color = "#0131c3"
	} else {
		color = model.DefaultColor
	}
	return
}

func getAlarmColor(Alarm string) (color string) {

	switch Alarm {
	case "红色":
		color = "#cb352f"
	case "橙色":
		color = "#ff9900"
	case "黄色":
		color = "#ffff00"
	case "蓝色":
		color = "#3265ff"
	default:
		color = model.DefaultColor
	}
	return
}
