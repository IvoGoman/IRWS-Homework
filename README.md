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

## Homework 2: Unigram Likelihood Model for Information Retrieval
Homework is split into a pure programming and a text processing assignment. 
For the programming assignment the tasks was to build a Unigram Likelihood Model for the 19News corpus, which is able to take ad-hoc queries. This part is implemented using Scala and the Spark Api.
For the second part of the Homework I used Python instead of calculating the different term weightings and metrics by hand. This also allows for later reuse.
