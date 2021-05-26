// Copyright 2021 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openapigen

import (
	"testing"
)

func TestTodoAPI(t *testing.T) {
	want := 4428
	result := GenerateSpec("https://jsonplaceholder.typicode.com/todos")

	if got := len(result); got != want {
		t.Errorf("GenerateSpec() length = %d, want %d", got, want)
	}
}

func TestOrderAPI(t *testing.T) {
	want := 5833
	result := GenerateSpec("https://emea-poc13-test.apigee.net/business-objects-api/orders")

	if got := len(result); got != want {
		t.Errorf("GenerateSpec() length = %d, want %d", got, want)
	}
}
