package handlers

import (
	vm "aagmon/rec-rest-server/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RcmndHandler struct {
	EmbeddingHandler *vm.EmbeddingHandler
}

type RcmndReq struct {
	History []string `json:"history"`
}

type RcmndRes struct {
	ViewedItems      []string `json:"viewed"`
	RecommendedItems []string `json:"recommended"`
	MissingItems     []string `json:"missing"`
}

type SimilarityReq struct {
	SimilarTo string `json:"similarto"`
	TopK      int    `json:"topk"`
}
type SimilarityRes struct {
	Original string   `json:"original"`
	Similar  []string `json:"similar"`
}

func NewRcmndHandler(factors_file string, items_file string, size int) *RcmndHandler {
	return &RcmndHandler{
		EmbeddingHandler: vm.NewEmbeddingHandler(factors_file, items_file, size),
	}
}

func (rec *RcmndHandler) GetMostSimilar(c *gin.Context) {

	time_start := time.Now()

	var request SimilarityReq
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	similar, err := rec.EmbeddingHandler.GetMostSimilar(request.SimilarTo, request.TopK)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := SimilarityRes{
		Original: request.SimilarTo,
		Similar:  similar,
	}
	elapsed := time.Since(time_start)
	log.Printf("GetMostSimilar took %s", elapsed)
	c.JSON(http.StatusOK, response)
}

func (rec *RcmndHandler) Recommend(c *gin.Context) {

	time_start := time.Now()

	var request RcmndReq
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	maxProgs := 5
	viewed_items := request.History
	if len(viewed_items) > maxProgs {
		viewed_items = viewed_items[len(viewed_items)-maxProgs:]
	}
	recommended, missing, err := rec.EmbeddingHandler.Recommend(viewed_items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	response := RcmndRes{
		ViewedItems:      viewed_items,
		RecommendedItems: recommended,
		MissingItems:     missing,
	}
	elapsed := time.Since(time_start)
	log.Printf("Recommend took %s", elapsed)
	c.JSON(http.StatusOK, response)

}
