package main

import _ "embed"

//go:embed third_party/e2fsprogs-1.46.5/mke2fs.arm64
var mke2fs []byte
