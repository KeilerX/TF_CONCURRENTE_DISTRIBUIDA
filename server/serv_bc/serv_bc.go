package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

type tmsg struct {
	Code string
	Addr string
	Blk  Block
	Bc   BlockChain
}

type nodeResult struct {
	Addr  string
	Valid int
}

// Las IP de los demás participantes acá, todos deberían usar el puerto 8000
var addrs = []string{"localhost:8010", "localhost:8011", "localhost:8012"}

var chInfo chan map[string]Block
var chServAns chan int
var chProcess chan bool
var chValidNode chan nodeResult
var localAddr string
var chain *BlockChain
var consUp bool
var mainSv bool

func main() {

	var lport string
	fmt.Print("Local Port: ")
	fmt.Scanf("%s\n", &lport)
	localAddr = "localhost:" + lport
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

	chain = InitBlockChain()
	consUp = false
	mainSv = false

	chServAns = make(chan int, 1)
	chValidNode = make(chan nodeResult, 1)
	chProcess = make(chan bool, 1)
	chInfo = make(chan map[string]Block)

	chProcess <- true
	go func() { chInfo <- map[string]Block{} }()
	server()

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
			consUp = true
			mainSv = true
			chain.AddBlock(string(msg.Blk.Data))
			newmsg := tmsg{"consenso", localAddr, *chain.Blocks[len(chain.Blocks)-1], BlockChain{}}
			fmt.Println("Sending this: ", *chain.Blocks[len(chain.Blocks)-1])

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
			if !consUp {
				consUp = true
				chain.AddBlock(string(msg.Blk.Data))
				newmsg := tmsg{"consenso", localAddr, *chain.Blocks[len(chain.Blocks)-1], BlockChain{}}

				for _, addr := range addrs {
					sendBlock(addr, newmsg)
				}
			}
			concensus(conn, msg)

		case "valid":
			chValidNode <- nodeResult{msg.Addr, 1}
		case "invalid":
			chValidNode <- nodeResult{msg.Addr, 0}
		case "getBackup":
			newmsg := tmsg{"setBackup", localAddr, Block{}, *chain}
			sendBlock(msg.Addr, newmsg)
		case "setBackup":
			*chain = msg.Bc
			if mainSv {
				chServAns <- 1
			}
			for _, block := range chain.Blocks {
				fmt.Printf("Previous Hash: %x\n", block.PrevHash)
				fmt.Printf("Data in Block: %s\n", block.Data)
				fmt.Printf("Hash: %x\n", block.Hash)
			}
			consUp = false
		}
	}
}

func concensus(conn net.Conn, msg tmsg) {
	info := <-chInfo
	if msg.Addr != "" {
		info[msg.Addr] = msg.Blk
	}
	if len(info) == len(addrs) {
		eqs := 0
		myblk := *chain.Blocks[len(chain.Blocks)-1]
		for _, iblk := range info {
			if string(myblk.Hash) == string(iblk.Hash) {
				eqs++
			}
		}
		if eqs > (len(addrs) - eqs) {
			fmt.Println("VALID!")
			for _, addr := range addrs {
				send("valid", addr)
			}
			if mainSv {
				chServAns <- 1
			}
			for _, block := range chain.Blocks {
				fmt.Printf("Previous Hash: %x\n", block.PrevHash)
				fmt.Printf("Data in Block: %s\n", block.Data)
				fmt.Printf("Hash: %x\n", block.Hash)
			}
			consUp = false
		} else {
			fmt.Println("INVALID!")
			for _, addr := range addrs {
				send("invalid", addr)
			}
			var validNode nodeResult
			noValids := true
			for i := 0; i < len(addrs); i++ {
				validNode = <-chValidNode
				if validNode.Valid == 1 {
					noValids = false
					send("getBackup", validNode.Addr)
					break
				}
			}
			if noValids {
				chServAns <- 0
			}
		}
		info = map[string]Block{}
	}
	go func() { chInfo <- info }()
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
	msg := tmsg{code, localAddr, Block{}, BlockChain{}}
	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("Can't dial", remoteAddr)
	} else {
		defer conn.Close()
		fmt.Println("Sending to", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}
