#!/bin/bash
set -ev

CURRENT_DIR=$(cd $(dirname $0); pwd)
echo $CURRENT_DIR
cd "$CURRENT_DIR"

protoc --version | grep 'libprotoc 24.2' || {
  echo "install protoc 24.2"
  PROTOC_ZIP=protoc-24.2-osx-x86_64.zip
  curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v24.2/$PROTOC_ZIP
  sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
  sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
  rm -f $PROTOC_ZIP
}
HostName=$(hostname)

if [[ $HostName != 'admindeMacBook-Pro-4.local' ]]; then
   md5 $(go env GOPATH)/bin/protoc-gen-go | grep 4994a5677761d18af2bcc03f51440250 || {
     echo "please reinstall protoc-gen-go:  curl 'http://qiniu.brightguo.com/sunmi/protoc-gen-go_mac13.0' -o protoc-gen-go && chmod +x protoc-gen-go && mv protoc-gen-go $(go env GOPATH)/bin"
     exit 1
   }
fi

protoc-gen-go-errors -version=true | grep 1.0.0 || {
    echo version too low, please reinstall: go install github.com/guoming0000/protoc-gen-go-gin/cmd/protoc-gen-go-errors@latest
    exit 2
}
protoc-gen-go-gin -version=true | grep 1.0.3 || {
    echo version too low, please reinstall: go install github.com/guoming0000/protoc-gen-go-gin/cmd/protoc-gen-go-gin@latest
    exit 3
}
protoc-gen-openapi -version=true | grep 1.0.0 || {
    echo version too low, please reinstall: go install github.com/guoming0000/protoc-gen-go-gin/cmd/protoc-gen-openapi@latest
    echo 4
}

which yq || {
   # https://github.com/mikefarah/yq
   # go install github.com/mikefarah/yq/v4@latest
   # brew install yq
   wget http://qiniu.brightguo.com/sunmi/yq -O /usr/local/bin/yq &&  chmod +x /usr/local/bin/yq
}

# code
protoc --go-errors_out=fe_ecode=./api/docs/fe_ecode:. api/proto/ecode.proto

protoc -I. -I api/proto/third_party --go-gin_out=. --go_out=. api/proto/space.proto
protoc -I. -I api/proto/third_party --go-gin_out=. --go_out=. api/proto/login.proto
protoc -I. -I api/proto/third_party --go-gin_out=. --go_out=. api/proto/dumplinks.proto


# 生成到openapi文件夹
for protoName in "space" "dumplinks" "login";
do
  go run api/docs/swagger.go api/proto/$protoName.proto api/docs/$protoName/$protoName.swagger.proto
  protoc --proto_path=. \
          --proto_path=./api/proto/third_party \
          --openapi_out=fq_schema_naming=true,default_response=false,output_mode=source_relative:. \
          api/docs/$protoName/$protoName.swagger.proto

  yq -Poj api/docs/$protoName/$protoName.swagger.yaml > api/docs/$protoName/$protoName.swagger.json
done

