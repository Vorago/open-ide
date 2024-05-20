package main

import (
	"flag"
	"fmt"
	"go.i3wm.org/i3/v4"
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

	fmt.Println(*maxDepth)
	fmt.Println(*rootDir)
	projectNames := searchProjects(*rootDir, *maxDepth)
	projectName := pickProject(projectNames)

	windowIds := searchWindow(projectName)
	if len(windowIds) == 0 {
		openProject(*idePath, *rootDir+projectName)
	} else if len(windowIds) == 1 {
		focusWindow(windowIds[0])
	} else {
		focusWindow(windowIds[0])
	}
}

func searchProjects(rootDir string, maxDepth int) []string {
	projects := make([]string, 0)
	_ = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
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

	return projects
}

func searchWindow(projectName string) []string {
	windowName := filepath.Base(projectName)
	byName, err := exec.Command("xdotool", "search", "--name", windowName).Output()
	if err != nil {
		return make([]string, 0)
	}

	byClass, err := exec.Command("xdotool", "search", "--class", "jetbrains-idea").Output()
	if err != nil {
		return make([]string, 0)
	}

	byNameArr := strings.Split(string(byName), "\n")
	byClassArr := strings.Split(string(byClass), "\n")

	return intersect(byNameArr, byClassArr)
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
			_, _ = io.WriteString(stdin, p+"\n")
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(out))
}

func focusWindow(windowId string) {
	_, err := i3.RunCommand(fmt.Sprintf("[id=\"%s\"] focus", windowId))
	if err != nil {
		panic(err)
	}
}

func openProject(idePath string, projectPath string) {
	_ = exec.Command(idePath, projectPath).Run()
}

func intersect(a []string, b []string) []string {
	res := make([]string, 0)
	for i := 0; i < len(a); i++ {
		if contains(b, a[i]) {
			res = append(res, a[i])
		}
	}

	return res
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
