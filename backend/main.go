package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type State struct {
    Name string `bson:"name" json:"name"`
    Code string `bson:"code" json:"code"`
}

func main() {
    // 1. 连接 MongoDB
    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.Background())
    if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
        log.Fatal(err)
    }
    coll := client.Database("statesdb").Collection("states")

    // 2. 定义 GraphQL Schema
    stateType := graphql.NewObject(graphql.ObjectConfig{
        Name: "State",
        Fields: graphql.Fields{
            "name": &graphql.Field{Type: graphql.String},
            "code": &graphql.Field{Type: graphql.String},
        },
    })
    rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
        "states": &graphql.Field{
            Type: graphql.NewList(stateType),
            Args: graphql.FieldConfigArgument{
                "q": &graphql.ArgumentConfig{Type: graphql.String},
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                filter := bson.M{}
                if q, ok := p.Args["q"].(string); ok && q != "" {
                    filter = bson.M{"name": bson.M{"$regex": q, "$options": "i"}}
                }
                cur, err := coll.Find(p.Context, filter)
                if err != nil {
                    return nil, err
                }
                var result []State
                if err := cur.All(p.Context, &result); err != nil {
                    return nil, err
                }
                return result, nil
            },
        },
    }}
    schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)})
    if err != nil {
        log.Fatal(err)
    }

    // 3. 包装 HTTP Handler
    graphqlHandler := handler.New(&handler.Config{
        Schema:   &schema,
        Pretty:   true,
        GraphiQL: true,
    })
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type"},
        AllowCredentials: true,
    })
    http.Handle("/graphql", c.Handler(graphqlHandler))

    // 4. 启动
    port := os.Getenv("PORT")
    if port == "" {
        port = "4000"
    }
    log.Println("Backend listening on :" + port + "/graphql")
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
