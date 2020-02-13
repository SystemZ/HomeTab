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
    path: string .
	sha256: string .
	phash: string .
	size: int .
	tagged: [uid] @reverse .
	assigned_to: [uid] @reverse .
	
 type File {
   name: string
   path: string
   sha256: string
   phash: string
   size: int
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

// FIXME use vars properly
func GraphSetTag(dg *dgo.Dgraph, tagname string, filename string) {
	query := `
	query {
      var(func: eq(name,"` + tagname + `")) {
	    Tag as uid
      }
      var(func: eq(name,"` + filename + `")) {
        File as uid
      }
	}`
	mu1 := &api.Mutation{
		SetNquads: []byte(`uid(Tag) <name> "` + tagname + `" .`),
	}
	mu2 := &api.Mutation{
		SetNquads: []byte(`uid(Tag) <dgraph.type> "Tag" .`),
	}
	mu3 := &api.Mutation{
		SetNquads: []byte(`uid(File) <tagged> uid(Tag) .`),
	}
	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{mu1, mu2, mu3},
		CommitNow: true,
	}

	// Update email only if matching uid found.
	ctx := context.Background()
	if _, err := dg.NewTxn().Do(ctx, req); err != nil {
		log.Fatal(err)
	}
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
	Path     string `json:"path,omitempty"`
	Phash    string `json:"phash,omitempty"`
}

type GraphTag struct {
	Id      int      `json:"id,omitempty"`
	Fid     int      `json:"fid,omitempty"`
	Name    string   `json:"name,omitempty"`
	Overall int      `json:"overall,omitempty"`
	DType   []string `json:"dgraph.type,omitempty"`
	Uid     string   `json:"uid,omitempty"`
}
