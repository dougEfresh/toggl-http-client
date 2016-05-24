#!/bin/bash

go  "$@"

for i in * ; do
if [ -f $i/doc.go ] ; then 
 (cd $i && go "$@")
fi
done
