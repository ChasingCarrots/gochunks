package gochunks

import (
	"io"

	"github.com/chasingcarrots/gobinary"
)

type ChunkReader struct {
	stream       *gobinary.StreamReader
	reader       gobinary.HighLevelReader
	header       chunkFileHeader
	atChunk      int
	nextPosition int64
}

func NewChunkReader(stream *gobinary.StreamReader) *ChunkReader {
	return &ChunkReader{
		stream: stream,
		reader: gobinary.MakeHighLevelReader(stream),
	}
}

func (cr *ChunkReader) Init() {
	cr.header.version = cr.reader.ReadUInt16()
	cr.header.numChunks = cr.reader.ReadUInt16()
	cr.nextPosition = cr.stream.Offset()
}

func (cr *ChunkReader) Version() uint16 {
	return cr.header.version
}

func (cr *ChunkReader) NumChunks() int { return int(cr.header.numChunks) }

func (cr *ChunkReader) NextChunk() *Chunk {
	if cr.atChunk >= cr.NumChunks() {
		return nil
	}
	cr.atChunk++
	cr.stream.Seek(cr.nextPosition, io.SeekStart)
	ch := Chunk{}
	ch.name = cr.reader.ReadString(int(cr.reader.ReadUInt16()))
	ch.version = cr.reader.ReadUInt16()
	ch.length = cr.reader.ReadUInt32()
	cr.nextPosition = cr.stream.Offset() + int64(ch.length)

	// setup chunk reader
	streamCopy := *cr.stream
	streamCopy.SeekCurrent()
	ch.stream = gobinary.MakeStreamReaderView(&streamCopy)
	ch.reader = gobinary.MakeHighLevelReader(&streamCopy)
	ch.stream.ViewHere()
	return &ch
}
