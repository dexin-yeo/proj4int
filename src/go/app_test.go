package app

import (
	"notion/generate"
	"notion/insert"
	"notion/update"

	"fmt"
	"io/ioutil"
	"sync"
	"testing"
)

func TestGenerate(t *testing.T) {
	args := []string{"2697606ce40348569b8a690c16d6f6ef"}

	err := generate.Generate(args)
	if err != nil {
		t.Error(err)
	}
}

func TestInsertFromCommand(t *testing.T) {
	args := []string{"a5574bd6fb084bd7a4bebae3330716be", "my_table", "tmp1", "18", "", "NO", "integer"}

	err := insert.Insert(args)
	if err != nil {
		t.Error(err)
	}

	args = []string{"a5574bd6fb084bd7a4bebae3330716be", "my_table", "tmp2", "19", "", "NO", "character varying", "66"}

	err = insert.Insert(args)
	if err != nil {
		t.Error(err)
	}
}

func TestInsertFromCSV(t *testing.T) {
	args := []string{"a5574bd6fb084bd7a4bebae3330716be", "my_table", "tmp1", "18", "", "NO", "integer"}

	err := insert.Insert(args)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	args := []string{"a5574bd6fb084bd7a4bebae3330716be", "my_table", "tmp1", "18", "", "NO", "integer"}
	err := update.Update(args)

	if err != nil {
		t.Error(err)
	}
}
