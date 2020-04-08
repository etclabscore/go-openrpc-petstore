package petstore

/*
Since Ethereum's RPC service library tolerates pretty arbitrary method styles,
filtering out args (eg context) as needed, we don't need to wrap the PetStore
type in any special kind of interface (but we will need to for go's standard rpc lib).
*/
