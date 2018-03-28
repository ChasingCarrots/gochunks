package gochunks

import (
	"io"

	"github.com/chasingcarrots/gobinary"
)

type ChunkWriter struct {
	stream        *gobinary.StreamWriter
	writer        gobinary.HighLevelWriter
	header        chunkFileHeader
	atChunk       int
	initialOffset int64

	chunkAvailable  bool
	lastChunkOffset int64
}

func NewChunkWriter(stream *gobinary.StreamWriter, version uint16) *ChunkWriter {
	return &ChunkWriter{
		lastChunkOffset: -1,
		stream:          stream,
		header:          chunkFileHeader{version: version},
		writer:          gobinary.MakeHighLevelWriter(stream),
	}
}

func (cw *ChunkWriter) Init() {
	cw.header.numChunks = 0
	cw.initialOffset = cw.stream.Offset()
	cw.writer.WriteUInt16(cw.header.version)
	// dummy value for the length, will be overwritten when all chunks have been
	// written.
	cw.writer.WriteUInt16(cw.header.numChunks)
}

func (cw *ChunkWriter) BeginChunk(name string, version uint16) ChunkHandle {
	if cw.lastChunkOffset >= 0 {
		panic("Cannot begin new chunk! You need to finish the chunk you started earlier!")
	}
	cw.stream.SeekCurrent()
	cw.writer.WriteUInt16(uint16(len(name)))
	cw.writer.WriteString(name)
	cw.writer.WriteUInt16(version)
	// length dummy value, will be overwritten when chunk is finished
	cw.writer.WriteUInt32(0)
	cw.lastChunkOffset = cw.stream.Offset()
	cw.header.numChunks++

	streamCopy := *cw.stream
	handle := ChunkHandle{
		stream: gobinary.MakeStreamWriterView(&streamCopy),
		writer: gobinary.MakeHighLevelWriter(&streamCopy),
	}
	handle.stream.ViewHere()
	return handle
}

func (cw *ChunkWriter) EndChunk(handle ChunkHandle) {
	if cw.lastChunkOffset < 0 {
		panic("No chunk has been started; there is no way to end one!")
	}
	if handle.stream.Base() != cw.lastChunkOffset {
		panic("Cannot end this chunk! This is not the last chunk that was started!")
	}
	currentOffset := handle.stream.GlobalOffset()
	cw.stream.Seek(cw.lastChunkOffset-4, io.SeekStart)
	chunkSize := currentOffset - cw.lastChunkOffset
	cw.writer.WriteUInt32(uint32(chunkSize))
	cw.stream.Seek(currentOffset, io.SeekStart)
	cw.lastChunkOffset = -1
}

func (cw *ChunkWriter) Finish() {
	cw.stream.Seek(cw.initialOffset+2, io.SeekStart)
	cw.writer.WriteUInt16(cw.header.numChunks)
}
