package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Terence1105/Tinder/cmd/tinder/app/dto"
	storageMock "github.com/Terence1105/Tinder/pkg/storage/mock"
	"github.com/Terence1105/Tinder/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TestHandlerSuite struct {
	suite.Suite

	ctrl        *gomock.Controller
	mockStorage *storageMock.MockTinderStorage
	router      *gin.Engine
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(TestHandlerSuite))
}

func (s *TestHandlerSuite) SetupSuite() {
	s.ctrl = gomock.NewController(s.T())
	s.mockStorage = storageMock.NewMockTinderStorage(s.ctrl)
	app := New(context.TODO(), s.mockStorage)
	s.router = app.Handler
}

func (s *TestHandlerSuite) TestAddSinglePersonAndMatch() {
	payload := &dto.AddSinglePersonAndMatchReq{
		Name:       "terence",
		Height:     180.0,
		Gender:     types.BOY,
		DateCounts: 1,
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/add-single-person-and-match", bytes.NewReader(body))

	s.mockStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(nil)

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *TestHandlerSuite) TestRemoveSinglePerson() {
	payload := &dto.RemoveSinglePersonReq{
		Name:   "terence",
		Gender: types.BOY,
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/remove-single-person", bytes.NewReader(body))

	s.mockStorage.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Return(nil)

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *TestHandlerSuite) TestQuerySinglePeople() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/query-single-people?counts=3", nil)

	s.mockStorage.EXPECT().GetPeople(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	s.mockStorage.EXPECT().DecrementDateCount(gomock.Any(), gomock.Any()).Return(0, nil)

	s.mockStorage.EXPECT().GetDateCount(gomock.Any(), gomock.Any()).Return("0", nil)

	s.mockStorage.EXPECT().RemovePerson(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

// TODO: write more test cases with QuerySinglePeople
