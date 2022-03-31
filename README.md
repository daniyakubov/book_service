
We want to build a program that will manage a library book stock.
First we’ll focus on the book inventory management database.

We create an index in ElasticSearch called “Books”
Each book should have the following attributes:
Title
Author’s name
Price (float, represents USD)
ebook available (boolean, this field is optional - meaning not all documents will have it)
Publish date

We decided that our search will not be case sensitive.

here we create the index with it's mappings and settings:

~~~

PUT books
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "fields": {"keyword": {"type": "keyword"}}, 
        "analyzer": "my_custom_analyzer"
      },
      "author": {
        "type": "text",
        "fields": {"keyword": {"type": "keyword"}}, 
        "analyzer": "my_custom_analyzer"
      },
      "price": {
        "type": "float"
      }, 
      "available": {
        "type": "boolean"
      }, 
      "date": {
        "type": "date"
      }
    }
  },
  "settings": {
    "analysis": {
      "analyzer": {
        "my_custom_analyzer": {
          "char_filter": [],
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "english_stemmer",
            "english_possessive_stemmer",
            "syn"
          ]
        }
      },
      "filter": {
        "english_stemmer": {
          "type": "porter_stem"
        },
        "english_possessive_stemmer": {
          "type": "stemmer",
          "language": "possessive_english"
        },
        "syn": {
          "type": "synonym",
          "synonyms": [
            "Mark Twain", "Samuel Langhorne Clemens"
          ]
        }
      }
    }
  }
}

~~~

We check whether the index was created using the cluster’s health API:

~~~
GET /_cluster/health/books
~~~

We index 1 test book:

~~~
PUT books/_doc/1
{
  "title": "in cupidatat",
  "price": 1490.05,
  "authorsName": "Ryan Estrada",
  "available": false,
  "date": "2015-01-01"
}
~~~
We check that the book was created properly using the GET api:

~~~
GET books/_doc/1
~~~

We set the book’s price to be 27.5 USD:

~~~
POST books/_update/1
{
    "doc" : {
        "price" : 27.5
    }
}
~~~

We check that the price was updated:
~~~
GET books/_doc/1
~~~

We index 20 more books:
~~~
POST books/_bulk
{"index":{"_id":1}}
{"title":"velit qui magna nostrud non reprehenderit","price":138.0108,"authorsName":"Trudy Carr","available":false,"date":"1993-05-20"}
{"index":{"_id":2}}
{"title":"quis sunt eu do est exercitation","price":143.5448,"authorsName":"Consuelo Pugh","available":true,"publish,date":"1976-05-17"}
{"index":{"_id":3}}
{"title":"deserunt nisi adipisicing ex est in","price":41.0763,"authorsName":"Meadows Green","available":true,"date":"1979-12-23"}
{"index":{"_id":4}}
{"title":"eiusmod non reprehenderit duis eu nostrud","price":183.8798,"authorsName":"Catalina Riddle","available":false,"date":"1975-06-15"}
{"index":{"_id":5}}
{"title":"et esse ad culpa voluptate ex","price":45.7831,"authorsName":"Lacey Lee","available":true,"date":"2011-02-09"}
{"index":{"_id":6}}
{"title":"magna laboris amet quis excepteur commodo","price":172.9611,"authorsName":"Mcpherson Matthews","available":true,"date":"1976-09-18"}
{"index":{"_id":7}}
{"title":"enim fugiat officia occaecat duis proident","price":179.1719,"authorsName":"Yvette Osborne","available":false,"date":"2014-02-21"}
{"index":{"_id":8}}
{"title":"tempor consectetur laboris ut laborum incididunt","price":115.8712,"authorsName":"Jennings Keller","available":false,"date":"1985-11-28"}
{"index":{"_id":9}}
{"title":"irure est fugiat quis reprehenderit esse","price":173.1195,"authorsName":"Walsh Norton","available":false,"date":"2013-12-19"}
{"index":{"_id":10}}
{"title":"sunt id sit duis ullamco veniam","price":49.5798,"authorsName":"Janelle Middleton","available":false,"date":"1989-02-07"}
{"index":{"_id":11}}
{"title":"ea reprehenderit consectetur duis aute fugiat","price":111.3198,"authorsName":"Traci Mclean","available":true,"date":"2009-03-27"}
{"index":{"_id":12}}
{"title":"elit consequat nisi est nulla est","price":37.9667,"authorsName":"Elisabeth Wong","available":true,"date":"2002-02-14"}
{"index":{"_id":13}}
{"title":"eu nisi deserunt voluptate est anim","price":159.3275,"authorsName":"Ana Joyce","available":false,"date":"2000-10-20"}
{"index":{"_id":14}}
{"title":"aliqua do eiusmod qui reprehenderit aliqua","price":43.4622,"authorsName":"Hines Reeves","available":true,"date":"1993-08-10"}
{"index":{"_id":15}}
{"title":"cillum incididunt esse do elit sunt","price":140.1977,"authorsName":"Myrtle Ray","available":true,"date":"1997-08-06"}
{"index":{"_id":16}}
{"title":"commodo ut amet veniam nisi proident","price":65.1081,"authorsName":"Glenda Vinson","available":false,"date":"2000-11-26"}
{"index":{"_id":17}}
{"title":"non minim commodo voluptate magna consectetur","price":104.3833,"authorsName":"Miles Mayo","available":true,"date":"1979-08-31"}
{"index":{"_id":18}}
{"title":"Lorem labore ad anim est ut","price":189.54,"authorsName":"Aguilar Downs","available":true,"date":"1974-11-17"}
{"index":{"_id":19}}
{"title":"eu deserunt commodo sint dolor do","price":97.6843,"authorsName":"Anna Byers","available":true,"date":"1976-01-31"}
{"index":{"_id":20}}
{"title":"tempor laboris ea ex proident fugiat","price":137.3987,"authorsName":"Francis Roberts","available":true,"date":"1996-12-10"}

~~~

We search for all the books that have an ebook version available. 
~~~
GET books/_search
{
  "query": {
    "constant_score": {
      "filter": {
        "term": {
          "available": true
        }
      }
    }
  }
}
~~~

We search for all the books that don’t have an ebook version available.

~~~
GET books/_search
{
  "query": {
    "constant_score": {
      "filter": {
        "term": {
          "available": false
        }
      }
    }
  }
}
~~~

We get all the books that are priced higher than 50 dollars:

~~~
GET books/_search
{
  "query": {
    "range": {
      "price": {
        "gte": 50
      }
    }
  }
}
~~~

We get all the books that are priced between 10 to 50 dollars or 100 to 200 dollars:

~~~
GET books/_search
{
  "query": {
    "constant_score": {
      "filter": {
        "bool": {
          "should": [
            {
              "range": {
                "price": {
                  "gte": 10,
                  "lte": 50
                }
              }
            },
            {
              "range": {
                "price": {
                  "gte": 100,
                  "lte": 200
                }
              }
            }
          ]
        }
      }
    }
  }
}

~~~

We detect how many books are available online and how many aren’t, at the same query:
~~~
GET books/_search
{
  "size": 0,
  "aggs": {
    "available": {
      "terms": { "field": "available" }
    }
  }
}
~~~

We get the median price: 
~~~
GET books/_search
{
    "size" : 0,
    "aggs" : {
        "median price" : {
            "percentiles" : {
                "field" : "price",
                "percents" : 50
            }
        }
    }
}
~~~

We get the upper 10% price:

~~~
GET books/_search
{
    "size" : 0,
    "aggs" : {
        "median price" : {
            "percentiles" : {
                "field" : "price",
                "percents" : 90
            }
        }
    }
}
~~~