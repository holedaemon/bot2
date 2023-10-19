package topster

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	_ "image/jpeg"

	"gotest.tools/v3/assert"
)

func TestTopster(t *testing.T) {
	url := os.Getenv("TOPSTER_ADDR")
	user := os.Getenv("TOPSTER_USER")
	assert.Assert(t, user != "" && url != "", "invalid env")

	c, err := New(url)
	assert.NilError(t, err, "creating topster client")

	chart, err := c.CreateChart(context.Background(), user, Title("Testing chart"))
	assert.NilError(t, err, "creating chart")

	raw, err := base64.StdEncoding.DecodeString(chart)
	assert.NilError(t, err, "decoding chart to []byte")

	fileName := fmt.Sprintf("topster%d.jpg", time.Now().Unix())
	file, err := os.Create(fileName)
	assert.NilError(t, err, "creating file")

	defer file.Close()

	file.Write(raw)
}
