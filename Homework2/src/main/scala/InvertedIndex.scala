package main

import java.io.File

import org.apache.spark._

import scala.collection.immutable.HashMap
import scala.collection.mutable
import scala.io.Source

/**
  * Created by ivo on 03.04.17.
  */
object InvertedIndex extends App {
  val root_path = "./data-local/20news-18828/"
  val conf = new SparkConf().setAppName(getClass.getName).setMaster("local")
  implicit val sc = new SparkContext(conf)

  val index = createIndex(root_path)

  def getListOfFilesRecursive(f: File): Array[File] = {
    val files = f.listFiles
    files ++ files.filter(_.isDirectory).flatMap(getListOfFilesRecursive(_))
  }

  def createIndex(p: String): Unit = {
    val fileList = this.getListOfFilesRecursive(new File(p))
    //    val file = new File(p)
    val termDocLists = for (f <- fileList if !f.isDirectory) yield {
      val name = f.getName
      val initialIndex = sc.textFile(f.getAbsolutePath).flatMap(_.split(" ")).map({
        val map = mutable.HashMap(name -> 1);
        (_, map)
      })
      initialIndex
    }

//    val merged = for (tdl <-termDocLists){
      val emptymap = mutable.HashMap[String, Int]()
      val mergeMap = (left:mutable.HashMap[String, Int], right:mutable.HashMap[String, Int]) => left++=right
      val aggregatedIndex = termDocLists.aggregateByKey(emptymap)(mergeMap, mergeMap)
//      val reducedIndex = initialIndex.reduceByKey((left, right) => {
//        left.merged(right)({ case ((k, v1), (_, v2)) => (k, v1 + v2) })
//      })
  }
}
