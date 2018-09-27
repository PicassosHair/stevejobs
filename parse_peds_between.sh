#!/bin/bash

BASEDIR=/root/pedsparser

. ${BASEDIR}/_rainbow.sh

echogreen "Start parsing PEDS data between $1 ... $2";

for (( y=$1; y<=$2; y++ ))
    do 
    echogreen "Start parsing year for $y"
    ${BASEDIR}/parse_peds_for.sh $y
done