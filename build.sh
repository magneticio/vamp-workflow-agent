#!/usr/bin/env bash

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

reset=`tput sgr0`
green=`tput setaf 2`
yellow=`tput setaf 3`

version="0.9.0"
target='target'
target_vamp=${target}'/vamp'
target_docker=${target}'/docker'
assembly_go='vamp.tar.gz'
docker_image_name="magneticio/vamp-workflow-agent:${version}"

cd ${dir}

function parse_command_line() {
    flag_help=0
    flag_list=0
    flag_clean=0
    flag_make=0
    flag_build=0

    for key in "$@"
    do
    case ${key} in
        -h|--help)
        flag_help=1
        ;;
        -l|--list)
        flag_list=1
        ;;
        -c|--clean)
        flag_clean=1
        ;;
        -m|--make)
        flag_make=1
        ;;
        -b|--build)
        flag_make=1
        flag_build=1
        ;;
        *)
        ;;
    esac
    done
}

function build_help() {
    echo "${green}Usage of $0:${reset}"
    echo "${yellow}  -h|--help   ${green}Help.${reset}"
    echo "${yellow}  -l|--list   ${green}List built Docker images.${reset}"
    echo "${yellow}  -r|--remove ${green}Remove Docker image.${reset}"
    echo "${yellow}  -m|--make   ${green}Build the binary and copy it to the Docker directories.${reset}"
    echo "${yellow}  -b|--build  ${green}Build Docker image.${reset}"
}

function go_build() {
    cd ${dir}
    bin='vamp-workflow-agent'
    export GOOS='linux'
    export GOARCH='amd64'
    echo "${green}building ${GOOS}:${GOARCH} ${yellow}${bin}${reset}"
    rm -rf ${target_vamp} && mkdir -p ${target_vamp}

    go get github.com/tools/godep
    godep restore
    go install
    CGO_ENABLED=0 go build -v -a -installsuffix cgo

    mv ${bin} ${target_vamp} && chmod +x ${target_vamp}/${bin}
}

function npm_make {
    cp ${dir}/package.json ${target_vamp}
    cd ${target_vamp}
    npm install
}

function docker_make {

    append_to=${dir}/${target_docker}/Dockerfile
    cat ${dir}/Dockerfile | grep -v ADD | grep -v ENTRYPOINT > ${append_to}

    echo "${green}appending common code to: ${append_to} ${reset}"
    function append() {
        printf "\n$1\n" >> ${append_to}
    }

    append "ADD ${assembly_go} /opt"
    append "ENTRYPOINT [\"/opt/vamp/vamp-workflow-agent\"]"
}

function vamp_archive {
    cd ${dir}/${target} && tar -zcf ${assembly_go} vamp
    mv ${dir}/${target}/${assembly_go} ${dir}/${target_docker} 2> /dev/null
}

function docker_build {
    echo "${green}building docker image: $1 ${reset}"
    docker build -t $1 $2
}

function docker_rmi {
    echo "${green}removing docker image: $1 ${reset}"
    docker rmi -f $1 2> /dev/null
}

function docker_image {
    echo "${green}built images:${yellow}"
    docker images | grep 'magneticio/vamp-workflow-agent'
}

function process() {
    rm -Rf ${dir}/${target} 2> /dev/null && mkdir -p ${dir}/${target_docker}

    if [ ${flag_make} -eq 1 ]; then
        docker_make
        go_build
        npm_make
        vamp_archive
    fi

    if [ ${flag_clean} -eq 1 ]; then
        docker_rmi ${docker_image_name}
    fi

    if [ ${flag_build} -eq 1 ]; then
        cd ${dir}/${target_docker}
        docker_build ${docker_image_name} .
    fi

    if [ ${flag_list} -eq 1 ]; then
        docker_image
    fi

    echo "${green}done.${reset}"
}

parse_command_line $@

if [ ${flag_help} -eq 1 ] || [[ $# -eq 0 ]]; then
    build_help
fi

if [ ${flag_list} -eq 1 ] || [ ${flag_clean} -eq 1 ] || [ ${flag_make} -eq 1 ] || [ ${flag_build} -eq 1 ]; then
    process
fi
