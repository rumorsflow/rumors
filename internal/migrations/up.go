package migrations

import (
	"fmt"
	"github.com/iagapie/rumors/pkg/litedb"
)

func Up(db *litedb.DB) error {
	files, err := things.ReadDir("schema")
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		data, err := things.ReadFile(fmt.Sprintf("schema/%s", file.Name()))
		if err != nil {
			return err
		}
		if _, err := db.NewQuery(string(data)).Execute(); err != nil {
			return err
		}
	}

	return nil
}
