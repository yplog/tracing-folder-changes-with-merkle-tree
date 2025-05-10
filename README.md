## Merkle Tree Folder Watcher

This is a simple Go program that watches a folder and calculates its Merkle root based on the content of its files. It also lets you verify Merkle proofs and visualize the tree structure in your terminal.

## Features

- Watch mode: Detects changes in folder content and updates the Merkle root
- Proof mode: Verifies whether a value is part of the Merkle tree
- Visualize mode: Prints the full Merkle tree structure in the terminal

## How to Run

```bash
git clone 

cd tracing-folder-changes-with-merkle-tree

go build -o ./build/tfcwmt
```

### Watch Mode

Watches a folder for any file changes:
```bash

./tfcwmt watch <file_path>
```

You should see output like this:

```bash
Watching folder: watched
Merkle Root initialized: <some_hash>
Merkle Root changed: <new_hash>
```
### Proof Mode

Verifies if a given value (file content) is part of the Merkle tree:

Note: Ensure the value matches the file content exactly (including newlines).

```bash
./tfcwmt proof <file_path> <value> <root_hash>

./tfcwmt proof watched 'hello' <root_hash>
```

You should see output like this:

```bash
Merkle proof steps:
 - <hash> (right)
 - <hash> (left)
Proof is valid: true
```

### Visualize Mode
Prints the Merkle tree structure to the terminal:

```bash
./tfcwmt visualize <file_path>
```
You should see output like this:

```bash
Merkle Tree Visualization:
├── [Leaf: <value>]  <hash>
└── [Leaf: <value>]  <hash>
```

Read the full blog post [here](https://yalinpala.dev/blog/tracking-folder-changes-with-merkle-trees-in-go).
