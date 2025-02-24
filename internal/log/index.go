package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
    offWidth uint64 = 4 // Record offset
    posWidth uint64 = 8 // Position in the store file
    entWidth = offWidth+posWidth   //12
)

type index struct{
    file *os.File
    mmap gommap.MMap
    size uint64
}

func newIndex(f *os.File, c Config) (*index,error) {
    idx := &index{
        file: f,
    }
    fi, err := os.Stat(f.Name())
    if err != nil {
        return nil,err
    }
    idx.size = uint64(fi.Size())
    // grow the file to the max index size before memory mapping the file
    if err = os.Truncate(
        f.Name(),int64(c.Segment.MaxIndexBytes),
    ); err != nil {
        return nil,err
    }

    if idx.mmap, err = gommap.Map(
        idx.file.Fd(),
        gommap.PROT_READ|gommap.PROT_WRITE,
        gommap.MAP_SHARED,
    );err != nil {
        return nil,err
    }
    return idx,nil
}

// Close function make sure memory mapped file
// has synced its data to the persisted file
// and that the presisted file has flushed
// its contents to stable storage
func (i *index)  Close() error {
    if err := i.mmap.Sync(gommap.MS_SYNC);err != nil {
        return err
    }
    if err := i.file.Sync(); err != nil {
        return err
    }
    if err := i.file.Truncate(int64(i.size));err != nil {
        return err
    }
    return i.file.Close()
}

//Append the given offset and position to the index
func (i *index) Write(off uint32,pos uint64) error{
    if uint64(len(i.mmap)) < i.size+entWidth {
        return io.EOF
    }
    enc.PutUint32(i.mmap[i.size:i.size+offWidth],off)
    enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth],pos)
    i.size += uint64(entWidth)
    return nil
}


func (i *index) Read(in int64) (out uint32,pos uint64,err error) {
    if i.size == 0{
        return 0,0,io.EOF
    }

    if in == -1 {
        out = uint32((i.size/ entWidth)-1)
    } else {
        out = uint32(in)
    }
    pos = uint64(out) * entWidth // if out is 1 then pos will be 12
    if i.size < pos + entWidth {    
        return 0,0,io.EOF
    }
    out = enc.Uint32(i.mmap[pos:pos+offWidth])
    pos = enc.Uint64(i.mmap[pos+offWidth:pos+entWidth])
    return out,pos,nil
}


func(i *index) Name() string{
    return i.file.Name()
}
