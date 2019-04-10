#!/bin/bash

for pool in null sync power2 reserved leakysync; do 
    GODEBUG=gctrace=1 ./main -pooltype $pool 2>&1 | perl -plE 's/.*?(?:\d+->(\d+)->\d+).*/$1/' | sort -rn | head -4 | tail -2
done
