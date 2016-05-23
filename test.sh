#!/bin/bash

go test  -v

for i in * ; do
if [ -f $i/doc.go ]  ; then
(cd $i && go test -v) || exit $?
fi
done
