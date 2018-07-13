// Licensed under the  MIT License (MIT)
// Copyright (c) 2016 Peter Waller <p@pwaller.net>

package main

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/builder/dockerignore"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/fileutils"
)

// getArchive returns the tarfile io.ReadCloser. It is a direct copy of the
// logic found in the official docker client.
// See <https://github.com/docker/docker/blob/78f2b8d8/api/client/build.go#L126-L172>.
func getArchive(contextDir, relDockerfile string) (io.ReadCloser, error) {
	var err error

	// And canonicalize dockerfile name to a platform-independent one
	relDockerfile = archive.CanonicalTarNameForPath(relDockerfile)

	f, err := os.Open(filepath.Join(contextDir, ".dockerignore"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var excludes []string
	if err == nil {
		excludes, err = dockerignore.ReadAll(f)
		if err != nil {
			return nil, err
		}
	}

	// If .dockerignore mentions .dockerignore or the Dockerfile
	// then make sure we send both files over to the daemon
	// because Dockerfile is, obviously, needed no matter what, and
	// .dockerignore is needed to know if either one needs to be
	// removed. The daemon will remove them for us, if needed, after it
	// parses the Dockerfile. Ignore errors here, as they will have been
	// caught by validateContextDirectory above.
	var includes = []string{"."}
	keepThem1, _ := fileutils.Matches(".dockerignore", excludes)
	keepThem2, _ := fileutils.Matches(relDockerfile, excludes)
	if keepThem1 || keepThem2 {
		includes = append(includes, ".dockerignore", relDockerfile)
	}

	return archive.TarWithOptions(contextDir, &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: excludes,
		IncludeFiles:    includes,
	})
}

// WriteCounter counts the bytes written to it.
type WriteCounter int

func (w *WriteCounter) Write(bs []byte) (int, error) {
	*w += WriteCounter(len(bs))
	return len(bs), nil
}

func main() {
	log.SetFlags(0)

	// Take a quick and dirty file count. This should be an over-estimate,
	// since it doesn't currently attempt to re-implement or reuse the
	// dockerignore logic.
	totalCount := 0
	totalStorage := int64(0)
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		totalCount++
		totalStorage += info.Size()
		return nil
	})

	// TODO(pwaller): Make these parameters?
	r, err := getArchive(".", "Dockerfile")
	if err != nil {
		log.Fatalf("Failed to make context: %v", err)
	}
	defer r.Close()

	// Keep mappings of paths/extensions to bytes/counts/times.
	dirStorage := map[string]int64{}
	dirFiles := map[string]int64{}
	dirTime := map[string]int64{}
	extStorage := map[string]int64{}

	// Counts of amounts seen so far.
	currentCount := 0
	currentStorage := int64(0)

	// Update the progress indicator at some frequency.
	const updateFrequency = 50 // Hz
	ticker := time.NewTicker(time.Second / updateFrequency)
	defer ticker.Stop()
	tick := ticker.C

	start := time.Now()
	last := time.Now()

	// Keep a count of how many bytes of Tar file were seen.
	writeCounter := WriteCounter(0)
	tf := tar.NewReader(io.TeeReader(r, &writeCounter))

	cr := []byte("\r")
	showUpdate := func(w io.Writer) {
		os.Stderr.Write(cr) // always to Stderr.
		fmt.Fprintf(w,
			"  %v / %v (%.0f / %.0f MiB) "+
				"(%.1fs elapsed)",
			currentCount,
			totalCount,
			float64(currentStorage)/1024/1024,
			float64(totalStorage)/1024/1024,
			time.Since(start).Seconds(),
		)
	}

	fmt.Println()
	fmt.Println("Scanning local directory (in tar / on disk):")
entries:
	for {
		header, err := tf.Next()
		switch err {
		case io.EOF:
			showUpdate(os.Stdout)
			fmt.Println(" .. completed")
			fmt.Println()
			break entries
		default:
			log.Fatalf("Error reading archive: %v", err)
			return
		case nil:
		}

		duration := time.Since(last).Nanoseconds()
		last = time.Now()

		dir := filepath.Dir(header.Name)
		size := header.FileInfo().Size()

		currentCount++
		currentStorage += size

		dirStorage[dir] += size
		dirTime[dir] += duration
		dirFiles[dir]++

		if !header.FileInfo().IsDir() {
			ext := filepath.Ext(strings.ToLower(header.Name))
			extStorage[ext] += size
		}

		select {
		case <-tick:
			showUpdate(os.Stderr)
		default:
		}
	}

	fmt.Printf(
		"Excluded by .dockerignore: "+
			"%d files totalling %.2f MiB\n",
		totalCount-currentCount,
		float64(totalStorage-currentStorage)/1024/1024)
	fmt.Println()
	fmt.Println("Final .tar:")
	// Epilogue.
	fmt.Printf(
		"  %v files totalling %.2f MiB (+ %.2f MiB tar overhead)\n",
		currentCount,
		float64(currentStorage)/1024/1024,
		float64(int64(writeCounter)-currentStorage)/1024/1024,
	)
	fmt.Printf("  Took %.2f seconds to build\n", time.Since(start).Seconds())
	fmt.Println()

	// Produce Top-N.
	topDirStorage := SortedBySize(dirStorage)
	topDirFiles := SortedBySize(dirFiles)
	topDirTime := SortedBySize(dirTime)
	topExtStorage := SortedBySize(extStorage)

	const N = 10

	fmt.Printf("Top %d directories by time spent:\n", N)
	for i := 0; i < N && i < len(topDirTime); i++ {
		entry := &topDirTime[i]
		fmt.Printf("%5d ms: %v\n", entry.Size/1000/1000, entry.Path)
	}
	fmt.Println()

	fmt.Printf("Top %d directories by storage:\n", N)
	for i := 0; i < N && i < len(topDirStorage); i++ {
		entry := &topDirStorage[i]
		fmt.Printf("%7.2f MiB: %v\n", float64(entry.Size)/1024/1024, entry.Path)
	}
	fmt.Println()

	fmt.Printf("Top %d directories by file count:\n", N)
	for i := 0; i < N && i < len(topDirFiles); i++ {
		entry := &topDirFiles[i]
		fmt.Printf("%5d: %v\n", entry.Size, entry.Path)
	}
	fmt.Println()

	fmt.Printf("Top %d file extensions by storage:\n", N)
	for i := 0; i < N && i < len(topExtStorage); i++ {
		entry := &topExtStorage[i]
		fmt.Printf("%7.2f MiB: %v\n", float64(entry.Size)/1024/1024, entry.Path)
	}
	fmt.Println()
}

// SortedBySize returns direcotries in m sorted by Size (biggest first).
func SortedBySize(m map[string]int64) []PathSize {
	bySize := BySize{}
	for dir, size := range m {
		bySize = append(bySize, PathSize{dir, size})
	}
	sort.Sort(bySize)
	return []PathSize(bySize)
}

// PathSize represents a directory with a size.
type PathSize struct {
	Path string
	Size int64
}

// BySize sorts by size (biggest first).
type BySize []PathSize

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size > a[j].Size }
