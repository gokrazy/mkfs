package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gokrazy/gokapi"
	"github.com/gokrazy/gokapi/ondeviceapi"
	"github.com/gokrazy/internal/rootdev"
)

func makeFilesystemNotWar() error {
	b, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return err
	}
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		mountpoint := parts[4]
		log.Printf("Found mountpoint %q", parts[4])
		if mountpoint == "/perm" {
			log.Printf("/perm file system already mounted, nothing to do")
			return nil
		}
	}

	// /perm is not a mounted file system. Try to create a file system.
	dev := rootdev.Partition(rootdev.Perm)
	log.Printf("No /perm mountpoint found. Creating file system on %s", dev)

	mkfs := exec.Command("/usr/local/bin/mke2fs", "-t", "ext4", dev)
	mkfs.Stdout = os.Stdout
	mkfs.Stderr = os.Stderr
	log.Printf("%v", mkfs.Args)
	if err := mkfs.Run(); err != nil {
		return fmt.Errorf("%v: %v", mkfs.Args, err)
	}

	// It is pointless to try and mount the file system here from within this
	// process, as gokrazy services are run in a separate mount namespace.
	// Instead, we trigger a reboot so that /perm is mounted early and
	// the whole system picks it up correctly.
	log.Printf("triggering reboot to mount /perm")
	cfg, err := gokapi.ConnectOnDevice()
	if err != nil {
		return err
	}
	cl := ondeviceapi.NewAPIClient(cfg)
	_, err = cl.UpdateApi.Reboot(context.Background(), &ondeviceapi.UpdateApiRebootOpts{})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := makeFilesystemNotWar(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	// tell gokrazy to not supervise this service, it’s a one-off:
	os.Exit(125)
}
