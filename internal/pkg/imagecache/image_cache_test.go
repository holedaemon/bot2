package imagecache

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zikaeroh/ctxlog"
	"gotest.tools/v3/assert"
)

func TestImageCache(t *testing.T) {
	c := New(
		WithTTL(time.Second * 30),
	)

	go c.Start()

	ctx := ctxlog.WithLogger(context.Background(), ctxlog.New(true))

	t.Run("GetJPEG", func(t *testing.T) {
		url := "https://holedaemon.net/images/yousuck.jpg"
		image, err := c.Get(ctx, url)
		assert.NilError(t, err, "getting image")

		err = writeFile(image, "jpg")
		assert.NilError(t, err, "writing file")

		time.Sleep(time.Second * 30)

		image, err = c.Get(ctx, url)
		assert.NilError(t, err, "getting image again")

		err = writeFile(image, "jpg")
		assert.NilError(t, err, "writing file again")
	})

	t.Run("GetPNG", func(t *testing.T) {
		url := "https://holedaemon.net/images/egopussy.png"
		image, err := c.Get(ctx, url)
		assert.NilError(t, err, "getting image")

		err = writeFile(image, "png")
		assert.NilError(t, err, "writing file")
	})

	t.Run("GetGif", func(t *testing.T) {
		url := "https://holedaemon.net/images/fortnite-1.gif"
		image, err := c.Get(ctx, url)
		assert.NilError(t, err, "getting image")

		err = writeFile(image, "gif")
		assert.NilError(t, err, "writing file")
	})
}

func writeFile(image *bytes.Buffer, ext string) error {
	fileName := fmt.Sprintf("out%d.%s", time.Now().Unix(), ext)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	file.Write(image.Bytes()) //nolint:errcheck
	return nil
}
