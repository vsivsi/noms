// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
   "io"
	// "net/http"
	"os"
	"time"
	"reflect"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/spec"
	// "github.com/attic-labs/noms/go/util/jsontonoms"
	"github.com/attic-labs/noms/go/util/progressreader"
	"github.com/attic-labs/noms/go/util/status"
	"github.com/dustin/go-humanize"
	flag "github.com/juju/gnuflag"
)

func NomsValueFromDecodedJSON(o interface{}, useStructList bool) types.Value {
	switch o := o.(type) {
	case string:
		return types.String(o)
	case bool:
		return types.Bool(o)
	case float64:
		return types.Number(o)
	case nil:
		return nil
	case []interface{}:
		var v types.Value
		if useStructList {
			items := make([]types.Value, 0, len(o))
			for _, v := range o {
				nv := NomsValueFromDecodedJSON(v, true)
				if nv != nil {
					items = append(items, nv)
				}
			}
			v = types.NewList(items...)
		} else {
			items := make([]types.Value, 0, len(o))
			for _, v := range o {
				nv := NomsValueFromDecodedJSON(v, true)
				if nv != nil {
					items = append(items, nv)
				}
			}
			v = types.NewSet(items...)
		}
		return v
	case map[string]interface{}:
		var v types.Value
		if useStructList {
			fields := make(types.StructData, len(o))
			for k, v := range o {
				nv := NomsValueFromDecodedJSON(v, true)
				if nv != nil {
					k := types.EscapeStructField(k)
					fields[k] = nv
				}
			}
			v = types.NewStruct("", fields)
		} else {
			kv := make([]types.Value, 0, len(o)*2)
			for k, v := range o {
				nv := NomsValueFromDecodedJSON(v, true)
				if nv != nil {
					kv = append(kv, types.String(k), nv)
				}
			}
			v = types.NewMap(kv...)
		}
		return v

	default:
		d.Chk.Fail("Nomsification failed.", "I don't understand %+v, which is of type %s!\n", o, reflect.TypeOf(o).String())
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <dataset> [<filename>]\n", os.Args[0])
		flag.PrintDefaults()
	}

	spec.RegisterDatabaseFlags(flag.CommandLine)
	flag.Parse(true)

	if len(flag.Args()) < 1 {
		d.CheckError(errors.New("expected dataset option"))
	}

	ds, err := spec.GetDataset(flag.Arg(0))
	d.CheckError(err)

	var file io.Reader

	structList := false

	if len(flag.Args()) >= 2 {

		if len(flag.Args()) == 3 {
			structList = true
		}

		filename := flag.Arg(1)
		if filename == "" {
			flag.Usage()
		}
		res, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Error opening file %s: %+v\n", filename, err)
		}
      file = res
		defer res.Close()
	} else {
		file = os.Stdin
   }

	var jsonObject interface{}
	start := time.Now()
	r := progressreader.New(file, func(seen uint64) {
		elapsed := time.Since(start).Seconds()
		rate := uint64(float64(seen) / elapsed)
		status.Printf("%s decoded in %ds (%s/s)...", humanize.Bytes(seen), int(elapsed), humanize.Bytes(rate))
	})
	err = json.NewDecoder(r).Decode(&jsonObject)
	if err != nil {
		log.Fatalln("Error decoding JSON: ", err)
	}
	status.Done()

	_, err = ds.CommitValue(NomsValueFromDecodedJSON(jsonObject, structList))
	d.PanicIfError(err)
	ds.Database().Close()
}
