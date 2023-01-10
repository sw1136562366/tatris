// Copyright 2022 Tatris Project Authors. Licensed under Apache-2.0.

// Package metadata is about the management of metadata (i.e. index)
package metadata

import (
	"encoding/json"
	"github.com/tatris-io/tatris/internal/protocol"
	"gotest.tools/v3/assert"
	"testing"
)

type testItem struct {
	Index protocol.Index
	Res   bool
}

func TestCheckParam(t *testing.T) {
	params := `[
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"keyword"}}}}},
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"text"}}}}},
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"INTEGER"}}}}},
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"long"}}}}},
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"FLOAT"}}}}},
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"double"}}}}},
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"BOOLEAN"}}}}},
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"date"}}}}},
		{"Res":true ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"dAtE"}}}}},
		{"Res":false},
		{"Res":false ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{}}},
		{"Res":false ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"keyword"},"age":{"type":"string"}}}}},
		{"Res":false ,"Index":{"settings":{"number_of_shards":3,"number_of_replicas":1},"mappings":{"properties":{"name":{"type":"bool"},"age":{"type":"int"}}}}}
	]`
	var items []testItem
	err := json.Unmarshal([]byte(params), &items)
	if err != nil {
		t.Error(err)
		return
	}
	for i, item := range items {
		err := checkParam(&item.Index)
		comparison := err == nil
		if !comparison {
			t.Logf("item %d error : %s", i, err)
		}
		assert.Equal(t, comparison, item.Res)
	}

}