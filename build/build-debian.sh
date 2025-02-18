#!/bin/bash
set -xe

PACKAGE="AppPlaygroundService"
TMP_DIR="tmp"
VERSION_STRING="$(cat VERSION)"
BINARY_NAME="AppPlaygroundService"
DEB_PACKAGE_DESCRIPTION="AppPlaygroundService"
GLOBAL_CONFIG_FILE="app-playground-service.yaml"
SYSTEMD_FILE="AppPlaygroundService.service"

rm -rf ${PACKAGE}-${VERSION_STRING}.deb

mkdir -p $TMP_DIR
mkdir -p "$TMP_DIR/usr/bin"
mkdir -p "$TMP_DIR/etc/ASUS"
mkdir -p "$TMP_DIR/etc/systemd/system"
cp -p "tmp/${PACKAGE}_${VERSION_STRING}" "$TMP_DIR/usr/bin/$BINARY_NAME"
cp -p "etc/$GLOBAL_CONFIG_FILE" "$TMP_DIR/etc/ASUS/$GLOBAL_CONFIG_FILE"
cp -p "systemd/$SYSTEMD_FILE" "$TMP_DIR/etc/systemd/system/$SYSTEMD_FILE"

fpm -t deb \
    -s dir \
    -C $TMP_DIR \
    --name $BINARY_NAME \
    --version $VERSION_STRING \
    --description "$DEB_PACKAGE_DESCRIPTION" \
    -p ${PACKAGE}-${VERSION_STRING}.deb \
    .

mv ${PACKAGE}-${VERSION_STRING}.deb $TMP_DIR/${PACKAGE}-${VERSION_STRING}.deb
rm -rf $TMP_DIR/etc
rm -rf $TMP_DIR/usr