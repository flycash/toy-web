package web

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StaticResourceHandlerOption func(h *StaticResourceHandler)

type StaticResourceHandler struct {
	dir string
	pathPrefix string
	extensionContentTypeMap map[string]string

	// 缓存静态资源的限制
	cache *lru.Cache
	maxFileSize int
}

type fileCacheItem struct {
	fileName string
	fileSize int
	contentType string
	data []byte
}

func NewStaticResourceHandler(dir string, pathPrefix string,
	options...StaticResourceHandlerOption) *StaticResourceHandler {
	res := &StaticResourceHandler{
		dir: dir,
		pathPrefix: pathPrefix,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe": "image/jpeg",
			"jpg": "image/jpeg",
			"png": "image/png",
			"pdf": "image/pdf",
		},
	}

	for _, o := range options {
		o(res)
	}
	return res
}
// WithFileCache 静态文件将会被缓存
// maxFileSizeThreshold 超过这个大小的文件，就被认为是大文件，我们将不会缓存
// maxCacheFileCnt 最多缓存多少个文件
// 所以我们最多缓存 maxFileSizeThreshold * maxCacheFileCnt
func WithFileCache(maxFileSizeThreshold int, maxCacheFileCnt int) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		c, err := lru.New(maxCacheFileCnt)
		if err != nil {
			fmt.Printf("could not create LRU, we won't cache static file")
		}
		h.maxFileSize = maxFileSizeThreshold
		h.cache = c
	}
}

func WithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		for ext, contentType := range extMap {
			h.extensionContentTypeMap[ext] = contentType
		}
	}
}

func (h *StaticResourceHandler) ServeStaticResource(c *Context)  {
	req := strings.TrimPrefix(c.R.URL.Path, h.pathPrefix)
	if item, ok := h.readFileFromData(req); ok {
		fmt.Printf("read data from cache...")
		h.writeItemAsResponse(item, c.W)
		return
	}
	path := filepath.Join(h.dir, req)
	f, err := os.Open(path)
	if err != nil {
		c.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	ext := getFileExt(f.Name())
	t, ok := h.extensionContentTypeMap[ext]
	if !ok {
		c.W.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	item := &fileCacheItem{
		fileSize: len(data),
		data: data,
		contentType: t,
		fileName: req,
	}

	h.cacheFile(item)
	h.writeItemAsResponse(item, c.W)

}

func (h *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if h.cache != nil && item.fileSize < h.maxFileSize {
		h.cache.Add(item.fileName, item)
	}
}

func (h *StaticResourceHandler) writeItemAsResponse(item *fileCacheItem, writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", item.contentType)
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", item.fileSize))
	_, _ = writer.Write(item.data)

}

func (h *StaticResourceHandler) readFileFromData(fileName string) (*fileCacheItem, bool) {
	if h.cache != nil {
		if item, ok := h.cache.Get(fileName); ok {
			return item.(*fileCacheItem), true
		}
	}
	return nil, false
}

func getFileExt(name string) string {
	index := strings.LastIndex(name, ".")
	if index == len(name) - 1{
		return ""
	}
	return name[index+1:]
}
