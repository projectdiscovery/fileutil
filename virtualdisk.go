package fileutil

import (
	"fmt"

	"github.com/projectdiscovery/chaos-middleware/exec"
)

const (
	BTRFS = "btrfs"
	EXT4  = "ext4"
)

type VirtualDiskLayout struct {
	FileSystem string
	Size       string
	FileDisk   string
	MountPoint string
	Truncate   bool
	Mount      bool
}

func CreateVirtualDisk(vlayout *VirtualDiskLayout) (string, error) {
	var fullOutput string
	if vlayout.Truncate {
		out, err := Truncate(vlayout.FileDisk, vlayout.Size)
		if err != nil {
			return out, err
		}
		fullOutput = fmt.Sprintf("%s\n", out)
	}

	if FileExists(vlayout.FileDisk) {
		out, err := Mkfs(vlayout.FileDisk, BTRFS)
		if err != nil {
			return out, err
		}
		fullOutput += fmt.Sprintf("%s\n", out)
	}

	return fullOutput, nil
}

func Truncate(path, size string) (string, error) {
	// Create empty file-disk
	truncateCmd := fmt.Sprintf("truncate -s %sG %s", size, path)
	output, err := exec.Run(truncateCmd)
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func Mkfs(path, fs string) (string, error) {
	// Format it to Btrfs
	formatCmd := fmt.Sprintf("mkfs -t %s %s", fs, path)
	output, err := exec.Run(formatCmd)
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func Mount(file, toFolder string) (string, error) {
	// Create destination folder if absent
	if err := CreateFolder(toFolder); err != nil {
		return "", err
	}

	// Create empty file-disk
	truncateCmd := fmt.Sprintf("mount %s %s", file, toFolder)
	output, err := exec.Run(truncateCmd)
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
