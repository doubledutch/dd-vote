#!/usr/bin/env bash

PROTOBUF3_DIR=~/protobuf3
pushd .
if [ -d "$PROTOBUF3_DIR" ] && [ -e "$PROTOBUF3_DIR/src/protoc" ]; then
    echo "Using cached protobuf3 build ..."
    cd $PROTOBUF3_DIR
else
    echo "Building protobuf3 from source ..."
    rm -rf $PROTOBUF3_DIR
    mkdir $PROTOBUF3_DIR

    # install some more dependencies required to build protobuf3
    apt-get install -y --no-install-recommends \
            curl \
            dh-autoreconf \
            unzip

    wget https://github.com/google/protobuf/archive/v3.0.0-beta-3.tar.gz -O protobuf3.tar.gz
    tar -xzf protobuf3.tar.gz -C $PROTOBUF3_DIR --strip 1
    rm protobuf3.tar.gz
    cd $PROTOBUF3_DIR
    ./autogen.sh
    ./configure --prefix=/usr
    make --jobs=$NUM_THREADS
fi
make install
popd
