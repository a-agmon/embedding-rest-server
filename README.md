# embedding-rest-server

This is a REST server implemented in GO (using GIN as a web server) that wraps users X items embeddings (that you can generate using your favorite matrix factorization or ML model) and serve similarity requests over HTTP-REST.
Although this is mainly for testing and debugging purposes, it runs a similarity search on an embedding really fast.
It returns an answer for a cosine similarity search over nearly 18K items, each represented by a vector of length 128 in about 40ms.
(Tested using the lastFM dataset)
Example for usage
```bash
curl [URL]:9090/mostsimilar \
-d '{"similarto":"Eminem", "topk":10}' \
-H 'Content-Type: application/json' -X GET | jq

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

