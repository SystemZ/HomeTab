package model

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
	"log"
)

func GraphInit() *dgo.Dgraph {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	//defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	GraphCleanDb(dg)
	GraphSetSchema(dg)
	return dg
}

func GraphCleanDb(dg *dgo.Dgraph) {
	err := dg.Alter(context.Background(), &api.Operation{DropOp: api.Operation_ALL})
	if err != nil {
		log.Printf("Error clearing DB: %v", err)
	}
}

func GraphSetSchema(dg *dgo.Dgraph) {
	op := &api.Operation{}
	op.Schema = `
	name: string @index(exact) .
	sha256: string .
	tagged: [uid] .
	assigned_to: [uid] .
	

 type File {
   name: string
   sha256
   tagged: [Tag]
 }

 type Tag {
   name: string
   assigned_to: [File]
 }

`

	ctx := context.Background()
	err := dg.Alter(ctx, op)
	if err != nil {
		log.Printf("Graph schema set error: %v", err)
	}
}

func GraphAddFile(dg *dgo.Dgraph, file GraphFile) string {
	file.DType = []string{"File"}
	// don't replace uid if already know it
	if len(file.Uid) < 1 {
		// placeholder for fetching DB ID when inserted
		file.Uid = "_:file"
	}
	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(file)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb
	ctx := context.Background()
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	return assigned.Uids["file"]
}

func GraphAddTag(dg *dgo.Dgraph, file GraphTag) string {
	file.DType = []string{"Tag"}
	// don't replace uid if already know it
	if len(file.Uid) < 1 {
		// placeholder for fetching DB ID when inserted
		file.Uid = "_:tag"
	}
	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(file)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb
	ctx := context.Background()
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	return assigned.Uids["tag"]
}

type GraphFile struct {
	Uid    string     `json:"uid,omitempty"`
	DType  []string   `json:"dgraph.type,omitempty"`
	Tagged []GraphTag `json:"tagged,omitempty"`
	//
	Fid      int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Size     int    `json:"size,omitempty"`
	Sha256   string `json:"sha256,omitempty"`
	Mime     string `json:"mime,omitempty"`
	ParentId int    `json:"parent_id,omitempty"`
}

type GraphTag struct {
	Id      int      `json:"id,omitempty"`
	Fid     int      `json:"fid,omitempty"`
	Name    string   `json:"name,omitempty"`
	Overall int      `json:"overall,omitempty"`
	DType   []string `json:"dgraph.type,omitempty"`
	Uid     string   `json:"uid,omitempty"`
}
