package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	tmp, err := os.MkdirTemp("", "gokrazy-mkfs-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	log.Printf("Writing self-contained mke2fs to %s", tmp)

	if err := ioutil.WriteFile(filepath.Join(tmp, "mke2fs"), mke2fs, 0755); err != nil {
		return err
	}
	mkfs := exec.Command(filepath.Join(tmp, "mke2fs"), "-t", "ext4", dev)
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
	httpPassword, err := readConfigFile("gokr-pw.txt")
	if err != nil {
		return fmt.Errorf("could read neither /perm/gokr-pw.txt, nor /etc/gokr-pw.txt, nor /gokr-pw.txt: %v", err)
	}

	port, err := ioutil.ReadFile("/etc/http-port.txt")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:"+string(port)+"/reboot", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("gokrazy:"+httpPassword)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		return fmt.Errorf("rebooting device: unexpected HTTP status code: got %d, want %d", got, want)
	}

	return nil
}

// readConfigFile reads configuration files from /perm /etc or / and returns
// trimmed content as string.
//
// TODO: de-duplicate this with gokrazy.go into a gokrazy/internal package
func readConfigFile(fileName string) (string, error) {
	str, err := ioutil.ReadFile("/perm/" + fileName)
	if err != nil {
		str, err = ioutil.ReadFile("/etc/" + fileName)
	}
	if err != nil && os.IsNotExist(err) {
		str, err = ioutil.ReadFile("/" + fileName)
	}

	return strings.TrimSpace(string(str)), err
}

func main() {
	if err := makeFilesystemNotWar(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	// tell gokrazy to not supervise this service, it’s a one-off:
	os.Exit(125)
}
