package repo

import (
	"archive/tar"
	"io"
	"os"
	"strings"
)

type PkgInfo struct {
	Filename    string
	Name        string
	Base        string
	Version     string
	Description string
}

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
			}
		}
	}

	return pkgInfo
}

func ReadRepo(path string) ([]PkgInfo, error) {
	var packages []PkgInfo

	file, err := os.Open(path + "/aurer.db")

	if err != nil {
		return packages, err
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
