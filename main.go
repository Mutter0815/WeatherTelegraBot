package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tucnak/telebot"
)

type WeatherResponse struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}
type GeocodingResponse struct {
	Name string  `json: "name"`
	Lat  float64 `json: "lat"`
	Lon  float64 `json: "lon"`
}

func getCoordinates(city, apiKey string) (float64, float64, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("не удалось получить координаты %s", resp.Status)
	}
	var geoResp []GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return 0, 0, err
	}
	if len(geoResp) == 0 {
		return 0, 0, fmt.Errorf("Город не найден")
	}
	return geoResp[0].Lat, geoResp[0].Lon, nil

}

func getWeather(lat, lon float64, apiKey string) (string, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%.2f&lon=%.2f&appid=%s", lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Не удалось получить погоду %s", resp.Status)
	}
	var weatherResp WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return "", err
	}
	weatherInfo := fmt.Sprintf("Погода в городе \nТемпература: %.1f \nОщущается как: %.1f \nВлажность: %d", weatherResp.Main.Temp-273.15, weatherResp.Main.FeelsLike-273.15, weatherResp.Main.Humidity)
	return weatherInfo, nil
}

func main() {

	weatherApiKey := "749caf75f9c260069ff03f13da598809"

	fmt.Println("hello")
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  "1514473225:AAHV3K9KFxusaF-Bs6c1Ai8EjI7aIVX1GZU",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}
	bot.Handle("/start", func(m *telebot.Message) {
		mes := fmt.Sprintf("Hello %d", bot.Me.Username)
		bot.Send(m.Sender, mes)
	})
	bot.Handle("/weather", func(m *telebot.Message) {
		city := m.Payload
		if city == "" {
			bot.Send(m.Sender, "Пожалуйста введите город. /weather город")
			return
		}
		lan, lot, err := getCoordinates(city, weatherApiKey)
		if err != nil {
			bot.Send(m.Sender, "Ошибка при получении координат")
		}
		weather, err := getWeather(lan, lot, weatherApiKey)
		if err != nil {
			bot.Send(m.Sender, "Ошибка при получении погоды")
		}
		bot.Send(m.Sender, city+"\n"+weather)

	})
	bot.Start()
}
