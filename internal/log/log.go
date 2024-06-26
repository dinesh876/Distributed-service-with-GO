package log

import (
	//"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	api "github.com/dinesh876/proglog/api/v1"
)

type Log struct{
    mu sync.Mutex
    Dir string
    Config Config

    activeSegement *segment
    segments []*segment
}

func NewLog(dir string,c Config) (*Log,error){
    if c.Segment.MaxStoreBytes == 0 {
        c.Segment.MaxStoreBytes = 1024
    }
    if c.Segment.MaxIndexBytes == 0 {
        c.Segment.MaxIndexBytes = 1024
    }

    l := &Log{
        Dir : dir,
        Config: c,
    }
    return l,l.setup()
}

func (l *Log) setup() error {
    files,err := ioutil.ReadDir(l.Dir)
    if err != nil {
        return err
    }
    var baseOffsets []uint64
    for _,file := range files{
        offStr := strings.TrimSuffix(
            file.Name(),
            path.Ext(file.Name()),
        )
        off,_ := strconv.ParseUint(offStr,10,0)
        baseOffsets = append(baseOffsets,off)
    }

    sort.Slice(baseOffsets, func(i,j int) bool {
        return baseOffsets[i] < baseOffsets[j]
    })

    for i:= 0;i<len(baseOffsets); i++ {
        if err = l.newSegment(baseOffsets[i]);err != nil {
            return err
        }
        // baseOffsets contains dup for index and store so we skip the dup
        i++
    }

    if l.segments == nil {
        if err  := l.newSegment(
            l.Config.Segment.InitialOffset,
        );err != nil {
            return err
        }
    }
    return nil
}


func (l *Log) Append(record *api.Record) (uint64,error) {
    l.mu.Lock()
    defer l.mu.Unlock()
    off, err := l.activeSegement.Append(record)
    if err != nil {
        return 0,err
    }
    if l.activeSegement.IsMaxed(){
        err = l.newSegment(off+1)
    }
    return off,err
}

func (l *Log) Read(off uint64) (*api.Record,error){
    l.mu.Lock()
    defer l.mu.Unlock()
    var s *segment
    for _,segement := range l.segments{
        if segement.baseOffset <= off && off < segement.nextOffset {
            s = segement
            break
        }
    }
    if s == nil || s.nextOffset <= off {
        return nil, api.ErrOffsetOutOfRange{Offset: off}
    }
    return s.Read(off)
}

// Close() iterates over the segements and closes them
// Remove() closes the log and then remove its data
// Reset() removes the log and then creates a new log to replace it

func(l *Log) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    for _,segement := range l.segments {
        if err := segement.Close(); err != nil {
            return err
        }
    }
    return nil
}

func (l *Log) Remove() error{
    if err := l.Close();err != nil {
        return err
    }
    return os.RemoveAll(l.Dir)
}

func (l *Log) Reset() error{
    if err := l.Remove(); err != nil {
        return err
    }
    return l.setup()
}

func (l *Log) LowestOffset() (uint64,error) {
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.segments[0].baseOffset, nil
} 

func (l *Log) HighestOffset() (uint64,error){
    l.mu.Lock()
    defer l.mu.Unlock()
    off := l.segments[len(l.segments)-1].nextOffset
    if off == 0 {
        return 0,nil
    }
    return off -1,nil
}

// removes all segements whose highest offset is lower than
// lowest. Because we don't have disks with infinite space,
// we'll periodically call Truncate() to remove old segments

func (l *Log) Truncate(lowest uint64) error{
    l.mu.Lock()
    defer l.mu.Unlock()
    var segements []*segment
    for _,s := range l.segments{
        if s.nextOffset <= lowest +1 {
            if err := s.Remove();err != nil {
            return err
            }
            continue
        }
        segements = append(segements,s)
    }
    l.segments = segements
    return nil
} 

func (l *Log) newSegment(off uint64) error {
    s,err := newSegment(l.Dir,off,l.Config)
    if err != nil {
        return err
    }
    l.segments = append(l.segments,s)
    l.activeSegement = s
    return nil
}

func (l *Log) Reader() io.Reader {
    l.mu.Lock()
    defer l.mu.Unlock()
    readers :=  make([]io.Reader,len(l.segments))
    for i,segment := range l.segments{
        readers[i] = &originReader{segment.store,0} 
    }
    return io.MultiReader(readers...)
}


type originReader struct {
    *store
    off int64
}

func (o *originReader) Read(p []byte) (int,error){
        n,err := o.ReadAt(p,o.off)
        o.off += int64(n)
        return n,err
}
