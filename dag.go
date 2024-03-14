package merkledag

import (
	"hash"
)

const (
	K          = 1 << 10
	BLOCK_SIZE = 256 * K
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// 将分片写入到KVStore中，并返回Merkle Root
	switch node.Type() {
	case FILE:
		return StoreFile(store, node, h)
	case DIR:
		return StoreDir(store, node, h)
	}
	return nil
}

func StoreFile(store KVStore, node File, h hash.Hash) []byte {
	t := []byte("blob")
	if node.Size() > BLOCK_SIZE {
		t = []byte("list")
	}

	// TODO: 实现将文件写入到 KVStore 中，并返回哈希值
	// 例如：store.Put(hash, node.Data())

	return h.Sum(nil), t
}

func StoreDir(store KVStore, dir DIR, h hash.Hash) []byte {
	// 定义树对象
	tree := Object{
		Links: make([]Link, 0),
		Data:  nil,
	}
	it := dir.It()
	for it.Next() {
		node := it.Node()
		if node.Type() == FILE {
			// 如果是文件，将文件添加到树的 Links 中
			fileHash, fileType := StoreFile(store, node.(File), h)
			tree.Links = append(tree.Links, Link{Name: node.Name(), Hash: fileHash, Size: node.Size()})
		} else if node.Type() == DIR {
			// 如果是目录，递归地将目录添加到树的 Links 中
			dirHash := StoreDir(store, node.(DIR), h)
			tree.Links = append(tree.Links, Link{Name: node.Name(), Hash: dirHash, Size: node.Size()})
		}
	}

	// TODO: 将树对象存储到 KVStore 中，并返回哈希值
	// 例如：store.Put(hash, Serialize(tree))

	return h.Sum(nil)
}
