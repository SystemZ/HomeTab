package model

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
	"log"
	"strconv"
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
	size: int .
	sha256: string .
	phash: string .
	similar: [uid] @reverse .

	tagged: [uid] @reverse .
	assigned_to: [uid] @reverse .

	mime: [uid] @reverse .
	mime_assigned_to: [uid] @reverse .
	
 type File {
   name: string
   path: string
   size: int
   sha256: string
   phash: string
   similar: [File]
   tagged: [Tag]
   mime: [Mime]
 }

 type Tag {
   name: string
   assigned_to: [File]
 }

 type Mime {
   name: string
   mime_assigned_to: [File]
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
		//SetNquads: []byte(`_:file <testz> "eeeez" .`),
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

// FIXME use vars properly
func GraphSetMime(dg *dgo.Dgraph, mime string, filename string) {
	query := `
	query {
      var(func: eq(name,"` + mime + `")) {
	    Mime as uid
      }
      var(func: eq(name,"` + filename + `")) {
        File as uid
      }
	}`
	mu1 := &api.Mutation{
		SetNquads: []byte(`uid(Mime) <name> "` + mime + `" .`),
	}
	mu2 := &api.Mutation{
		SetNquads: []byte(`uid(Mime) <dgraph.type> "Mime" .`),
	}
	mu3 := &api.Mutation{
		SetNquads: []byte(`uid(File) <mime> uid(Mime) .`),
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

// FIXME use vars properly
func GraphSetDistance(txn *dgo.Txn, filename1 string, filename2 string, distance int) {
	distanceStr := strconv.Itoa(distance)
	query := `
	query {
      var(func: eq(name,"` + filename1 + `")) {
        File1 as uid
      }
	  var(func: eq(name,"` + filename2 + `")) {
        File2 as uid
      }
	}`
	mu := &api.Mutation{
		SetNquads: []byte(`uid(File1) <similar> uid(File2) (distance=` + distanceStr + `) .`),
	}
	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{mu},
		//CommitNow: true,
	}

	ctx := context.Background()
	txn.Do(ctx, req)
	// Update email only if matching uid found.
	//ctx := context.Background()
	//if _, err := dg.NewTxn().Do(ctx, req); err != nil {
	//	log.Fatal(err)
	//}
}

func GraphSearchPhash(dg *dgo.Dgraph) []GraphFile {
	const q = `
	{
	  Files(func: type("File")) @filter(has(phash)) {
	    uid
		name
		phash
	  }
	}
	`
	resp, err := dg.NewTxn().Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}

	var decode struct {
		Files []GraphFile
	}
	if err := json.Unmarshal(resp.GetJson(), &decode); err != nil {
		log.Fatal(err)
	}
	return decode.Files
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
