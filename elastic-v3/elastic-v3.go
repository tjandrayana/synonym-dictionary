package main

import (
	"fmt"
	"log"
	"time"

	elastic "gopkg.in/olivere/elastic.v3"
)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type DigitalOrderFinish struct {
	OrderId           int64     `json:"order_id"`
	UserId            int64     `json:"user_id"`
	Product           string    `json:"product"`
	PhoneNumber       string    `json:"phone_number"`
	UserEmail         string    `json:"user_email"`
	TopicConsume      string    `json:"-"`
	TopicPublish      string    `json:"-"`
	IpAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	PromoCode         string    `json:"promo_code"`
	Cashback          float64   `json:"cashback"`
	PgId              int       `json:"pg_id"`
	AmountCut         float64   `json:"amount_cut"`
	TopPoints         float64   `json:"lp_amount"`
	TransactionStart  time.Time `json:"transaction_start"`
	TransactionFinish time.Time `json:"transaction_finish"`
	RegisterTime      int64     `json:"register_time"`
	ClientNumber      string    `json:"client_number"`
	SalesPrice        float64   `json:"sales_price"`
	VoucherCode       string    `json:"voucher_code"`
	FirstTime         bool      `json:"first_time"`

	Created time.Time             `json:"created,omitempty"`
	Suggest *elastic.SuggestField `json:"suggest_field,omitempty"`
}

func main() {

	// create a client and add Compose connection strings
	client, err := elastic.NewClient(
		elastic.SetURL("http://192.168.100.16:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal(err)
	}
	// create a variable that stores the result
	// of the executed cluster health query and prints the result
	health, err := client.ClusterHealth().Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("<------ Cluster Health ------>\n%+v\n", health)

	// errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can also provide your own HTTP client here.
	// client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }

	// Trace request and response details like this
	//client.SetTracer(log.New(os.Stdout, "", 0))

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://192.168.100.16:9200").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("digital").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("digital").Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Index a tweet (using JSON serialization)
	transaction1 := DigitalOrderFinish{
		OrderId:           1,
		UserId:            2001,
		Product:           "SIMPATI50",
		PhoneNumber:       "622138531111",
		UserEmail:         "coba@gmail.com",
		IpAddress:         "127.168.0.1",
		UserAgent:         "MozilaFX....",
		PromoCode:         "HEMAT50",
		Cashback:          0,
		PgId:              1,
		TransactionStart:  time.Now().Add(-10 * time.Minute),
		TransactionFinish: time.Now(),
		ClientNumber:      "622138531111",
		SalesPrice:        50000,
		FirstTime:         true,
	}
	put1, err := client.Index().
		Index("digital").
		Type("orderfinish").
		Id("1").
		BodyJson(transaction1).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed digital %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// // Get tweet with specified ID
	// get1, err := client.Get().
	// 	Index("twitter").
	// 	Type("tweet").
	// 	Id("1").
	// 	Do()
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// if get1.Found {
	// 	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	// }
	//
	// // Flush to make sure the documents got written.
	// _, err = client.Flush().Index("twitter").Do()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Search with a term query
	// termQuery := elastic.NewTermQuery("user", "olivere")
	// searchResult, err := client.Search().
	// 	Index("twitter").   // search in index "twitter"
	// 	Query(termQuery).   // specify the query
	// 	Sort("user", true). // sort by "user" field, ascending
	// 	From(0).Size(10).   // take documents 0-9
	// 	Pretty(true).       // pretty print request and response JSON
	// 	Do()                // execute
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	//
	// // searchResult is of type SearchResult and returns hits, suggestions,
	// // and all kinds of other information from Elasticsearch.
	// fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	//
	// // Each is a convenience function that iterates over hits in a search result.
	// // It makes sure you don't need to check for nil values in the response.
	// // However, it ignores errors in serialization. If you want full control
	// // over iterating the hits, see below.
	// var ttyp Tweet
	// for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
	// 	t := item.(Tweet)
	// 	fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	// }
	// // TotalHits is another convenience function that works even when something goes wrong.
	// fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())
	//
	// // Here's how you iterate through results with full control over each step.
	// if searchResult.Hits.TotalHits > 0 {
	// 	fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)
	//
	// 	// Iterate through results
	// 	for _, hit := range searchResult.Hits.Hits {
	// 		// hit.Index contains the name of the index
	//
	// 		// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
	// 		var t Tweet
	// 		err := json.Unmarshal(*hit.Source, &t)
	// 		if err != nil {
	// 			// Deserialization failed
	// 		}
	//
	// 		// Work with tweet
	// 		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	// 	}
	// } else {
	// 	// No hits
	// 	fmt.Print("Found no tweets\n")
	// }
	//
	// // Update a tweet by the update API of Elasticsearch.
	// // We just increment the number of retweets.
	// script := elastic.NewScript("ctx._source.retweets += num").Param("num", 1)
	// update, err := client.Update().Index("twitter").Type("tweet").Id("1").
	// 	Script(script).
	// 	Upsert(map[string]interface{}{"retweets": 0}).
	// 	Do()
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)
	//
	// // ...
	//
	// // Delete an index.
	// deleteIndex, err := client.DeleteIndex("twitter").Do()
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// if !deleteIndex.Acknowledged {
	// 	// Not acknowledged
	// }

}
