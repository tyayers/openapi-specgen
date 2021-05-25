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
