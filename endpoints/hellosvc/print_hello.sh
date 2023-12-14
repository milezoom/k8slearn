#!/bin/bash

grpcurl -plaintext -emit-defaults -import-path . \
    -proto "/root/projects/k8slearn/hellosvc/contract/hellosvc.proto" \
    -d @ localhost:6001 hellosvc.HelloService/PrintHello <<EOM
{
    "name": "World"
}
EOM