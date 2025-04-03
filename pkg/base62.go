package pkg

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// EncodeID 将ID转换为短链接key
// 参数：
//   - id: 需要转换的ID
//
// 返回：
//   - string: 生成的短链接key
func EncodeID(id int64) string {
	if id == 0 {
		return string(base62Chars[0])
	}
	var result []byte
	for id > 0 {
		// 将id对62取余，得到一个0-61之间的数字
		// 将这个数字转换为base62Chars中的一个字符
		// 将这个字符添加到result中
		// 将id除以62
		result = append([]byte{base62Chars[id%62]}, result...)
		id /= 62
	}
	return string(result)
}
