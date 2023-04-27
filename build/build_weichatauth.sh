#!/bin/sh
echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../weichatauth
npm install
npm run build
cd ../build

echo remove last package if exist
if [ -e package/web/weichatauth ]; then
  rm -rf package/web/weichatauth
fi

mv ../weichatauth/build ./package/web/weichatauth

echo weichatauth build over.
