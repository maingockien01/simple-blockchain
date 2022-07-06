# simple-blockchain

A simple blockchain implemented by Go Lang.

This project is made for academic purpose. Please feel free to check and give me some comments.

# Architecture
- Network layer: `udp` folder
- Protocol layer: `protocol` folder
- Peer layer: files `peer/peer*.go` in charge of handle request, reply
- Blockchain layer
        - `chain.go` in charge of chain
        - `consensus.go` in charge of doing consensus
        - `mining.go` in charge of managing mining pool

# Consensus
- Files: `consensus.go` and `chain.go` and `peerConsensus.go`

- Algorithm:
Peer will send stats to all known peers
Retry 1 time if stats timeout
When all stats return results
Using `Consensus` struct in `consensus.go` to add the stats reply and process the consensus


After getting the chain with agreed contacts (that returned same chain: height and hash),
using mutex lock to lock syncing mode
append new block or rebuild if encounter misfit block
for each block,
        try to get block from a peer in agreed contact list
        if block is not added, move to the next peer in the list

        if none block from all peers agreed fits the current chain, rebuild from 0

After done with rebuilding block, verify the whole chain again from 0 to highest chain.

unlock lock for syncing mode

# Running
to build `make all`
to run, use command `go run main.go [current_host] [port_to_run]`
or `./blockchain [current_host] [port_to_run]`

Example, I'm currently on `grouse.cs.umanitoba.ca` and want to run on port `8999`. To run: `./blockchain grouse.cs.umanitoba.ca 8999`

# Mining
- My peer wont start mining untill there is at least 1 words to do 
- do add word protocol to start the mining