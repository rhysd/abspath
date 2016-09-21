package abspath

import (
	"fmt"
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

func TestGetwd(t *testing.T) {
	expected, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	actual, err := Getwd()
	if err != nil {
		t.Fatalf("Getwd() unexpectedly returns an error although os.Getwd() doesn't return an error: %s", err.Error())
	}
	if actual.String() != expected {
		t.Fatalf("Expected Getwd() to return '%s' but actually did '%s'", expected, actual.String())
	}
}

func TestHomeDir(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	h, err := HomeDir()
	if err != nil {
		t.Fatal(err)
	}

	if h.String() != u.HomeDir {
		t.Fatalf("'%s' is expected as home directory, but actually '%s'", u.HomeDir, h.String())
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
	actual := b
	expected := filepath.Ext(filepath.FromSlash("/foo/bar.poyo"))
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}
}

func TestHasPrefix(t *testing.T) {
	a, _ := FromSlash("/foo/bar.poyo")
	b := filepath.FromSlash("/foo")
	actual := a.HasPrefix(b)
	expected := filepath.HasPrefix(filepath.FromSlash("/foo/bar.poyo"), b)
	if actual != expected {
		t.Errorf("Expected %v but actually %v", expected, actual)
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

	c := a.Join()
	actual = c.String()
	expected = filepath.Join(filepath.FromSlash("/foo/bar"))
	if actual != expected {
		t.Errorf("Expected %s but actually %s", expected, actual)
	}

	d := a.Join("piyo")
	actual = d.String()
	expected = filepath.Join(filepath.FromSlash("/foo/bar"), "piyo")
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
	return filepath.Base(a.underlying)
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

func BenchmarkLongPath(b *testing.B) {
	p := ""
	for i := 0; i < 10; i++ {
		for c := 'a'; c <= 'z'; c++ {
			p += fmt.Sprintf("/%c%c%c_%c%c%c", c, c, c, c, c, c)
		}
	}
	s := filepath.FromSlash(p)
	a, _ := New(s)

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

func Example() {
	var a AbsPath
	var err error

	// From a string representing absolute path.
	// If it doesn't stand for absolute path, an error will be returned.
	a, err = New("/absolute/path")
	if err != nil {
		panic(err)
	}
	fmt.Println(a.String())

	// Expanded to $PWD/relative_path
	a, err = ExpandFrom("relative_path")
	if err != nil {
		panic(err)
	}
	fmt.Println(a.String())

	// Expanded to $HOME/relative_path
	a, err = ExpandFrom("~/relative_path")
	if err != nil {
		panic(err)
	}
	fmt.Println(a.String())

	// Slashes in the string will be replaced by a file separator.
	// In Windows, below can be 'C:\absolute\path'
	a, err = FromSlash("/absoltue/path")
	if err != nil {
		panic(err)
	}
	fmt.Println(a.String())

	// Similar to ExpandFrom().  But all slashes will be replaced with a file sperator.
	a, err = ExpandFromSlash("relative/path")
	if err != nil {
		panic(err)
	}
	fmt.Println(a.String())

	// Check the path exists/is file/is directory.
	if _, err := os.Stat(a.String()); err == nil {
		fmt.Printf("'%s' exists", a.String())
	}
	if s, err := os.Stat(a.String()); err == nil && s.IsDir() {
		fmt.Printf("'%s' is a directory", a.String())
	}
	if s, err := os.Stat(a.String()); err == nil && !s.IsDir() {
		fmt.Printf("'%s' is a file", a.String())
	}
}
