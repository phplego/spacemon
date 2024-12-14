#!/bin/bash

go build

if [[ $? -eq 0 ]]
then
    ./spacemon --dry $@
fi

