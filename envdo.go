package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var passhroughEnvVars = []string{
	"LANG",
	"LC_CTYPE",
	"TERM",
	"SHELL",
	"PATH",
	"PWD",
	"HOME",
	"USER",
	"LOGNAME",
}

func readFiles(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var allFiles []string

	for _, f := range files {
		if f.Mode().IsDir() {
			if strings.HasPrefix(f.Name(), ".") {
				continue
			}

			subdirFiles, err := readFiles(path + "/" + f.Name())
			if err != nil {
				return nil, err
			}

			allFiles = append(allFiles, subdirFiles...)

			continue
		}

		allFiles = append(allFiles, path+"/"+f.Name())
	}

	return allFiles, nil
}

func profileName(base, path string) string {
	return strings.TrimPrefix(path, base+"/")
}

func absProfilePath(path, profileName string) (string, error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	if stat.Mode() == os.ModeSymlink {
		path, err = os.Readlink(path)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s/%s", path, profileName), nil
}

func passtroughEnvWithValues() []string {
	var evs []string

	for _, name := range passhroughEnvVars {
		evs = append(evs, fmt.Sprintf("%s=%s", name, os.Getenv(name)))
	}

	return evs
}

func fetchEnv(path, profile string, recipient *string) ([]string, error) {
	absPath, err := absProfilePath(path, profile)
	if err != nil {
		return nil, err
	}

	var reader io.Reader

	ext := filepath.Ext(absPath)
	if ext == ".gpg" {
		var b []byte
		var gpgArgs []string

		if recipient != nil {
			gpgArgs = append(gpgArgs, []string{"-r", *recipient}...)
		}

		gpgArgs = append(gpgArgs, []string{"-d", "--quiet", absPath}...)

		buf := bytes.NewBuffer(b)

		err := executeCmdWithWriter("gpg2", gpgArgs, true, nil, buf)
		if err != nil {
			return nil, err
		}

		reader = buf
	} else {
		f, err := os.Open(absPath)
		if err != nil {
			return nil, err
		}

		reader = f

		defer f.Close()
	}

	var lines []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix("#", line) {
			continue
		}

		lines = append(lines, strings.Replace(line, `"`, ``, -1))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func executeCmd(command string, args []string, preserveEnv bool, env []string) error {
	return executeCmdWithWriter(command, args, preserveEnv, env, os.Stdout)
}

func executeCmdWithWriter(command string, args []string, preserveEnv bool, env []string, writer io.Writer) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = writer

	// Default variables to pass through.
	cmd.Env = passtroughEnvWithValues()

	if preserveEnv {
		cmd.Env = os.Environ()
	}

	cmd.Env = append(cmd.Env, env...)

	return cmd.Run()
}
