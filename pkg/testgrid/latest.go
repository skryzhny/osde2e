package testgrid

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"cloud.google.com/go/storage"
)

const (
	latestBuildFileName = "latest-build.txt"
)

func (t *TestGrid) getLatestBuild(ctx context.Context) (int, error) {
	key := t.latestBuildKey()
	rdr, err := t.bucket.Object(key).NewReader(ctx)
	if err != nil {
		// build num is 0 if not yet set
		if err == storage.ErrObjectNotExist {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to request '%s': %v", key, err)
	}
	defer rdr.Close()

	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		return 0, fmt.Errorf("failed to tranfer '%s': %v", key, err)
	}

	buildNum, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse build number from '%s': %v", key, err)
	}
	return int(buildNum), nil
}

func (t *TestGrid) setLatestBuild(ctx context.Context, buildNum int) error {
	buildNumData := []byte(strconv.Itoa(buildNum))

	key := t.latestBuildKey()
	w := t.bucket.Object(key).NewWriter(ctx)
	if _, err := w.Write(buildNumData); err != nil {
		return fmt.Errorf("failed writing latest build in '%s': %v", key, err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed updating latest build in '%s': %v", key, err)
	}
	return nil
}

func (t *TestGrid) latestBuildKey() string {
	return filepath.Join(t.prefix, latestBuildFileName)
}
