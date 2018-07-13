#!/bin/sh

protoc --proto_path=. --twirp_out=./go/ --go_out=./go/ ./*.proto

protoc --proto_path=. --swift_opt=Visibility=Public --swift_out=./swift/ ./*.proto

java -jar wire-compiler-2.3.0-RC1-jar-with-dependencies.jar --proto_path=. --java_out=./java/

touch ./ts/api.js
pbjs -t static-module -p $GOPATH/src -p ../vendor -p ./ -w commonjs -o ./ts/api.js ./*.proto | pbts -o ./ts/api.d.ts ./ts/api.js