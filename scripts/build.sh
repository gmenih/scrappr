#! /bin/bash

mkdir -p dist
mkdir -p dist/tmp

cp go.{sum,mod} ./dist/tmp
cp -r ./src ./dist/tmp
mv ./dist/tmp/src/*.go ./dist/tmp
cd ./dist/tmp
zip ../index.zip -r *
cd -
rm -r ./dist/tmp
