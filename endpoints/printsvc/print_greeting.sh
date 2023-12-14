#!/bin/bash

grpcurl -plaintext -emit-defaults -import-path . \
    -proto "/root/projects/k8slearn/printsvc/contract/printsvc.proto" \
    -d "{}" localhost:6002 printsvc.PrintService/PrintGreeting