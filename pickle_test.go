package pickle

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestPickle(t *testing.T) {

	os.RemoveAll("./test/output")

	// When I run function Gen and specifty a TOML input file and a output path
	// The directory should be created for the given output path
	// I should get code generated
	// And the code should compile

	in, err := filepath.Abs("./test/config/functions.toml")
	if err != nil {
		t.Error(err)
	}
	out, err := filepath.Abs("./test/output")
	if err != nil {
		t.Error(err)
	}
	Gen(in, out)

	f1, err1 := ioutil.ReadFile("./test/assert/greeter/main.go")

	if err1 != nil {
		log.Fatal(err1)
	}

	// per comment, better to not read an entire file into memory
	// this is simply a trivial example.
	f2, err2 := ioutil.ReadFile("./test/output/greeter/main.go")

	if err1 != nil {
		log.Fatal(err2)
	}

	if !bytes.Equal(f1, f2) {
		t.Errorf("Files are not equal, got: %v, want: %v.", f1, f2)
	}

	f1, err1 = ioutil.ReadFile("./test/assert/greeter/go.mod")

	if err1 != nil {
		log.Fatal(err1)
	}

	// per comment, better to not read an entire file into memory
	// this is simply a trivial example.
	f2, err2 = ioutil.ReadFile("./test/output/greeter/go.mod")

	if err1 != nil {
		log.Fatal(err2)
	}

	if !bytes.Equal(f1, f2) {
		t.Errorf("Files are not equal, got: %v, want: %v.", f1, f2)
	}

	f1, err1 = ioutil.ReadFile("./test/assert/docker-compose.yaml")

	if err1 != nil {
		log.Fatal(err1)
	}

	// per comment, better to not read an entire file into memory
	// this is simply a trivial example.
	f2, err2 = ioutil.ReadFile("./test/output/docker-compose.yaml")

	if err1 != nil {
		log.Fatal(err2)
	}

	if !bytes.Equal(f1, f2) {
		t.Errorf("Files are not equal, got: %v, want: %v.", f1, f2)
	}
}
