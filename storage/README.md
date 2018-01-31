## Storage design

### Background

Repositories are a data storage abstraction that provide more functionality than simple CRUD.  
We are going to use the idea of a repository to ensure that the immutability of the blockchain is maintained, and provide functions that make reading blocks easier.

#### Domain

As this is a blockchain implementation, we need to store blocks!  
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

As far as the outside wolrd cares, there are two things to read out of the database:  
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

##### Finding UTXOs

In order to find Unspent Transaction Outputs, we iterate *back* through the entire chain starting from the best block.  
This is already pretty easy, since each block contains the hash of the previous block.  

The `BestBlock()` function lets us find the top of the chain.  
The `BlockByHash(hash)` function lets us find the next block.

##### Providing blocks for peers

When syncing with peers that are at a lower version, the peers will request headers up to the latest, starting at a specific index/hash.  
In order to find blocks up to the latest, we iterate *forward* through the chain starting from the requested block  
This is more difficult, since blocks don't know the hash of the next block.

The `BlockByIndex(index uint64)` function lets us find the starting block, as well the next blocks.

##### Saving blocks

The simplest use-case; blocks are saved when they are downloaded from a peer during sync, or when a node broadcasts a new block.

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
                    
   -> Body Bucket   -> 0x004f2ea8 -> Body   {Index:1}  
                    -> 0x009cb1e8 -> Body   {Index:2}
```

This makes saving and reading Blocks and Headers pretty simple, satisfying two of the Repository functions!  

We still need a way to find Blocks and Headers by Index...  
It makes sense to store hashes by index, but that idea doesn't fit in either of the buckets we have so far; so lets make a new bucket for hashes!

```
DB -> Hash Bucket   -> 1          -> 0x004f2ea8
                    -> 2          -> 0x009cb1e8
```

With that, we should have all the parts necessary for all our repository functions!
