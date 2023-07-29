package dataset

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tburke/iso8211"
)

type File struct {
	Path string
}

type Dataset struct {
	Id          string
	Description string
	Files       []File
}

func GetS57Datasets(path string) ([]Dataset, error) {
	datasets := make([]Dataset, 0)
	err := filepath.WalkDir(path, func(fp string, entry fs.DirEntry, err error) error {
		if entry != nil {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if strings.ToUpper(info.Name()) == "CATALOG.031" {
				dataset := Dataset{Id: "test", Description: ""}
				f, err := os.Open(fp)
				if err != nil {
					return err
				}
				var l iso8211.LeadRecord
				l.Read(f)
				var d iso8211.DataRecord
				d.Lead = &l
				for d.Read(f) == nil {
					if d.Fields[1].SubFields[5] == "BIN" {
						fileName := fmt.Sprintf("%s", d.Fields[1].SubFields[2])
						if strings.Contains(fileName, ".000") {
							file := strings.ReplaceAll(fileName, "\\", string(os.PathSeparator))
							file = filepath.Join(filepath.Dir(fp), file)
							dataset.Files = append(dataset.Files, File{Path: file})
						}

					}

				}
				datasets = append(datasets, dataset)

			}
		} else {
			return errors.New(fmt.Sprintf("Invalid path:%s", path))
		}

		return nil
	})
	if err != nil {
		return datasets, err
	}
	return datasets, nil
}
