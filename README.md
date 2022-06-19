# Embedding REST Server

This is a REST service implemented in GO (using GIN as a web server) that wraps embeddings of items (or words) that are typically used for recommendations.  Alhough it was mainly created for testing and debugging purposes, it enables a quick and efficient way of serving REST based similarity requests using embeddings generated by your favorite matrix factorization or ML model.  You simply point the server to an *embedding file* and an *item name-id mapping* and it can start serving requests as below.  
On a reasonable scale, it runs a similarity search on an embedding really fast.
It returns an answer for a similarity search over nearly 18K items, each represented by a vector of length 128, in about 30ms.
(Tested using the lastFM dataset)
Example for usage
```bash
curl http://127.0.0.1:8080/similar?to=Eminem&topk=10 |jq
```

```json
{
  "original": "Eminem",
  "similar": [
    "Eminem",
    "Ludacris",
    "Jay-Z",
    "Kanye West",
    "2Pac",
    "50 Cent",
    "Juno Reactor vs. Don Davis",
    "The Brothers Gutworm",
    "Black Eyed Peas",
    "Timbaland"
  ]

```


Projects/packages that are used here - 
1. *Argsort* - for fast matrix sorting  -  https://github.com/mkmik/argsort/blob/v1.1.0/argsort.go
2. *Facts* - for recommend requests based on colaborative filtering
3. *gonum* - for fast matrix and vector operations
4. *go learn* - for the Cosine distance and Dot functions
5. and *Gin*

## Running

The server can be easily compiled as 

```bash
go build -o bin/server ./server 
```

and then run while pointing to the embedding files 

```bash
./bin/server ./bin/server.config.yaml 
```

Config file should follow this structure

The supported strcuture of the embedding file should be in the following CSV format:
```csv
id, vector_element_1, vector_element_2, vector_element_3 ......
```
The supported strcuture of the embedding map file should be in the following CSV format:
```csv
id, item_name
```
**note that it is case sensitive** 





