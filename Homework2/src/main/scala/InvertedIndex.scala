package main

import java.io.File
import java.util.Properties

import edu.stanford.nlp.process.DocumentPreprocessor
import org.apache.spark._

import scala.collection.immutable.{HashMap, TreeMap}
import scala.collection.mutable
import scala.io.Source

/**
  * Created by ivo on 03.04.17.
  */
class InvertedIndex {
  val root_path = "./data-local/20news-18828/comp.graphics/"
  val conf = new SparkConf().setAppName(getClass.getName).setMaster("local")
  implicit val sc = new SparkContext(conf)

  val (documentFrequencies, indextree) = createIndex(root_path)
  val totalTokens = documentFrequencies.foldLeft(0)((agg, kv) => agg + kv._2)
//  calculateMercier("That Computer")

  def getListOfFilesRecursive(f: File): Array[File] = {
    val files = f.listFiles
    files ++ files.filter(_.isDirectory).flatMap(getListOfFilesRecursive(_))
  }

  def createIndex(p: String): (mutable.HashMap[String, Int], TreeMap[String, mutable.HashMap[String, Int]]) = {
    val fileList = this.getListOfFilesRecursive(new File(p))
    val docFrequencies = mutable.HashMap[String, Int]()
    val termDocLists = for (f <- fileList if !f.isDirectory) yield {
      val name = f.getName
      val initialIndex = sc.textFile(f.getAbsolutePath).flatMap(_.replaceAll("[\\W+]", " ").split(" ")).map({
        docFrequencies += ((name, 1))
        docFrequencies.put(name, docFrequencies.getOrElse(name, 0) + 1)
        val tuple = (name, 1)
        (_, tuple)
      })
      val tempIndex = initialIndex.reduceByKey((x, y) => (x._1, (x._2 + y._2)))
      tempIndex
    }
    val combinedTDL = termDocLists.reduce((left, right) => left ++ right)
    val emptymap = mutable.HashMap[String, Int]()
    val fillMap = (left: mutable.HashMap[String, Int], right: (String, Int)) => left += (right)
    val mergeMap = (left: mutable.HashMap[String, Int], right: mutable.HashMap[String, Int]) => left ++= right
    val aggComTDL = combinedTDL.aggregateByKey(emptymap)(fillMap, mergeMap)

    var indextree = TreeMap.empty[String, mutable.HashMap[String, Int]]
    indextree = indextree ++ aggComTDL.collect()
    (docFrequencies, indextree)
  }


  def calculateMercier(query: String): Unit = {
    val tokens = query.replaceAll("[\\W+]", " ").split(" ")
    val queryResult = for (t <- tokens; if indextree.contains(t)) yield {
      val result = indextree.get(t).get
      val invertedResult = result.map(row => (row._1 -> (t -> row._2)))
      val totalFrequency = result.foldLeft(0)((v, kv2) => v + kv2._2)
      (t, result, totalFrequency)
    }

    val keys = queryResult.foldLeft(Set[String]())((x, kv)=>x ++ kv._2.keySet)
    val fullListOfTermDocFreq = mutable.HashMap[String,Double]()
    for (qR <- queryResult) {
      for (key <- keys){
        val lm = qR._2.getOrElse(key, 0)/documentFrequencies.getOrElse(key,0).toDouble
        val gm = qR._3 / totalTokens.toDouble
        fullListOfTermDocFreq.put(key, fullListOfTermDocFreq.getOrElse(key,1.0).toDouble * score(lm,gm))
      }
      }
    def score(LM:Double, GM:Double) :Double= 0.5 * LM + 1 - 0.5 * GM
    val ranking  = fullListOfTermDocFreq.toList.sortWith((x,y)=>x._2>y._2)
    for( i <- 0 to 10){
      println(s"Rank $i: ${ranking(i)}")
    }

  }


}
