#!/bin/bash

. ./_rainbow.sh

echogreen "Start diffing PEDS data between $1 ... $2";

for (( y=$1; y<=$2; y++ ))
    do 
    echogreen "Start parsing year for $y"
    ./diff_peds_for.sh $y
done