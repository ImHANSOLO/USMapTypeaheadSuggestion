package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/graphql-go/graphql"
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
    uri := os.Getenv("MONGO_URI") 
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal(err)
    }
    if err := client.Ping(ctx, readpref.Primary()); err != nil {
        log.Fatal(err)
    }
    coll := client.Database("statesdb").Collection("states")

    stateType := graphql.NewObject(graphql.ObjectConfig{
        Name: "State",
        Fields: graphql.Fields{
            "name": &graphql.Field{Type: graphql.String},
            "code": &graphql.Field{Type: graphql.String},
        },
    })
    schema, err := graphql.NewSchema(graphql.SchemaConfig{
        Query: graphql.NewObject(graphql.ObjectConfig{
            Name: "RootQuery",
            Fields: graphql.Fields{
                "states": &graphql.Field{
                    Type: graphql.NewList(stateType),
                    Args: graphql.FieldConfigArgument{
                        "q": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
                    },
                    Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                        q := p.Args["q"].(string)
                        filter := bson.M{"name": bson.M{"$regex": "^" + strings.ToLower(q), "$options": "i"}}
                        cur, err := coll.Find(ctx, filter)
                        if err != nil {
                            return nil, err
                        }
                        var out []State
                        _ = cur.All(ctx, &out)
                        return out, nil
                    },
                },
            },
        }),
    })
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
        var req struct{ Query string }
        _ = json.NewDecoder(r.Body).Decode(&req)
        res := graphql.Do(graphql.Params{Schema: schema, RequestString: req.Query})
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(res)
    })

    log.Println("Backend listening on :4000/graphql")
    log.Fatal(http.ListenAndServe(":4000", nil))
}
