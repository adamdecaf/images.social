package upload

import (
	"os"
	"path"
	"testing"
)

func TestImageDetection(t *testing.T) {
	cases := []struct {
		name string
		tpe  ImageType
	}{
		{"hat.jpg", ImageType("jpeg")},
		{"hat.png", ImageType("png")},
		{"hat.gif", ImageType("gif")},
	}
	for _, v := range cases {
		f, err := os.Open(path.Join("../testdata/", v.name))
		if err != nil {
			t.Errorf("error reading test file, err=%v", err)
		}
		tpe, err := Detect(f)
		if err != nil {
			t.Errorf("error detecting on file '%s'", v.name)
		}
		if tpe != v.tpe {
			t.Errorf("detected type '%s' not equal to expected '%s'", tpe.Ext(), v.tpe.Ext())
		}
	}
}
