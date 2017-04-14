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
Homework 2 has multiple mistakes hidden in depths of unreadable monkey-written scala code. It is your special task to find these mistakes and make the author personally liable for his wrongdoings.
Hakuna Matata!
