#!/bin/bash

go fmt

for i in * ; do
[ -f $i/doc.go ] && (cd $i  && go fmt)
done
