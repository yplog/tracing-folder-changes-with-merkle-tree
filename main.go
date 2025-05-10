package main

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func readFiles(dir string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(dir, path)
		files[relPath] = content

		return nil
	})

	return files, err
}

func filesToBlocks(files map[string][]byte) [][]byte {
	var keys []string
	for k := range files {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var blocks [][]byte
	for _, k := range keys {
		blocks = append(blocks, files[k])
	}

	return blocks
}

func runWatchMode(folder string) {
	fmt.Println("Watching folder:", folder)

	var lastRoot string
	var initialized bool

	for {
		files, err := readFiles(folder)
		if err != nil {
			log.Println("Error reading files:", err)
			continue
		}

		blocks := filesToBlocks(files)
		tree := NewMerkleTree(blocks)
		newRoot := hex.EncodeToString(tree.Root.Hash)

		if !initialized {
			fmt.Println("Merkle Root initialized:", newRoot)
			initialized = true
		} else if newRoot != lastRoot {
			fmt.Println("Merkle Root changed:", newRoot)
		}

		lastRoot = newRoot
		time.Sleep(2 * time.Second)
	}
}

func runProofMode(folder string, value string, expectedRoot string) {
	files, err := readFiles(folder)
	if err != nil {
		log.Fatal("Error reading files:", err)
	}

	blocks := filesToBlocks(files)
	tree := NewMerkleTree(blocks)

	proof, ok := tree.GenerateProof([]byte(value))
	if !ok {
		fmt.Println("No proof found for the given value.")
		return
	}

	fmt.Println("Merkle proof steps:")
	for _, p := range proof {
		fmt.Printf(" - %s (%s)\n", hex.EncodeToString(p.Hash), ternary(p.IsLeft, "left", "right"))
	}

	rootHash, err := hex.DecodeString(expectedRoot)
	if err != nil {
		log.Fatal("Invalid root hash provided.")
	}

	valid := VerifyProof([]byte(value), proof, rootHash)
	fmt.Println("Proof is valid:", valid)
}

func runVisualizeMode(folder string) {
	files, err := readFiles(folder)
	if err != nil {
		log.Fatal("Error reading files:", err)
	}
	blocks := filesToBlocks(files)
	tree := NewMerkleTree(blocks)

	fmt.Println("Merkle Tree Visualization:")
	printMerkleTree(tree.Root, "", true)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  go run . watch <folder>")
		fmt.Println("  go run . proof <folder> <value> <expected-root-hash>")
		fmt.Println("  go run . visualize <folder>")
		return
	}

	mode := os.Args[1]
	switch mode {
	case "watch":
		runWatchMode(os.Args[2])
	case "proof":
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run . proof <folder> <value> <expected-root-hash>")
			return
		}
		runProofMode(os.Args[2], os.Args[3], os.Args[4])
	case "visualize":
		if len(os.Args) != 3 {
			fmt.Println("Usage: go run . visualize <folder>")
			return
		}
		runVisualizeMode(os.Args[2])
	default:
		fmt.Println("Unknown mode:", mode)
	}
}
