package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/tidwall/gjson"
	"gopkg.in/ini.v1"
)

/*
*
step:
1,read options of ini file
2,load region code
3,load region today weather info
4,load 'one line' info
5,assemble 'special day' result
5,assemble all info
6,push result to template message
end
*/

func main() {

	cfg, _ := ini.Load("config.ini")
	allOption := new(AllOption)
	_ = cfg.MapTo(allOption)

	allResult := new(AllResult)

	//weather info
	weatherInfos := GetWeatherInfos(allOption)
	allResult.WeatherInfos = weatherInfos

	//hitokoto info
	hitokotoInfo := getHitokotoInfo(allOption)
	allResult.Hitokoto = hitokotoInfo

	//count-down
	countDownInfos := getCountDownInfos(allOption)
	allResult.CountDowns = countDownInfos

	//count
	countInfos := getCountInfos(allOption)
	allResult.Counts = countInfos

	sendTempMessage(allOption, allResult)

}

func sendTempMessage(allOption *AllOption, allResult *AllResult) {

	dataMap := map[string]any{}
	dataMap["hitokoto"] = map[string]string{
		"value": fmt.Sprintf("%s(%s)", gjson.Get(allResult.Hitokoto, "hitokoto"), gjson.Get(allResult.Hitokoto, "from")),
		"color": randomHashColor(),
	}

	lang := carbon.NewLanguage()
	lang.SetLocale("zh-CN")

	c := carbon.SetLanguage(lang)
	now := c.Now(carbon.Shanghai)
	dataMap["now"] = map[string]string{
		"value": fmt.Sprintf("%s %s", now.ToDateString(carbon.Shanghai), now.ToWeekString(carbon.Shanghai)),
		"color": randomHashColor(),
	}

	for i := 0; i < len(allResult.WeatherInfos); i++ {

		m := map[string]string{
			"value": fmt.Sprintf(
				"%s：%s℃ %s %s %s级",
				allOption.Weather.Region[i],
				gjson.Get(allResult.WeatherInfos[i], "now.temp"),
				gjson.Get(allResult.WeatherInfos[i], "now.text"),
				gjson.Get(allResult.WeatherInfos[i], "now.windDir"),
				gjson.Get(allResult.WeatherInfos[i], "now.windScale"),
			),
			"color": randomHashColor(),
		}

		dataMap[fmt.Sprintf("weather%d", i+1)] = m
	}

	for i := 0; i < len(allResult.Counts); i++ {

		m := map[string]string{}
		m["value"] = allResult.Counts[i]
		m["color"] = randomHashColor()

		dataMap[fmt.Sprintf("count%d", i+1)] = m
	}

	for i := 0; i < len(allResult.CountDowns); i++ {

		m := map[string]string{}
		m["value"] = allResult.CountDowns[i]
		m["color"] = randomHashColor()

		dataMap[fmt.Sprintf("countDown%d", i+1)] = m
	}

	tempMap := map[string]any{
		"template_id": allOption.Wechat.TemplateId,
		"topcolor":    randomHashColor(),
		"data":        dataMap,
	}

	getAccessTokenUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", allOption.Wechat.AppId, allOption.Wechat.AppSecret)
	tokenStr := SendSimpleGet(getAccessTokenUrl)

	sendMagUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", gjson.Get(tokenStr, "access_token"))

	for i := 0; i < len(allOption.Wechat.User); i++ {

		tempMap["touser"] = allOption.Wechat.User[i]
		j, _ := json.Marshal(tempMap)
		SendSimplePost(sendMagUrl, string(j))

	}

}

func randomHashColor() string {
	codeArr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	hash := "#"

	for i := 0; i < 6; i++ {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(16)
		hash += codeArr[randNum]
	}

	return hash

}

func getCountInfos(allOption *AllOption) []string {

	cb := carbon.SetTimezone(carbon.Shanghai)
	today := cb.Now(carbon.Shanghai)

	var countInfos []string

	for i := 0; i < len(allOption.Day.Count); i++ {
		dc := allOption.Day.Count[i]
		dcc := cb.Parse(dc, carbon.Shanghai)

		day := today.DiffAbsInDays(dcc)

		countInfos = append(countInfos, fmt.Sprintf("%s%d天", allOption.Day.CountTitle[i], day))

	}

	return countInfos

}

func getCountDownInfos(allOption *AllOption) []string {

	cb := carbon.SetTimezone(carbon.Shanghai)
	today := cb.Now(carbon.Shanghai)

	var countDownInfos []string

	for i := 0; i < len(allOption.Day.CountDown); i++ {
		dc := allOption.Day.CountDown[i]
		dcc := cb.Parse(fmt.Sprintf("%s", dc), carbon.Shanghai)
		if dcc.IsPast() {
			dcc = dcc.SetYear(today.Year() + 1)
		}

		day := dcc.DiffAbsInDays(today)

		countDownInfos = append(countDownInfos, fmt.Sprintf("%s%d天", allOption.Day.CountDownTitle[i], day))

	}

	return countDownInfos

}

func getHitokotoInfo(allOption *AllOption) string {
	hitokotoUrl := "https://v1.hitokoto.cn"

	if allOption.Hitokoto.Types != nil {

		parseUrl, _ := url.Parse(hitokotoUrl)

		params, _ := url.ParseQuery(parseUrl.RawQuery)
		for i := 0; i < len(allOption.Hitokoto.Types); i++ {
			params.Add("c", allOption.Hitokoto.Types[i])
		}
		parseUrl.RawQuery = params.Encode()
		hitokotoUrl = parseUrl.String()
	}

	hitokotoResultStr := SendSimpleGet(hitokotoUrl)
	return hitokotoResultStr
}

func GetWeatherInfos(allOption *AllOption) []string {

	var weatherRegionInfos []string
	for i := 0; i < len(allOption.Weather.Region); i++ {
		var r = allOption.Weather.Region[i]
		weatherRegionCodeUrl := fmt.Sprintf("https://geoapi.qweather.com/v2/city/lookup?location=%s&key=%s", r, allOption.Weather.Key)
		weatherRegionCodeResultStr := SendSimpleGet(weatherRegionCodeUrl)
		code := gjson.Get(weatherRegionCodeResultStr, "location.0.id")

		weatherRegionInfoUrl := fmt.Sprintf("https://devapi.qweather.com/v7/weather/now?key=%s&location=%s", allOption.Weather.Key, code.String())
		weatherRegionInfo := SendSimpleGet(weatherRegionInfoUrl)
		weatherRegionInfos = append(weatherRegionInfos, weatherRegionInfo)
	}

	return weatherRegionInfos

}

// SendSimplePost run a simple post request
func SendSimplePost(url string, param string) string {

	rq, _ := http.NewRequest("POST", url, strings.NewReader(param))
	rp, _ := http.DefaultClient.Do(rq)
	defer rp.Body.Close()
	rpBody, _ := io.ReadAll(rp.Body)

	return string(rpBody)
}

// SendSimpleGet run a simple get request
func SendSimpleGet(getUrl string) string {

	parseUrl, _ := url.Parse(getUrl)

	params, _ := url.ParseQuery(parseUrl.RawQuery)
	parseUrl.RawQuery = params.Encode()

	fmt.Printf("\n" + parseUrl.String())

	rq, _ := http.NewRequest("GET", parseUrl.String(), nil)
	rq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	rp, _ := http.DefaultClient.Do(rq)
	defer rp.Body.Close()
	rpBody, _ := io.ReadAll(rp.Body)

	return string(rpBody)
}

type AllResult struct {
	Hitokoto     string
	Title        string
	WeatherInfos []string
	Counts       []string
	CountDowns   []string
}

// AllOption struct
type AllOption struct {
	Weather  Weather
	Wechat   Wechat
	Day      Day
	Hitokoto Hitokoto
}

type Hitokoto struct {
	Types []string `ini:"types"`
}

// Weather api option
type Weather struct {
	Key    string   `ini:"key"`
	Region []string `ini:"region"`
}

// Wechat api option
type Wechat struct {
	AppId      string   `ini:"app-id"`
	AppSecret  string   `ini:"app-secret"`
	TemplateId string   `ini:"template-id"`
	User       []string `ini:"user"`
}

// Day option
type Day struct {
	CountDown      []string `ini:"count-down"`
	CountDownTitle []string `ini:"count-down-title"`
	Count          []string `ini:"count"`
	CountTitle     []string `ini:"count-title"`
}
