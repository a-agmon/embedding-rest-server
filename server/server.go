package main

import (
	"aagmon/rec-rest-server/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	e := gin.New()
	e.GET("/recommend", GetRecHandler().Recommend)
	e.GET("/mostsimilar", GetRecHandler().GetMostSimilar)
	e.GET("/", func(ctx *gin.Context) {
		log.Println("Get Request:" + ctx.FullPath())
	})

	log.Fatal(e.Run(":9090"))
}

func GetRecHandler() *handlers.RcmndHandler {

	vec_file := "/Users/alonagmon/MyData/work/golang-projects/vectors_model/factors.csv"
	item_file := "/Users/alonagmon/MyData/work/golang-projects/vectors_model/artists.csv"
	vec_size := 129

	return handlers.NewRcmndHandler(vec_file, item_file, vec_size)

}
