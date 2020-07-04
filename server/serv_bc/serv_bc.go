package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

type FbDoc struct {
	Data      string    `json:"Data"`
	PrevHash  []byte    `json:"PrevHash"`
	Hash      []byte    `json:"Hash"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func deleteCollection(ctx context.Context, client *firestore.Client,
	ref *firestore.CollectionRef, batchSize int) error {

	for {
		// Get a batch of documents
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha1.Sum(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func AddBlock(new Block) {

	_, _, err = client.Collection("pacientes"+lport).Add(ctx, map[string]interface{}{
		"Data":      string(new.Data),
		"PrevHash":  new.PrevHash,
		"Hash":      new.Hash,
		"CreatedAt": time.Now(),
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}

func InitChain() {
	chain = CreateBlock("Genesis", []byte{})
	_, _, err = client.Collection("pacientes"+lport).Add(ctx, map[string]interface{}{
		"Data":      "Genesis",
		"PrevHash":  chain.PrevHash,
		"Hash":      chain.Hash,
		"CreatedAt": time.Now().UTC(),
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}

}

func getBackup(addr string) {

	deleteCollection(ctx, client, client.Collection("pacientes"+lport), 10)
	fmt.Println("pacientes" + addr[10:])
	iter := client.Collection("pacientes" + addr[10:]).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		_, _, err2 = client.Collection("pacientes"+lport).Add(ctx, doc.Data())
		if err2 != nil {
			log.Fatalf("Failed adding alovelace: %v", err2)
		}
	}
	getChain()
}

type tmsg struct {
	Code string
	Addr string
	Blk  Block
}

type nodeResult struct {
	Addr  string
	Valid int
}

// Las IP de los demás participantes acá, todos deberían usar el puerto 8000
var addrs = []string{"localhost:8010", "localhost:8011", "localhost:8012", "localhost:8013"}

var chInfo chan map[string]Block
var chServAns chan int
var chProcess chan bool
var chValidNode chan nodeResult
var localAddr string
var lport string
var chain *Block
var newBlk *Block
var consUp bool
var mainSv bool
var valids int
var invalids int
var ctx = context.Background()
var sa = option.WithCredentialsFile("../go-react-pcd-firebase-adminsdk-dedz3-72cae0df1b.json")
var app, err = firebase.NewApp(ctx, nil, sa)
var client, err2 = app.Firestore(ctx)

func main() {

	fmt.Print("Local Port: ")
	fmt.Scanf("%s\n", &lport)
	localAddr = "localhost:" + lport

	// Use a service account
	if err != nil {
		log.Fatalln(err)
	}
	if err2 != nil {
		log.Fatalln(err2)
	}
	defer client.Close()

	/*
		for {
			var tmpPort string
			fmt.Print("External Port: ")
			fmt.Scanf("%s\n", &tmpPort)
			if tmpPort != "-1" {
				addrs = append(addrs, "localhost:"+tmpPort)
			} else {
				break
			}
		}
	*/
	for i, ad := range addrs {
		if ad == localAddr {
			addrs[i] = addrs[len(addrs)-1]
			addrs[len(addrs)-1] = ""
			addrs = addrs[:len(addrs)-1]
			break
		}
	}

	consUp = false
	mainSv = false
	getChain()

	chServAns = make(chan int, 1)
	chValidNode = make(chan nodeResult, 1)
	chProcess = make(chan bool, 1)
	chInfo = make(chan map[string]Block, len(addrs))

	chProcess <- true
	go func() { chInfo <- map[string]Block{} }()
	server()

}
func getChain() {
	iter := client.Collection("pacientes"+lport).OrderBy("CreatedAt", firestore.Desc).Limit(1).Documents(ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		InitChain()
	} else {
		jsonbody, err2 := json.Marshal(doc.Data())
		if err2 != nil {
			// do error check
			fmt.Println(err2)
			return
		}
		tempdoc := FbDoc{}
		if err2 := json.Unmarshal(jsonbody, &tempdoc); err2 != nil {
			// do error check
			fmt.Println(err2)
			return
		}
		chain = &Block{
			Data:     []byte(tempdoc.Data),
			Hash:     []byte(tempdoc.Hash),
			PrevHash: []byte(tempdoc.PrevHash),
		}
	}
}
func server() {
	if ln, err := net.Listen("tcp", localAddr); err != nil {
		log.Panicln("Can't start listener on", localAddr)
	} else {
		defer ln.Close()
		fmt.Println("Listeing on", localAddr)
		for {
			if conn, err := ln.Accept(); err != nil {
				log.Println("Can't accept", conn.RemoteAddr())
			} else {
				go handle(conn)
			}
		}
	}
}
func handle(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var msg tmsg
	if err := dec.Decode(&msg); err != nil {
		log.Println("Can't decode from", conn.RemoteAddr())
	} else {
		fmt.Println(msg)
		switch msg.Code {
		case "serv":
			<-chProcess

			chValidNode = make(chan nodeResult, 1)
			valids = 0
			invalids = 0
			consUp = true
			mainSv = true
			newBlk = CreateBlock(string(msg.Blk.Data), chain.Hash)
			newmsg := tmsg{"consenso", localAddr, *newBlk}
			fmt.Println("Sending this: ", *chain)

			for _, addr := range addrs {
				sendBlock(addr, newmsg)
			}

			concensus(conn, tmsg{})
			svAns := <-chServAns
			conn2, err2 := net.Dial("tcp", msg.Addr)
			if err2 == nil {
				defer conn2.Close()
				fmt.Fprint(conn2, strconv.Itoa(svAns)+"\n")
			}
			chProcess <- true

		case "consenso":
			fmt.Println(consUp)
			if !consUp {
				valids = 0
				invalids = 0
				consUp = true
				chValidNode = make(chan nodeResult, 1)
				newBlk = CreateBlock(string(msg.Blk.Data), chain.Hash)
				newmsg := tmsg{"consenso", localAddr, *newBlk}

				for _, addr := range addrs {
					sendBlock(addr, newmsg)
				}
			}
			concensus(conn, msg)

		case "valid":
			if (valids + invalids) == len(addrs) {
				close(chValidNode)
				chValidNode = make(chan nodeResult, 1)
			}
			chValidNode <- nodeResult{msg.Addr, 1}
			valids++
		case "invalid":
			if (valids + invalids) == len(addrs) {
				close(chValidNode)
				chValidNode = make(chan nodeResult, 1)
			}
			chValidNode <- nodeResult{msg.Addr, 0}
			invalids++
		}
	}
}

func concensus(conn net.Conn, msg tmsg) {
	fmt.Println("Inicio consenso")
	info := <-chInfo
	fmt.Println(msg.Addr)
	if msg.Addr != "" {
		info[msg.Addr] = msg.Blk
	}
	if len(info) == len(addrs) {
		eqs := 0
		myblk := *newBlk
		for _, iblk := range info {
			if string(myblk.Hash) == string(iblk.Hash) {
				eqs++
			}
		}
		if eqs > (len(addrs) - eqs) {
			fmt.Println("VALID!")
			AddBlock(*newBlk)

			for _, addr := range addrs {
				go send("valid", addr)
			}
			if mainSv {
				chServAns <- 1
			}
			fmt.Println("PREV BLOCK:")
			fmt.Printf("Previous Hash: %x\n", chain.PrevHash)
			fmt.Printf("Data in Block: %s\n", chain.Data)
			fmt.Printf("Hash: %x\n", chain.Hash)
			fmt.Println("NEW BLOCK:")
			fmt.Printf("Previous Hash: %x\n", newBlk.PrevHash)
			fmt.Printf("Data in Block: %s\n", newBlk.Data)
			fmt.Printf("Hash: %x\n", newBlk.Hash)
			chain = newBlk
			consUp = false
		} else {
			fmt.Println("INVALID!")
			for _, addr := range addrs {
				go send("invalid", addr)
			}
			var validNode nodeResult
			noValids := true

			for i := 0; i < len(addrs); i++ {

				validNode = <-chValidNode
				if validNode.Valid == 1 {
					noValids = false
					getBackup(validNode.Addr)
					consUp = false
					chServAns <- 1
					break
				}
			}
			if noValids {
				consUp = false
				chServAns <- 0
			}
		}
		info = map[string]Block{}
	}
	chInfo <- info
}
func sendBlock(remoteAddr string, msg tmsg) {
	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("Can't dial", remoteAddr)
	} else {
		defer conn.Close()
		fmt.Println("Sending to", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}
func send(code string, remoteAddr string) {
	msg := tmsg{code, localAddr, Block{}}
	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("Can't dial", remoteAddr)
	} else {
		defer conn.Close()
		fmt.Println("Sending to", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}
