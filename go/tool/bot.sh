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

go2root(){
    cd $(dirname $0)
    cd ../
}

run(){
    reflect
    echo $PWD
    go=$(check_bin go)
    $go run $MAIN
}

reflect(){
    pkgreflect=$(check_bin pkgreflect)
    $pkgreflect $MODULE
}

usage(){
    echo "Usage: $0 <run|reflect>"
}

case $1 in
    run)
        go2root
        run
        ;;
    reflect)
        go2root
        reflect 
        ;;
    *)
        usage
        ;;
esac
