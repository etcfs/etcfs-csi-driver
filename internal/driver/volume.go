package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"syscall"
)

// volumeIDPattern is deliberately narrow. A volume ID arrives from the API
// server and is turned into a filesystem path, so it is a trust boundary: this
// rejects separators, relative components and leading dots outright rather
// than trying to sanitise them afterwards. Kubernetes-generated names
// ("pvc-<uuid>") and hand-written static ones both fit.
var volumeIDPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,252}$`)

// volumeDir resolves a volume ID to its directory under the filesystem root.
func volumeDir(root, volumeID string) (string, error) {
	if volumeID == "" {
		return "", fmt.Errorf("volume ID is required")
	}
	if !volumeIDPattern.MatchString(volumeID) || volumeID == ".." {
		return "", fmt.Errorf("invalid volume ID %q: expected a single path component matching %s",
			volumeID, volumeIDPattern)
	}
	return filepath.Join(root, volumeID), nil
}

// isMountPoint reports whether path is the root of a mount, by comparing its
// device number with its parent's.
//
// The node plugin checks this before publishing. Without it, a host where the
// EtcFS daemon is down still has an ordinary empty directory at the mount
// path, and every pod scheduled there would get a private local directory that
// looks like a working shared filesystem — data written into it is invisible
// to every other node and lost with the host.
func isMountPoint(path string) (bool, error) {
	var st, parent syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		return false, err
	}
	if err := syscall.Stat(filepath.Dir(path), &parent); err != nil {
		return false, err
	}
	return st.Dev != parent.Dev, nil
}

// inodeOf returns the filesystem's own inode number for a path, which is what
// the quota records in etcd are keyed by.
func inodeOf(path string) (uint64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("stat %s: no inode information", path)
	}
	return st.Ino, nil
}
