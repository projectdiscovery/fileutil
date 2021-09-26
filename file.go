package fileutil

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/karrick/godirwalk"
)

// FileExists checks if the file exists in the provided path
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// FolderExists checks if the folder exists
func FolderExists(foldername string) bool {
	info, err := os.Stat(foldername)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
}

func DeleteFilesOlderThan(folder string, maxAge time.Duration, callback func(string)) error {
	startScan := time.Now()
	return godirwalk.Walk(folder, &godirwalk.Options{
		Unsorted:            true,
		FollowSymbolicLinks: false,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if osPathname == "" {
				return nil
			}
			if de.IsDir() {
				return nil
			}
			fileInfo, err := os.Stat(osPathname)
			if err != nil {
				return nil
			}
			if fileInfo.ModTime().Add(maxAge).Before(startScan) {
				os.RemoveAll(osPathname)
				if callback != nil {
					callback(osPathname)
				}
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})
}

// DownloadFile to specified path
func DownloadFile(filepath string, url string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func CreateFolders(paths ...string) error {
	for _, path := range paths {
		if err := CreateFolder(path); err != nil {
			return err
		}
	}

	return nil
}

func CreateFolder(path string) error {
	return os.MkdirAll(path, 0700)
}

// HasStdin determines if the user has piped input
func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	isPipedFromChrDev := (mode & os.ModeCharDevice) == 0
	isPipedFromFIFO := (mode & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

func ReadFileWithReader(r io.Reader) (chan string, error) {
	out := make(chan string)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}

func ReadFileWithReaderAndBufferSize(r io.Reader, maxCapacity int) (chan string, error) {
	out := make(chan string)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(r)
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}

func ReadFile(filename string) (chan string, error) {
	if !FileExists(filename) {
		return nil, errors.New("file doesn't exist")
	}
	out := make(chan string)
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		if err != nil {
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}

func ReadFileWithBufferSize(filename string, maxCapacity int) (chan string, error) {
	if !FileExists(filename) {
		return nil, errors.New("file doesn't exist")
	}
	out := make(chan string)
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		if err != nil {
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}
