go build

if [[ $? -eq 0 ]]
then
    ./spacemon $@
fi

