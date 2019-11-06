#!/bin/bash
geth_path="/home/mimota/go/src/github.com/ethereum/go-ethereum/build/bin"
eth_export_path="/home/mimota/ethenv"
datadir_path="/home/mimota/data"
data_size_record_path="/home/mimota/mimota-0.1-size.log"
time_record_path="/home/mimota/mimota-0.1-time.log"

rm $data_size_record_path
rm $time_record_path

for loop in 1 2 3 4 5 6 7 8 9 10 20 30 40 50 60 70 80 90 100
do
    start=$(date +%s%N)
    start_ms=${start:0:16}

    ${geth_path}/geth import ${eth_export_path}/ethexport${loop} -datadir $datadir_path

    end=$(date +%s%N)12
    end_ms=${end:0:16}
    #echo "cost time is:"
    printf "${loop}," >> ${time_record_path}
    echo "scale=6;($end_ms - $start_ms)/1000000" | bc >> ${time_record_path}

	  printf "${loop}," >> $data_size_record_path
    du -s $datadir_path >> $data_size_record_path

done