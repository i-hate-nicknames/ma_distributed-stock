mkdir -p $PWD/bin
GOBIN=$PWD/bin go install ./cmd/{warehouse,stock}

