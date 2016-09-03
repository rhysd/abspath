package abspath

import (
	"fmt"
	"os/user"
	"path/filepath"
)

type AbsPath struct {
	underlying string
}

type NotAbsolutePathError struct {
	specified string
}

func (err *NotAbsolutePathError) Error() string {
	return fmt.Sprintf("Not an absolute path: %s", err.specified)
}

func New(from string) (AbsPath, error) {
	if !filepath.IsAbs(from) {
		return AbsPath{""}, &NotAbsolutePathError{from}
	}
	return AbsPath{filepath.Clean(from)}, nil
}

func ExpandFrom(maybe_relative string) (AbsPath, error) {
	if filepath.IsAbs(maybe_relative) {
		return AbsPath{filepath.Clean(maybe_relative)}, nil
	}

	if maybe_relative == "" {
		return AbsPath{""}, fmt.Errorf("Empty path cannot be expanded")
	}

	if maybe_relative[0] == '~' {
		u, err := user.Current()
		if err != nil {
			return AbsPath{""}, err
		}
		return AbsPath{filepath.Join(u.HomeDir, maybe_relative[1:])}, nil
	}

	p, err := filepath.Abs(maybe_relative)
	if err != nil {
		return AbsPath{""}, err
	}
	return AbsPath{p}, nil
}

func FromSlash(s string) (AbsPath, error) {
	return New(filepath.FromSlash(s))
}

func (a AbsPath) Base() AbsPath {
	return AbsPath{filepath.Base(a.underlying)}
}

func (a AbsPath) Dir() AbsPath {
	return AbsPath{filepath.Dir(a.underlying)}
}

func (a AbsPath) EvalSymlinks() (AbsPath, error) {
	s, err := filepath.EvalSymlinks(a.underlying)
	if err != nil {
		return AbsPath{""}, err
	}
	s, err = filepath.Abs(s)
	if err != nil {
		return AbsPath{""}, err
	}
	return AbsPath{s}, nil
}

func (a AbsPath) Ext() string {
	return filepath.Ext(a.underlying)
}

func (a AbsPath) Join(elem ...string) AbsPath {
	return AbsPath{filepath.Join(a.underlying, filepath.Join(elem...))}
}

func (a AbsPath) Match(pattern string) (bool, error) {
	return filepath.Match(pattern, a.underlying)
}

func (a AbsPath) Rel(targpath string) (string, error) {
	return filepath.Rel(a.underlying, targpath)
}

func (a AbsPath) Split() (dir AbsPath, file string) {
	d, f := filepath.Split(a.underlying)
	return AbsPath{d}, f
}

func (a AbsPath) ToSlash() string {
	return filepath.ToSlash(a.underlying)
}

func (a AbsPath) VolumeName() string {
	return filepath.VolumeName(a.underlying)
}

func (a AbsPath) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(a.underlying, walkFn)
}

func (a AbsPath) String() string {
	return a.underlying
}
