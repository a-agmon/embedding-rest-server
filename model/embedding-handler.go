package model

import (
	utl "aagmon/rec-rest-server/utils"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"gonum.org/v1/gonum/mat"
	//"github.com/khaibin/go-cosinesimilarity"
)

type EmbeddingHandler struct {
	Embedding   *Embeddings
	VectorModel *VectorModel
	ItemsMatrix [][]float64
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

	// create the matrix of all the items
	items_matrix := make([][]float64, len(e.Factors))
	for key, vector := range e.Factors {
		items_matrix[key] = vector
	}

	log.Printf("Embedding handler created")

	return &EmbeddingHandler{
		Embedding:   e,
		VectorModel: m,
		ItemsMatrix: items_matrix,
	}
}

func (handler *EmbeddingHandler) GetMostSimilar(item string, topk int) ([]string, error) {

	itemID, ok := handler.Embedding.Items2ID[item]
	if !ok {
		return nil, errors.New("item not found")
	}
	itemVec, ok := handler.Embedding.Factors[itemID]
	if !ok {
		return nil, errors.New("the item's id was found but its vector was not - something is very wrong")
	}
	cosineScores := make([]float64, len(handler.ItemsMatrix))
	vecSize := len(itemVec)

	time_start := time.Now()
	for i, matrixVec := range handler.ItemsMatrix {
		vectorX := mat.NewDense(vecSize, 1, matrixVec)
		vectorY := mat.NewDense(vecSize, 1, itemVec)
		result := Distance(vectorX, vectorY)
		cosineScores[i] = result
	}
	elapsed := time.Since(time_start)
	log.Printf("Cosine similarity matrix computed in %v", elapsed)

	sorted_scores := utl.Sort(sort.Float64Slice(cosineScores))
	top_items := make([]string, len(sorted_scores[:topk]))
	for i, itemID := range sorted_scores[:topk] {
		top_items[i], _ = handler.Embedding.GetItemNameByID(itemID)
	}

	return top_items, nil
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
		return nil, itemsNotfound, errors.New("non of the seen items was recognized")
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

// Dot computes dot value of vectorX and vectorY.
func Dot(vectorX *mat.Dense, vectorY *mat.Dense) float64 {
	subVector := new(mat.Dense)
	subVector.MulElem(vectorX, vectorY)
	result := mat.Sum(subVector)

	return result
}

// Distance computes Cosine distance.
// It will return distance which represented as 1-cos() (ranged from 0 to 2).
func Distance(vectorX *mat.Dense, vectorY *mat.Dense) float64 {
	dotXY := Dot(vectorX, vectorY)
	lengthX := math.Sqrt(Dot(vectorX, vectorX))
	lengthY := math.Sqrt(Dot(vectorY, vectorY))

	cos := dotXY / (lengthX * lengthY)

	return 1 - cos
}
