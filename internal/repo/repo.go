package repo

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type PkgInfo struct {
	Filename    string
	Name        string
	Base        string
	Version     string
	Description string
	BuildDate   int
}

const REPO_DB_FILE = "aurer.db.tar"

func parsePkgInfo(info string) PkgInfo {
	lines := strings.Split(info, "\n")

	pkgInfo := PkgInfo{}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if len(line) == 0 {
			continue
		}

		if line[0] == '%' && line[len(line)-1] == '%' {
			section := line[1 : len(line)-1]

			i++
			line = lines[i]

			value := ""

			for len(line) != 0 {
				value = value + line

				i++

				line = lines[i]
			}

			switch section {
			case "FILENAME":
				pkgInfo.Filename = value
			case "NAME":
				pkgInfo.Name = value
			case "BASE":
				pkgInfo.Base = value
			case "VERSION":
				pkgInfo.Version = value
			case "DESC":
				pkgInfo.Description = value
			case "BUILDDATE":
				pkgInfo.BuildDate, _ = strconv.Atoi(value)
			}
		}
	}

	return pkgInfo
}

func ReadRepo(path string) ([]PkgInfo, error) {
	packages := make([]PkgInfo, 0)

	repo_path := path + "/" + REPO_DB_FILE

	file, err := os.Open(repo_path)

	if err != nil {
		out, err := exec.Command("repo-add", repo_path).CombinedOutput()

		log.Println(string(out))

		if err != nil {
			return packages, err
		}

	}

	tr := tar.NewReader(file)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return packages, err
		}

		info := header.FileInfo()

		if info.IsDir() {
			continue
		}

		b, err := io.ReadAll(tr)

		if err != nil {
			return packages, err
		}

		pkgInfo := parsePkgInfo(string(b))

		packages = append(packages, pkgInfo)
	}

	return packages, nil
}

func DeletePackage(repoPath string, pkgName string) error {
	log.Printf("Removing package %v\n", pkgName)

	repo, err := ReadRepo(repoPath)

	if err != nil {
		return err
	}

	for _, pkg := range repo {
		if pkg.Name == pkgName {
			cmd := exec.Command("repo-remove", repoPath+"/"+REPO_DB_FILE, pkgName)

			output, err := cmd.CombinedOutput()

			log.Println(string(output))

			if err != nil {
				return err
			}

			err = os.Remove(repoPath + "/" + pkg.Filename)

			if err != nil {
				return err
			}

			log.Printf("Removed package %v", pkgName)

			return nil
		}
	}

	return nil
}
