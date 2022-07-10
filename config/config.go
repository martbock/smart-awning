package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type gpioConfig struct {
	PinExtend  int
	PinRetract int
	PinStop    int
}

type httpConfig struct {
	Username string
	Password string
}

var (
	GPIO gpioConfig
	HTTP httpConfig
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err)
	}

	pinExtend, _ := strconv.Atoi(os.Getenv("PIN_EXTEND"))
	pinRetract, _ := strconv.Atoi(os.Getenv("PIN_RETRACT"))
	pinStop, _ := strconv.Atoi(os.Getenv("PIN_STOP"))

	GPIO = gpioConfig{
		PinExtend:  pinExtend,
		PinRetract: pinRetract,
		PinStop:    pinStop,
	}

	HTTP = httpConfig{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	log.Println("Loaded .env file")
	log.Println("GPIO config: ", GPIO)
	log.Println("HTTP username: ", HTTP.Username)
}
