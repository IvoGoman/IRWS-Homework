package main

import java.io.File

import org.apache.spark._

import scala.collection.immutable.TreeMap
import scala.collection.mutable

/**
  * Created by ivo on 03.04.17.
  */
class InvertedIndex(implicit sc: SparkContext) extends Serializable {
  val root_path = "./data-local/20news-18828/comp.graphics/"
  val (documentFrequencies, indextree) = createIndex(root_path)
  val totalTokens = documentFrequencies.foldLeft(0)((agg, kv) => agg + kv._2)

  /* Function that will return all files in the subdirectories of the given rootdirectory*/
  def getListOfFilesRecursive(f: File): Array[File] = {
    val files = f.listFiles
    files ++ files.filter(_.isDirectory).flatMap(getListOfFilesRecursive(_))
  }

  /*Function that creates the inverted index over the corpus of documents*/
  def createIndex(p: String): (mutable.HashMap[String, Int], TreeMap[String, mutable.HashMap[String, Int]]) = {
    //    get the list of all files under the root
    val fileList = this.getListOfFilesRecursive(new File(p))
    val docFrequencies = mutable.HashMap[String, Int]()
    //    iterate through each file which is not a directory
    val termDocLists = for (f <- fileList if !f.isDirectory) yield {
      //      get the name of the file, which will be used as the DocId
      val name = f.getName
      //      Index is initalized by creating an inverted index of form (word->(docID, 1)
      val initialIndex = sc.textFile(f.getAbsolutePath).flatMap(_.replaceAll("\"[!?,;:<>()\\[\\]']\"", " ").split("\\s+")).map({
        //        Create the count of tokens for each document, needed for Local Model
        val tuple = (name, 1)
        (_, tuple)
      })
      //      Reduce the temporary index of the file by same words and add their counts
      val tempIndex = initialIndex.reduceByKey((x, y) => (x._1, x._2 + y._2))
      docFrequencies.put(name, docFrequencies.getOrElse(name, 0) + tempIndex.aggregate(0)((agg, value) => agg + value._2._2, (left: Int, right: Int) => left + right))
      tempIndex
    }

    //    combine the inverted indices of the different documents into one big one
    val combinedTDL = termDocLists.reduce((left, right) => left ++ right)
    //    Aggregate the DocID and Frequency of the different documents and the words into a hashmap
    val emptymap = mutable.HashMap[String, Int]()
    val fillMap = (left: mutable.HashMap[String, Int], right: (String, Int)) => left += (right)
    val mergeMap = (left: mutable.HashMap[String, Int], right: mutable.HashMap[String, Int]) => left ++= right
    val aggComTDL = combinedTDL.aggregateByKey(emptymap)(fillMap, mergeMap)
    //  Convert the RDD into a RedBlack Tree for easier retrieval
    var indextree = TreeMap.empty[String, mutable.HashMap[String, Int]]
    indextree = indextree ++ aggComTDL.collect()
    (docFrequencies, indextree)
  }


  def getRankedResults(query: String, lambda: Double): Unit = {
    //    tokenize the query
    val tokens = query.replaceAll("[\\W+]", " ").split(" ")
    //    get the entries of the inverted index for each word
    val queryResult = for (t <- tokens;
                           if indextree.contains(t)) yield {
      val result = indextree.get(t).get
      //      val invertedResult = result.map(row => (row._1 -> (t -> row._2)))
      //      calculate the total frequency of the tokens over all documents in which they occur
      val totalFrequency = result.foldLeft(0)((v, kv2) => v + kv2._2)
      (t, result, totalFrequency)
    }
    //  Get all document keys from the list of documents retrieved from the inverted index
    val keys = queryResult.foldLeft(Set[String]())((x, kv) => x ++ kv._2.keySet)
    val fullListOfTermDocFreq = mutable.HashMap[String, Double]()
    //    For each DocID get either the frequency of which the query term occured in the documents or 0
    for (qR <- queryResult) {
      //        create the global unigram model for the query term
      val gm = qR._3 / totalTokens.toDouble
      for (key <- keys) {
        //        create local unigram model for the query term
        val lm = qR._2.getOrElse(key, 0) / documentFrequencies.getOrElse(key, 0).toDouble
        //        for each DocID add the Jelinek-Mercer Score
        fullListOfTermDocFreq.put(key, fullListOfTermDocFreq.getOrElse(key, 1.0).toDouble * calcJelinekMercier(lm, gm, lambda))
      }
    }

    def calcJelinekMercier(LM: Double, GM: Double, lambda: Double): Double = lambda * LM + (1 - lambda) * GM

    //    order the ranked results descending and print out the top 10
    val ranking = fullListOfTermDocFreq.toList.sortWith((x, y) => x._2 > y._2)
    var i = 10
    if (ranking.size < 10) {
      i = ranking.size
    }
    for (i <- 0 to i) {
      println(s"Rank $i: ${
        ranking(i)
      }")
    }

  }
}

object InvertedIndex extends App {
  val conf = new SparkConf().setAppName(getClass.getName).setMaster("local")
  implicit val sc = new SparkContext(conf)
  val myIndex = new InvertedIndex
  myIndex.getRankedResults("Computer Science Football Game", 0.5)
}