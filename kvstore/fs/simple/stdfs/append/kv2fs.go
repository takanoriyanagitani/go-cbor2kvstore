package append2fs

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"

	kv "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore"
	kf "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple"
)

const (
	FileModeDefaultFile os.FileMode = 0644
	FileModeDefaultDir  os.FileMode = 0755
)

type FileSync func(*os.File) error

type CreateDir func(string) error

type CborExt string

const CborExtDefault CborExt = "cbor"

func CreateDirNewFromMode(m os.FileMode) CreateDir {
	return func(dirname string) error {
		return os.MkdirAll(dirname, m)
	}
}

var FileSyncNop FileSync = func(_ *os.File) error { return nil }

type FsStoreAppend struct {
	Dirname string
	os.FileMode
	CreateDir
	FileSync
	CborExt
}

func FsStoreAppendNewDefault(dirname string) FsStoreAppend {
	return FsStoreAppend{
		Dirname:   dirname,
		FileMode:  FileModeDefaultFile,
		CreateDir: CreateDirNewFromMode(FileModeDefaultDir),
		FileSync:  FileSyncNop,
		CborExt:   CborExtDefault,
	}
}

func (a FsStoreAppend) WriteData(data []byte, f *os.File) error {
	_, e := f.Write(data)
	if nil != e {
		return e
	}
	return a.FileSync(f)
}

// Non-atomically(no rename) write(append) the data to the path.
func (a FsStoreAppend) WriteToPathNonAtomic(
	_ context.Context,
	pat string,
	dat []byte,
) error {
	var joined string = filepath.Join(a.Dirname, pat)
	var withExt string = joined + "." + string(a.CborExt)
	var dirname string = filepath.Dir(joined)

	e := a.CreateDir(dirname)
	if nil != e {
		return e
	}

	f, e := os.OpenFile(
		withExt,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		a.FileMode,
	)
	if nil != e {
		return e
	}

	var closeErr error = nil
	var closer func() = sync.OnceFunc(func() {
		closeErr = f.Close()
	})
	defer closer()

	e = a.WriteData(dat, f)
	if nil != e {
		closer()
		return errors.Join(e, closeErr)
	}

	closer()
	return closeErr
}

func (a FsStoreAppend) ToRawStore() kf.RawStore {
	return func(ctx context.Context, p kf.Path, v kv.Val) error {
		var pat string = string(p)
		var dat []byte = v
		return a.WriteToPathNonAtomic(ctx, pat, dat)
	}
}
