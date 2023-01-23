#!/bin/bash

export PROJECT_NAME=$1 # 项目名字
export REGISTRY_URL=$2 # 仓库地址
export THREAD=$3 # 最大并行构建线程数

build_names=('auth-rpc-token-store' 'auth-rpc-token-enhancer' 'user-rpc-sys' \
'video-rpc-sys' 'chat-rpc-sys' \
'auth-api' 'user-api' 'video-api' 'chat-api' \
'mq-nsq-consumer')

function go_build() {
  if [ "$1" -ef "" ]; then
    return 0
  fi

  path='./app/service'

  eval "$(echo "$1" | awk '{split($0,array,"-");for(i in array)print "arr["i"]="array[i]}')"
  # shellcheck disable=SC2154
  length=${#arr[@]}


  filename=""

  if [ "$length" = 2 ]; then
#   Api 服务
    path="${path}""/""${arr[1]}/api"
    filename="${arr[1]}"
  elif [ "$length" -eq 3 ] && [ "${arr[1]}" = "mq" ]; then
#   Mq 服务
    path="${path}""/""${arr[1]}""/""${arr[2]}""/""${arr[3]}"
    filename="${arr[3]}"
  else
#   Rpc 服务
    for((i=1;i<=length;i++))
    do
        path="${path}""/""${arr[i]}"
    done
    if [ "$length" -ge 3 ]; then
      filename="$filename""${arr[1]}""."
      for((j=3;j<=length;j++))
      do
        if [ "$j" = "$length" ]; then
          filename="$filename""${arr[j]}"
        else
          filename="$filename""${arr[j]}""_"
        fi
      done
    fi
  fi
  echo "go build -ldflags=""\"-s -w\""" -o ""$path""/""$PROJECT_NAME-""$1"" ""$path""/""$filename"".go"
  go build -ldflags="-s -w" -o "$path""/""$PROJECT_NAME-""$1" "$path""/""$filename"".go"
  mkdir "$path""/""manifest"
  cp -r ./manifest/config "$path""/""manifest/config"

  return 1
}

function docker_build() {
    if [ "$1" -ef "" ]; then
      return 0
    fi

    path='./app/service'

    eval "$(echo "$1" | awk '{split($0,array,"-");for(i in array)print "arr["i"]="array[i]}')"
    # shellcheck disable=SC2154
    length=${#arr[@]}

    for((i=1;i<=length;i++))
    do
      path="${path}""/""${arr[i]}"
    done
    docker build -t "$REGISTRY_URL""/""$PROJECT_NAME""/""$1" "${path}"
    return 1
}

function docker_push() {
    if [ "$1" -ef "" ]; then
      return 0
    fi

    path='./app/service'

    eval "$(echo "$1" | awk '{split($0,array,"-");for(i in array)print "arr["i"]="array[i]}')"
    # shellcheck disable=SC2154
    length=${#arr[@]}

    for((i=1;i<=length;i++))
    do
      path="${path}""/""${arr[i]}"
    done
    docker push "$REGISTRY_URL""/""$PROJECT_NAME""/""$1"
    return 1
}

[ -e /tmp/fd1 ] || mkfifo /tmp/fd1
exec 3<>/tmp/fd1
rm -rf /tmp/fd1

for ((i=1;i<=THREAD;i++))
do
  echo >&3
done

if [ ! -d "/www/service/""$PROJECT_NAME" ];then
  mkdir "/www/service/""$PROJECT_NAME"
fi

cd /www/service/"$PROJECT_NAME" || exit

echo "start building..."
for build_name in "${build_names[@]}"
do
  read -r -u3
  {
    echo "$build_name building go..."
    go_build "${build_name}"
    echo "$build_name build go successfully"

    echo "$build_name building docker..."
    docker_build "${build_name}"
    echo "$build_name build docker successfully"

    echo "$build_name pushing docker..."
    docker_push "${build_name}"
    echo "$build_name pushing docker successfully"

    echo >&3
  } &
done

wait

exec 3<&-
exec 3>&-