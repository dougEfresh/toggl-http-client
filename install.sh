#!/bin/bash

go install 

for i in * ; do
[ -f $i/doc.go ] && (cd $i && go install)
done
