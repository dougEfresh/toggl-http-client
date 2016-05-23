#!/bin/bash

go install 

for i in * ; do
[ -f $i/doc.go ]  && go install
done