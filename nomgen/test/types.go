// This file was generated by nomgen.
// To regenerate, run `go generate` in this package.

package main

import (
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

// MyTestSet

type MyTestSet struct {
	s types.Set
}

type MyTestSetIterCallback (func(p types.UInt32) (stop bool))

func NewMyTestSet() MyTestSet {
	return MyTestSet{types.NewSet()}
}

func MyTestSetFromVal(p types.Value) MyTestSet {
	return MyTestSet{p.(types.Set)}
}

func (s MyTestSet) NomsValue() types.Set {
	return s.s
}

func (s MyTestSet) Equals(p MyTestSet) bool {
	return s.s.Equals(p.s)
}

func (s MyTestSet) Ref() ref.Ref {
	return s.s.Ref()
}

func (s MyTestSet) Empty() bool {
	return s.s.Empty()
}

func (s MyTestSet) Len() uint64 {
	return s.s.Len()
}

func (s MyTestSet) Has(p types.UInt32) bool {
	return s.s.Has(p)
}

func (s MyTestSet) Iter(cb MyTestSetIterCallback) {
	s.s.Iter(func(v types.Value) bool {
		return cb(types.UInt32FromVal(v))
	})
}

func (s MyTestSet) Insert(p ...types.UInt32) MyTestSet {
	return MyTestSet{s.s.Insert(s.fromElemSlice(p)...)}
}

func (s MyTestSet) Remove(p ...types.UInt32) MyTestSet {
	return MyTestSet{s.s.Remove(s.fromElemSlice(p)...)}
}

func (s MyTestSet) Union(others ...MyTestSet) MyTestSet {
	return MyTestSet{s.s.Union(s.fromStructSlice(others)...)}
}

func (s MyTestSet) Subtract(others ...MyTestSet) MyTestSet {
	return MyTestSet{s.s.Subtract(s.fromStructSlice(others)...)}
}

func (s MyTestSet) Any() types.UInt32 {
	return types.UInt32FromVal(s.s.Any())
}

func (s MyTestSet) fromStructSlice(p []MyTestSet) []types.Set {
	r := make([]types.Set, len(p))
	for i, v := range p {
		r[i] = v.s
	}
	return r
}

func (s MyTestSet) fromElemSlice(p []types.UInt32) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v
	}
	return r
}

// MapOfTestStructToSetOfBool

type MapOfTestStructToSetOfBool struct {
	m types.Map
}

type MapOfTestStructToSetOfBoolIterCallback (func(k TestStruct, v SetOfBool) (stop bool))

func NewMapOfTestStructToSetOfBool() MapOfTestStructToSetOfBool {
	return MapOfTestStructToSetOfBool{types.NewMap()}
}

func MapOfTestStructToSetOfBoolFromVal(p types.Value) MapOfTestStructToSetOfBool {
	return MapOfTestStructToSetOfBool{p.(types.Map)}
}

func (m MapOfTestStructToSetOfBool) NomsValue() types.Map {
	return m.m
}

func (m MapOfTestStructToSetOfBool) Equals(p MapOfTestStructToSetOfBool) bool {
	return m.m.Equals(p.m)
}

func (m MapOfTestStructToSetOfBool) Ref() ref.Ref {
	return m.m.Ref()
}

func (m MapOfTestStructToSetOfBool) Empty() bool {
	return m.m.Empty()
}

func (m MapOfTestStructToSetOfBool) Len() uint64 {
	return m.m.Len()
}

func (m MapOfTestStructToSetOfBool) Has(p TestStruct) bool {
	return m.m.Has(p.NomsValue())
}

func (m MapOfTestStructToSetOfBool) Get(p TestStruct) SetOfBool {
	return SetOfBoolFromVal(m.m.Get(p.NomsValue()))
}

func (m MapOfTestStructToSetOfBool) Set(k TestStruct, v SetOfBool) MapOfTestStructToSetOfBool {
	return MapOfTestStructToSetOfBoolFromVal(m.m.Set(k.NomsValue(), v.NomsValue()))
}

// TODO: Implement SetM?

func (m MapOfTestStructToSetOfBool) Remove(p TestStruct) MapOfTestStructToSetOfBool {
	return MapOfTestStructToSetOfBoolFromVal(m.m.Remove(p.NomsValue()))
}

func (m MapOfTestStructToSetOfBool) Iter(cb MapOfTestStructToSetOfBoolIterCallback) {
	m.m.Iter(func(k, v types.Value) bool {
		return cb(TestStructFromVal(k), SetOfBoolFromVal(v))
	})
}

// ListOfInt32

type ListOfInt32 struct {
	l types.List
}

type ListOfInt32IterCallback (func (p types.Int32) (stop bool))

func NewListOfInt32() ListOfInt32 {
	return ListOfInt32{types.NewList()}
}

func ListOfInt32FromVal(p types.Value) ListOfInt32 {
	return ListOfInt32{p.(types.List)}
}

func (l ListOfInt32) NomsValue() types.List {
	return l.l
}

func (l ListOfInt32) Equals(p ListOfInt32) bool {
	return l.l.Equals(p.l)
}

func (l ListOfInt32) Ref() ref.Ref {
	return l.l.Ref()
}

func (l ListOfInt32) Len() uint64 {
	return l.l.Len()
}

func (l ListOfInt32) Empty() bool {
	return l.Len() == uint64(0)
}

func (l ListOfInt32) Get(idx uint64) types.Int32 {
	return types.Int32FromVal(l.l.Get(idx))
}

func (l ListOfInt32) Slice(idx uint64, end uint64) ListOfInt32 {
	return ListOfInt32{l.l.Slice(idx, end)}
}

func (l ListOfInt32) Set(idx uint64, v types.Int32) ListOfInt32 {
	return ListOfInt32{l.l.Set(idx, v)}
}

func (l ListOfInt32) Append(v ...types.Int32) ListOfInt32 {
	return ListOfInt32{l.l.Append(l.fromElemSlice(v)...)}
}

func (l ListOfInt32) Insert(idx uint64, v ...types.Int32) ListOfInt32 {
	return ListOfInt32{l.l.Insert(idx, l.fromElemSlice(v)...)}
}

func (l ListOfInt32) Remove(idx uint64, end uint64) ListOfInt32 {
	return ListOfInt32{l.l.Remove(idx, end)}
}

func (l ListOfInt32) RemoveAt(idx uint64) ListOfInt32 {
	return ListOfInt32{(l.l.RemoveAt(idx))}
}

func (l ListOfInt32) fromElemSlice(p []types.Int32) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v
	}
	return r
}

// MapOfStringToFloat64

type MapOfStringToFloat64 struct {
	m types.Map
}

type MapOfStringToFloat64IterCallback (func(k types.String, v types.Float64) (stop bool))

func NewMapOfStringToFloat64() MapOfStringToFloat64 {
	return MapOfStringToFloat64{types.NewMap()}
}

func MapOfStringToFloat64FromVal(p types.Value) MapOfStringToFloat64 {
	return MapOfStringToFloat64{p.(types.Map)}
}

func (m MapOfStringToFloat64) NomsValue() types.Map {
	return m.m
}

func (m MapOfStringToFloat64) Equals(p MapOfStringToFloat64) bool {
	return m.m.Equals(p.m)
}

func (m MapOfStringToFloat64) Ref() ref.Ref {
	return m.m.Ref()
}

func (m MapOfStringToFloat64) Empty() bool {
	return m.m.Empty()
}

func (m MapOfStringToFloat64) Len() uint64 {
	return m.m.Len()
}

func (m MapOfStringToFloat64) Has(p types.String) bool {
	return m.m.Has(p)
}

func (m MapOfStringToFloat64) Get(p types.String) types.Float64 {
	return types.Float64FromVal(m.m.Get(p))
}

func (m MapOfStringToFloat64) Set(k types.String, v types.Float64) MapOfStringToFloat64 {
	return MapOfStringToFloat64FromVal(m.m.Set(k, v))
}

// TODO: Implement SetM?

func (m MapOfStringToFloat64) Remove(p types.String) MapOfStringToFloat64 {
	return MapOfStringToFloat64FromVal(m.m.Remove(p))
}

func (m MapOfStringToFloat64) Iter(cb MapOfStringToFloat64IterCallback) {
	m.m.Iter(func(k, v types.Value) bool {
		return cb(types.StringFromVal(k), types.Float64FromVal(v))
	})
}

// TestStruct

type TestStruct struct {
	m types.Map
}

func NewTestStruct() TestStruct {
	return TestStruct{
		types.NewMap(types.NewString("$name"), types.NewString("TestStruct")),
	}
}

func TestStructFromVal(v types.Value) TestStruct {
	return TestStruct{v.(types.Map)}
}

// TODO: This was going to be called Value() but it collides with root.value. We need some other place to put the built-in fields like Value() and Equals().
func (s TestStruct) NomsValue() types.Map {
	return s.m
}

func (s TestStruct) Equals(p TestStruct) bool {
	return s.m.Equals(p.m)
}

func (s TestStruct) Ref() ref.Ref {
	return s.m.Ref()
}

func (s TestStruct) Title() types.String {
	return types.StringFromVal(s.m.Get(types.NewString("title")))
}

func (s TestStruct) SetTitle(p types.String) TestStruct {
	return TestStructFromVal(s.m.Set(types.NewString("title"), p))
}

// SetOfBool

type SetOfBool struct {
	s types.Set
}

type SetOfBoolIterCallback (func(p types.Bool) (stop bool))

func NewSetOfBool() SetOfBool {
	return SetOfBool{types.NewSet()}
}

func SetOfBoolFromVal(p types.Value) SetOfBool {
	return SetOfBool{p.(types.Set)}
}

func (s SetOfBool) NomsValue() types.Set {
	return s.s
}

func (s SetOfBool) Equals(p SetOfBool) bool {
	return s.s.Equals(p.s)
}

func (s SetOfBool) Ref() ref.Ref {
	return s.s.Ref()
}

func (s SetOfBool) Empty() bool {
	return s.s.Empty()
}

func (s SetOfBool) Len() uint64 {
	return s.s.Len()
}

func (s SetOfBool) Has(p types.Bool) bool {
	return s.s.Has(p)
}

func (s SetOfBool) Iter(cb SetOfBoolIterCallback) {
	s.s.Iter(func(v types.Value) bool {
		return cb(types.BoolFromVal(v))
	})
}

func (s SetOfBool) Insert(p ...types.Bool) SetOfBool {
	return SetOfBool{s.s.Insert(s.fromElemSlice(p)...)}
}

func (s SetOfBool) Remove(p ...types.Bool) SetOfBool {
	return SetOfBool{s.s.Remove(s.fromElemSlice(p)...)}
}

func (s SetOfBool) Union(others ...SetOfBool) SetOfBool {
	return SetOfBool{s.s.Union(s.fromStructSlice(others)...)}
}

func (s SetOfBool) Subtract(others ...SetOfBool) SetOfBool {
	return SetOfBool{s.s.Subtract(s.fromStructSlice(others)...)}
}

func (s SetOfBool) Any() types.Bool {
	return types.BoolFromVal(s.s.Any())
}

func (s SetOfBool) fromStructSlice(p []SetOfBool) []types.Set {
	r := make([]types.Set, len(p))
	for i, v := range p {
		r[i] = v.s
	}
	return r
}

func (s SetOfBool) fromElemSlice(p []types.Bool) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v
	}
	return r
}
