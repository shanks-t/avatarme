package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

type AppStatus struct {
	Status      string
	Environment string
	Version     string
	Gravitar    string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "App environment (development|prodction")
	flag.Parse()

	fmt.Println("running...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := getGravitarURL("  shanks@dude.com  ")

		currentStatus := AppStatus{
			Status:      "available",
			Environment: cfg.env,
			Version:     version,
			Gravitar:    url,
		}

		js, err := json.MarshalIndent(currentStatus, "", "\t")
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil)
	if err != nil {
		log.Println(err)
	}
}

func getGravitarURL(email string) string {
	hash := createHashFromEmail(email)
	url := "https://www.gravatar.com/avatar/" + hash + "?d=identicon"

	return url
}

func createHashFromEmail(email string) string {
	formattedEmail := trimAndFormatEmailString(email)

	hash := md5.Sum([]byte(formattedEmail))

	hashToString := (hex.EncodeToString(hash[:]))

	return hashToString
}

func trimAndFormatEmailString(email string) string {
	trimmed := strings.TrimSpace(email)
	trimmedAndLowerCase := strings.ToLower(trimmed)

	return trimmedAndLowerCase
}
