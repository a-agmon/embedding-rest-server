package model

import (
	"errors"
	"fmt"
	"log"
)

type EmbeddingHandler struct {
	Embedding   *Embeddings
	VectorModel *VectorModel
}

func NewEmbeddingHandler(factorsFile string, itemsFile string, factorsSize int) *EmbeddingHandler {

	e := NewEmbedding()
	e.LoadVectors(factorsFile, factorsSize)
	e.LoadItems(itemsFile)

	confidence := 40.0
	reg := 0.001
	m, err := NewVectorModel(e.Factors, confidence, reg)
	if err != nil {
		log.Fatalf("Error creating embedding handler %v", err)
	}

	return &EmbeddingHandler{
		Embedding:   e,
		VectorModel: m,
	}
}

func (handler *EmbeddingHandler) Recommend(old_items []string) ([]string, []string, error) {

	itemIds := make(map[int]bool)
	itemsNotfound := make([]string, 0)
	for index, item := range old_items {

		if itemId, ok := handler.Embedding.Items2ID[item]; ok {
			itemIds[itemId] = true
		} else {
			s := fmt.Sprintf("%d (%s)", index, item)
			itemsNotfound = append(itemsNotfound, s)
		}
	}
	if len(itemIds) < 1 {
		return nil, itemsNotfound, errors.New("Non of the seen items was recognized")
	}

	recItems, err := handler.VectorModel.Recommend(itemIds, 5)
	if err != nil {
		return nil, itemsNotfound, errors.New("Error creating reccomendation: " + err.Error())
	}
	strRecItems := make([]string, len(recItems))
	for index, recItem := range recItems {
		recItemID := recItem.DocumentID
		recItem, _ := handler.Embedding.GetItemNameByID(recItemID)
		strRecItems[index] = recItem
	}
	return strRecItems, itemsNotfound, nil
}
