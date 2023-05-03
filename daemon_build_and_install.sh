cd daemon && go build -o ../spacemond

if [[ $? -eq 0 ]]
then
    cd ..
    ./service-uninstall.sh
    ./service-install.sh
fi

