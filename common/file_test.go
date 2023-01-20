package common_test

import (
	"demo-store/common"
	"testing"
)

func TestDirExists(t *testing.T) {

	if ok := common.DirExists("../common"); !ok {
		t.Errorf("Unexpected error: dir /users/: got %v want %v,", ok, "true")
	}
}

func TestDirNotExists(t *testing.T) {

	if ok := common.DirExists("../somedir"); ok {
		t.Errorf("Unexpected error: dir /users/: got %v want %v,", ok, "false")
	}
}

func TestToJsonSuccess(t *testing.T) {

	data := "somedata"
	res, _ := common.ToJson(data)
	expected := "\"somedata\""
	if res != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res, expected)
	}
}

func TestToJsonFailed(t *testing.T) {

	d := make(chan int)
	_, err := common.ToJson(d)
	if err == nil {
		t.Errorf("handler returned unexpected json: got nill want %v", err)
	}
}
