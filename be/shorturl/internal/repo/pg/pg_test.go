package pg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/repo"
)

var tsDuration = time.Duration(200 * time.Millisecond)

func connect(t *testing.T) *Repo {
	db, err := Connect(Config{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
		SSLMode:  "disable",
	})
	require.NoError(t, err)

	err = Down(db)
	require.NoError(t, err)
	err = Up(db)
	require.NoError(t, err)

	return &Repo{
		db: db,
	}

}

func TestRepo(t *testing.T) {

	testRepo := connect(t)

	testCases := []struct {
		name string
		req  repo.NewRecord
		err  error
	}{
		{
			name: "ok",
			req: repo.NewRecord{
				Owner:   "admin",
				Key:     "key1",
				LongURL: "long1",
			},
			err: nil,
		},
		{
			name: "err",
			req: repo.NewRecord{
				Owner:   "admin",
				Key:     "key1",
				LongURL: "long1",
			},
			err: fmt.Errorf("err"),
		},
		{
			name: "ok 2",
			req: repo.NewRecord{
				Owner:   "admin",
				Key:     "key2",
				LongURL: "long2",
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testRepo.Create(context.Background(), tc.req)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testRepo := connect(t)

	err := testRepo.Create(context.Background(), repo.NewRecord{
		Owner:   "owner",
		Key:     "key1",
		LongURL: "long",
	})
	require.NoError(t, err)

	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "ok",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testRepo.Delete(context.Background(), "key1")
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

		})
	}
}

func TestUpdate(t *testing.T) {
	testRepo := connect(t)

	err := testRepo.Create(context.Background(), repo.NewRecord{
		Owner:   "owner",
		Key:     "key1",
		LongURL: "long",
	})
	require.NoError(t, err)

	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "ok",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testRepo.Update(context.Background(), "key1", "long2")
			if tc.err == nil {
				require.NoError(t, err)
				rec, err := testRepo.GetByKey(context.Background(), "key1")
				require.NoError(t, err)
				require.Equal(t, "long2", rec.LongURL)
			} else {
				require.Error(t, err)
			}

		})
	}
}

func TestGetByKey(t *testing.T) {
	testRepo := connect(t)

	err := testRepo.Create(context.Background(), repo.NewRecord{
		Owner:   "owner",
		Key:     "key1",
		LongURL: "long1",
	})
	require.NoError(t, err)

	testCases := []struct {
		name   string
		key    string
		err    error
		record *repo.URLRecord
	}{
		{
			name: "ok",
			key:  "key1",
			record: &repo.URLRecord{
				Key:       "key1",
				Owner:     "owner",
				LongURL:   "long1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			err: nil,
		},
		{
			name: "empty",
			err:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := testRepo.GetByKey(context.Background(), tc.key)
			if tc.err == nil {
				require.NoError(t, err)
				if tc.record != nil {
					fmt.Printf("%+v \n", rec)
					assert.WithinDuration(t, tc.record.CreatedAt, rec.CreatedAt, tsDuration)
					assert.WithinDuration(t, tc.record.UpdatedAt, rec.UpdatedAt, tsDuration)
					rec.CreatedAt = tc.record.CreatedAt
					rec.UpdatedAt = tc.record.UpdatedAt
					require.Equal(t, tc.record, rec)
				} else {
					require.Nil(t, rec)
				}

			} else {
				require.Error(t, err)
			}

		})
	}
}

func TestList(t *testing.T) {
	testRepo := connect(t)

	err := testRepo.Create(context.Background(), repo.NewRecord{
		Owner:   "owner",
		Key:     "key1",
		LongURL: "long1",
	})
	require.NoError(t, err)

	testCases := []struct {
		name   string
		owner  string
		err    error
		record *repo.URLRecord
	}{
		{
			name:  "ok",
			owner: "owner",
			record: &repo.URLRecord{
				Key:       "key1",
				Owner:     "owner",
				LongURL:   "long1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recs, err := testRepo.ListByOwner(context.Background(), tc.owner)
			if tc.err == nil {
				require.NoError(t, err)
				if tc.record != nil {
					rec := recs[0]
					fmt.Printf("%+v \n", rec)
					assert.WithinDuration(t, tc.record.CreatedAt, rec.CreatedAt, tsDuration)
					assert.WithinDuration(t, tc.record.UpdatedAt, rec.UpdatedAt, tsDuration)
					rec.CreatedAt = tc.record.CreatedAt
					rec.UpdatedAt = tc.record.UpdatedAt
					require.Equal(t, tc.record, rec)
				}

			} else {
				require.Error(t, err)
			}

		})
	}
}
