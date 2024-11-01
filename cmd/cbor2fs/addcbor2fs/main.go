package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strconv"

	ck "github.com/takanoriyanagitani/go-cbor2kvstore"

	util "github.com/takanoriyanagitani/go-cbor2kvstore/util"

	ap "github.com/takanoriyanagitani/go-cbor2kvstore/app/cbor2iter2kvstore"

	ic "github.com/takanoriyanagitani/go-cbor2kvstore/iter/cbor2array"
	ca "github.com/takanoriyanagitani/go-cbor2kvstore/iter/cbor2array/amacker"

	sc "github.com/takanoriyanagitani/go-cbor2kvstore/ser/cbor"
	ac "github.com/takanoriyanagitani/go-cbor2kvstore/ser/cbor/amacker"

	kv "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore"
	kf "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple"

	sa "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple/any2buf"
	sp "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple/pathcheck"

	fa "github.com/takanoriyanagitani/go-cbor2kvstore/kvstore/fs/simple/stdfs/append"
)

var BytesToPathBufDefault kf.BytesToPathBuf = kf.BytesToPathBufDefault
var BytesToPathDefault kf.BytesToPath = BytesToPathBufDefault.ToBytesToPath()

var RowToCborBuf sc.CborToBuffer = ac.CborToBufferDefault
var RowToCborVal kv.RowToVal = RowToCborBuf.ToCborToBytes().ToRowToVal()

var AnyToBuf sa.AnyToBuffer = sa.AnyToBufferDefaultNew()
var AnyToDirToBuf kf.AnyToDirNameToBuf = AnyToBuf.AsAnyToDirnameToBuf()

var PatChecker sp.PathChecker = sp.PathCheckerDefault
var Path2buf kf.PathToBuffer = PatChecker.ToPathToBuffer()
var Cpath kf.CreatePath = Path2buf.ToCreatePath()

const (
	DirnameColIndexDefault  uint32 = 0
	FilenameColIndexDefault uint32 = 0
)

func ParseUint32(alt uint32, s string) uint32 {
	u, e := strconv.ParseUint(s, 10, 32)
	switch e {
	case nil:
		return uint32(u)
	default:
		return alt
	}
}

func GetEnv32uNew(alt uint32) func(envkey string) uint32 {
	return util.Compose(
		os.Getenv,
		util.Curry(ParseUint32)(alt),
	)
}

var EnvKeyToDirIdx func(string) uint32 = GetEnv32uNew(DirnameColIndexDefault)
var EnvKeyToFileIdx func(string) uint32 = GetEnv32uNew(FilenameColIndexDefault)

type IoConfig struct {
	io.Reader
}

func (i IoConfig) ToCborToArrIter() ca.CborToArrayIter {
	return ca.CborToArrayIterNew(
		bufio.NewReader(i.Reader),
	)
}

func (i IoConfig) ToCborToArrays() ic.CborToArrays {
	return i.ToCborToArrIter().AsCborToArrays()
}

type Store struct {
	kf.RawStore
	kf.BytesToPath
	kf.RowToKey
	kf.RowToVal
}

func (s Store) ToKvStore() kv.Kvstore {
	return s.RawStore.ToKvStore(s.BytesToPath)
}

func (s Store) ToRowSaver() ck.RowSaver {
	return kf.FsStore{
		Kvstore:  s.ToKvStore(),
		RowToKey: s.RowToKey,
		RowToVal: s.RowToVal,
	}.ToRowSaver()
}

type App struct {
	IoConfig
	Store
}

func (a App) ToCborToKvs() ap.CborToIterToKvstore {
	return ap.CborToIterToKvstore{
		CborToArrays: a.IoConfig.ToCborToArrays(),
		RowSaver:     a.Store.ToRowSaver(),
	}
}

func (a App) SaveAll(ctx context.Context) error {
	return a.ToCborToKvs().SaveAll(ctx)
}

func RawStoreFromDirnameDefault(dirname string) kf.RawStore {
	return fa.FsStoreAppendNewDefault(dirname).ToRawStore()
}

type AnyToName struct {
	kf.AnyToDirNameToBuf
	DirnameSourceColumnIndex  uint32
	FilenameSourceColumnIndex uint32
	kf.CreatePath
}

func (a AnyToName) ToRowToDirAny() kf.RowToDirAny {
	return kf.RowToDirAnyFromColIdx(a.DirnameSourceColumnIndex)
}

func (a AnyToName) ToRowToFileAny() kf.RowToDirAny {
	return kf.RowToDirAnyFromColIdx(a.FilenameSourceColumnIndex)
}

func (a AnyToName) ToRowToDirToBuf() kf.RowToDirNameToBuf {
	return a.AnyToDirNameToBuf.ToRowToDirToBuf(a.ToRowToDirAny())
}

func (a AnyToName) ToRowToFileToBuf() kf.RowToDirNameToBuf {
	return a.AnyToDirNameToBuf.ToRowToDirToBuf(a.ToRowToFileAny())
}

func (a AnyToName) ToRowToDirName() kf.RowToDirName {
	return a.ToRowToDirToBuf().ToRowToDirName()
}

func (a AnyToName) ToRowToFileName() kf.RowToFileName {
	return kf.RowToFileName(a.ToRowToFileToBuf().ToRowToDirName())
}

func (a AnyToName) ToRowToKey() kf.RowToKey {
	return kf.RowToFsKey{
		RowToDirName:  a.ToRowToDirName(),
		RowToFileName: a.ToRowToFileName(),
		CreatePath:    a.CreatePath,
	}.ToRowToKey()
}

type Config struct {
	OutputDirname             string
	DirnameSourceColumnIndex  uint32
	FilenameSourceColumnIndex uint32
}

func ConfigNewEnv() Config {
	var outdir string = os.Getenv("ENV_OUTPUT_DIR_NAME")

	return Config{
		OutputDirname:             outdir,
		DirnameSourceColumnIndex:  EnvKeyToDirIdx("ENV_COL_IX4DIR"),
		FilenameSourceColumnIndex: EnvKeyToDirIdx("ENV_COL_IX4FILE"),
	}
}

func rdr2fs(ctx context.Context, rdr io.Reader, cfg Config) error {
	a2n := AnyToName{
		AnyToDirNameToBuf:         AnyToDirToBuf,
		DirnameSourceColumnIndex:  cfg.DirnameSourceColumnIndex,
		FilenameSourceColumnIndex: cfg.FilenameSourceColumnIndex,
		CreatePath:                Cpath,
	}
	var rstore kf.RawStore = RawStoreFromDirnameDefault(cfg.OutputDirname)

	store := Store{
		RawStore:    rstore,
		BytesToPath: BytesToPathDefault,
		RowToKey:    a2n.ToRowToKey(),
		RowToVal:    kf.RowToVal(RowToCborVal),
	}

	icfg := IoConfig{Reader: rdr}

	app := App{
		IoConfig: icfg,
		Store:    store,
	}

	return app.SaveAll(ctx)
}

func stdin2fs(ctx context.Context, cfg Config) error {
	return rdr2fs(ctx, os.Stdin, cfg)
}

func sub(ctx context.Context) error {
	return stdin2fs(ctx, ConfigNewEnv())
}

func main() {
	e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
