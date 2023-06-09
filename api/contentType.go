package api

import (
	"fmt"
	"net/http"
	"os"
)

/*
文件类型识别：
	1.根据文件的扩展名来判断
		比如 .jpg 表示 JPEG 图像，.mp3 表示 MP3 音频等
	2.根据文件内容的特征来判断
		比如 JPEG 图像的前几个字节是 FF D8 FF，MP3 音频的前几个字节是 ID3 等
	3.Go 标准库：net/http DetectContentType()
		它读取文件内容的前 512 个字节内容，返回一个 MIME 类型字符串，例如 image/jpeg
		它使用了 mimesniff 算法²⁴，根据一组预定义的规则来匹配文件内容的特征和对应的 MIME 类型
		既不依赖于文件扩展名，也不需要完整地读取文件内容，因此既快速又准确
*/

func GetContentType(f string) {
	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		return
	}
	buff := make([]byte, 512)
	n, err := file.Read(buff)
	if err != nil {
		return
	}
	// 实际上，如果字节数超过 512，该函数也只会使用前 512 个字节
	contentType := http.DetectContentType(buff[:n])
	fmt.Println(contentType)
}
