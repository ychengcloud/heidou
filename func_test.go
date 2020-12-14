package heidou

import (
	"testing"
)

func TestPascal(t *testing.T) {
	id := pascal("ID")
	t.Log(id)
	if id != "ID" {
		t.Error()
	}
	id = pascal("user_id")
	t.Log(id)
	if id != "UserID" {
		t.Error()
	}
}

func TestCamel(t *testing.T) {
	id := camel("ID")
	t.Log(id)
	if id != "id" {
		t.Error()
	}

	id = camel("Id")
	t.Log(id)
	if id != "id" {
		t.Error()
	}

	id = camel("user_id")
	t.Log(id)
	if id != "userID" {
		t.Error()
	}
}