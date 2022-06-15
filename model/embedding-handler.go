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

	num_vectors := len(handler.Embedding.Factors)
	vec_size := len(handler.Embedding.Factors[0])

	v := make([]float64, num_vectors*vec_size)
	A := mat.NewDense(num_vectors, vec_size, v)
	for i, vec := range handler.Embedding.Factors {
		A.SetRow(i, vec)
	}

	time_start := time.Now()
	b := mat.NewVecDense(vec_size, itemVec)
	results_arr := make([]float64, num_vectors)
	results_vec := mat.NewVecDense(num_vectors, results_arr)
	// this was the method, I was looking for.
	results_vec.MulVec(A, b)
	elapsed := time.Since(time_start)
	log.Printf("Cosine similarity matrix computed in %v", elapsed)
	//fmt.Printf("Dist:%v\n", results_vec.AtVec(0))

	top_items_k := num_vectors - topk
	cosine_scores := make([]float64, num_vectors)
	for i := range cosine_scores {
		cosine_scores[i] = results_vec.AtVec(i)
	}
	sorted_scores := utl.Sort(sort.Float64Slice(cosine_scores))
	top_items := make([]string, len(sorted_scores[top_items_k:]))
	for i, itemID := range sorted_scores[top_items_k:] {
		top_items[i], _ = handler.Embedding.GetItemNameByID(itemID)
		log.Printf("Rec Item: %v", top_items[i])
	}
	/// --------------

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

func cosine(a []float64, b []float64) (cosine float64, err error) {
	count := 0
	length_a := len(a)
	length_b := len(b)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(b[k], 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(a[k], 2)
			continue
		}
		sumA += a[k] * b[k]
		s1 += math.Pow(a[k], 2)
		s2 += math.Pow(b[k], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("vectors should not be null (all zeros)")
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}
