#!/bin/bash

go  "$@"

for i in * ; do
[ -f $i/doc.go ] && (cd $i && go "$@")
done
