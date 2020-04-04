module github.com/whilei/openrpc-go-petstore

go 1.13

require github.com/ethereum/go-ethereum v1.9.12

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/etclabscore/openrpc-go-document v0.0.6
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/rs/xhandler v0.0.0-20170707052532-1eb70cf1520d // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
)

replace github.com/ethereum/go-ethereum => github.com/meowsbits/core-geth v1.10.0-core.0.20200404194932-7af048e1a645

replace github.com/etclabscore/openrpc-go-document => github.com/meowsbits/openrpc-go-document v0.0.6

replace github.com/alecthomas/jsonschema => github.com/meowsbits/jsonschema v0.0.2

replace github.com/etclabscore/go-jsonschema-traverse => github.com/meowsbits/go-jsonschema-traverse v0.0.4

replace github.com/gregdhill/go-openrpc => github.com/meowsbits/go-openrpc v0.0.1
