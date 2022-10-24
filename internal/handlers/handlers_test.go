package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/storage"
	mocksStorage "github.com/nastradamus39/ya_practicum_go_advanced/internal/storage/mocks"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	storage *mocksStorage.Mockstore
}

func (s *HandlersTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.storage = mocksStorage.NewMockstore(s.ctrl)

	storage.Storage = s.storage
}

// TestCreateShortURLHandler создание короткого url
func (s *HandlersTestSuite) TestCreateShortURLHandler() {
	s.storage.EXPECT().Save(gomock.Any()).Return(nil).Times(1)

	request := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader("http://jwlqct1udntv.com/xr0cz5fshffj/pimnbpv/otw2im3fudstqi1"),
	)
	w := httptest.NewRecorder()

	CreateShortURLHTTPHandler(w, request)

	result := w.Result()

	urlResult, err := ioutil.ReadAll(result.Body)
	assert.Equal(s.T(), "/580c5ab5ef6a4f27b3da9956ae192f4f", string(urlResult))
	assert.Equal(s.T(), 201, result.StatusCode)

	require.NoError(s.T(), err)
	err = result.Body.Close()
	require.NoError(s.T(), err)
}

// TestGetShortURLHandler возвращает полный url по короткому
func (s *HandlersTestSuite) TestGetShortURLHandler() {
	s.storage.EXPECT().FindByHash(gomock.Any()).Return(true, &types.URL{
		UUID:     "uuid",
		Hash:     "580c5ab5ef6a4f27b3da9956ae192f4f",
		URL:      "https://ya.ru?x=y",
		ShortURL: "https://localhost/580c5ab5ef6a4f27b3da9956ae192f4f",
	}, nil).AnyTimes()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("hash", "580c5ab5ef6a4f27b3da9956ae192f4f")

	request := httptest.NewRequest(http.MethodGet, "/580c5ab5ef6a4f27b3da9956ae192f4f", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	GetShortURLHTTPHandler(w, request)

	result := w.Result()

	_, err := ioutil.ReadAll(result.Body)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusTemporaryRedirect, result.StatusCode)

	err = result.Body.Close()
	require.NoError(s.T(), err)
}

// TestAPICreateShortURLHandler Api для создания короткого урла
func (s *HandlersTestSuite) TestAPICreateShortURLHandler() {
	s.storage.EXPECT().Save(gomock.Any()).Return(nil).Times(1)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		strings.NewReader(`{"correlation_id" : "as7d6as8d68as67dausghdjahsgd","original_url" : "http://yandex.ru?x=1&y=2"}`),
	)

	w := httptest.NewRecorder()

	APICreateShortURLHTTPHandler(w, request)

	result := w.Result()

	_, err := ioutil.ReadAll(result.Body)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusCreated, result.StatusCode)

	err = result.Body.Close()
	require.NoError(s.T(), err)
}

// TestAPICreateShortURLBatchHandler Api для создания короткого урла
func (s *HandlersTestSuite) TestAPICreateShortURLBatchHandler() {
	s.storage.EXPECT().SaveBatch(gomock.Any()).Return(nil).Times(1)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten/batch",
		strings.NewReader(`[{"correlation_id" : "as7d6as8d68as67dausghdjahsgd", "original_url" : "http://yandex.ru?x=1&y=2"}]`),
	)

	w := httptest.NewRecorder()

	APICreateShortURLBatchHTTPHandler(w, request)

	result := w.Result()

	_, err := ioutil.ReadAll(result.Body)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusCreated, result.StatusCode)

	err = result.Body.Close()
	require.NoError(s.T(), err)
}

// TestGetUserURLSHandler возвращает все сокращенные урлы пользователя
func (s *HandlersTestSuite) TestGetUserURLSHandler() {
	s.storage.EXPECT().FindByUUID(gomock.Any()).Return(nil, nil).Times(1)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/user/urls",
		nil,
	)

	w := httptest.NewRecorder()

	GetUserURLSHTTPHandler(w, request)

	result := w.Result()

	_, err := ioutil.ReadAll(result.Body)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNoContent, result.StatusCode)

	err = result.Body.Close()
	require.NoError(s.T(), err)
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
