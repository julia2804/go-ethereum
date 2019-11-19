#!/bin/bash
for loop in 100 200 300 400 500
do
    /home/mimota/go/src/github.com/ethereum/go-ethereum/build/bin/geth export ethexport${loop} 0 ${loop}000 -datadir ethereum_data
done