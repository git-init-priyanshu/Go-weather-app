package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type envVariable struct {
	OpenWeatherMapApiKey string `json:"OpenWeaterApiKey"`
}
type weaterhData struct {
	Name string `json:name`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadEnv(filename string) (envVariable, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return envVariable{}, err
	}

	var env envVariable
	fmt.Print("bytes", bytes)
	err = json.Unmarshal(bytes, &env)
	if err != nil {
    fmt.Print("error while Unmarshalling")
		return envVariable{}, err
	}
	fmt.Print("env", env)
	return env, nil
}
func getWeather(city string) (weaterhData, error) {
	env, err := loadEnv("envFile")
	if err != nil {
		return weaterhData{}, err
	}

	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + env.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weaterhData{}, err
	}

	defer res.Body.Close()

	var data weaterhData

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return weaterhData{}, err
	}

	return data, nil

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		env, err := loadEnv("./envFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Print("Api key=", env)
		w.Write([]byte("Hello from weather-api"))
	})

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := getWeather(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)

}
