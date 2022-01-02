package main

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"math"

// 	cli "github.com/urfave/cli/v2"
// )

// type CLI struct {
// 	App *cli.App
// }

// func CreateNewCli() *CLI {

// 	app := &cli.App{
// 		Name:  "A completely Orginal DS using back-linked hashes",
// 		Usage: "CRUD with a Blockchain",
// 		Commands: []*cli.Command{
// 			// {
// 			// 	Name:    "add",
// 			// 	Aliases: []string{"a"},
// 			// 	Usage:   "add a Tx to the Blockahin",
// 			// 	Action: func(c *cli.Context) error {
// 			// 		data := c.Args().First()
// 			// 		bc.AddBlock(data)
// 			// 		return nil
// 			// 	},
// 			// },
// 			{
// 				Name:    "createblockchain",
// 				Aliases: []string{"cbc"},
// 				Usage:   "Create a New Blockchain",
// 				Action: func(c *cli.Context) error {
// 					_ = CreateBlockchain(c.Args().First())

// 					return nil
// 				},
// 			},
// 			{
// 				Name:    "printchain",
// 				Aliases: []string{"p"},
// 				Usage:   "Print curent chain",
// 				Action: func(c *cli.Context) error {
// 					Printchain()
// 					return nil
// 				},
// 			},
// 			{
// 				Name:    "balance",
// 				Aliases: []string{"b"},
// 				Usage:   "Balance of the User",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{Name: "address", Aliases: []string{"a"}},
// 				},
// 				Action: func(c *cli.Context) error {

// 					address := c.String("address")
// 					balance := GetBalance(NewBlockchain(), address)
// 					fmt.Printf("Balance of %s is %d\n", address, balance)
// 					return nil
// 				},
// 			},
// 			{
// 				Name:    "send",
// 				Aliases: []string{"s"},
// 				Usage:   "Send from one account to another",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{Name: "from", Aliases: []string{"f"}},
// 					&cli.StringFlag{Name: "to", Aliases: []string{"t"}},
// 					&cli.IntFlag{Name: "amount", Aliases: []string{"a"}},
// 				},
// 				Action: func(c *cli.Context) error {

// 					from := c.String("from")

// 					to := c.String("to")

// 					amount := c.Int("amount")

// 					Send(NewBlockchain(), amount, from, to)

// 					return nil
// 				},
// 			},
// 		},
// 	}
// 	return &CLI{app}
// }

// func Printchain() {
// 	// TODO: Fix this
// 	bc := NewBlockchain()

// 	bci := bc.Iterator()

// 	for {
// 		block := bci.Next()

// 		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
// 		fmt.Printf("Hash: %x\n", block.Hash)
// 		//pow := NewProofOfWork(block)
// 		//fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

// 		fmt.Println("Transcations")

// 		for idx, tx := range block.Transactions {

// 			fmt.Printf("	Transaction No %d with ID = %s\n", idx, hex.EncodeToString(tx.ID))

// 			for inIdx, in := range tx.Vin {

// 				fmt.Printf("	Input Transcation %d with txId=%x,Vout=%d,script=%s\n", inIdx, in.TxId, in.Vout, in.ScriptSig)
// 			}
// 			for outIdx, in := range tx.Vout {

// 				fmt.Printf("	Output Transcation %d with Vout=%d,script=%s\n", outIdx, in.Value, in.ScriptPubKey)
// 			}

// 		}

// 		fmt.Println()

// 		if len(block.PrevHash) == 0 {
// 			break
// 		}
// 	}
// }

// func GetBalance(bc *Blockchain, address string) int {
// 	_, _, acc := bc.FindUTXOs(address, math.MaxInt32)
// 	return acc
// }

// func Send(bc *Blockchain, amount int, from, to string) {
// 	tx := bc.NewTransaction(from, to, amount)
// 	bc.AddBlock([]*Transaction{tx})
// 	fmt.Println("Money sent")
// }
