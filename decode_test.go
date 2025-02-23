package main

import (
	"reflect"
	"testing"
)

type User struct {
	Name  string `csv:"name"`
	Email string `csv:"email"`
}

func TestUnmarshall(t *testing.T) {
	data := [][]string{
		{"name", "email"},
		{"mario", "m@c.com"},
		{"luigi", "l@c.com"},
	}
	want := []User{
		{"mario", "m@c.com"},
		{"luigi", "l@c.com"},
	}

	var vv []User
	err := Unmarshal(data, &vv)
	if err != nil {
		t.Fatalf("did not expect to fail: %s\n", err)
	}

	if !reflect.DeepEqual(want, vv) {
		t.Fatalf("got %+v, want %+v\n", vv, want)
	}
}

func TestCheckValidHeader(t *testing.T) {
	header := []string{
		"name", "email",
	}
	v := []User{}

	err := checkValidHeader(header, &v)
	if err != nil {
		t.Fatalf("expect no error with %+v and %+v: %s\n", header, v, err)
	}
}

func TestDecoderInit(t *testing.T) {
	d := newDecoder()
	header := []string{"name", "email"}
	vv := []User{}
	if err := d.init(header, &vv); err != nil {
		t.Fatalf("should not fail: %q\n", err)
	}
	if d.mapper[0] != "Name" {
		t.Fatalf("expected mapper 0 to be name: %s\n", d.mapper[0])
	}
	if d.mapper[1] != "Email" {
		t.Fatalf("expected mapper 0 to be name: %s\n", d.mapper[0])
	}
}
