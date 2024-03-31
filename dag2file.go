
package merkledag

// Hash2File 根据hash和path，返回对应的文件内容
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 从 KVStore 中获取哈希对应的对象
	objBytes, err := store.Get(hash)
	if err != nil {
		return nil
	}

	// 反序列化对象
	obj := Deserialize(objBytes)

	// 如果路径为空，则返回整个对象的内容
	if path == "" {
		return obj.Data
	}

	// 按路径查找文件内容
	currentObj := obj
	pathComponents := strings.Split(path, "/")
	for _, component := range pathComponents {
		if component == "" {
			continue
		}

		found := false
		for _, link := range currentObj.Links {
			if link.Name == component {
				// 如果是文件，则返回文件内容
				if link.Size != 0 {
					fileObjBytes, err := store.Get(link.Hash)
					if err != nil {
						return nil
					}
					return Deserialize(fileObjBytes).Data
				}
				// 如果是目录，则继续向下查找
				hash = link.Hash
				objBytes, err := store.Get(hash)
				if err != nil {
					return nil
				}
				currentObj = Deserialize(objBytes)
				found = true
				break
			}
		}

		// 如果路径中的某个组件未找到，则返回错误
		if !found {
			return nil
		}
	}

	// 如果路径找到了，但是对应的文件内容为空，则返回错误
	return nil
}
