package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	db "mamikost/db/sqlc"
	"mamikost/middleware"
	mockdb "mamikost/services/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestGetListCategory(t *testing.T) {

	token, _ := middleware.GenerateJWT("ivan@gmail.com")

	//fake data
	categories := []*db.Category{
		{
			CateID:   1,
			CateName: "Premium",
		},
		{
			CateID:   2,
			CateName: "Reguler",
		},
	}

	rowExpected := 2

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllCategories(gomock.Any()).
					Times(1).
					Return(categories, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, len(categories), rowExpected)
			},
		},
		{
			"InternalError",
			gin.H{},
			func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllCategories(gomock.Any()).
					Times(1).
					Return(nil, fmt.Errorf("internal error"))
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestHttpServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := "/api/category/"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(data))
			require.NoError(t, err)
			request.Header.Set("Authorization", fmt.Sprint("Bearer "+token))

			server.Router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestCreateCategory(t *testing.T) {
	username := "ivan"
	token, _ := middleware.GenerateJWT(username)
	//fake data
	cate := &db.Category{
		CateID:   0,
		CateName: "Premium",
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"category_name": cate.CateName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(cate, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "BadRequestBody",
			body: gin.H{
				"category": cate.CateName,
			},
			buildStubs: func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "DuplicateCategoryname",
			body: gin.H{
				"category_name": cate.CateName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, &pgconn.PgError{ConstraintName: "category_name_uq"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"category_name": cate.CateName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, fmt.Errorf("internal error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestHttpServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := "/api/category/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

			server.Router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestFindCategory(t *testing.T) {
	username := "ivan"
	token, _ := middleware.GenerateJWT(username)
	cate := &db.Category{
		CateID:   1,
		CateName: "Premium",
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategoryByID(gomock.Any(), cate.CateID).
					Times(1).
					Return(cate, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			"CategoryNotFound",
			gin.H{},
			func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategoryByID(gomock.Any(), cate.CateID).
					Times(1).
					Return(nil, nil)
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			"InternalError",
			gin.H{},
			func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategoryByID(gomock.Any(), cate.CateID).
					Times(1).
					Return(nil, fmt.Errorf("internal error"))
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestHttpServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := "/api/category/" + strconv.Itoa(int(cate.CateID))
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(data))
			require.NoError(t, err)
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

			server.Router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	username := "kang"
	token, _ := middleware.GenerateJWT(username)
	category := &db.Category{
		CateID:   1,
		CateName: "PremiumA",
	}

	args := db.UpdateCategoryParams{
		CateID:   category.CateID,
		CateName: category.CateName,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"category_id":   category.CateID,
				"category_name": category.CateName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCategory(gomock.Any(), args).
					Times(1).
					Return(category, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			"CategoryNotFound",
			gin.H{},
			func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil)
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			"InternalError",
			gin.H{},
			func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCategory(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, fmt.Errorf("internal error"))
			},
			func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestHttpServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := "/api/category/" + strconv.Itoa(int(category.CateID))
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
			require.NoError(t, err)
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

			server.Router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
