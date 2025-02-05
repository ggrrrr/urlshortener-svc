package app

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/urlshortener-svc/be/common/application"
	"github.com/ggrrrr/urlshortener-svc/be/common/roles"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/keygenerator"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/repo"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/models"
)

func TestGetURL(t *testing.T) {

	mockRepo := new(repo.MockRepo)

	testApp := &Application{
		repo: mockRepo,
	}

	tests := []struct {
		name   string
		key    string
		prepFn func(t *testing.T)
		err    error
		resp   string
	}{
		{
			name: "ok",
			key:  "key1",
			prepFn: func(t *testing.T) {
				mockRepo.On("GetByKey", "key1").Return(&repo.URLRecord{LongURL: "long_url1"}, nil).Once()
			},
			err:  nil,
			resp: "long_url1",
		},
		{
			name: "not found",
			key:  "not_found",
			prepFn: func(t *testing.T) {
				mockRepo.On("GetByKey", "not_found").Return(nil, application.NewNotFound()).Once()
			},
			err:  application.NewNotFound(),
			resp: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepFn(t)
			resp, err := testApp.GetLongURL(context.Background(), tc.key)
			if tc.err == nil {
				require.NoError(t, err)
				assert.Equal(t, tc.resp, resp)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {

	mockRepo := new(repo.MockRepo)

	adminCtx := roles.InjectUser(context.Background(), roles.AuthenticatedUser{Username: "admin"})

	testApp := &Application{
		repo: mockRepo,
	}

	tests := []struct {
		name   string
		ctx    context.Context
		req    models.DeleteShortURL
		prepFn func(t *testing.T)
		err    error
	}{
		{
			name: "ok",
			ctx:  adminCtx,
			req: models.DeleteShortURL{
				Key: "key1",
			},
			prepFn: func(t *testing.T) {
				mockRepo.On("GetByKey", "key1").Return(&repo.URLRecord{LongURL: "long_url1", Owner: "admin"}, nil).Once()
				mockRepo.On("Delete", "key1").Return(nil).Once()
			},
			err: nil,
		},
		{
			name: "ok 403 owner not match",
			ctx:  adminCtx,
			req: models.DeleteShortURL{
				Key: "key403",
			},
			prepFn: func(t *testing.T) {
				mockRepo.On("GetByKey", "key403").Return(&repo.URLRecord{LongURL: "long_url1", Owner: "notadmin"}, nil).Once()
			},
			err: application.NewForbidden("not allowed"),
		},
		{
			name: "ok 401 owner not match",
			ctx:  context.Background(),
			req: models.DeleteShortURL{
				Key: "key401",
			},
			prepFn: func(t *testing.T) {
			},
			err: application.NewUnauthorized(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepFn(t)
			err := testApp.Delete(tc.ctx, tc.req)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {

	mockRepo := new(repo.MockRepo)

	adminCtx := roles.InjectUser(context.Background(), roles.AuthenticatedUser{Username: "admin"})

	testApp := &Application{
		repo: mockRepo,
	}

	tests := []struct {
		name   string
		ctx    context.Context
		req    models.UpdateShortURL
		prepFn func(t *testing.T)
		err    error
	}{
		{
			name: "ok",
			ctx:  adminCtx,
			req: models.UpdateShortURL{
				Key:     "key1",
				LongURL: "long1",
			},
			prepFn: func(t *testing.T) {
				mockRepo.On("GetByKey", "key1").Return(&repo.URLRecord{LongURL: "long_url1", Owner: "admin"}, nil)
				mockRepo.On("Update", "key1", "long1").Return(nil)
			},
			err: nil,
		},
		{
			name: "ok 403 owner not match",
			ctx:  adminCtx,
			req: models.UpdateShortURL{
				Key:     "key403",
				LongURL: "long1",
			},
			prepFn: func(t *testing.T) {
				mockRepo.On("GetByKey", "key403").Return(&repo.URLRecord{LongURL: "long_url1", Owner: "notadmin"}, nil)
			},
			err: application.NewForbidden("not allowed"),
		},
		{
			name: "ok 401 owner not match",
			ctx:  context.Background(),
			req: models.UpdateShortURL{
				Key:     "key401",
				LongURL: "long1",
			},
			prepFn: func(t *testing.T) {
			},
			err: application.NewUnauthorized(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepFn(t)
			err := testApp.Update(tc.ctx, tc.req)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestCreate(t *testing.T) {

	mockRepo := new(repo.MockRepo)
	mockGenerator := new(keygenerator.MockGenerator)

	adminCtx := roles.InjectUser(context.Background(), roles.AuthenticatedUser{Username: "admin"})

	testApp := &Application{
		repo:      mockRepo,
		generator: mockGenerator,
	}

	tests := []struct {
		name   string
		ctx    context.Context
		req    models.CreateShortURL
		prepFn func(t *testing.T)
		err    error
		key    *models.Key
	}{
		{
			name: "ok",
			ctx:  adminCtx,
			req: models.CreateShortURL{
				LongURL: "long1",
			},
			prepFn: func(t *testing.T) {
				mockGenerator.On("Generate").Return("key1").Once()
				mockRepo.On("Create", repo.NewRecord{Owner: "admin", Key: "key1", LongURL: "long1"}).Return(nil).Once()
			},
			err: nil,
			key: &models.Key{Key: "key1"},
		},
		{
			name: "ok 401",
			ctx:  context.Background(),
			req: models.CreateShortURL{
				LongURL: "long1",
			},
			prepFn: func(t *testing.T) {
			},
			err: application.NewUnauthorized(),
			key: nil,
		},
		{
			name: "ok 500",
			ctx:  context.Background(),
			req: models.CreateShortURL{
				LongURL: "long1",
			},
			prepFn: func(t *testing.T) {
				mockGenerator.On("Generate").Return("key1")
				mockRepo.On("Create", repo.NewRecord{Owner: "admin", Key: "key1", LongURL: "long1"}).Return(fmt.Errorf("some error")).Once()
			},
			err: application.NewSystemError("", nil),
			key: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepFn(t)
			key, err := testApp.Create(tc.ctx, tc.req)
			if tc.err == nil {
				require.NoError(t, err)
				require.Equal(t, tc.key, key)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestListForOwner(t *testing.T) {

	testRepo := new(repo.MockRepo)
	ts := time.Now()
	adminCtx := roles.InjectUser(context.Background(), roles.AuthenticatedUser{Username: "admin"})

	testApp := &Application{
		repo: testRepo,
	}

	tests := []struct {
		name   string
		ctx    context.Context
		prepFn func(t *testing.T)
		err    error
		list   []*models.ShortURLRecord
	}{
		{
			name: "ok",
			ctx:  adminCtx,
			prepFn: func(t *testing.T) {
				testRepo.On("ListByOwner", "admin").Return([]*repo.URLRecord{
					{
						Key:       "key1",
						Owner:     "admin",
						LongURL:   "long1",
						CreatedAt: ts,
						UpdatedAt: ts,
					},
				}, nil)
			},
			err: nil,
			list: []*models.ShortURLRecord{
				{
					Key:       "key1",
					LongURL:   "long1",
					CreatedAt: ts,
					UpdatedAt: ts,
				},
			},
		},
		{
			name: "err 401",
			ctx:  context.Background(),
			prepFn: func(t *testing.T) {
			},
			err: application.NewUnauthorized(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepFn(t)
			list, err := testApp.ListForOwner(tc.ctx)
			if tc.err == nil {
				require.NoError(t, err)
				assert.Equal(t, tc.list, list)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestRetryLoop(t *testing.T) {
	mockRepo := new(repo.MockRepo)
	mockGenerator := new(keygenerator.MockGenerator)
	testApp := Application{
		repo:      mockRepo,
		generator: mockGenerator,
	}

	mockGenerator.On("Generate").Return("key1").Once()
	mockGenerator.On("Generate").Return("key2").Once()
	mockRepo.On("Create", repo.NewRecord{Owner: "admin", LongURL: "long_url1", Key: "key1"}).Return(fmt.Errorf("some error")).Once()
	mockRepo.On("Create", repo.NewRecord{Owner: "admin", LongURL: "long_url1", Key: "key2"}).Return(nil).Once()

	key, err := testApp.retryLoop(context.Background(), repo.NewRecord{Owner: "admin", LongURL: "long_url1"}, 10)

	require.NoError(t, err)
	require.Equal(t, "key2", key)
}
