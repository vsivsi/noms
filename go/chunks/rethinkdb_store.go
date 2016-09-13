// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package chunks

import (
	"fmt"
	"sync"

	"github.com/attic-labs/noms/go/constants"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/hash"
	"github.com/golang/snappy"
	flag "github.com/juju/gnuflag"
	r "gopkg.in/dancannon/gorethink.v2"
)

const (
	rethinkSysTable         = "sys"
	rethinkChunkTable       = "chunks"
	rethinkChunkPrefixConst = "/chunk/"
)

type rethinkVersionDoc struct {
	ID      []byte `gorethink:"id"`
	Version []byte `gorethink:"Version"`
}

type rethinkRootDoc struct {
	ID   []byte `gorethink:"id"`
	Root []byte `gorethink:"Root"`
}

type rethinkChunkDoc struct {
	ID   []byte `gorethink:"id"`
	Data []byte `gorethink:"Data"`
}

type RethinkDBStoreFlags struct {
	dumpStats bool
}

var (
	rethinkFlags           = RethinkDBStoreFlags{false}
	rethinkFlagsRegistered = false
)

func RegisterRethinkDBFlags(flags *flag.FlagSet) {
	if !rethinkFlagsRegistered {
		rethinkFlagsRegistered = true
		flags.BoolVar(&rethinkFlags.dumpStats, "rethink-dump-stats", false, "print get/has/put counts on close")
	}
}

func NewRethinkStoreUseFlags(url, db, ns string) *RethinkStore {
	return newRethinkStore(newRethinkBackingStore(url, db, rethinkFlags.dumpStats), []byte(ns), true)
}

func NewRethinkStore(url, db, ns string, dumpStats bool) *RethinkStore {
	return newRethinkStore(newRethinkBackingStore(url, db, dumpStats), []byte(ns), true)
}

func newRethinkStore(store *internalRethinkStore, ns []byte, closeBackingStore bool) *RethinkStore {
	copyNsAndAppend := func(suffix string) (out []byte) {
		out = make([]byte, len(ns)+len(suffix))
		copy(out[copy(out, ns):], []byte(suffix))
		return
	}
	return &RethinkStore{
		internalRethinkStore: store,
		rootKey:              []byte(""),
		versionKey:           []byte(""),
		chunkPrefix:          copyNsAndAppend(rethinkChunkPrefixConst),
		closeBackingStore:    closeBackingStore,
	}
}

type RethinkStore struct {
	*internalRethinkStore
	rootKey           []byte
	versionKey        []byte
	chunkPrefix       []byte
	closeBackingStore bool
	versionSetOnce    sync.Once
}

func (l *RethinkStore) Root() hash.Hash {
	d.Chk.True(l.internalRethinkStore != nil, "Cannot use RethinkStore after Close().")
	return l.rootByKey(l.rootKey)
}

func (l *RethinkStore) UpdateRoot(current, last hash.Hash) bool {
	d.Chk.True(l.internalRethinkStore != nil, "Cannot use RethinkStore after Close().")
	l.versionSetOnce.Do(l.setVersIfUnset)
	return l.updateRootByKey(l.rootKey, current, last)
}

func (l *RethinkStore) Get(ref hash.Hash) Chunk {
	d.Chk.True(l.internalRethinkStore != nil, "Cannot use RethinkStore after Close().")
	return l.getByKey(l.toChunkKey(ref), ref)
}

func (l *RethinkStore) Has(ref hash.Hash) bool {
	d.Chk.True(l.internalRethinkStore != nil, "Cannot use RethinkStore after Close().")
	return l.hasByKey(l.toChunkKey(ref))
}

func (l *RethinkStore) Version() string {
	d.Chk.True(l.internalRethinkStore != nil, "Cannot use RethinkStore after Close().")
	return l.versByKey(l.versionKey)
}

func (l *RethinkStore) Put(c Chunk) {
	d.Chk.True(l.internalRethinkStore != nil, "Cannot use RethinkStore after Close().")
	l.versionSetOnce.Do(l.setVersIfUnset)
	l.putByKey(l.toChunkKey(c.Hash()), c)
}

func (l *RethinkStore) PutMany(chunks []Chunk) (e BackpressureError) {
	d.Chk.True(l.internalRethinkStore != nil, "Cannot use RethinkStore after Close().")
	l.versionSetOnce.Do(l.setVersIfUnset)
	// numBytes := 0
	// b := new(Rethink.Batch)

	// TODO: This is an initial implementation using individual Puts
	// This can almost certainly be batched into a single insert
	for _, c := range chunks {
		l.putByKey(l.toChunkKey(c.Hash()), c)
		// data := snappy.Encode(nil, c.Data())
		// numBytes += len(data)
		// b.Put(l.toChunkKey(c.Hash()), data)
	}
	// l.putBatch(b, numBytes)
	return
}

func (l *RethinkStore) Close() error {
	if l.closeBackingStore {
		l.internalRethinkStore.Close()
	}
	l.internalRethinkStore = nil
	return nil
}

func (l *RethinkStore) toChunkKey(r hash.Hash) []byte {
	digest := r.DigestSlice()
	out := make([]byte, len(l.chunkPrefix), len(l.chunkPrefix)+len(digest))
	copy(out, l.chunkPrefix)
	return append(out, digest...)
}

func (l *RethinkStore) setVersIfUnset() {
	// The Rethink query handles this case...
	// exists, err := l.session.Has(l.versionKey, nil)
	// d.Chk.NoError(err)
	// if !exists {
	l.setVersByKey(l.versionKey)
	//}
}

type internalRethinkStore struct {
	session                                *r.Session
	db                                     string
	sys                                    r.Term
	chunks                                 r.Term
	getCount, hasCount, putCount, putBytes int64
	dumpStats                              bool
}

func newRethinkBackingStore(url, db string, dumpStats bool) *internalRethinkStore {
	d.PanicIfTrue(url == "", "url cannot be empty")
	d.PanicIfTrue(db == "", "db cannot be empty")

	session, err := r.Connect(r.ConnectOpts{
		Address:  url,
		Database: db,
	})
	d.Chk.NoError(err, "opening connection %s in internalRethinkStore", url)

	// Create requested DB if it doesn't exist
	_, err = r.Branch(
		r.DBList().Contains(db),
		r.Expr(map[string]interface{}{"dbs_created": 0}),
		r.DBCreate(db)).Run(session)
	d.Chk.NoError(err, "conditionally creating requested DB %s in internalRethinkStore", db)

	// Create system table if it doesn't exist
	_, err = r.Branch(
		r.TableList().Contains(rethinkSysTable),
		r.Expr(map[string]interface{}{"tables_created": 0}),
		r.TableCreate(rethinkSysTable)).Run(session)

	d.Chk.NoError(err, "conditionally creating system table %s in internalRethinkStore", rethinkSysTable)

	// Create chunk table if it doesn't exist
	_, err = r.Branch(
		r.TableList().Contains(rethinkChunkTable),
		r.Expr(map[string]interface{}{"tables_created": 0}),
		r.TableCreate(rethinkChunkTable)).Run(session)

	d.Chk.NoError(err, "conditionally creating chunk table %s in internalRethinkStore", rethinkChunkTable)

	return &internalRethinkStore{
		session:   session,
		db:        db,
		sys:       r.Table(rethinkSysTable),
		chunks:    r.Table(rethinkChunkTable),
		dumpStats: dumpStats,
	}
}

func (l *internalRethinkStore) rootByKey(key []byte) hash.Hash {
	cursor, err := l.sys.Get(key).Run(l.session)
	d.Chk.NoError(err)
	var doc rethinkRootDoc
	if cursor.Next(&doc) {
		return hash.Parse(string(doc.Root))
	} else {
		d.Chk.NoError(cursor.Err())
		return hash.Hash{}
	}
}

func (l *internalRethinkStore) updateRootByKey(key []byte, current, last hash.Hash) bool {
	// TODO: do something with the response below?
	_, err := l.sys.Get(key).Replace(
		r.Branch(
			r.Row.IsEmpty().Or(r.Row.Field("Root").Eq([]byte(last.String()))),
			r.Expr(rethinkRootDoc{
				ID:   key,
				Root: []byte(current.String()),
			}),
			r.Row),
		r.ReplaceOpts{Durability: "hard", ReturnChanges: true}).RunWrite(l.session)
	d.Chk.NoError(err)
	return true
}

func (l *internalRethinkStore) getByKey(key []byte, ref hash.Hash) Chunk {
	cursor, err := l.chunks.Get(key).Run(l.session)
	d.Chk.NoError(err)
	l.getCount++

	if cursor.IsNil() {
      return EmptyChunk
   } else {
      var doc rethinkChunkDoc
      err = cursor.One(&doc)
		fmt.Println("Doc! ", doc.ID, len(doc.Data))
		data, err := snappy.Decode(nil, doc.Data)
		d.Chk.NoError(err)
		return NewChunkWithHash(ref, data)
   }

	//
	// var doc rethinkChunkDoc
	// if cursor.Next(&doc) {
	// 	fmt.Println("Doc! ", doc.ID, len(doc.Data))
	// 	data, err := snappy.Decode(nil, doc.Data)
	// 	d.Chk.NoError(err)
	// 	return NewChunkWithHash(ref, data)
	// } else {
	// 	d.Chk.NoError(cursor.Err())
	// 	return EmptyChunk
	// }
}

func (l *internalRethinkStore) hasByKey(key []byte) bool {
	cursor, err := l.chunks.Get(key).IsEmpty().Not().Run(l.session)
	d.Chk.NoError(err)
	var exists bool
	err = cursor.One(&exists)
	d.Chk.NoError(err)
	l.hasCount++
	return exists
}

func (l *internalRethinkStore) versByKey(key []byte) string {
	var res string
	cursor, err := l.sys.Get(key).Field("Version").Run(l.session)
	defer cursor.Close()
	d.Chk.NoError(err)
	if cursor.Next(&res) {
		return res
	} else {
		d.Chk.NoError(cursor.Err())
		return constants.NomsVersion
	}
}

func (l *internalRethinkStore) setVersByKey(key []byte) {
	// TODO: do something with the response below?
	_, err := l.sys.Get(key).Replace(
		r.Branch(
			r.Row.IsEmpty(),
			r.Expr(
				rethinkVersionDoc{
					ID:      key,
					Version: []byte(constants.NomsVersion),
				}),
			r.Row,
		)).RunWrite(l.session)
	d.Chk.NoError(err)
}

func (l *internalRethinkStore) putByKey(key []byte, c Chunk) {
	data := snappy.Encode(nil, c.Data())
	// TODO: do something with the response below?
	_, err := l.chunks.Insert(rethinkChunkDoc{
		ID:   key,
		Data: data,
	}).RunWrite(l.session)
	d.Chk.NoError(err)
	l.putCount++
	l.putBytes += int64(len(data))
}

// Not currently needed. See TODO for PutMany...

// func (l *internalRethinkStore) putBatch(b *Rethink.Batch, numBytes int) {
// 	err := l.db.Write(b, nil)
// 	d.Chk.NoError(err)
// 	l.putCount += int64(b.Len())
// 	l.putBytes += int64(numBytes)
// }

func (l *internalRethinkStore) Close() error {
	err := l.session.Close()
	if l.dumpStats {
		fmt.Println("--Rethink Stats--")
		fmt.Println("GetCount: ", l.getCount)
		fmt.Println("HasCount: ", l.hasCount)
		fmt.Println("PutCount: ", l.putCount)
		fmt.Println("Average PutSize: ", l.putBytes/l.putCount)
	}
	return err
}

func NewRethinkStoreFactory(url, db string, dumpStats bool) Factory {
	return &RethinkStoreFactory{url, db, dumpStats, newRethinkBackingStore(url, db, dumpStats)}
}

func NewRethinkStoreFactoryUseFlags(url, db string) Factory {
	return NewRethinkStoreFactory(url, db, rethinkFlags.dumpStats)
}

type RethinkStoreFactory struct {
	url       string
	db        string
	dumpStats bool
	store     *internalRethinkStore
}

func (f *RethinkStoreFactory) CreateStore(ns string) ChunkStore {
	d.Chk.True(f.store != nil, "Cannot use RethinkStoreFactory after Shutter().")
	return newRethinkStore(f.store, []byte(ns), false)
}

func (f *RethinkStoreFactory) Shutter() {
	f.store.Close()
	f.store = nil
}
