third_party/e2fsprogs-1.46.5/mke2fs.amd64: Dockerfile
	docker build --rm -t e2fsprogs .
	docker run --rm -v $$(pwd)/third_party/e2fsprogs-1.46.5:/tmp/bins e2fsprogs cp -r /bins/ /tmp/

clean:
	rm -f third_party/e2fsprogs-1.46.5/*
