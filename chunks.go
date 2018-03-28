package gochunks

import (
	"fmt"

	"github.com/chasingcarrots/gobinary"
)

type chunkFileHeader struct {
	version   uint16
	numChunks uint16
}

type Chunk struct {
	name    string
	version uint16
	length  uint32

	stream gobinary.StreamReaderView
	reader gobinary.HighLevelReader
}

func (c *Chunk) Name() string {
	return c.name
}

func (c *Chunk) Version() uint16 {
	return c.version
}

func (c *Chunk) CheckVersion(version uint16) error {
	if version < c.version {
		return fmt.Errorf("trying to read chunk %s with version %d, but the current version is only %d", c.name, c.version, version)
	}
	return nil
}

func (c *Chunk) Reader() *gobinary.HighLevelReader {
	return &c.reader
}

func (c *Chunk) Stream() *gobinary.StreamReaderView {
	return &c.stream
}

func (c *Chunk) Length() uint32 {
	return c.length
}
