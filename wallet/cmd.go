package main

import (
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"nshmadhani.com/notcoin/blockchain/block"
	"nshmadhani.com/notcoin/blockchain/util"
	WalletCTX "nshmadhani.com/notcoin/blockchain/wallet"
)

var App *grumble.App

func InitShell() {
	App = grumble.New(&grumble.Config{
		Name:                  "Notacoin Wallet",
		Description:           "Multiple Account Wallet for Notacoin",
		HistoryFile:           "/tmp/notacoin_wallet.hist",
		Prompt:                "wallet »»» ",
		PromptColor:           color.New(color.FgGreen, color.Bold),
		HelpHeadlineColor:     color.New(color.FgGreen),
		HelpHeadlineUnderline: true,
		HelpSubCommands:       true,
		Flags: func(f *grumble.Flags) {
			f.String("d", "directory", "DEFAULT", "set an alternative root directory path")
			f.Bool("v", "verbose", false, "enable verbose mode")
		},
	})

	App.AddCommand(conncectToNode())
	App.AddCommand(createWallet())
	App.AddCommand(printWallet())
	App.AddCommand(printwallets())
	App.AddCommand(loadWallet())
	App.AddCommand(updateChain())
	App.AddCommand(newTransaction())
	App.AddCommand(printChainCmd())

}

func printChainCmd() *grumble.Command {
	runCommand := &grumble.Command{
		Name: "printchain",
		Help: "Prints the Blockchain",

		Run: func(c *grumble.Context) error {

			hashs := Manager.Node.Chain.GetBlockchainHashs()
			for _, h := range hashs {
				c.App.Println(hex.EncodeToString(h))
			}
			return nil
		},
	}
	return runCommand
}

func conncectToNode() *grumble.Command {
	cmd := &grumble.Command{
		Name: "connect",
		Help: "Connect to a Node",
		Args: func(a *grumble.Args) {
			a.String("nodeAddr", "Addr for Node")
		},
		Run: func(c *grumble.Context) error {

			nAddr := c.Args.String("nodeAddr")

			if Manager == nil {
				Manager = LoadManager(nAddr)
			} else {
				Manager.PeerAddr = nAddr
				Manager.Updatechain()
			}

			return nil
		},
	}
	return cmd
}

func createWallet() *grumble.Command {
	cmd := &grumble.Command{
		Name: "createwallet",
		Help: "Create a New Wallet",
		Args: func(a *grumble.Args) {
			a.String("alias", "Name for Wallet")
		},
		Run: func(c *grumble.Context) error {

			alias := c.Args.String("alias")

			wallet := WalletCTX.NewWallet(alias)

			Manager.Add(*wallet)

			return nil
		},
	}
	return cmd
}

func printwallets() *grumble.Command {
	cmd := &grumble.Command{
		Name: "printwallets",
		Help: "Prints All Wallets",
		Run: func(c *grumble.Context) error {
			for alias, address := range Manager.Aliases {
				_, _, balance := Manager.Node.Chain.FindUTXOs(util.PubHashFromAddress([]byte(address)), 0)
				fmt.Printf("%s, %s, %d\n", alias, address, balance)
			}
			return nil
		},
	}
	return cmd
}

func printWallet() *grumble.Command {
	cmd := &grumble.Command{
		Name: "printwallet",
		Help: "Print Wallet Details",
		Args: func(a *grumble.Args) {
			a.String("id", "Identifier for wallet")
		},
		Run: func(c *grumble.Context) error {

			id := c.Args.String("id")

			byteAddress := []byte(id)

			var wallet *WalletCTX.Wallet

			if len(id) == 0 {
				fmt.Println("No Account Mentioned")
			} else {

				wallet = Manager.FindWallet(id)
				if wallet != nil {
					byteAddress = []byte(wallet.GetAddress())
				}
			}

			utxoSet := block.UTXO{Manager.Node.Chain}

			_, balance := utxoSet.FindUTXO(util.PubHashFromAddress(byteAddress), 0)

			fmt.Printf("Balance is %d\n", balance)

			return nil
		},
	}
	return cmd
}

func loadWallet() *grumble.Command {
	cmd := &grumble.Command{
		Name: "loadWallet",
		Help: "Load a Wallet from a .wt File",
		Args: func(a *grumble.Args) {
			a.String("file", "Wallet File")
		},
		Run: func(c *grumble.Context) error {

			walletFile := c.Args.String("file")
			walletFile, _ = filepath.Abs(walletFile)

			wallet := WalletFromFile(walletFile)
			Manager.Add(wallet)

			fmt.Println("Wallet Added")

			return nil
		},
	}
	return cmd
}
func updateChain() *grumble.Command {
	cmd := &grumble.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Help:    "Update Blockchain",
		Run: func(c *grumble.Context) error {
			Manager.Updatechain()
			return nil
		},
	}
	return cmd
}

func newTransaction() *grumble.Command {
	cmd := &grumble.Command{
		Name: "transaction",
		Help: "Create a Transaction",
		Args: func(a *grumble.Args) {
			a.String("f", "From Alias")
			a.String("t", "to Address")
			a.Int("a", "Amount to Send")
		},
		Run: func(c *grumble.Context) error {

			from := c.Args.String("f")
			to := c.Args.String("t")
			amount := c.Args.Int("a")

			tx := CreateNewTransaction(*Manager, from, to, amount)

			Manager.SendTX(tx, Manager.PeerAddr)

			Manager.Updatechain()

			fmt.Println("Transaction Has been sent to the node")
			return nil
		},
	}
	return cmd
}
