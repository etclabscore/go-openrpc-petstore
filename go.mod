module github.com/whilei/openrpc-go-petstore

go 1.13

require github.com/ethereum/go-ethereum v1.9.12

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/alecthomas/jsonschema v0.0.2
	github.com/etclabscore/openrpc-go-document v0.0.10
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-openapi/spec v0.19.7
	github.com/gregdhill/go-openrpc v0.0.1
	github.com/rs/xhandler v0.0.0-20170707052532-1eb70cf1520d // indirect
)

replace github.com/etclabscore/openrpc-go-document => /home/ia/dev/etclabscore/openrpc-go-document

replace github.com/alecthomas/jsonschema => github.com/meowsbits/jsonschema v0.0.2

replace github.com/etclabscore/go-jsonschema-traverse => github.com/meowsbits/go-jsonschema-traverse v0.0.4

replace github.com/gregdhill/go-openrpc => github.com/meowsbits/go-openrpc v0.0.1
