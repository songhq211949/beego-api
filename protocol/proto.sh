#!/bin/bash

outPath=../proto
fileArray=(WSBaseReqProto WSBaseResProto WSMessageResProto WSUserResProto)

for i in ${fileArray[@]};
do
    echo "generate cli protocol go code: ${i}.proto"
    protoc --go_out=$outPath ./$i.proto
done
