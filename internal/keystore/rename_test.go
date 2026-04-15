package keystore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRename_MovesEntry(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	require.NoError(t, err)

	require.NoError(t, s.Set("old-svc", "key-abc"))
	require.NoError(t, s.Rename("old-svc", "new-svc", false))

	_, err = s.Get("old-svc")
	assert.Error(t, err, "old name should no longer exist")

	val, err := s.Get("new-svc")
	require.NoError(t, err)
	assert.Equal(t, "key-abc", val)
}

func TestRename_MissingSource(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	require.NoError(t, err)

	err = s.Rename("ghost", "new-svc", false)
	assert.ErrorContains(t, err, "not found")
}

func TestRename_DestinationExistsNoOverwrite(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	require.NoError(t, err)

	require.NoError(t, s.Set("svc-a", "key-1"))
	require.NoError(t, s.Set("svc-b", "key-2"))

	err = s.Rename("svc-a", "svc-b", false)
	assert.ErrorContains(t, err, "already exists")
}

func TestRename_OverwriteExistingDestination(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	require.NoError(t, err)

	require.NoError(t, s.Set("svc-a", "key-1"))
	require.NoError(t, s.Set("svc-b", "key-2"))

	require.NoError(t, s.Rename("svc-a", "svc-b", true))

	val, err := s.Get("svc-b")
	require.NoError(t, err)
	assert.Equal(t, "key-1", val)

	_, err = s.Get("svc-a")
	assert.Error(t, err)
}

func TestRename_TagsAreIndependent(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	require.NoError(t, err)

	require.NoError(t, s.Set("svc-a", "key-1"))
	require.NoError(t, s.SetTags("svc-a", []string{"prod", "critical"}))
	require.NoError(t, s.Rename("svc-a", "svc-b", false))

	tags, err := s.GetTags("svc-b")
	require.NoError(t, err)
	assert.Equal(t, []string{"critical", "prod"}, tags)
}

func TestRename_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)
	s, err := New(path)
	require.NoError(t, err)

	require.NoError(t, s.Set("old-svc", "key-xyz"))
	require.NoError(t, s.Rename("old-svc", "new-svc", false))

	s2, err := New(path)
	require.NoError(t, err)

	val, err := s2.Get("new-svc")
	require.NoError(t, err)
	assert.Equal(t, "key-xyz", val)

	_, err = s2.Get("old-svc")
	assert.Error(t, err)
}
