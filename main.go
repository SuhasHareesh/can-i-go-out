package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Colors struct {
    Reset   string
    Red     string
    Green   string
    Yellow  string
    Blue    string
    Magenta string
    Gold    string
}

func NewColors() Colors {
    return Colors{
        Reset:   "\033[0m",
        Red:     "\033[31m",
        Green:   "\033[1m\033[32m",
        Yellow:  "\033[33m",
        Blue:    "\033[34m",
        Magenta: "\033[35m",
        Gold:    "\033[1m\033[38;5;214m",
    }
}

type Weather struct {
    Location struct {
        Name string `json:"name"`
        Region string `json:"region"`
        Country string `json:"country"`
        LocalTime int64 `json:"localtime_epoch"`
        TimeZone string `json:"tz_id"`
    } `json:"location"`
    Current struct {
        TempC float64 `json:"temp_c"`
        FeelC float64 `json:"feelslike_c"`
        TempF float64 `json:"temp_f"`
        FeelF float64 `json:"feelslike_f"`
        Condition struct {
            Text string `json:"text"`
        } `json:"condition"`
    } `json:"current"`
    Forecast struct {
        ForecastDay []struct {
            Date string `json:"date"`
            Hour []struct {
                TimeOfDay int64 `json:"time_epoch"`
                TempC float64 `json:"temp_c"`
                TempF float64 `json:"temp_f"`
                RainChance float64 `json:"chance_of_rain"`
                Condition struct {
                    Text string `json:"text"`
                } `json:"condition"`
            }`json:"hour"`
        }`json:"forecastday"`
    } `json:"forecast"`
}

func callWeatherAPI(loc string, day int) {
    const API_KEY = "fb74f2ccc2a64b548b1163707250403"
    const API_URL = "http://api.weatherapi.com/v1/forecast.json?"

    res, err := http.Get(fmt.Sprint(API_URL + "key=" + API_KEY + "&q=" + loc + "&days=3" + "&aqi=no&alerts=no"))
    if err != nil {
        panic(err)
    }

    defer res.Body.Close()

    if res.StatusCode != 200 {
        panic("Weather API not Available")
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        panic(err)
    }

    var weather Weather
    err = json.Unmarshal(body, &weather)
    if err != nil {
        panic(err)
    }

    colors := NewColors()
    current := weather.Current
    location := weather.Location
    forecastHours := weather.Forecast.ForecastDay[day].Hour

    fmt.Printf(colors.Gold + "%s, %s :" + colors.Reset, location.Name, location.Region)
    fmt.Printf(" %.0fF (Feels like %.0fF), %.0fC (Feels like %.0fC), %s\n", current.TempF, current.FeelF, current.TempC, current.FeelC, current.Condition.Text)
    fmt.Printf("%s%s%s\n", colors.Green, weather.Forecast.ForecastDay[day].Date, colors.Reset);

    // localTime := time.Unix(weather.Location.LocalTime, 0)

    loctime, err := time.LoadLocation(location.TimeZone)
    if err != nil {
        fmt.Println("⚠️ Warning: Could not load timezone, using UTC.")
        loctime = time.UTC
    }

    for _, hour := range forecastHours {
        utcTime := time.Unix(hour.TimeOfDay, 0)
        localTime := utcTime.In(loctime)

        now := time.Now().In(loctime)

        if localTime.Before(now) { // ✅ Skip past times properly
            continue
        }

        if hour.RainChance > 40 {
            fmt.Printf(
                "%s%s - %.0fF, %.0fC, %.0f%%, %s%s\n",
                colors.Red,
                localTime.Format("15:00"),
                hour.TempF,
                hour.TempC,
                hour.RainChance,
                hour.Condition.Text,
                colors.Reset,
                )
        } else {
            fmt.Printf(
                "%s%s%s - %.0fF, %.0fC, %.0f%%, %s%s%s\n",
                colors.Magenta,
                localTime.Format("15:00"),
                colors.Reset,
                hour.TempF,
                hour.TempC,
                hour.RainChance,
                colors.Yellow,
                hour.Condition.Text,
                colors.Reset,
                )
        }

    }

}

func main() {
    location := flag.String("l", "Gainesville", "Location for weather lookup")
	days := flag.Int("d", 0, "Number of forecast days (1-2)")

	flag.Parse()

	// Call weather API with user inputs
	callWeatherAPI(*location, *days)
}
