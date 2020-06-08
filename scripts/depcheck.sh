function depcheck {
    if [[ -z $(which $1 2> /dev/null) ]]; then
        echo "$1 missing, please install"
        exit 1
    fi
}
