package bot

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

// ImageCache caches images for continual use.
type ImageCache struct {
	images map[string][]byte
	mu     sync.Mutex
}

// NewImageCache creates a new image cache.
func NewImageCache() *ImageCache {
	return &ImageCache{
		images: make(map[string][]byte),
	}
}

// Download retrieves an image from the given URL
// and caches it.
func (i *ImageCache) Download(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.images[url] = b
	return nil
}

// Get returns a buffered image if it exists.
func (i *ImageCache) Get(name string) *bytes.Buffer {
	i.mu.Lock()
	defer i.mu.Unlock()

	old, ok := i.images[name]
	if !ok {
		return nil
	}

	new := make([]byte, len(old))
	copy(new, old)

	return bytes.NewBuffer(new)
}
