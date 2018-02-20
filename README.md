# Oxidize: golang

Learn blockchain development, one bit at a time!

This project is meant for people interested in understanding how blockchain works by building one from scratch.

## Installation
```bash
go get github.com/tclchiam/oxidize-go
glide install
```

## Design Docs

Each high level package will eventually have docs describing the system design at that level.  
Each package is designed to work as a stand alone package (WIP)

### Packages

- [ ] [Blockchain](https://github.com/tclchiam/oxidize-go/tree/master/blockchain)
- [ ] [Identity](https://github.com/tclchiam/oxidize-go/tree/master/identity)
- [ ] [RPC](https://github.com/tclchiam/oxidize-go/tree/master/rpc)
- [x] [Storage](https://github.com/tclchiam/oxidize-go/tree/master/storage)
- [ ] [Node](https://github.com/tclchiam/oxidize-go/tree/master/node)
- [ ] [Wallet](https://github.com/tclchiam/oxidize-go/tree/master/wallet)

## Contributing

Pull requests are welcome!

I have a rough plan on what I want to do with the project in the future, but nothing is set in stone.  
I have a tentative road map in Trello, there is a link on the GitHub page.

### Disclaimer 

*This is not meant to be used in production*

I like test driving, but much of the code was implemented without writing tests first. 
Moving forward, I plan on test driving where it makes sense (I don't love testing view logic)
