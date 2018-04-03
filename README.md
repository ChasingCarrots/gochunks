# Go Chunks
`gochunks` is a very bare-bones Go library to read and write chunk-based binary files. A *chunk* is a section within a binary file that consists of a header plus a payload. The header contains the chunk name, its version, and the size of its payload in bytes. Chunks are stored consecutively in a *chunk file*.

## Chunk File Format
The format of a chunk file is given by a header followed by the number of chunks specified in the chunk file header.

### Chunk File Header
The header is formatted like this:
```
uint16 version
uint16 numChunks
```
<<<<<<< HEAD
The version field may be updated to allow for future changes to the system.
=======
>>>>>>> b696b74272ee7b00b10d2f277edb11e84a7fef39

### Chunks
Chunks consist of a header followed by a payload. The header is formatted like this:
```
uint16 nameLength (in bytes)
string name (utf8, of the byte length given above)
uint16 version
uint32 payloadSize (in bytes)
```
<<<<<<< HEAD
It is followed by `payloadSize` many raw bytes whose interpretation is up to the program handling the chunks. The semantics of the `version` field are up to you: While the name identifies the type of chunk, the version may be used to use differentiate between different versions of a chunk's payload.
=======
It is followed by `payloadSize` many raw bytes whose interpretation is up to the program handling the chunks.
>>>>>>> b696b74272ee7b00b10d2f277edb11e84a7fef39

## Writing Chunks
Chunk files are written using the `ChunkWriter` struct:
```golang
file, _ := os.Create("YourFileNameHere")
saver := gochunks.NewChunkWriter(gobinary.NewStreamWriter(file), 1 /* version */)

saver.Init()

// Begins a new chunk.
chunkHandle := saver.BeginChunk("MyChunk", 1 /* version */)

// Acquire a writer for this chunk. Alternatively, use writer.Stream() to get
// a *gobinary.StreamWriterView.
var writer *gobinary.HighLevelWriter
writer = chunkHandle.Writer()

// Use the writer to write your payload.
writer.WriteUInt32(1337)

// Ends the current chunk. This is mandatory and you cannot begin a new chunk
// without first ending the current one.
saver.EndChunk(chunkHandle)

saver.Finish()
```

## Reading Chunks
Reading (predictably) uses the `ChunkReader` struct:
```golang
file, _ := os.Open("YourFileNameHere")
reader := gochunks.NewChunkReader(gobinary.NewStreamReader(file))
reader.Init()

for i := 0; i < reader.NumChunks(); i++ {
    chunk := reader.NextChunk()
    // ensure that the version of the chunk does not exceed 1
    err := chunk.CheckVersion(1)
    if (err != nil) {
        // your error handling here
        panic(err)
    }

    // Acquire a reader for this chunk. Alternatively, use chunk.Stream() and
    // chunk.Length() to do it manually. Note that currently nothing stops you
    // from reading outside of chunk bounds.
    reader := chunk.Reader()
    if reader.ReadUInt32() != 1337 {
        panic("Something went terribly wrong!")
    }

    // no need to close the chunk
}

```