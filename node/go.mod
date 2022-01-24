module nshmadhani.com/notcoin/node

go 1.13

require (
	github.com/desertbit/grumble v1.1.1
	github.com/fatih/color v1.10.0
	github.com/gin-gonic/gin v1.7.1
	github.com/urfave/cli/v2 v2.3.0
	nshmadhani.com/notcoin/blockchain v0.0.0-00010101000000-000000000000
)

replace nshmadhani.com/notcoin/blockchain => ../blockchain
