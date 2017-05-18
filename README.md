# IRWS-Homework
Repository including all the programming assignements given throughout the course of Information Retrieval and Web Search at the University of Mannheim during the Spring Term 2017

## Homework 1: Minimum Edit Distance
For Homework 1 the (Damerau-)Levensthein distance has been implemented both with dynamic programming and recursions.

There are several flags to customize what happens at runtime
  - original: left side of the comparison
  - compare: right side of the comparison
  - recursive: true if a recursive version shall be used
  - damerau: true if the Damerau-Levensthein Distance shall be used
  - weigths: true if custom weights for transposition/ replacement shall be used
  
For the implementation golang was used, there are a couple of tests to show sample output and benchmark tests to see the difference in runtime between recursive and dynamic programming versions.

## Homework 2: Vector Space and Probabilistic Retrieval
- __Term weighting__:
Compute TF-IDF for a toy document collection with different definitions for TF and IDF and rank the documents given a query with cosine similarity.
- __Distance/similarity metrics__:
Ranking of documents given a query and 'raw Euclidean distance', 'normalized Euclidean distance' and 'cosine similarity'
- __Optimizing vector space model__:
Given a toy collection of TF-IDF vectors perform random projections to reduce computation costs. Do a pre-clustering of the documents using a given set of leader vectors. Finally retrieve top 5 documents for a query vector using the random projection vectors and leader vectors with clusters.
- __Classic probabilistic retrieval__:
Given a query rank documents with 'Binary independence model', 'Two-Poisson model', 'BM25'
- __Unigram Likelihood Model for Information Retrieval__:
For the programming assignment the tasks was to build a query likelihood model based on a unigram Likelihood Model for the 20 News corpus, which is able to take ad-hoc queries and rank the documents by relevance based on the unigram model. This part is implemented using Scala and the Spark Api.

## Homework 3: Semantic Retrieval, Text Clustering, and IR Evaluation
- __Latent Semantic Indexing__: Computing the similarity of latent vectors for a toy collection of documents and a query
- __Text Clustering__: Using 'K-Means' and 'Single Pass Clustering' to cluster a toy collection of TF-IDF vectors
- __IR Evaluation__: Calculating precision, recall, F1, P@k, R-precision, average precision and mean average precision for a toy collection of retrievals and their relevance rating
- __Semantic Retrieval with Word-Embeddings__: Implementation of a simple retrieval engine based on aggregation of word embeddings using the pretrained 'GloVe' word embeddings and a random subsample of 500 documents from the '20 News Groups dataset'
