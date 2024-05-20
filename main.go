package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	maxDepth := flag.Int("depth", 1, "max directory depth to look for projects")
	rootDir := flag.String("codeDir", "", "directory with code")
	idePath := flag.String("ideCommand", "", "ide binary location")

	flag.Parse()

	*maxDepth += strings.Count(*rootDir, string(os.PathSeparator))

	projectNames := searchProjects(*rootDir, *maxDepth)
	projectName := pickProject(projectNames)

	openProject(*idePath, *rootDir+projectName)
}

func searchProjects(rootDir string, maxDepth int) []string {
	projects := make([]string, 0)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.Count(path, string(os.PathSeparator)) > maxDepth {
			return fs.SkipDir
		}

		if d.Name() == "node_modules" {
			return fs.SkipDir
		}

		if d.Name() == ".git" {
			projects = append(projects, strings.TrimPrefix(strings.TrimSuffix(path, "/.git"), rootDir))
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return projects
}

func projectRegexp(windowName string) string {
	return fmt.Sprintf("^%s( – |$)", windowName)
}

func pickProject(projects []string) string {
	cmd := exec.Command("rofi", "-dmenu")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		for _, p := range projects {
			_, err = io.WriteString(stdin, p+"\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(out))
}

func openProject(idePath string, projectPath string) {
	err := exec.Command(idePath, projectPath).Run()
	if err != nil {
		log.Fatal(err)
	}
}
