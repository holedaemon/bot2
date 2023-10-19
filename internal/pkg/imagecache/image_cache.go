package imagecache

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/holedaemon/bot2/internal/pkg/httpx"
	"github.com/jellydator/ttlcache/v3"
	"github.com/zikaeroh/ctxlog"
	"go.uber.org/zap"
)

var ErrContentType = errors.New("cache: Content-Type of image is not recognized")

// Option configures a Cache.
type Option func(*Cache)

// WithTTL sets the TTL on the cache.
func WithTTL(dur time.Duration) Option {
	return func(c *Cache) {
		c.ttl = dur
	}
}

// WithHTTPClient specifies the HTTP Client to use for image downloading.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Cache) {
		c.client = client
	}
}

type Cache struct {
	ttl    time.Duration
	cache  *ttlcache.Cache[string, []byte]
	client *http.Client
}

func New(opts ...Option) *Cache {
	c := &Cache{
		ttl: time.Hour * 1,
	}

	for _, o := range opts {
		o(c)
	}

	c.cache = ttlcache.New[string, []byte](
		ttlcache.WithTTL[string, []byte](c.ttl),
	)

	if c.client == nil {
		c.client = httpx.DefaultClient
	}

	return c
}

// Start begins deleting expired items from the cache.
func (c *Cache) Start() {
	c.cache.Start()
}

func (c *Cache) download(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if !httpx.IsOK(res.StatusCode) {
		return httpx.ErrStatus
	}

	image, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	mime := http.DetectContentType(image)
	switch mime {
	case "image/gif", "image/jpeg", "image/png":
	default:
		return ErrContentType
	}

	c.cache.Set(url, image, ttlcache.DefaultTTL)
	return nil
}

// Get returns a buffered image from the cache.
// The image will be downloaded if it doesn't already exist within.
func (c *Cache) Get(ctx context.Context, url string) (*bytes.Buffer, error) {
	item := c.cache.Get(url)
	if item == nil {
		ctxlog.Debug(ctx, "image does not exist in cache, downloading...", zap.String("url", url))
		if err := c.download(ctx, url); err != nil {
			return nil, err
		}

		return c.Get(ctx, url)
	}

	old := item.Value()
	new := make([]byte, len(old))
	copy(new, old)

	buf := bytes.NewBuffer(new)
	return buf, nil
}
