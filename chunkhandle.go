package gochunks

import "github.com/chasingcarrots/gobinary"

type ChunkHandle struct {
	stream gobinary.StreamWriterView
	writer gobinary.HighLevelWriter
}

func (ch *ChunkHandle) Stream() *gobinary.StreamWriterView {
	return &ch.stream
}

func (ch *ChunkHandle) Writer() *gobinary.HighLevelWriter {
	return &ch.writer
}
