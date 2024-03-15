package main

import (
	"flag"
	"fmt"
	"github.com/jxeng/shortcut"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type FilePair struct {
	Name string
	Path string
	Ico  bool
	Rdp  bool
}

func listFiles(cwd string) (files map[string]FilePair, err error) {
	fileInfo, err := os.ReadDir(cwd)
	if err != nil {
		log.Fatal(err)
	}

	files = make(map[string]FilePair)
	for _, file := range fileInfo {
		baseName := filepath.Base(file.Name())
		ext := filepath.Ext(baseName)
		nameWithoutExt := strings.TrimSuffix(baseName, ext)

		prev, ok := files[nameWithoutExt]
		if ok {
			files[nameWithoutExt] = FilePair{
				Name: nameWithoutExt,
				Path: cwd,
				Ico:  prev.Ico || ext == ".ico",
				Rdp:  prev.Rdp || ext == ".rdp",
			}
		} else {
			files[nameWithoutExt] = FilePair{
				Name: nameWithoutExt,
				Path: cwd,
				Ico:  ext == ".ico",
				Rdp:  ext == ".rdp",
			}
		}
	}

	return files, nil
}

func getLnkName(file FilePair) string {
	if nameTemplate != "" {
		t := template.Must(template.New("lnkName").Parse(nameTemplate))
		var b strings.Builder
		err := t.Execute(&b, file)
		if err != nil {
			return fmt.Errorf("error creating lnk name: %s", err).Error()
		}
		return b.String()
	} else {
		return file.Name
	}
}

// Create a .lnk file
func createLnkFile(lnkBase string, file FilePair) (string, error) {
	if !(file.Ico && file.Rdp) {
		missing := "ico"
		if !file.Rdp {
			missing = "rdp"
		}

		return file.Name, fmt.Errorf("both ico and rdp files are required, %s is missing", missing)
	}

	lnkName := getLnkName(file)
	var lnk string
	if lnkBase == "" {
		lnk = file.Path + "\\" + lnkName + ".lnk"
	} else {
		lnk = lnkBase + "\\" + lnkName + ".lnk"
	}

	target := file.Path + "\\" + file.Name + ".rdp"
	icon := file.Path + "\\" + file.Name + ".ico"
	r := shortcut.Shortcut{
		ShortcutPath:     lnk,
		Target:           target,
		IconLocation:     icon,
		Arguments:        "",
		Description:      file.Name,
		Hotkey:           "",
		WindowStyle:      "",
		WorkingDirectory: "",
	}
	return lnk, shortcut.Create(r)
}

var path string
var lnkPath string
var nameTemplate string

func init() {
	flag.StringVar(&path, "p", "", "Path to the directory")
	flag.StringVar(&lnkPath, "l", "", "Path to the lnk being created")
	flag.StringVar(&nameTemplate, "n", "",
		"Name template for the lnk being created. "+
			"\nDefault is the same as the rdp file name, without the extension. "+
			"\nExample: -n \"{{.Name}}-shortcut\" will create a shortcut with the name of the rdp file, followed by \"-shortcut\". "+
			"\nAvailable fields: {{.Name}}, {{.Path}}, {{.Ico}}, {{.Rdp}}.")
	flag.Parse()
}
func main() {
	println("Simple RDP Shortcut Creator for RemoteAppTool")

	var workingDir string
	if path != "" {
		workingDir, _ = filepath.Abs(path)
		fmt.Println("Using path: ", workingDir)
	} else {
		workingDir, _ = os.Getwd()
		fmt.Println("Using current working directory: ", workingDir)
	}

	files, err := listFiles(workingDir)
	if err != nil {
		log.Fatal(err)
	}

	lnkBase, err := filepath.Abs(lnkPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Using lnk path: ", lnkBase)

	fmt.Printf("Files in %s:\n", workingDir)
	for _, file := range files {
		lnk, err := createLnkFile(lnkBase, file)
		if err != nil {
			fmt.Println("error creating: ", lnk, err)
		} else {
			fmt.Println("Created: ", lnk)
		}
	}
}
