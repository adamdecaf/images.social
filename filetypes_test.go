package main

import (
	"os"
	"path"
	"testing"
)

func TestImageDetection(t *testing.T) {
	cases := []struct{
		name string
		tpe ImageType
	} {
		{ "hat.jpg", JPEG },
		{ "hat.png", PNG },
		{ "hat.gif", GIF },
	}
	for _, v := range cases {
		f, err := os.Open(path.Join("testdata/", v.name))
		if err != nil {
			t.Errorf("error reading test file, err=%v", err)
		}
		tpe := Detect(f)
		if tpe != v.tpe {
			t.Errorf("detected type '%s' not equal to expected '%s'", tpe.Ext(), v.tpe.Ext())
		}
	}
}
