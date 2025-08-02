package mock

import (
	"archiveFiles/internal/entity"
	"archiveFiles/internal/repo"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

var _ repo.IhttpClient = (*MockHttpClient)(nil)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) DownloadByLink(url string, num int, l *zap.Logger) *entity.DownloadResult {
	args := m.Called(url, num, l)
	return args.Get(0).(*entity.DownloadResult)
}
