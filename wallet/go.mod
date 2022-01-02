module nshmdhani.com/notcoin/wallet

go 1.13

require (
	github.com/boltdb/bolt v1.3.1
	github.com/desertbit/grumble v1.1.1
	github.com/fatih/color v1.11.0
	github.com/gin-gonic/gin v1.7.1 // indirect
	github.com/go-playground/validator/v10 v10.6.0 // indirect

	github.com/go-resty/resty/v2 v2.4.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ugorji/go v1.2.5 // indirect
	github.com/urfave/cli v1.22.5
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	nshmadhani.com/notcoin/blockchain v0.0.0-00010101000000-000000000000
)

replace nshmadhani.com/notcoin/blockchain => ../blockchain
