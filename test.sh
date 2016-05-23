#!/bin/bash

go install 

for i in * ; do
[ -f $i/doc.go ]  && go test -v
done
