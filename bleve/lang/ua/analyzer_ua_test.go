//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package ua

import (
	"reflect"
	"testing"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
)

func TestUkrainianAnalyzer(t *testing.T) {
	tests := []struct {
		input  []byte
		output analysis.TokenStream
	}{
		// digits safe
		{
			input: []byte("text 1000"),
			output: analysis.TokenStream{
				&analysis.Token{
					Term: []byte("text"),
				},
				&analysis.Token{
					Term: []byte("1000"),
				},
			},
		},
		{
			input: []byte("Разом з тим про силу електромагнітної енергії мали уявлення ще"),
			output: analysis.TokenStream{
				&analysis.Token{
					Term: []byte("вміст"),
				},
				&analysis.Token{
					Term: []byte("сил"),
				},
				&analysis.Token{
					Term: []byte("електромагнітн"),
				},
				&analysis.Token{
					Term: []byte("енерг"),
				},
				&analysis.Token{
					Term: []byte("мав"),
				},
				&analysis.Token{
					Term: []byte("представлен"),
				},
			},
		},
		{
			input: []byte("Але знання це зберігалося в таємниці"),
			output: analysis.TokenStream{
				&analysis.Token{
					Term: []byte("знан"),
				},
				&analysis.Token{
					Term: []byte("це"),
				},
				&analysis.Token{
					Term: []byte("збер"),
				},
				&analysis.Token{
					Term: []byte("таем"),
				},
			},
		},
	}

	cache := registry.NewCache()
	analyzer, err := cache.AnalyzerNamed(AnalyzerName)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		actual := analyzer.Analyze(test.input)
		if len(actual) != len(test.output) {
			t.Fatalf("expected length: %d, got %d", len(test.output), len(actual))
		}
		for i, tok := range actual {
			if !reflect.DeepEqual(tok.Term, test.output[i].Term) {
				t.Errorf("expected term %s (% x) got %s (% x)", test.output[i].Term, test.output[i].Term, tok.Term, tok.Term)
			}
		}
	}
}
