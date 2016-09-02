package abspath

import (
	"fmt"
	"path/filepath"
)

type AbsPath string

type NotAbsolutePathError struct {
	specified string
}

func (err *NotAbsolutePathError) Error() string {
	return fmt.Sprintf("Not an absolute path: %s", err.specified)
}

func New(from string) (AbsPath, error) {
	if !filepath.IsAbs(from) {
		return AbsPath(""), &NotAbsolutePathError{from}
	}
	return AbsPath(from), nil
}

func ExpandFrom(maybe_relative string) (AbsPath, error) {
	if filepath.IsAbs(maybe_relative) {
		return AbsPath(filepath.Clean(maybe_relative)), nil
	}

	if maybe_relative == "" {
		return AbsPath(""), fmt.Errorf("Empty path cannot be expanded")
	}

	if maybe_relative[0] == '~' {
		u, err := user.Current()
		if err != nil {
			return AbsPath(""), err
		}
		return AbsPath(filepath.Join(u.HomeDir, s[1:])), nil
	}

	p, err := filepath.Abs()
	if err != nil {
		return AbsPath(""), err
	}
	return AbsPath(p), nil
}

func FromSlash(s string) (AbsPath, error) {
	return New(filepath.FromSlash(s))
}

func (a AbsPath) Base(path AbsPath) AbsPath {
	return AbsPath(filepath.Base(string(AbsPath)))
}

func (a AbsPath) Dir(path AbsPath) AbsPath {
	return AbsPath(filepath.Dir(string(AbsPath)))
}

func (a AbsPath) EvalSymlinks(path AbsPath) (AbsPath, error) {
	s, err := filepath.EvalSymlinks(string(path))
	if err != nil {
		return AbsPath(""), err
	}
	a, err := filepath.Abs(s)
	if err != nil {
		return AbsPath(""), err
	}
	return AbsPath(a), nil
}

func (a AbsPath) Ext(path AbsPath) string {
	return filepath.Ext(string(path))
}

func (a AbsPath) HasPrefix(prefix AbsPath) bool {
	return filepath.HasPrefix(string(a), prefix)
}

func (a AbsPath) Join(elem ...string) AbsPath {
	return absoltuePath(filepath.Join(string(a), elem...))
}

func (a AbsPath) Match(pattern string) (bool, error) {
	return filepath.Match(pattern, string(a))
}

func (a AbsPath) Rel(targpath string) (AbsPath, error) {
	return filepath.Rel(string(a), targpath)
}

func (a AbsPath) Split() (dir AbsPath, file string) {
	d, f := filepath.Split(string(a))
	return AbsPath(d), f
}

func (a AbsPath) ToSlash() string {
	return filepath.ToSlash(string(a))
}

func (a AbsPath) VolumeName() string {
	return filepath.VolumeName(string(a))
}

func (a AbsPath) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(string(a), walkFn)
}
