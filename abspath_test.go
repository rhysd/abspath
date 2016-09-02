package abspath

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	testcases := []struct {
		input    string
		expected string
	}{
		{"/path/to/file", "/path/to/file"},
		{"/", "/"},
		{"/path includes whitespaces", "/path includes whitespaces"},
	}

	for _, c := range testcases {
		a, err := New(c.input)
		if err != nil {
			t.Error(err)
		}
		if string(a) != c.expected {
			t.Errorf("Expected %s but actually %s", c.expected, a)
		}
	}

	errorcases := []string{
		"relative_path",
		"",
	}
	for _, e := range errorcases {
		_, err := New(e)
		if err == nil {
			t.Errorf("Error was expected for input '%s'", e)
		}
	}
}

func abs(s string) string {
	r, err := filepath.Abs(s)
	if err != nil {
		panic(err)
	}
	return r
}

func TestExpandFrom(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	for _, c := range []struct {
		input    string
		expected string
	}{
		{"/path/to/file", "/path/to/file"},
		{"relative_path", abs("relative_path")},
		{"./foo", abs("./foo")},
		{"../foo", abs("../foo")},
		{"foo/bar", abs("foo/bar")},
		{"~/.vimrc", filepath.Join(u.HomeDir, ".vimrc")},
		{"~/foo/bar", filepath.Join(u.HomeDir, "foo/bar")},
	} {
		a, err := ExpandFrom(c.input)
		if err != nil {
			t.Error(err)
		}
		if string(a) != c.expected {
			t.Errorf("Expected %s but actually %s", c.expected, a)
		}
	}

	for _, e := range []string{
		"",
	} {
		_, err := ExpandFrom(e)
		if err == nil {
			t.Errorf("Error was expected for input '%s'", e)
		}
	}
}

func TestFromSlash(t *testing.T) {
	for _, c := range []struct {
		input    string
		expected string
	}{
		{"/path/to/file", filepath.Clean(filepath.FromSlash("/path/to/file"))},
		{"/path////to//file/", filepath.Clean(filepath.FromSlash("/path////to//file/"))},
	} {
		a, err := ExpandFrom(c.input)
		if err != nil {
			t.Error(err)
		}
		if string(a) != c.expected {
			t.Errorf("Expected %s but actually %s", c.expected, a)
		}
	}

	for _, e := range []string{
		"path/to/file",
		"",
		"relative_path",
	} {
		_, err := FromSlash(e)
		if err == nil {
			t.Errorf("Error was expected for input '%s'", e)
		}
	}
}

// Note:
// Only tests a representative case because it passes to functions defined in path/filepath package.

func TestBase(t *testing.T) {
	a, _ := New("/foo/bar.poyo")
	b := a.Base()
	actual := string(b)
	expected := filepath.Base("/foo/bar.poyo")
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestDir(t *testing.T) {
	a, _ := New("/foo/bar")
	b := a.Dir()
	actual := string(b)
	expected := filepath.Dir("/foo/bar")
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestEvalSymlinks(t *testing.T) {
	a, _ := ExpandFrom("testdata/sym-link")
	b, err := a.EvalSymlinks()
	if err != nil {
		t.Error(err)
	}
	actual := string(b)
	expected, _ := filepath.Abs("testdata/test-file")
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}

	c, _ := New("/foo/bar")
	_, err = c.EvalSymlinks()
	if err == nil {
		t.Errorf("Not existing file path must cause an error!")
	}
}

func TestExt(t *testing.T) {
	a, _ := New("/foo/bar.poyo")
	b := a.Ext()
	actual := string(b)
	expected := filepath.Ext("/foo/bar.poyo")
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestJoin(t *testing.T) {
	a, _ := New("/foo/bar")
	b := a.Join("tsurai", "darui")
	actual := string(b)
	expected := filepath.Join("/foo/bar", "tsurai", "darui")
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestMatch(t *testing.T) {
	a, _ := New("/foo/bar")
	b, err := a.Match("/*/*")
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Errorf("It should match to the pattern")
	}
}

func TestRel(t *testing.T) {
	a, _ := New("/a")
	s, err := a.Rel("/b/c")
	if err != nil {
		t.Fatal(err)
	}
	expected, _ := filepath.Rel("/a", "/b/c")
	if s != expected {
		t.Errorf("Expected %s but actually %s", expected, s)
	}
}

func TestSplit(t *testing.T) {
	a, _ := New("/foo/bar.poyo")
	d, f := a.Split()
	d2, f2 := filepath.Split("/foo/bar.poyo")
	if string(d) != d2 {
		t.Errorf("Expected %s but actually %s", d2, d)
	}
	if f != f2 {
		t.Errorf("Expected %s but actually %s", f2, f)
	}
}

func TestToSlash(t *testing.T) {
	a, _ := FromSlash("/foo/bar.poyo")
	s := a.ToSlash()
	expected := filepath.Clean(filepath.FromSlash("/foo/bar.poyo"))
	if s != expected {
		t.Errorf("Expected %s but actually %s", expected, s)
	}
}

func TestVolumeName(t *testing.T) {
	a, _ := New("C:/foo/bar.poyo")
	v := a.VolumeName()
	expected := filepath.VolumeName("C:/foo/bar.poyo")
	if v != expected {
		t.Errorf("Expected %s but actually %s", expected, v)
	}
}

func TestWalk(t *testing.T) {
	w := func(p string, info os.FileInfo, err error) error {
		return nil
	}

	a, _ := ExpandFrom(".")
	err := a.Walk(w)
	if err != nil {
		t.Fatal(err)
	}
}
