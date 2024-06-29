all: _gokrazy/extrafiles_amd64.tar _gokrazy/extrafiles_arm64.tar _gokrazy/extrafiles_arm.tar

third_party/e2fsprogs-1.47.1/mke2fs.amd64: Dockerfile
	docker build --rm -t e2fsprogs .
	docker run --rm -v $$(pwd)/third_party/e2fsprogs-1.47.1:/tmp/bins e2fsprogs cp -r /bins/ /tmp/

_gokrazy/extrafiles_amd64.tar: third_party/e2fsprogs-1.47.1/mke2fs.amd64 third_party/e2fsprogs-1.47.1/e2fsck.amd64
	mkdir -p _gokrazy
	tar cf $@ $^ "--transform=s,third_party/e2fsprogs-1.47.1/\([^.]*\)\..*,usr/local/bin/\1,g"

_gokrazy/extrafiles_arm64.tar: third_party/e2fsprogs-1.47.1/mke2fs.arm64 third_party/e2fsprogs-1.47.1/e2fsck.arm64
	mkdir -p _gokrazy
	tar cf $@ $^ "--transform=s,third_party/e2fsprogs-1.47.1/\([^.]*\)\..*,usr/local/bin/\1,g"

_gokrazy/extrafiles_arm.tar: third_party/e2fsprogs-1.47.1/mke2fs.arm third_party/e2fsprogs-1.47.1/e2fsck.arm
	mkdir -p _gokrazy
	tar cf $@ $^ "--transform=s,third_party/e2fsprogs-1.47.1/\([^.]*\)\..*,usr/local/bin/\1,g"

clean:
	rm -f third_party/e2fsprogs-1.47.1/*
