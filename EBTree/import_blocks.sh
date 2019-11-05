#!/bin/bash
for loop in 100 200 300 400 500
do
    /home/mimota/go/src/github.com/ethereum/go-ethereum/build/bin/geth import /home/mimota/ethexport${loop} -datadir /home/mimota/data0
    du -sh /home/mimota/data0 >> import_chain.log

done