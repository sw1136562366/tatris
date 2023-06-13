//  Copyright (c) 2020 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package index

import (
	"reflect"
	"testing"
)

func TestIndexFieldDict(t *testing.T) {
	cfg, cleanup := CreateConfig("TestIndexFieldDict")
	defer func() {
		err := cleanup()
		if err != nil {
			t.Log(err)
		}
	}()

	idx, err := OpenWriter(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cerr := idx.Close()
		if cerr != nil {
			t.Fatal(cerr)
		}
	}()

	doc := &FakeDocument{
		NewFakeField("_id", "1", true, false, false),
		NewFakeField("name", "test", false, false, true),
	}
	b := NewBatch()
	b.Update(testIdentifier("1"), doc)
	err = idx.Batch(b)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}

	doc = &FakeDocument{
		NewFakeField("_id", "2", true, false, false),
		NewFakeField("name", "test test test", false, false, true),
		NewFakeField("desc", "eat more rice", false, true, true),
		NewFakeField("prefix", "bob cat cats catting dog doggy zoo", false, true, true),
	}
	b2 := NewBatch()
	b2.Update(testIdentifier("2"), doc)
	err = idx.Batch(b2)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}

	indexReader, err := idx.Reader()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = indexReader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	dict, err := indexReader.DictionaryIterator("name", nil, nil, nil)
	if err != nil {
		t.Errorf("error creating reader: %v", err)
	}
	defer func() {
		err = dict.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	termCount := 0
	curr, err := dict.Next()
	for err == nil && curr != nil {
		termCount++
		if curr.Term() != "test" {
			t.Errorf("expected term to be 'test', got '%s'", curr.Term())
		}
		curr, err = dict.Next()
	}
	if termCount != 1 {
		t.Errorf("expected 1 term for this field, got %d", termCount)
	}

	dict2, err := indexReader.DictionaryIterator("desc", nil, nil, nil)
	if err != nil {
		t.Fatalf("error creating reader: %v", err)
	}
	defer func() {
		err = dict2.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	termCount = 0
	terms := make([]string, 0)
	curr, err = dict2.Next()
	for err == nil && curr != nil {
		termCount++
		terms = append(terms, curr.Term())
		curr, err = dict2.Next()
	}
	if termCount != 3 {
		t.Errorf("expected 3 term for this field, got %d", termCount)
	}
	expectedTerms := []string{"eat", "more", "rice"}
	if !reflect.DeepEqual(expectedTerms, terms) {
		t.Errorf("expected %#v, got %#v", expectedTerms, terms)
	}
	// test start and end range
	dict3, err := indexReader.DictionaryIterator("desc", nil, []byte("fun"), []byte("nice"))
	if err != nil {
		t.Errorf("error creating reader: %v", err)
	}
	defer func() {
		err = dict3.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	termCount = 0
	terms = make([]string, 0)
	curr, err = dict3.Next()
	for err == nil && curr != nil {
		termCount++
		terms = append(terms, curr.Term())
		curr, err = dict3.Next()
	}
	if termCount != 1 {
		t.Errorf("expected 1 term for this field, got %d", termCount)
	}
	expectedTerms = []string{"more"}
	if !reflect.DeepEqual(expectedTerms, terms) {
		t.Errorf("expected %#v, got %#v", expectedTerms, terms)
	}

	// test use case for prefix
	kBeg := []byte("cat")
	kEnd := incrementBytes(kBeg)
	dict4, err := indexReader.DictionaryIterator("prefix", nil, kBeg, kEnd)
	if err != nil {
		t.Errorf("error creating reader: %v", err)
	}
	defer func() {
		err = dict4.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	termCount = 0
	terms = make([]string, 0)
	curr, err = dict4.Next()
	for err == nil && curr != nil {
		termCount++
		terms = append(terms, curr.Term())
		curr, err = dict4.Next()
	}
	if termCount != 3 {
		t.Errorf("expected 3 term for this field, got %d", termCount)
	}
	expectedTerms = []string{"cat", "cats", "catting"}
	if !reflect.DeepEqual(expectedTerms, terms) {
		t.Errorf("expected %#v, got %#v", expectedTerms, terms)
	}
}

func incrementBytes(in []byte) []byte {
	rv := make([]byte, len(in))
	copy(rv, in)
	for i := len(rv) - 1; i >= 0; i-- {
		rv[i]++
		if rv[i] != 0 {
			// didn't overflow, so stop
			break
		}
	}
	return rv
}
