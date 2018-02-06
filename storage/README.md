## Storage design

This article documents the design of our persistence system for storing Blocks in a Blockchain.

If you are not familiar with the concepts of a Blockchain or Blocks, I recommend you checkout any of the wonderful articles on Blockchain that are becoming increasingly popular; 
such as [this](https://medium.com/s/welcome-to-blockchain/everything-you-need-to-know-about-blockchain-but-were-too-embarrassed-to-ask-b3cee3e918f8) 
or [this](https://medium.freecodecamp.org/the-authoritative-guide-to-blockchain-development-855ab65b58bc)

### Background

Every node in a Blockchain system needs a way to store information about blocks. 
This could be in memory storage, or something persistent like a database.

Repositories are a data storage abstraction that provide more functionality than simple CRUD.  
We are going to use the idea of a repository to ensure that the immutability of the Blockchain is maintained, and provide functions that make reading blocks easier.  
Use-cases for these functions are described below.

#### Domain

As this is a Blockchain implementation, we need to store blocks!  
Block anatomy:

```go
type BlockHeader struct {
	Index            uint64
	PreviousHash     *Hash
	Timestamp        uint64
	TransactionsHash *Hash
	Nonce            uint64
	Hash             *Hash
}

type Block struct {
 	Header       *BlockHeader
 	Transactions []*Transaction
}
```

#### API

As far as the outside world cares, there are two things to read out of the database:  
Blocks, and Headers.  

```go
type BlockRepository interface {
	Close() error
	BestBlock() (*Block, error)

	BlockByHash(hash *Hash) (*Block, error)
	BlockByIndex(index uint64) (*Block, error)

	SaveBlock(*Block) error
}

type HeaderRepository interface {
	Close() error
	BestHeader() (*BlockHeader, error)

	HeaderByHash(hash *Hash) (*BlockHeader, error)
	HeaderByIndex(index uint64) (*BlockHeader, error)

	SaveHeader(*BlockHeader) error
}
```

Contract:  
 - Functions should only return error if there was an issue reading or writing
 - Read functions will return nil if the entity isn't found
 - Write functions are destructive, and will override if an entity exists (luckily, this should happen if entities are stored by hash)

#### Use cases

Here we define common business problems we see while using a blockchain system.  
The repository functions are defined to help simplify the interactions with our persistence strategy.

##### Saving blocks/headers

The simplest use-case; blocks/headers are saved when they are downloaded from a peer during sync, or when a node broadcasts a new block.
The `SaveHeader(*BlockHeader)` function lets us save a header.  
The `SaveBlock(*Block)` function lets us save a block.  

##### Advertising node version

When a node discovers a peer, it sends a Version message; describing the protocol version, as well as the *best block* that the node has to offer.  
The `BestBlock()` function lets us find the top of the chain.  

You can find detailed information about Nodes and Peers [here](https://github.com/tclchiam/oxidize-go/tree/master/node) **TODO**  
You can find detailed information about [gRPC](https://grpc.io) and how we use it to communicate between nodes [here](https://github.com/tclchiam/oxidize-go/tree/master/rpc) **TODO**

##### Syncing blocks from peers

When syncing with peers that are at a *higher* version, the node will request headers from peers, starting at the node's *best header*.  
In order to find the header to start from, we find the current *best header* in our node's local storage.  
The `BestHeader()` function lets us find most recent header.

When the node downloads and validates a header, it will save it.  
The `SaveHeader(*BlockHeader)` function lets us save a header.  

This sets up a situation where the *best header* can have a higher index than the *best block*.  
When that happens, the node requests the full block from a different peer.

The `BestBlock()` function lets us find the most recent block.

When the node downloads and validates the block, it will save it.  
The `SaveBlock(*Block)` function lets us save a block.

##### Providing blocks for peers

When syncing with peers that are at a *lower* version, the peers will request headers up to the latest, starting at a specific index/hash.  
In order to find blocks up to the latest, we iterate *forward* through the chain starting from the requested header  
This is difficult, since headers don't know the hash of the next block.

The `HeaderByIndex(index uint64)` function lets us find the starting header, as well the next headers.
The `BlockByIndex(index uint64)` function lets us find the starting block, as well the next blocks.

##### Finding UTXOs

In order to find Unspent Transaction Outputs, we iterate *back* through the entire chain starting from the *best block*.  
This is already pretty easy, since each block contains the hash of the previous block.  

The `BestBlock()` function lets us find the top of the chain.  
The `BlockByHash(hash)` function lets us find the next block.

#### Bolt

Bolt DB is a key/value store - for Go - with a simple hierarchy.  
The database file's top level is a key/value set where the values are Buckets, and the keys are the Bucket's name  
Buckets are a key/value set, each key must be unique.  
Together, this makes Bolt a *two* level key/value store.

### Storage

Under the hood, we are going to separate a Block into its Header and Body.  
Header matches the struct shown above:

```go
type BlockHeader struct {
	Index            uint64
	PreviousHash     *Hash
	Timestamp        uint64
	TransactionsHash *Hash
	Nonce            uint64
	Hash             *Hash
}
```

Body is the part of the Block that is not the Header; in this case, just the transactions:

```go
type BlockBody struct {
 	Transactions []*Transaction
}
```

We are going to store the Body and Header, indexed by the block's hash  
With Bolt, the keys in a bucket must be unique; 
so in order to be able to store both the body and the header by the hash, 
we need to put them in different buckets.

```
DB -> Header Bucket -> 0x004f2ea8 -> Header {Index:1}  
                    -> 0x009cb1e8 -> Header {Index:2}  
                    -> 0x0026f8a6 -> Header {Index:3}  
                    
   -> Body Bucket   -> 0x004f2ea8 -> Body   {Index:1}  
                    -> 0x009cb1e8 -> Body   {Index:2}
```

The idea of a hashing function is to produce a unique representation of an arbitrary piece of data, so we should never have a collision where two different blocks/headers have the same hash.   
Our hashing function produces byte arrays with a length of 32.

This makes saving and reading Blocks and Headers pretty simple, satisfying six of the Repository functions!  
- Saving a header
- Saving a block
- Reading blocks by hash
- Reading blocks by index
- Reading headers by hash
- Reading headers by index

We still need a way to find Blocks and Headers by Index...  
It makes sense to store hashes by index; but we also have to keep in mind that there can be more headers stored than blocks.  
Since Bolt just sees keys as []byte - and uint64 can be represented as 8 bytes - 
we can store the index/hash pair in the same buckets as blocks/headers (respectively) without worrying about collisions.

```
DB -> Header Bucket -> 0x00000001 -> 0x004f2ea8
                    -> 0x00000002 -> 0x009cb1e8
                    -> 0x00000003 -> 0x0026f8a6
                    
   -> Body Bucket   -> 0x00000001 -> 0x004f2ea8
                    -> 0x00000002 -> 0x009cb1e8
```

With that, we should have all the parts necessary for all our repository functions!
