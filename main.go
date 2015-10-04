package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	config      *oauth2.Config
	globalFlags struct {
		ClientSecretFile string
		Port             int
	}
)

func init() {
	flag.StringVar(&globalFlags.ClientSecretFile, "secretFile",
		"client_secret.json", "Path to the Google Drive client secret file")
	flag.IntVar(&globalFlags.Port, "port", 80, "The Port to listen on")
	flag.Parse()
}

func main() {
	secret, err := ioutil.ReadFile(globalFlags.ClientSecretFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err = google.ConfigFromJSON(secret, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	router := httprouter.New()
	router.DELETE("/delete/*filepath", Delete)
	router.GET("/auth/new", NewUser)
	router.GET("/browse/*filepath", Browse)
	router.GET("/read/*filepath", Read)
	router.POST("/auth/validate", Validate)
	router.POST("/add/*filepath", Add)
	router.GET("/publish/*filepath", Publish)

	log.Println("Listening on port " + strconv.Itoa(globalFlags.Port))
	http.ListenAndServe(":"+strconv.Itoa(globalFlags.Port), router)
}
