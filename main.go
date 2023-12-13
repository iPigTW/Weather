package main

import (
	"io"
	"net/http"
	"os"

	"github.com/valyala/fastjson"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}
func main() {
	key := os.Getenv("API_KEY")
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
            "title": "Home",
        })
	})
	app.Post("/", func(c *fiber.Ctx) error {
		var err error
		var data []byte
		var body struct {
			City string
		}
		c.BodyParser(&body)
		if data, err = getWeather(body.City, key); err != nil {
			return err
		}
		json, err := fastjson.Parse(string(data))
		weathers:=json.GetArray("weather")
		weather := weathers[0]
		main := weather.GetStringBytes("main")
		return c.Render("index", fiber.Map{
            "Weather": string(main),
			"City": body.City,
			"Description": string(weather.GetStringBytes("description")),
        })
	})
	app.Listen(":3000")
}
func getWeather(city string, key string) ([]byte, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + key + "&units=metric"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
