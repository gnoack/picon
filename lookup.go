// Package picon implements picon lookup.
//
// To use this, see https://kinzler.com/picons/ftp/ for a repository
// of icons which you can download to your ~/.picons directory. The
// structure should be:
//
//	.picons/
//	    local/...
//	    misc/...
//	    unknown/...
//
// etc.
//
// This package is hacky and the API is unstable.
package picon

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Lookup returns the fully qualified filename with the best picon
// file for the given email.
func Lookup(email string) (filename string, ok bool) {
	cs, err := parseAddress(email)
	if err != nil {
		return "", false
	}
	var todo []string
	for i := len(cs); i > 0; i-- {
		path := filepath.Join(cs[0:i]...)

		if i < len(cs) {
			path = filepath.Join(path, "unknown")
		}
		todo = append(todo, path)
	}
	todo = append(todo, filepath.Join("MISC", cs[len(cs)-1]))
	todo = append(todo, "MISC/noface")

	dbs := getDBs()
	for _, path := range todo {
		for _, db := range dbs {
			for _, name := range []string{"face.png", "face.jpg", "face.jpeg", "face.gif", "face.xpm", "face.xbm"} {
				path := filepath.Join(db, path, name)
				_, err := os.Stat(path)
				if err != nil {
					continue
				}
				return path, true
			}
		}
	}
	return "", false
}

func parseAddress(email string) ([]string, error) {
	n, domain, ok := strings.Cut(email, "@")
	if !ok {
		return nil, errors.New("invalid address")
	}
	return append(reverse(strings.Split(domain, ".")), n), nil
}

// orderedSubdirs is the default subdirectory iteration order;
// compare https://kinzler.com/picons/ftp/faq.html#lookup
var orderedSubdirs = []string{
	"local", "users", "usenix", "misc", "domains", "unknown",
}

func isWellKnownSubdir(name string) bool {
	for _, o := range orderedSubdirs {
		if name == o {
			return true
		}
	}
	return false
}

func orderedSubdirNames(baseDir string) []string {
	var res []string

	dbs, err := os.ReadDir(baseDir)
	switch {
	case os.IsNotExist(err):
		return nil
	case os.IsPermission(err):
		// Workaround for Landlock eCryptfs issue, see
		// https://lore.kernel.org/linux-security-module/c1c9c688-c64d-adf2-cc96-dc2aaaae5944@digikod.net/
		// We can not list directories, but we can access them directly,
		// so as a fallback we assume that the well-known subdirectories exist.
		res = orderedSubdirs
	case err != nil:
		return nil
	default:
		for _, e := range dbs {
			if !e.IsDir() {
				continue
			}
			// Well-known subdirs are appended last.
			if isWellKnownSubdir(e.Name()) {
				continue
			}
			res = append(res, e.Name())
		}
		res = append(res, orderedSubdirs...)
	}

	for i, r := range res {
		res[i] = filepath.Join(baseDir, r)
	}
	return res
}

func getDBs() []string {
	var res []string
	for _, base := range []string{
		"/usr/share/picons", // where it's installed on Debian
		filepath.Join(os.Getenv("HOME"), ".picons"),
	} {
		res = append(res, orderedSubdirNames(base)...)
	}
	return res
}

func reverse(s []string) []string {
	for i := 0; i < len(s)/2; i++ {
		tmp := s[i]
		s[i] = s[len(s)-1-i]
		s[len(s)-1-i] = tmp
	}
	return s
}
