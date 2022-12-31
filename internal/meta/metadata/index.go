// Copyright 2022 Tatris Project Authors. Licensed under Apache-2.0.

// Package metadata is about the management of metadata (i.e. index)
package metadata

import (
	json "encoding/json"
	"github.com/tatris-io/tatris/internal/meta"
	"github.com/tatris-io/tatris/internal/meta/metadata/storage"
)

var metaStore storage.MetaStore

func init() {
	metaStore, _ = storage.Open()
}

func Create(idx *meta.Index) error {
	json, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	return metaStore.Set(fillKey(idx.Name), json)
}

func Get(idxName string) (*meta.Index, error) {
	if b, err := metaStore.Get(fillKey(idxName)); err != nil {
		return nil, err
	} else if b == nil {
		return nil, nil
	} else {
		idx := new(meta.Index)
		if err := json.Unmarshal(b, idx); err != nil {
			return nil, err
		}
		return idx, nil
	}
}

func fillKey(name string) string {
	return "/index/" + name
}