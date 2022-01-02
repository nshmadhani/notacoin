package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	TransactonCTX "nshmadhani.com/notcoin/blockchain/tx"
)

func StartRPC(port string) {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	fmt.Println("RPC Server is Running :" + port)

	r.GET("/blockchain", func(c *gin.Context) {
		//This will send the Bloclchain db file
		a, _ := filepath.Abs(Node.Chain.DB.Path())
		//fmt.Println(a)
		c.File(a)

	})

	r.POST("/tx/new", func(c *gin.Context) {
		transactionData := c.PostForm("tx")
		_ = []byte(c.PostForm("senderAddress"))

		fmt.Println(c.PostForm("senderAddress"))

		amount, err := strconv.ParseInt(c.PostForm("amount"), 10, 0)

		fmt.Println(amount)

		txBytes, err := hex.DecodeString(transactionData)
		tx := TransactonCTX.DeserializeTX(txBytes)

		if err != nil {
			log.Panic(err)
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		}

		if Node.Chain.VerifyTx(tx) {

			cbtx := TransactonCTX.NewCoinBaseTx(Node.Wallet.GetAddress(), Node.Wallet.Alias, 0)

			Node.Chain.MineBlock([]*TransactonCTX.Transaction{tx, cbtx})

			c.JSON(200, gin.H{
				"message": "Verified",
			})

		} else {
			fmt.Println("Verification Failed")
			c.JSON(400, gin.H{
				"message": "Could not verify",
			})
		}

	})
	go r.Run(":" + port)
}
