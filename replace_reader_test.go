package reader

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestReplaceReader(t *testing.T) {
	str := `j'aimerai vous dire que je n'aime pas le collage
	des caractères -- ensemble. Changez ca tout de suite, je n'en veux plus de tous ces --
	et faites le maintenant`

	reader := strings.NewReader(str)

	rpl := NewReplaceReader(reader, "--", "++")

	var res []byte
	var err error
	res, err = ioutil.ReadAll(rpl)
	if err != nil {
		t.Fatal(err)
	}

	expected := `j'aimerai vous dire que je n'aime pas le collage
	des caractères ++ ensemble. Changez ca tout de suite, je n'en veux plus de tous ces ++
	et faites le maintenant`

	if string(res) != expected {
		t.Fatalf("\nexpected : \n%v \n received : \n%v", expected, string(res))
	}

}

func TestReplaceReaderMoreComplicated(t *testing.T) {
	str := `j'aimerai vous dire que je n'aime pas le collage
	des caractères -- ensemble. Changez ca tout de suite, je n'en veux plus de tous ces --
	et faites le maintenant`

	reader := strings.NewReader(str)

	rpl := NewReplaceReader(reader, "dire", "dite")

	var res []byte
	var err error
	res, err = ioutil.ReadAll(rpl)
	if err != nil {
		t.Fatal(err)
	}

	expected := `j'aimerai vous dite que je n'aime pas le collage
	des caractères -- ensemble. Changez ca tout de suite, je n'en veux plus de tous ces --
	et faites le maintenant`

	if string(res) != expected {
		t.Fatalf("\nexpected : \n%v \n received : \n%v", expected, string(res))
	}

}

func TestReplaceReaderMoreMoreComplicated(t *testing.T) {
	str := `j'aimerai vous dire que je n'aime pas le collage
	des caractères -- ensemble. Changez ca tout de suite, je n'en veux plus de tous ces --
	et faites le maintenant`

	reader := strings.NewReader(str)

	rpl := NewReplaceReader(reader, "collage", "fait de coller")

	var res []byte
	var err error
	res, err = ioutil.ReadAll(rpl)
	if err != nil {
		t.Fatal(err)
	}

	expected := `j'aimerai vous dire que je n'aime pas le fait de coller
	des caractères -- ensemble. Changez ca tout de suite, je n'en veux plus de tous ces --
	et faites le maintenant`

	if string(res) != expected {
		t.Fatalf("\nexpected : \n%v \n received : \n%v", expected, string(res))
	}

}
