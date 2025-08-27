module github.com/eggybyte-technology/yao-portal

go 1.24.5

require (
	github.com/eggybyte-technology/yao-verse-shared v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.16.2
	github.com/gorilla/websocket v1.5.3
	golang.org/x/time v0.9.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/holiman/uint256 v1.3.2 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/eggybyte-technology/yao-verse-shared => ../shared
