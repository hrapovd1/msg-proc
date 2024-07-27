package dbstorage

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDBStorage(t *testing.T) {
	ds, err := NewDBStorage("", log.Default())
	require.NoError(t, err)
	assert.IsType(t, &DBStorage{}, ds)
}
