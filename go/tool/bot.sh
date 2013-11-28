#!/bin/bash

MAIN='bot.go'
MODULE='./ircbot/module/'

check_bin(){
    bin=$(which $1)
    if [[ -z "$bin" ]]; then
        "Need $1, check the document !!"
        exit 1
    fi
    echo $bin
}

run(){
    reflect
    cd $(dirname $(dirname $0))
    go=$(check_bin go)
    $go run $MAIN
}

reflect(){
    cd $(dirname $(dirname $0))
    pkgreflect=$(check_bin pkgreflect)
    $pkgreflect $MODULE
}

usage(){
    echo "Usage: $0 <run|reflect>"
}

case $1 in
    run)
        run
        ;;
    reflect)
        reflect 
        ;;
    *)
        usage
        ;;
esac
