# Flexpool SOLO (WIP)
A lightweight portable self-hosted SOLO Ethereum Mining Pool.

Developed with ❤️ to make it as fast as Flexpool.

> WARNING: This project is currently under heavy development

# Why?

There's a bunch of SOLO pools outside, and sadly none of them are transparent enough to have any trust. Also, there were incidents when SOLO pools have just stolen blocks. At Flexpool, we decided to develop a self-hosted open-source solution to all this mess.

### Portability

The Flexpool Solo is a highly portable software, and the only dependency you need is the Ethereum node. Also, it uses an embedded lightning-fast key-value LevelDB database, so there is no need to set up an external databases.

### Speed

The Flexpool SOLO is written in the Go Programming Language, which makes it very speedy and more profitable accordingly.

# Installation

TBD

# TODO

### Core Mining Engine
- [x] Worker Authentication & Work Receiver
- [x] Share verification
- [x] Block submission

### Analytics
- [x] Accept hashrate reports
- [x] Collect worker statistics
- [x] Collect mined blocks and best shares
- [ ] Implement querying & API
- TBD

### Front End
- [ ] Base Website
- [ ] General & Per worker Statistics
- [ ] Network statistics
- [ ] Node Health monitoring 
- TBD

### Documentation
TBD

### Other
- [ ] Write tests for everything

# License

GNU Affero General Public License v3.0
