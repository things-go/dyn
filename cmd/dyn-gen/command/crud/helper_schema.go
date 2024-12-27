package crud

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"ariga.io/atlas/sql/schema"
	"github.com/things-go/ens"
	"github.com/things-go/ens/driver"
)

func LoadDriver(URL string) (driver.Driver, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	d, err := driver.LoadDriver(u.Scheme)
	if err != nil {
		return nil, err
	}
	return d, nil
}

type source struct {
	// sql file
	InputFile []string
	Schema    string
	// database url
	URL     string
	Tables  []string
	Exclude []string
}

func getSchema(c *source) (*ens.Schema, error) {
	if c.URL != "" {
		u, err := url.Parse(c.URL)
		if err != nil {
			return nil, err
		}
		d, err := driver.LoadDriver(u.Scheme)
		if err != nil {
			return nil, err
		}
		return d.InspectSchema(context.Background(), &driver.InspectOption{
			URL: c.URL,
			InspectOptions: schema.InspectOptions{
				Mode:    schema.InspectTables,
				Tables:  c.Tables,
				Exclude: c.Exclude,
			},
		})
	}
	if len(c.InputFile) > 0 {
		d, err := driver.LoadDriver(c.Schema)
		if err != nil {
			return nil, err
		}
		inputFile, err := filepath.Glob(c.InputFile[0])
		if err != nil {
			return nil, err
		}
		if len(inputFile) == 1 && inputFile[0] == c.InputFile[0] {
			inputFile = c.InputFile
		}
		entities := make([]*ens.EntityDescriptor, 0, 256)
		for _, filename := range inputFile {
			tmpSc, err := func() (*ens.Schema, error) {
				content, err := os.ReadFile(filename)
				if err != nil {
					return nil, err
				}
				return d.InspectSchema(context.Background(), &driver.InspectOption{
					URL:            "",
					Data:           string(content),
					InspectOptions: schema.InspectOptions{},
				})
			}()
			if err != nil {
				slog.Warn("🧐 parse failed !!!", slog.String("file", filename), slog.Any("error", err))
				continue
			}
			entities = append(entities, tmpSc.Entities...)
		}
		return &ens.Schema{
			Name:     "",
			Entities: entities,
		}, nil
	}

	return nil, errors.New("at least one of [url input] is required")
}
