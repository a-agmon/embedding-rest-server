package model

import (
	
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
)

type Embeddings struct {
	Factors  map[int][]float64
	Items2ID map[string]int
}

func NewEmbedding() *Embeddings {
	return &Embeddings{}
}

func strArrayToFloat(strAr []string, arrSize int) []float64 {
	size := len(strAr)
	if size != arrSize {
		log.Fatalf("recieved array of size %v instead of 128", size)
	}
	floatArr := make([]float64, arrSize)
	for i, v := range strAr {
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Fatal("cannot convert value to float 64")
		}
		floatArr[i] = floatVal
	}
	return floatArr
}

func (e *Embeddings) GetItemNameByID(id int) (string, error) {
	for k, v := range e.Items2ID {
		if v == id {
			return k, nil
		}
	}
	return "", errors.New("item not found")
}

func (e *Embeddings) LoadItems(fileName string) error {

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("cannot read items file: %v", err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	itemsMap := make(map[string]int)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("cannot read items file: %v", err)
		}
		itemID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatalf("cannot convert item id to int: %v", err)
		}

		itemName := record[1]
		itemsMap[itemName] = itemID
	}
	e.Items2ID = itemsMap
	return nil
}


func (e *Embeddings) LoadVectors(factorsFile string, embSize int) error {

	f, err := os.Open(factorsFile)
	if err != nil {
		log.Fatalf("cannot read factors file")
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	vectorsMap := make(map[int][]float64)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading records from file %v", err)
		}

		floatVector := strArrayToFloat(record, embSize+1)
		vector_id := int(floatVector[0])
		vector_data := floatVector[1:]
		vectorsMap[vector_id] = vector_data
	}
	e.Factors = vectorsMap
	return nil
}
