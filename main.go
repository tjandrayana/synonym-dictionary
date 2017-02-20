package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"golang.org/x/net/context"

	elastic "gopkg.in/olivere/elastic.v5"
)

type DigitalOrderFinish struct {
	OrderID           int64     `json:"order_id"`
	UserID            int64     `json:"user_id"`
	Product           string    `json:"product"`
	PhoneNumber       string    `json:"phone_number"`
	UserEmail         string    `json:"user_email"`
	TopicConsume      string    `json:"-"`
	TopicPublish      string    `json:"-"`
	IpAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	PromoCode         string    `json:"promo_code"`
	Cashback          float64   `json:"cashback"`
	PgID              int       `json:"pg_id"`
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
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	// client, err := elastic.NewClient()
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }

	client, err := elastic.NewClient(
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	// Ping the Elasticsearch server to get e.g. the version number
	// info, code, err := client.Ping().Do(ctx)
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// fmt.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("digital").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("digital").Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// bulkInReq := elastic.BulkIndexRequest{
	// 	index: "twitter",
	// 	typ:   "tweet",
	// 	id:    "id1",
	// 	ttl:   "1m",
	// }

	// a := elastic.NewBulkIndexRequest()
	// x := a.OpType("index").Index("index1").Type("tweet").Id("1").TTL("1m").Doc(Tweet{User: "olivere", Created: time.Date(2014, 1, 18, 23, 59, 58, 0, time.UTC)})
	// //
	//
	// put1, err := elastic.NewBulkService(client).
	// 	Index("twitter").
	// 	Type("tweet").
	// 	// Id("3").
	// 	Add(x).
	// 	Do(ctx)
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// fmt.Printf("Indexed tweet %s \n", put1)

	// Index a tweet (using JSON serialization)
	transaction1 := DigitalOrderFinish{
		OrderID:           1,
		UserID:            2001,
		Product:           "SIMPATI00000",
		PhoneNumber:       "622138531111",
		UserEmail:         "coba@gmail.com",
		IpAddress:         "127.168.0.1",
		UserAgent:         "MozilaFX....",
		PromoCode:         "HEMAT50000",
		Cashback:          0,
		PgID:              1,
		TransactionStart:  time.Now().Add(-10 * time.Minute),
		TransactionFinish: time.Now(),
		ClientNumber:      "622138531111",
		SalesPrice:        150000,
		FirstTime:         true,
	}

	transaction2 := DigitalOrderFinish{
		OrderID:           2,
		UserID:            2001,
		Product:           "SIMPATI100!!!!!!!",
		PhoneNumber:       "622138531111",
		UserEmail:         "coba@gmail.com",
		IpAddress:         "127.168.0.1",
		UserAgent:         "MozilaFX....",
		PromoCode:         "HEMAT50",
		Cashback:          0,
		PgID:              1,
		TransactionStart:  time.Now().Add(-10 * time.Minute),
		TransactionFinish: time.Now(),
		ClientNumber:      "622138531111",
		SalesPrice:        50000,
		FirstTime:         true,
	}
	put1, err := client.Index().
		Index("digital").
		Type("orderfinish").
		TTL("1m").
		// Id("AVpVL_2TZhO9OinkNoNJ").
		BodyJson(transaction1).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	put2, err := client.Index().
		Index("digital").
		Type("orderfinish").
		Id("2").
		TTL("1m").
		BodyJson(transaction2).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed transaction %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	fmt.Println("Put 2 = ", put2, "\n")

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("digital").
		Type("orderfinish").
		Id("1").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("digital").Do(ctx)
	if err != nil {
		panic(err)
	}

	// Search with a term query
	// termQuery := elastic.NewTermQuery("SIMPATI50", "product")
	searchResult, err := client.Search().
		Index("digital").
		Type("orderfinish"). // search in index "twitter"
		// Query(termQuery).    // specify the query
		// Sort("product", true). // sort by "user" field, ascending
		// From(0).Size(10). // take documents 0-9
		// Pretty(true).     // pretty print request and response JSON
		Do(ctx) // execute

	if err != nil {
		// Handle error
		panic(err)
	}
	//
	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.Hits)
	//
	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var dof DigitalOrderFinish
	for _, item := range searchResult.Each(reflect.TypeOf(dof)) {
		if t, ok := item.(DigitalOrderFinish); ok {
			fmt.Printf("Order ID %d is User ID %d that used Promo: %s\n", t.OrderID, t.UserID, t.Product)
		}
	}
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
	// update, err := client.Update().Index("twitter").Type("tweet").Id("1").
	// 	// Script("ctx._source.retweets += num").
	// 	// ScriptParams(map[string]interface{}{"num": 1}).
	// 	Upsert(map[string]interface{}{"retweets": 0}).
	// 	Do(ctx)
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)
	//
	// // ...
	//
	// // Delete an index.
	// deleteIndex, err := client.DeleteIndex("twitter").Do(ctx)
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// if !deleteIndex.Acknowledged {
	// 	// Not acknowledged
	// }
}
