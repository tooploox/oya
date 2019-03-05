#!/usr/bin/env bash

OYA_URL="https://github.com/tooploox/oya"
OYA_TAGS_URL="https://api.github.com/repos/tooploox/oya/tags"

check_sops() {
    printf " Sops -> "
    if hash sops 2>/dev/null; then
        printf "OK ($(sops --version))\n"
    else
        printf "FAIL\n"
        echo "Oya requires SOPS (Secrets OPerationS) to manage secret files."
        echo "Visit https://github.com/mozilla/sops for install instuctions."
        exit 1
    fi
}

check_deps() {
    check_sops
}

get_os() {
    unamestr=`uname`
    if [[ "$unamestr" == 'Linux' ]]; then
        os='linux'
        return 0
    elif [[ "$unamestr" == 'Darwin' ]]; then
        os='darwin'
        return 0
    else
        echo "Sorry Oya doesn't support ${unamestr} yet! :("
        exit 1
    fi
}

get_arch() {
    archstr=`uname -m`
    if [[ "$archstr" == 'x86_64' || "$archstr" == 'x86' ]]; then
        arch='386'
    elif [[ "$archstr" == 'AMD64' ]]; then
         arch='amd64'
    else
        echo "Sorry Oya doesn't support ${archstr} architecture yet! :("
        exit 1
    fi
}

find_latest() {
    version="$(fetch_versions | tail -n 1)"
    if [[ -n ${version} ]];then
        echo "${version}"
        return 0
    fi
}

fetch_versions() {
    curl -s $OYA_TAGS_URL |
        \awk -v RS=',' -v FS='"' '$2=="name"{print $4}' |
        sort -t. -k 1,1n -k 2,2n -k 3,3n -k 4,4n -k 5,5n 
}

install_version() {
    _ver=$1
    _os=$2
    _arch=$3
    _url="${OYA_URL}/releases/download/${_ver}/oya_${version}_${_os}_${_arch}.gz"
    _sum_url="${OYA_URL}/releases/download/${_ver}/oya_${version}_SHA256SUMS"
    get_and_check "${_url}" "${_sum_url}" && return
}

get_and_check() {
    _url=$1
    _sum_url=$2
    _tmp_dir=`mktemp -d`
    _archive="oya.gz"
    _fileName="oya"
    _bin_dir="/Users/bart/work/tooploox/tmp/bin"
    _args=""
    get_package "$_url" "$_tmp_dir" "$_archive" || return $?
    verify_sha "$_sum_url" "$_tmp_dir" "$_archive" || return $?
    gunzip ${_tmp_dir}/${_archive} || return $?
    chmod +x ${_tmp_dir}/${_fileName} || return $?
    if [ ! -w "$_bin_dir" ]; then
        _args="sudo "
    fi
    ${_args}mv ${_tmp_dir}/${_fileName} ${_bin_dir}
}

get_package() {
    _url=$1; _dir=$2; _file=$3
    curl -sSL ${_url} > ${_dir}/${_file} ||
        {
            echo "error downloading oya package $url"
            return 1
        }
}

verify_sha() {
    _url=$1; _dir=$2; _file=$3
    curl -sSL ${_url} > ${_dir}/shasums ||
        {
            echo "error downloading oya package SHA sums $url"
            return 1
        }
    _mysum=$(shasum -a 256 ${_dir}/${_file} | awk '{print $1}')
    _result=$(cat ${_dir}/shasums | grep $_mysum)
    if [ -z "$_result" ]; then
        echo "downloaded archive checksum failed"
        return 1
    fi
}

oya_install() {
    get_os
    get_arch
    check_deps
    if [ -z "$2" ];then
        version=$(find_latest)
    else
        version=$2
    fi
    install_version $version $os $arch
}

oya_install "$@"
