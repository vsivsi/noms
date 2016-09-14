// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package chunks

import (
	"bytes"
	"testing"

	"github.com/attic-labs/testify/suite"
)

func TestRethinkStoreTestSuite(t *testing.T) {
	suite.Run(t, &RethinkStoreTestSuite{})
}

type RethinkStoreTestSuite struct {
	ChunkStoreTestSuite
	factory Factory
	url     string
	db      string
	*RethinkStore
}

func (suite *RethinkStoreTestSuite) SetupTest() {
	suite.url = "localhost"
	suite.db = "gotest"
	suite.factory = NewRethinkStoreFactory(suite.url, suite.db, false)
	suite.RethinkStore = suite.factory.CreateStore("name").(*RethinkStore)
	suite.putCountFn = func() int {
		return int(suite.RethinkStore.putCount)
	}

	suite.Store = suite.RethinkStore
}

func (suite *RethinkStoreTestSuite) TearDownTest() {
	suite.RethinkStore._Teardown()
	suite.factory.Shutter()
}

func (suite *RethinkStoreTestSuite) TestReservedKeys() {
	// Apparently, the following:
	//  s := []byte("")
	//  s = append(s, 1, 2, 3)
	//  f := append(s, 10, 20, 30)
	//  g := append(s, 4, 5, 6)
	//
	// Results in both f and g being [1, 2, 3, 4, 5, 6]
	// This was happening to us here, so ldb.chunkPrefix was "/chunk/" and ldb.rootKey was "/chun" instead of "/root"
	l := suite.factory.CreateStore("").(*RethinkStore)
	suite.True(bytes.HasSuffix(l.rootKey, []byte(rethinkRootKeyConst)))
	suite.True(bytes.HasSuffix(l.versionKey, []byte(rethinkVersionKeyConst)))
	suite.True(bytes.HasSuffix(l.chunkPrefix, []byte(rethinkChunkPrefixConst)))
}
