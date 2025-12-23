#!/bin/bash
echo "will generate code in gen from the proto file"
buf dep update
buf lint
buf generate
