package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"nshmadhani.com/notcoin/blockchain/block"
	"nshmadhani.com/notcoin/blockchain/util"
)

var Node *block.Node
var App *grumble.App

func init() {
	App = grumble.New(&grumble.Config{
		Name:                  "Notacoin Node",
		Description:           "Runs a Full Node for Notacoin",
		HistoryFile:           "/tmp/notacoing_node.hist",
		Prompt:                "node »»» ",
		PromptColor:           color.New(color.FgGreen, color.Bold),
		HelpHeadlineColor:     color.New(color.FgGreen),
		HelpHeadlineUnderline: true,
		HelpSubCommands:       true,

		Flags: func(f *grumble.Flags) {
			f.String("d", "directory", "DEFAULT", "set an alternative root directory path")
			f.Bool("v", "verbose", false, "enable verbose mode")
		},
	})

	App.SetPrintASCIILogo(func(a *grumble.App) {
		a.Println("        		       $$  								  $$  		  ")
		a.Println("					   $$       									  ")
		a.Println("	$$$$$$$  $$$$$$  $$$$$$    $$$$$$$   $$$$$$$  $$$$$$  $$  $$$$$$$ ")
		a.Println("	$$   $$ $$    $$   $$      $$       $$       $$    $$ $$  $$   $$ ")
		a.Println("	$$   $$ $$    $$   $$      $$$$$$$  $$ 		 $$    $$ $$  $$   $$ ")
		a.Println("	$$   $$ $$    $$   $$  $$  $$   $$  $$       $$    $$ $$  $$   $$ ")
		a.Println("	$$   $$  $$$$$$     $$$$   $$$$$$$   $$$$$$$  $$$$$$  $$  $$   $$ ")

	})
	RunCMD()
	AddPeerCmd()
	RunRPC()
	ReindexUTXO()
	PrintChainCmd()
	ListPeersCmd()
	whoamiCmd()
}

func PrintChainCmd() {
	runCommand := &grumble.Command{
		Name: "printchain",
		Help: "Run the Notcoin Node",

		Run: func(c *grumble.Context) error {

			if len(Node.NetworkAddr) == 0 {
				c.App.PrintError(errors.New("Init the Node First"))
			} else {

				hashs := Node.Chain.GetBlockchainHashs()
				for _, h := range hashs {
					c.App.Println(hex.EncodeToString(h))
				}
			}
			return nil
		},
	}
	App.AddCommand(runCommand)
}

func RunCMD() {
	runCommand := &grumble.Command{
		Name: "run",
		Help: "Run the Notcoin Node",
		Args: func(a *grumble.Args) {
			a.String("id", "NodeID")
			a.String("db", "blockchaindb file")
			a.String("w", "Wallet File")
		},
		Run: func(c *grumble.Context) error {

			Node = block.NewNode(c.Args.String("id"),
				c.Args.String("db"),
				c.Args.String("w"))

			go Node.Run()

			return nil
		},
	}
	App.AddCommand(runCommand)
}

func AddPeerCmd() {
	runCommand := &grumble.Command{
		Name: "addpeer",
		Help: "Add a New Peer",
		Args: func(a *grumble.Args) {
			a.String("pAddr", "Peer Address")
		},
		Run: func(c *grumble.Context) error {

			Node.AddPeer((c.Args.String("pAddr")))

			return nil
		},
	}
	App.AddCommand(runCommand)
}

func ListPeersCmd() {
	runCommand := &grumble.Command{
		Name: "peers",
		Help: "Shows all Peers",
		Run: func(c *grumble.Context) error {

			for _, peer := range Node.Peers {
				fmt.Println(peer)
			}
			return nil
		},
	}
	App.AddCommand(runCommand)
}

//****************************************RPC*****************************

func RunRPC() {
	runCommand := &grumble.Command{
		Name: "run-rpc",
		Help: "Runs the RPC Server",
		Args: func(a *grumble.Args) {
			a.String("port", "Port for RPC Server")
		},
		Run: func(c *grumble.Context) error {

			if len(Node.NetworkAddr) == 0 {
				return errors.New("Please Start the node first")
			}
			port := c.Args.String("port")
			if strings.Trim(port, " ") == "" {
				port = "3000"
			}

			StartRPC(port)

			return nil
		},
	}
	App.AddCommand(runCommand)
}

/*************************************************************************************************
*											BLOCKCHAIN COMMANDS
**************************************************************************************************/

func ReindexUTXO() {

	runCommand := &grumble.Command{
		Name: "reindex-utxo",
		Help: "Reindexes the UTXOSet",
		Args: func(a *grumble.Args) {
		},
		Run: func(c *grumble.Context) error {

			utxoSet := block.NewUTXO(Node.Chain)
			utxoSet.Reindex()
			return nil
		},
	}
	App.AddCommand(runCommand)
}

func MempoolCmd() {

	runCommand := &grumble.Command{
		Name: "mempool",
		Help: "Prints TxID of all transactions in Mempool",
		Args: func(a *grumble.Args) {
		},
		Run: func(c *grumble.Context) error {

			mempoolTx := Node.GetMempool()

			for txId := range mempoolTx {
				fmt.Println(txId)
			}

			return nil
		},
	}
	App.AddCommand(runCommand)
}

/**************************************************************************************************
									INFO COMMANDS
***************************************************************************************************/

func whoamiCmd() {

	runCommand := &grumble.Command{
		Name: "whoami",
		Help: "Print Information about the Current Node",
		Args: func(a *grumble.Args) {
		},
		Run: func(c *grumble.Context) error {

			if Node == nil {
				return errors.New("Please start server first")
			}

			fmt.Printf("Address is %s\n", Node.NetworkAddr)

			utxoSet := block.NewUTXO(Node.Chain)
			_, balance := utxoSet.FindUTXO(util.PubHashFromAddress([]byte(Node.Wallet.GetAddress())), 0)
			fmt.Printf("%s, %s, %d\n", Node.Wallet.Alias, Node.Wallet.GetAddress(), balance)

			return nil
		},
	}
	App.AddCommand(runCommand)
}
