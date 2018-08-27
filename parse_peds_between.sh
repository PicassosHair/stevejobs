#!/bin/bash

echo "[INFO] Start parsing PEDS data between $1 ... $2";

for (( y=$1; y<=$2; y++ ))
    do 
    echo "[INFO] Start parsing year for $y"
    ./parse_peds_for.sh $y
done