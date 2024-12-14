#!/bin/bash

cd daemon && go build -o ../spacemond

if [[ $? -eq 0 ]]
then
    cd ..
    ./spacemond $@
fi

