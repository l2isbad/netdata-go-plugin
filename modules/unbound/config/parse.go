package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/*
Files can be included using the include: directive.
It can appear anywhere, it accepts a single file name as argument.

Processing continues as if the text  from  the included file was copied into the config file at that point.

If also using chroot, using full path names for the included files works, relative pathnames for the included names
work if the directory where the daemon is started equals its chroot/working directory or is specified before
the include statement with  directory:  dir. Wildcards can be used to include multiple files, see glob(7).


Unbound stop processing and exits on any error:
 - syntax error
 - recursive include
*/

type option struct{ name, value string }

func isOptionUsed(opt option) bool {
	switch opt.name {
	case
		"include",
		"statistics-cumulative",
		"control-enable",
		"control-interface",
		"control-port",
		"control-use-cert",
		"control-use-key-file",
		"control-use-cert-file":
		return true
	}
	return false
}

// Parse
func Parse(entryPath string) (*UnboundConfig, error) {
	options, err := parse(entryPath, nil)
	if err != nil {
		return nil, err
	}
	return fromOptions(options), nil
}

func parse(filename string, visited map[string]bool) ([]option, error) {
	if visited == nil {
		visited = make(map[string]bool)
	}
	visited[filename] = true

	f, err := open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var options []option
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		opt, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("file '%s', error on parsing line '%s': %v", filename, line, err)
		}

		if !isOptionUsed(opt) {
			continue
		}

		if opt.name != "include" {
			options = append(options, opt)
			continue
		}

		opts, err := handleInclude(opt.value, visited)
		if err != nil {
			return nil, err
		}
		options = append(options, opts...)
	}
	return options, nil
}

func handleInclude(include string, visited map[string]bool) ([]option, error) {
	filenames, err := filenamesFromInclude(include)
	if err != nil {
		return nil, err
	}
	var options []option
	for _, filename := range filenames {
		opts, err := parse(filename, visited)
		if err != nil {
			return nil, err
		}
		options = append(options, opts...)
	}
	return options, nil
}

func filenamesFromInclude(include string) ([]string, error) {
	if isGlobPattern(include) {
		return filepath.Glob(include)
	}
	return []string{include}, nil
}

func parseLine(line string) (option, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return option{}, errors.New("bad syntax")
	}
	key, value := cleanKeyValue(parts[0], parts[1])
	return option{name: key, value: value}, nil
}

func cleanKeyValue(key, value string) (string, string) {
	if i := strings.IndexByte(value, '#'); i > 0 {
		value = value[:i-1]
	}
	key = strings.TrimSpace(key)
	value = strings.Trim(strings.TrimSpace(value), "\"'")
	return key, value
}

func isGlobPattern(value string) bool {
	magicChars := `*?[`
	if runtime.GOOS != "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(value, magicChars)
}

func open(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("'%s' is not a regular file", filename)
	}
	return f, nil
}
