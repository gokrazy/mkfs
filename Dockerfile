FROM debian:bullseye

RUN apt-get update && apt-get install -y crossbuild-essential-armel crossbuild-essential-arm64 build-essential curl

RUN mkdir /src /bins && curl -s https://mirrors.edge.kernel.org/pub/linux/kernel/people/tytso/e2fsprogs/v1.46.5/e2fsprogs-1.46.5.tar.gz | tar --strip-components=1 -C /src -xzf -

ENV SOURCE_DATE_EPOCH 1600000000

RUN cd /src && ./configure CFLAGS='-O2 -static' LDFLAGS=-static CC=arm-linux-gnueabi-gcc --host=arm-linux-gnueabi && make -j$(nproc) && cp ./misc/mke2fs /bins/mke2fs.arm && make distclean
RUN cd /src && ./configure CFLAGS='-O2 -static' LDFLAGS=-static CC=aarch64-linux-gnu-gcc --host=aarch64-linux-gnu && make -j$(nproc) && cp ./misc/mke2fs /bins/mke2fs.arm64 && make distclean
RUN cd /src && ./configure CFLAGS='-O2 -static' LDFLAGS=-static && make -j$(nproc) && cp ./misc/mke2fs /bins/mke2fs.amd64 && make distclean
