package main

import (
	"aagmon/rec-rest-server/handlers"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	configFile := os.Args[1]
	serverConfig, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	log.Printf("config loaded: %+v", serverConfig)
	requestHandler := GetRecHandler(serverConfig)

	e := gin.New()
	e.GET("/recommend", requestHandler.Recommend)
	e.GET("/similar", requestHandler.GetMostSimilarURL)

	listenerAddress := fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port)
	log.Fatal(e.Run(listenerAddress))
}

func GetRecHandler(serverConfig *ServerConfig) *handlers.RcmndHandler {

	vec_file := serverConfig.EmbeddingFile
	item_file := serverConfig.ItemsFile
	vec_size := serverConfig.EmbeddingSize

	return handlers.NewRcmndHandler(vec_file, item_file, vec_size)

}
