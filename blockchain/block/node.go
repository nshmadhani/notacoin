package block

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"nshmadhani.com/notcoin/blockchain/network"
	TransactionCTX "nshmadhani.com/notcoin/blockchain/tx"
	"nshmadhani.com/notcoin/blockchain/wallet"
)

type Node struct {
	Chain       *Blockchain
	Server      network.Server
	Peers       []string
	Wallet      wallet.Wallet
	NetworkAddr string
	blocksToAdd *BlockQueue

	mempool map[string]*TransactionCTX.Transaction

	mu *sync.Mutex
}

func (n *Node) Run() {

	n.Server.Start(func(c *net.TCPConn) {

		n.mu.Lock()

		defer n.mu.Unlock()

		request, _ := ioutil.ReadAll(c)

		commandType := string(request[0:10])

		commandType = strings.Trim(commandType, "0")

		fmt.Println("Received " + commandType)

		c.Close()

		if commandType == "Version" {
			versionCommand := network.DeserializeVersion(request[10:])
			n.Version(versionCommand)
		} else if commandType == "Verack" {
			//FIXME: What to do with me???
			//verackCommand := network.DeserializeVerack(request[10:])
			//n.Verack(verackCommand)
		} else if commandType == "GetBlock" {
			getBlockCommand := network.DeserializeGetBlock(request[10:])
			n.GetBlockCmd(getBlockCommand)
		} else if commandType == "Inv" {

			invCommnd := network.DeserializeInv(request[10:])
			n.Inv(invCommnd)

		} else if commandType == "GetData" {

			gdCommnd := network.DeserializeGetData(request[10:])
			n.GetData(gdCommnd)

		} else if commandType == "Block" {
			blockCommand := network.DeserializeBlock(request[10:])
			n.BlockCMD(blockCommand)
		} else if commandType == "TX" {
			txCommand := network.DeserializeTX(request[10:])
			n.TXCommand(txCommand)
		} else if commandType == "Peers" {
			peersCommand := network.DeserializePeers(request[10:])
			n.PeersCommand(peersCommand)
		} else if commandType == "ReqPeers" {
			reqPeerCommand := network.DeserializeReqPeers(request[10:])
			n.ReqPeersCommand(reqPeerCommand)
		}

	})

}

/******************************************************************************************
*									HANDLE COMMANDS
******************************************************************************************/

func (n *Node) Version(version network.Version) {

	fmt.Println("Received Version from ", version.AddrFrom)

	bestHeight := n.Chain.GetBestHeight()
	peerBestHeight := version.BestHeight

	peerAddr := version.AddrFrom

	if !n.isPeer(peerAddr) {
		n.addPeer(peerAddr)
	}

	if bestHeight > peerBestHeight {
		n.SendVersion(peerAddr)
	} else if bestHeight < peerBestHeight {
		n.SendGetBlock(peerAddr)
	} else {
		fmt.Println("All Matched with", version.AddrFrom)
	}
	n.SendReqPeers(peerAddr)
}

func (n *Node) GetBlockCmd(geblock network.GetBlock) {
	n.SendInv(geblock.AddrFrom)
}

func (n *Node) Inv(invCommnd network.Inv) {

	//peerAddr := invCommnd.AddrFrom

	if invCommnd.Type == "block" {

		for _, hash := range invCommnd.Items {

			block := n.Chain.GetBlock(hash)
			if len(block.Hash) == 0 {
				n.SendGetData(hash, "block", invCommnd.AddrFrom)
			}

		}

	} else if invCommnd.Type == "tx" {

		for _, txID := range invCommnd.Items {
			if !n.inMempool(hex.EncodeToString(txID)) {
				n.SendGetData(txID, "tx", invCommnd.AddrFrom)
			}
		}

	}

}

func (n *Node) GetData(gd network.GetData) {

	fmt.Print("Received GetData for ")
	if gd.Type == "block" {

		blockHash := gd.Item

		fmt.Println("block:", hex.EncodeToString(blockHash), "from", gd.AddrFrom)

		block := n.Chain.GetBlock(blockHash)

		if len(block.Hash) != 0 {
			n.SendBlock(block.Serialize(), gd.AddrFrom)
		} else {
			fmt.Println("Block Not Found")
		}

	} else if gd.Type == "tx" {
		txID := hex.EncodeToString(gd.Item)
		if n.inMempool(txID) {
			n.SendTxCmd(txID, gd.AddrFrom)
		}
	}

}

func (n *Node) BlockCMD(blockCmd network.Block) {
	block := DeserializeBlock(blockCmd.Item)

	if block == nil {
		fmt.Println("Malformed Request")
		return
	}

	myHeight := n.Chain.GetBestHeight()

	//Reject a Block
	if myHeight >= block.Height {
		foundBlock := n.Chain.GetBlock(block.Hash)
		if len(foundBlock.Hash) != 0 { //1. Block Already there in the Chain
			fmt.Println("Already have the Block", hex.EncodeToString(block.Hash))
			return
		} else { //2. Fork in the chain
			log.Panic(fmt.Errorf("New Fork at Height:%d , Hash=%s", block.Height, hex.EncodeToString(block.Hash)))
		}

	} else { //Accept a block

		n.blocksToAdd.PushBlock(block)

		for {
			rootBlock := n.blocksToAdd.Peek()
			if len(rootBlock.Hash) == 0 {
				fmt.Println("	Peek was nil")
				break
			}
			//ERROR: There can be at times that the process is finished but the new or top nodes are not added
			if rootBlock.Height == n.Chain.GetBestHeight()+1 {
				fmt.Println("Added Block", hex.EncodeToString(rootBlock.Hash))
				blks := n.blocksToAdd.PopBlock()
				n.Chain.AddBlock(blks)
			} else {
				fmt.Println("Height Not Matched", rootBlock.Height, n.Chain.GetBestHeight(), blockCmd.AddrFrom, rootBlock.Hash)
				break
			}
		}
	}
}

func (n *Node) TXCommand(txCMD network.TX) {

	tx := TransactionCTX.DeserializeTX(txCMD.Transaction)

	if n.inMempool(hex.EncodeToString(tx.ID)) {
		fmt.Println("Already have the Tx=" + hex.EncodeToString(tx.ID))
		return
	}
	txId := hex.EncodeToString(tx.ID)
	n.mempool[txId] = tx

	fmt.Printf("Verification is: %s\n", n.Chain.VerifyTx(tx))

	for _, peer := range n.Peers {
		if peer != n.NetworkAddr && peer != txCMD.AddrFrom {
			n.SendTxInv(txId, peer)
		}
	}

	if len(n.mempool) > 0 {

	MineTxs:

		var txs []*TransactionCTX.Transaction

		for id := range n.mempool {
			if n.Chain.VerifyTx(n.mempool[id]) {
				txs = append(txs, n.mempool[id])
			}
		}

		if len(txs) == 0 {
			fmt.Println("All Mempool Tx are invalid")
			return
		}

		txs = append(txs, TransactionCTX.NewCoinBaseTx(n.Wallet.GetAddress(), "HEYA", 0))

		n.Chain.MineBlock(txs)

		for _, tx := range txs {
			txID := hex.EncodeToString(tx.ID)
			delete(n.mempool, txID)
		}

		for _, peer := range n.Peers {
			if peer != n.NetworkAddr {
				n.SendInv(peer)
			}
		}

		if len(n.mempool) > 0 {
			goto MineTxs
		}

	}

}
func (n *Node) PeersCommand(peersCMD network.Peers) {
	extenalPeerList := peersCMD.Peers
	for idx, _ := range extenalPeerList {
		peerAddr := extenalPeerList[idx]
		if n.isPeer(peerAddr) {
			continue
		}
		n.addPeer(peerAddr)
		if peerAddr != peersCMD.AddrFrom {
			n.AddPeer(peerAddr)
		}
	}
}
func (n *Node) ReqPeersCommand(reqpeersCMD network.ReqPeers) {
	n.SendPeers(reqpeersCMD.AddrFrom) //Send ma Peers
	if !n.isPeer(reqpeersCMD.AddrFrom) {
		n.SendReqPeers(reqpeersCMD.AddrFrom) //Ask for his peers
	}
}

/******************************************************************************************
*									SEND COMMANDS
*******************************************************************************************/

func SendCommand(peerAddr string, command []byte) {
	client := network.NewClient(peerAddr)
	client.Send(command)
	fmt.Println("Sent Command", string(command[:10]), "to", peerAddr)
}

func (n *Node) SendVersion(peerAddr string) {
	version := network.Version{1, n.Chain.GetBestHeight(), n.NetworkAddr}
	SendCommand(peerAddr, version.Serialize())
}
func (n Node) SendGetBlock(peerAddr string) {
	getBlock := network.GetBlock{n.NetworkAddr}
	SendCommand(peerAddr, getBlock.Serialize())
}
func (n Node) SendInv(peerAddr string) {
	blockHashs := n.Chain.GetBlockchainHashs()
	inv := network.Inv{blockHashs, "block", n.NetworkAddr}
	SendCommand(peerAddr, inv.Serialize())

}
func (n Node) SendTxInv(txID, peerAddr string) {

	var idBytes [][]byte
	aTx, _ := hex.DecodeString(txID)
	idBytes = append(idBytes, aTx)
	inv := network.Inv{idBytes, "tx", n.NetworkAddr}
	SendCommand(peerAddr, inv.Serialize())
}

func (n Node) SendGetData(hash []byte, typeOfData string, peerAddr string) {
	gd := network.GetData{hash, typeOfData, n.NetworkAddr}
	SendCommand(peerAddr, gd.Serialize())
}
func (n Node) SendBlock(serializedBlock []byte, peerAddr string) {
	gd := network.Block{serializedBlock, n.NetworkAddr}
	SendCommand(peerAddr, gd.Serialize())
}
func (n Node) SendReqPeers(peerAddr string) {
	gd := network.ReqPeers{n.NetworkAddr}
	SendCommand(peerAddr, gd.Serialize())
}
func (n Node) SendPeers(peerAddr string) {
	peersCmd := network.Peers{n.Peers, n.NetworkAddr}
	SendCommand(peerAddr, peersCmd.Serialize())
}
func (n Node) SendTxCmd(txId, peerAddr string) {
	txCmd := network.TX{n.mempool[txId].Serialize(), peerAddr}
	SendCommand(peerAddr, txCmd.Serialize())
}

/********************************************************************************************
									NODE UTIL FUNCTIONS
*********************************************************************************************/

func (n Node) isPeer(peerAddr string) bool {

	for _, peer := range n.Peers {

		if peer == peerAddr {
			return true
		}

	}
	return false

}

func (n *Node) addPeer(peerAddr string) {
	fmt.Println("New Incoming Peer", peerAddr)
	n.Peers = append(n.Peers, peerAddr)
}

func (n Node) inMempool(txID string) bool {
	return n.mempool[txID] != nil
}

func (n *Node) AddPeer(peerAddr string) {
	n.SendVersion(peerAddr)
}

func (n *Node) GetMempool() []string {
	txIdArr := make([]string, len(n.mempool))

	for i := range n.mempool {
		txIdArr = append(txIdArr, i)
	}
	return txIdArr

}

func NewNode(nodeID, dbFile, walletFile string) *Node {

	//FIXME: View the comments at the end
	bc := NewBlockchain(dbFile)

	w := wallet.GetWallet(walletFile)

	w.Alias = nodeID

	a, _ := strconv.Atoi(nodeID)

	server, nAddr := network.NewServer(a)

	mempool := make(map[string]*TransactionCTX.Transaction)
	bq := NewBlockQueue()

	var mu sync.Mutex

	peers := make([]string, 0)
	peers = append(peers, nAddr) //is it okay to say, i am my own peer, this i dont have to call addPeer is sending and receving a Version

	return &Node{bc, server, peers, w, nAddr, &bq, mempool, &mu}
}

/**************************************************************************************

So the genesis block was giving me two different hashes, this was due to the fact that even
if the transaction was the same in where it was sending noney and how much was sent, the
Serialized copy of txCopy was different in a few byte when it was serialized, the hex value
of the copy will be the same for both tx, but the bytes were different in my checker

So then what to do? well i tried to do something with how the NewNode func works, i commented the
wallet part and then it was all normal, and even with the Wallet Init code if the Wallet file is not created then we do
no get a different txID, so i simply chnage the order, first blockchain is Init then the Wallet and it all
seems to be working

but how did the writing of a Wallet file, interfere with the ecoding of a transaction that was not even
created as of wrting the wallet file.

I will call this a paranormal bug, as i cannot explain what is happening


**************************************************************************************/
