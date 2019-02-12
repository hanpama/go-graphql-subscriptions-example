package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"github.com/hanpama/graphql-with-go/playground"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type root struct {
}

func (*root) Hello() string {
	return "Hello, world!"
}

func (rv *root) Time(ctx context.Context) chan string {
	clock := make(chan string)
	count := 0

	go func() {
		defer close(clock)
		for {
			if count > 10 {
				return
			}
			count++

			time.Sleep(1 * time.Second)
			nowString := time.Now().Format(time.UnixDate)

			select {
			case <-ctx.Done():
				println("Done!")
				return

			default:
				clock <- nowString
				println(nowString)
			}
		}
	}()
	return clock
}

func main() {
	s := `
		schema {
			query: Query
			subscription: Subscription
		}
		type Query {
			hello: String!
		}
		type Subscription {
			time: String!
		}
	`

	schema := graphql.MustParseSchema(s, &root{})
	http.HandleFunc("/graphql", graphqlws.NewHandlerFunc(
		schema, &relay.Handler{Schema: schema},
	))
	http.HandleFunc("/playground", playground.HandlePlayground)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
