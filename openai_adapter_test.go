package main

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestFetchCompletionStream(t *testing.T) {
	expected := []string{"these", "tokens", "are", "returned"}
	tokenResponses := make([]string, 0)

	client := &OpenAIAdapter{APIKey: "mykey"}

	respCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := client.FetchCompletionStream(CreateRequest("prompt"), respCh)
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()

	if !reflect.DeepEqual(expected, tokenResponses) {
		t.Errorf("Expected %v\nGot %v\n", expected, tokenResponses)
	}

}
