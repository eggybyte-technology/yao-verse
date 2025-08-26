module github.com/eggybyte-technology/yao-chronos

go 1.24.5

require (
	github.com/eggybyte-technology/yao-verse-shared v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.16.2
)

require (
	github.com/holiman/uint256 v1.3.2 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)

replace github.com/eggybyte-technology/yao-verse-shared => ../shared
