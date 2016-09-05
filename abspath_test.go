package abspath

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
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
		if !strings.HasPrefix(err.Error(), "Not an absolute path: ") {
			t.Errorf("Unexpected kind of error: %s", err.Error())
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

func TestExists(t *testing.T) {
	a, _ := ExpandFrom("unknown_file")
	if a.Exists() {
		t.Errorf("Unknown file was detected as exist")
	}

	f, err := os.OpenFile("_test_file", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	f.Close()
	defer os.Remove("_test_file")
	b, _ := ExpandFrom("_test_file")
	if !b.Exists() {
		t.Errorf("Existing file was not found")
	}
}

func TestIsDir(t *testing.T) {
	a, _ := ExpandFrom("unknown_dir")
	if a.IsDir() {
		t.Errorf("Unknown directory was detected as exist")
	}

	f, err := os.OpenFile("_test_file", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	f.Close()
	defer os.Remove("_test_file")
	b, _ := ExpandFrom("_test_file")
	if b.IsDir() {
		t.Errorf("Existing file was detected as directory")
	}

	err = os.Mkdir("_test_dir", os.ModeDir|os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer os.Remove("_test_dir")
	c, _ := ExpandFrom("_test_dir")
	if !c.IsDir() {
		t.Errorf("Existing directory was not found")
	}
}

func TestIsFile(t *testing.T) {
	a, _ := ExpandFrom("unknown_file")
	if a.IsFile() {
		t.Errorf("Unknown file was detected as exist")
	}

	f, err := os.OpenFile("_test_file", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	f.Close()
	defer os.Remove("_test_file")
	b, _ := ExpandFrom("_test_file")
	if !b.IsFile() {
		t.Errorf("Existing file was not detected")
	}

	err = os.Mkdir("_test_dir", os.ModeDir|os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer os.Remove("_test_dir")
	c, _ := ExpandFrom("_test_dir")
	if c.IsFile() {
		t.Errorf("Existing directory was not found as file")
	}
}

func TestEquals(t *testing.T) {
	a, _ := FromSlash("/foo/bar")
	b, _ := FromSlash("/foo/bar")
	if a != b {
		t.Errorf("Expected '%s' == '%s'", a, b)
	}

	c, _ := FromSlash("/foo/bar")
	d, _ := FromSlash("/piyo/poyo")
	if c == d {
		t.Errorf("Expected '%s' != '%s'", c, d)
	}
}

func TestStat(t *testing.T) {
	f, err := os.OpenFile("_test_file", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	f.Close()
	defer os.Remove("_test_file")
	a, _ := ExpandFrom("_test_file")
	s, err := a.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if s.IsDir() {
		t.Errorf("_test_file is not an directory")
	}
}

func funcPassByValue(a AbsPath) string {
	return a.String()
}

func funcPassByRef(a *AbsPath) string {
	return a.String()
}

func BenchmarkPassBy(b *testing.B) {
	a, _ := FromSlash("/path/to/entry")

	b.Run("Pass by value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			funcPassByValue(a)
		}
	})

	b.Run("Pass by reference", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			funcPassByRef(&a)
		}
	})
}

func funcRawBase(a *AbsPath) string {
	return filepath.Base(a.String())
}

func BenchmarkBaseMethod(b *testing.B) {
	a, _ := FromSlash("/path/to/entry")
	s := "/path/to/entry"

	b.Run("Fastest raw Base() function", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filepath.Base(s)
		}
	})

	b.Run("Equivalent Base() function", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			funcRawBase(&a)
		}
	})

	b.Run("Original Base() method", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a.Base()
		}
	})
}
