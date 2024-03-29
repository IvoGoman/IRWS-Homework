from pprint import pprint
import numpy as np
from nltk.stem.porter import PorterStemmer
from nltk.stem import WordNetLemmatizer
from nltk.corpus import wordnet as wn
import nltk
import math

def preprocess_query(query):
    """Preprocessing of the corpus, filter for nouns and adjectives and lemmatize"""
    # stop = set(stopwords.words('english'))
    tags = {'NN', 'NNS', 'NNP', 'NNP', 'NNPS', 'JJ', 'JJR', 'JJS'}
    wordnet_lemmatizer = WordNetLemmatizer()
    # for i in range(len(query)):
    query = [(word.lower(), convert(tag)) for (word, tag) in nltk.pos_tag(nltk.word_tokenize(query)) if tag in tags]
    query = [wordnet_lemmatizer.lemmatize(w, t) for (w, t) in query ]
    return query

def preprocess(docs):
    """Preprocessing of the corpus, filter for nouns and adjectives and lemmatize"""
    # stop = set(stopwords.words('english'))
    tags = {'NN', 'NNS', 'NNP', 'NNP', 'NNPS', 'JJ', 'JJR', 'JJS'}
    for i in range(len(docs)):
        docs[i] = [(word.lower(), convert(tag)) for (word, tag) in nltk.pos_tag(nltk.word_tokenize(docs[i])) if tag in tags]
    return lemmatize_docs(docs)

def lemmatize_docs(docs):
    """Lemmatize the terms of the corpus"""
    wordnet_lemmatizer = WordNetLemmatizer()
    for i in range(len(docs)):
        docs[i] = [wordnet_lemmatizer.lemmatize(w, t) for (w, t) in docs[i]]
    return docs

def convert(tag):
    """Convert tag from treebank to wordnet format"""
    if is_noun(tag):
        return wn.NOUN
    if is_adjective(tag):
        return wn.ADJ

def is_noun(tag):
    """True if tag corresponds to treebank noun tags"""
    return tag in ['NN', 'NNS', 'NNP', 'NNPS']

def is_adjective(tag):
    """True if tag corresponds to treebank adjective tags"""
    return tag in ['JJ', 'JJR', 'JJS']

def stem_docs(docs):
    """Stems the tokens in the corpus with Porter Stemmer"""
    porter_stemmer = PorterStemmer()
    for i in range(len(docs)):
        docs[i] = [porter_stemmer.stem(t) for t in docs[i]]

def calc_tf(doc):
    """Calculates the tf scores based on the corpus"""
    tf = {}
    for term in doc:
        if term not in tf:
            tf[term] = doc.count(term)
    return tf

def calc_tf_log(doc):
    """Calculates the tf scores based on the corpus"""
    tf = calc_tf(doc)
    max_tf = tf[max(tf, key=tf.get)]
    tf_log = {}
    for key, val in tf.items():
        tf_log[key] = (1 + math.log(val)) / (1 + math.log(max_tf))
    return tf_log

def calc_idf(docs):
    """Calculates the idf scores based on the corpus"""
    terms = set()
    for doc in docs:
        for term in doc:
            terms.add(term)
    idf = {}
    for term in terms:
        term_count = 0
        doc_count = 0
        for doc in docs:
            doc_count += 1
            if term in doc:
                term_count += 1
        idf[term] = doc_count/term_count
    return idf

def calc_tdf(docs):
    """Calculates the tdf scores based on the corpus"""
    terms = set()
    for doc in docs:
        for term in doc:
            terms.add(term)
    tdf = {}
    for term in terms:
        doc_count = 0
        for doc in docs:
            doc_count += 1 if term in doc else 0
        tdf[term] = doc_count
    return tdf

def calc_idf_two(docs):
    """Calculates the df scores based on the corpus"""
    terms = set()
    for doc in docs:
        for term in doc:
            terms.add(term)
    idf = {}
    for term in terms:
        term_count = 0
        doc_count = 0
        for doc in docs:
            doc_count += 1
            if term in doc:
                term_count += 1
        idf[term] = max(0, ((doc_count-term_count)/term_count))
    return idf

def calc_tf_idf(idf, tf):
    """Calculate the TF-IDF score based on tf and idf components"""
    tfidf = {}
    for key, val in tf.items():
        tfidf[key] = val * idf[key]
    return tfidf

def manual_correction(docs):
    # """Some manual correction to pos-tag filtered words"""
    docs[0]=[x.replace('orcs', 'orc') for x in docs[0]]
    docs[0].remove("hundred")
    docs[0].append("surround")
    docs[1].append("body")
    docs[2].append("closer")
    docs[2].append("radiat")
    docs[2]=[x.replace('orcs', 'orc') for x in docs[2]]
    docs[3]=[x.replace('orcs', 'orc') for x in docs[3]]
    docs[3].remove("only")
    docs[3].append("deadly")
    return docs

def calc_cos_sim(query, tfidf):
    dotproduct = 0
    logx = 0
    logy = 0
    for term in query:
        if term in tfidf:
            dotproduct += tfidf[term]
    logx = math.sqrt(len(query))
    for key, val in tfidf.items():
        logy += math.pow(val, 2)
    logy = math.sqrt(logy)
    return dotproduct / (logx * logy)

def calc_euclidean_sim(query, tfidf):
    distance = 0
    query = set(query)
    c = 0
    for q in query:
        distance = distance + math.pow(1-tfidf[q], 2) if q in tfidf.keys() else distance + math.pow(1-0, 2)
    for k, v in tfidf.items():
        if k not in query:
            distance += math.pow(0-v, 2)
    return math.sqrt(distance)

def unitnorm(vector):
    unitnorm = 0
    for k,v in vector.items():
        unitnorm += v
    return unitnorm

def calc_euclidean_norm_dist(query, tfidf):
    distance = 0
    query = set(query)
    qlen = len(query)
    dlen = unitnorm(tfidf)
    c = 0
    for q in query:
        distance = distance + math.pow((1/qlen)-(tfidf[q]/dlen), 2) if q in tfidf.keys() else distance + math.pow((1/qlen), 2)
    for k, v in tfidf.items():
        if k not in query:
            distance += math.pow(v/dlen, 2)
    return math.sqrt(distance)

def calc_binary_ind(query, docs, tdf):
    bign = len(docs)
    weight = {}
    for d in range(len(docs)):
        score = 0
        for t in query:
            score = math.log10((bign-tdf[t]+0.5) / (tdf[t]+0.5)) if t in docs[d] else math.log10((0.5 + bign) / (0.5))
        weight[d] = score
    return weight

def calc_two_poisson(query, docs, tdf, k):
    ranking = {}
    bign = len(docs)
    for d in range(len(docs)):
        score = 0
        term_count = 0
        for term in query:
            for word in docs[d]:
                term_count += 1 if term == word else 0
            w = math.log10((bign-tdf[term]+0.5) / (tdf[term]+0.5)) if term in set(docs[d]) else math.log10((0.5 + bign) / (0.5))
            score += ((term_count * (k+1))/(term_count + k)) * w
        ranking[d] = score
    return ranking

def calc_bm_eleven(query, docs, tdf, k):
    ranking = {}
    len_avg = 0
    bign = len(docs)
    for doc in docs:
        len_avg += len(doc)
    len_avg = len_avg/len(docs)
    for d in range(len(docs)):
        score = 0
        term_count = 0
        for term in query:
            for word in docs[d]:
                term_count += 1 if term == word else 0
            w = math.log10((bign-tdf[term]+0.5) / (tdf[term]+0.5)) if term in set(docs[d]) else math.log10((0.5 + bign) / (0.5))
            score += ((term_count * (k+1))/(term_count + k * (len(docs[d])/len_avg))) * w
        ranking[d] = score
    return ranking

def calc_bm_twentyfive(query, docs, tdf, k, b):
    ranking = {}
    len_avg = 0
    bign = len(docs)
    for doc in docs:
        len_avg += len(doc)
    len_avg = len_avg/len(docs)
    for d in range(len(docs)):
        score = 0
        term_count = 0
        for term in query:
            for word in docs[d]:
                term_count += 1 if term == word else 0
            w = math.log10((bign-tdf[term]+0.5) / (tdf[term]+0.5)) if term in set(docs[d]) else math.log10((0.5 + bign) / (0.5))
            score += ((term_count * (k+1))/(term_count + k * (len(docs[d])/len_avg)*b+k*(1-b))) * w
        ranking[d] = score
    return ranking

def main():
    d1 = """Frodo and Sam were trembling in the darkness, surrounded in darkness by hundreds of blood-thirsty orcs. Sam was certain these beasts were about to taste the scent of their flesh."""
    d2 = """The faceless black beast then stabbed Frodo. He felt like every nerve in his body was hurting. Suddenly, he thought of Sam and his calming smile. Frodo had betrayed him."""
    d3 = """Frodo's sword was radiating blue, stronger and stronger every second. Orcs were getting closer. And these weren't just regular orcs either, Uruk-Hai were among them. Frodo had killed regular orcs before, but he had never stabbed an Uruk-Hai, not with the blue stick."""
    d4 = """Sam was carrying a small lamp, shedding some blue light. He was afraid that orcs might spot him, but it was the only way to avoid deadly pitfalls of Mordor."""
    query = preprocess_query("Sam and blue orc")
    query2 = preprocess_query("Sam and Frodo love blue orcs")
    query2 = [x.replace('orcs', 'orc') for x in query2]
    docs = [d1, d2, d3, d4]
    docs = preprocess(docs)
    docs = manual_correction(docs)
    pprint(docs)
    print('IDF scores')
    idf = calc_idf(docs)
    idf2 = calc_idf_two(docs)
    pprint(idf)
    cosinesim = {}
    cosinesim2 = {}
    euclideansim = {}
    euclideansimnorm = {}
    for d in range(len(docs)):
        print('************************************************************')
        print('Document %d' % (d + 1))
        print('************************************************************')
        tf = calc_tf(docs[d])
        tfidf = calc_tf_idf(idf, tf)
        cosinesim[d] = calc_cos_sim(query, tfidf)
        euclideansim[d] = calc_euclidean_sim(query, tfidf)
        euclideansimnorm[d] = calc_euclidean_norm_dist(query, tfidf)
        print('TF scores')
        pprint(tf)
        print('TF-IDF scores')
        pprint(tfidf)
        print('************************************************************')
        tf_log = calc_tf_log(docs[d])
        tfidf2 = calc_tf_idf(idf2, tf_log)
        cosinesim2[d] = calc_cos_sim(query, tfidf2)
        print('TF scores')
        pprint(tf_log)
        print('TF-IDF scores')
        pprint(tfidf2)

    tdf = calc_tdf(docs)
    binaryindep = calc_binary_ind(query2, docs, tdf)
    twopoisson = calc_two_poisson(query2, docs, tdf, 1)
    bmeleven = calc_bm_eleven(query2, docs, tdf, 1)
    bmtwentyfive = calc_bm_twentyfive(query2, docs, tdf, 1, 0.75)

    print('************************************************************')
    print("Query: Sam and blue orc")
    print('************************************************************')
    sorted_similarity = sorted(cosinesim.items(), key=lambda item: item[1], reverse=True)
    print('Cosine Similarity Ranking')
    pprint(sorted_similarity)
    sorted_similarity = sorted(cosinesim2.items(), key=lambda item: item[1], reverse=True)
    print('Cosine Similarity Ranking')
    pprint(sorted_similarity)
    sorted_similarity = sorted(euclideansim.items(), key=lambda item: item[1], reverse=True)
    print('Euclidean Similarity Ranking')
    pprint(sorted_similarity)
    sorted_similarity = sorted(euclideansimnorm.items(), key=lambda item: item[1], reverse=True)
    print('Normalized Euclidean Similarity Ranking')
    pprint(sorted_similarity)
    print('************************************************************')
    print('Query: Sam and Frodo love blue orcs')
    print('************************************************************')
    print('Binary Independence Model')
    pprint(sorted(binaryindep.items(), key=lambda item: item[1], reverse=True))
    print('************************************************************')
    print('Two Poisson Model')
    pprint(sorted(twopoisson.items(), key=lambda item: item[1], reverse=True))
    print('************************************************************')
    print('BM 11 Model')
    pprint(sorted(bmeleven.items(), key=lambda item: item[1], reverse=True))
    print('************************************************************')
    print('BM 25 Model')
    pprint(sorted(bmtwentyfive.items(), key=lambda item: item[1], reverse=True))
main()
