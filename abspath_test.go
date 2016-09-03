package abspath

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

type TestCase struct {
	input    string
	expected string
}

func tc(i, e string) TestCase {
	return TestCase{filepath.FromSlash(i), filepath.FromSlash(e)}
}

func TestNew(t *testing.T) {
	for _, c := range []TestCase{
		tc("/path/to/file", "/path/to/file"),
		tc("/", "/"),
		tc("/path includes whitespaces", "/path includes whitespaces"),
	} {
		a, err := New(c.input)
		if err != nil {
			t.Error(err)
		}
		if a.String() != c.expected {
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

	for _, c := range []TestCase{
		tc("/path/to/file", "/path/to/file"),
		tc("relative_path", abs("relative_path")),
		tc(".", abs(".")),
		tc("./foo", abs("./foo")),
		tc("../foo", abs("../foo")),
		tc("foo/bar", abs("foo/bar")),
		tc("~/.vimrc", filepath.Join(u.HomeDir, ".vimrc")),
		tc("~/foo/bar", filepath.Join(u.HomeDir, "foo/bar")),
	} {
		a, err := ExpandFrom(c.input)
		if err != nil {
			t.Error(err)
			continue
		}
		if a.String() != c.expected {
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
	for _, c := range []TestCase{
		{"/path/to/file", filepath.Clean(filepath.FromSlash("/path/to/file"))},
		{"/path////to//file/", filepath.Clean(filepath.FromSlash("/path////to//file/"))},
	} {
		a, err := FromSlash(c.input)
		if err != nil {
			t.Error(err)
			continue
		}
		if a.String() != c.expected {
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

func TestExpandFromSlash(t *testing.T) {
	for _, c := range []TestCase{
		{"relative/path", abs(filepath.Clean(filepath.FromSlash("relative/path")))},
		{"./foo/bar", abs(filepath.Clean(filepath.FromSlash("./foo/bar")))},
		{"../foo/bar", abs(filepath.Clean(filepath.FromSlash("../foo/bar")))},
		{".", abs(filepath.Clean(filepath.FromSlash(".")))},
	} {
		a, err := ExpandFromSlash(c.input)
		if err != nil {
			t.Error(err)
			continue
		}
		if a.String() != c.expected {
			t.Errorf("Expected %s but actually %s", c.expected, a)
		}
	}
}

// Note:
// Only tests a representative case because it passes to functions defined in path/filepath package.

func TestBase(t *testing.T) {
	a, _ := FromSlash("/foo/bar.poyo")
	b := a.Base()
	actual := b.String()
	expected := filepath.Base(filepath.FromSlash("/foo/bar.poyo"))
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestDir(t *testing.T) {
	a, _ := FromSlash("/foo/bar")
	b := a.Dir()
	actual := b.String()
	expected := filepath.Dir(filepath.FromSlash("/foo/bar"))
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestEvalSymlinks(t *testing.T) {
	a, _ := ExpandFromSlash("testdata/sym-link")
	b, err := a.EvalSymlinks()
	if err != nil {
		t.Error(err)
	}
	actual := b.String()
	expected, _ := filepath.Abs(filepath.FromSlash("testdata/test-file"))
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}

	c, _ := FromSlash("/foo/bar")
	_, err = c.EvalSymlinks()
	if err == nil {
		t.Errorf("Not existing file path must cause an error!")
	}
}

func TestExt(t *testing.T) {
	a, _ := FromSlash("/foo/bar.poyo")
	b := a.Ext()
	actual := string(b)
	expected := filepath.Ext(filepath.FromSlash("/foo/bar.poyo"))
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestJoin(t *testing.T) {
	a, _ := FromSlash("/foo/bar")
	b := a.Join("tsurai", "darui")
	actual := b.String()
	expected := filepath.Join(filepath.FromSlash("/foo/bar"), "tsurai", "darui")
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestMatch(t *testing.T) {
	a, _ := FromSlash("/foo/bar")
	b, err := a.Match(filepath.FromSlash("/*/*"))
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Errorf("It should match to the pattern")
	}
}

func TestRel(t *testing.T) {
	a, _ := FromSlash("/a")
	s, err := a.Rel(filepath.FromSlash("/b/c"))
	if err != nil {
		t.Fatal(err)
	}
	expected, _ := filepath.Rel(filepath.FromSlash("/a"), filepath.FromSlash("/b/c"))
	if s != expected {
		t.Errorf("Expected %s but actually %s", expected, s)
	}
}

func TestSplit(t *testing.T) {
	a, _ := FromSlash("/foo/bar.poyo")
	d, f := a.Split()
	d2, f2 := filepath.Split(filepath.FromSlash("/foo/bar.poyo"))
	if d.String() != d2 {
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
