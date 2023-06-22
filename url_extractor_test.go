package main

import (
	"reflect"
	"testing"
)

func Test_extractURLs(t *testing.T) {
	testData := " <a href=\"https://www.example.com\">Visit example.com</a> \n " +
		"<a href=\"https://www.test.com\">Visit example.com</a>" +
		"<div></div>"

	want := map[string]struct{}{
		"https://www.test.com":    {},
		"https://www.example.com": {},
	}

	got, err := extractURLs([]byte(testData))
	if err != nil {
		t.Errorf("err [%v]", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wanted %v got %v", want, got)
	}
}
