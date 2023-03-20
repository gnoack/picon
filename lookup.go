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

func getDBs() []string {
	var res []string
	for _, base := range []string{
		"/usr/share/picons", // where it's installed on Debian
		filepath.Join(os.Getenv("HOME"), ".picons"),
	} {
		dbs, err := os.ReadDir(base)
		if os.IsNotExist(err) {
			continue
		} else if os.IsPermission(err) {
			// Potential Landlock eCryptfs issue, see
			// https://lore.kernel.org/linux-security-module/c1c9c688-c64d-adf2-cc96-dc2aaaae5944@digikod.net/
			//
			// Let's assume the default iteration order
			// https://kinzler.com/picons/ftp/faq.html#lookup
			for _, e := range []string{
				"local", "users", "usenix", "misc", "domains", "unknown",
			} {
				res = append(res, filepath.Join(base, e))
			}
			continue
		} else if err != nil {
			return nil
		}
		for _, e := range dbs {
			if !e.IsDir() {
				continue
			}
			res = append(res, filepath.Join(base, e.Name()))
		}
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
