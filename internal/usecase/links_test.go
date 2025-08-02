package usecase_test

import (
	"archive/zip"
	"archiveFiles/internal/entity"
	httpclient "archiveFiles/internal/repo/httpClient"
	"archiveFiles/internal/repo/mock"
	"archiveFiles/internal/usecase/links"
	builderror "archiveFiles/pkg/buildError"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func NewTestLogger(t *testing.T) *zap.Logger {
	return zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel))
}

func TestCreateTask(t *testing.T) {

	l := NewTestLogger(t)

	// repo
	httpClient := httpclient.New("image/jpeg", 51200)

	// usecase
	u := links.New(httpClient, 3, 2)
	t.Run("test1", func(t *testing.T) {

		task1, err := u.CreateTask(l)
		assert.NoError(t, err)
		assert.Equal(t, task1.ID, uint64(1))

		task2, err := u.CreateTask(l)
		assert.NoError(t, err)
		assert.Equal(t, task2.ID, uint64(2))

		_, err = u.CreateTask(l)
		assert.ErrorIs(t, err, builderror.ActiveTasksLimitError{MaxAllowed: 2})
	})
}

func TestProcessLinks(t *testing.T) {
	l := NewTestLogger(t)

	// repo
	httpClient := new(mock.MockHttpClient)

	httpClient.On("DownloadByLink", "url-1", 0, l).Return(
		&entity.DownloadResult{
			Filename: "file0.jpg",
			Content:  []byte{1, 2, 3},
			FileNum:  0,
		},
	)

	httpClient.On("DownloadByLink", "url-2", 1, l).Return(
		&entity.DownloadResult{
			Filename: "file1.jpg",
			Content:  []byte{1, 2, 3},
			FileNum:  1,
		},
	)

	// usecase
	u := links.New(httpClient, 3, 3)

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	links := []entity.Link{}
	links = append(links,
		entity.Link{
			URL: "url-1",
		},
		entity.Link{
			URL: "url-2",
		},
	)

	errMsg, files, err := u.ProcessLinks(zipWriter, links, l)
	assert.NoError(t, err)
	assert.Equal(t, len(errMsg), 0)
	assert.Equal(t, len(files), 2)
}

// test max links error
func TestErrorProcessLinks(t *testing.T) {
	l := NewTestLogger(t)

	// repo
	httpClient := new(mock.MockHttpClient)

	httpClient.On("DownloadByLink", "url-1", 0, l).Return(
		&entity.DownloadResult{
			Filename: "file0.jpg",
			Content:  []byte{1, 2, 3},
			FileNum:  0,
		},
	)

	httpClient.On("DownloadByLink", "url-2", 1, l).Return(
		&entity.DownloadResult{
			Filename: "file1.jpg",
			Content:  []byte{1, 2, 3},
			FileNum:  1,
		},
	)

	// usecase
	// max links == 1,
	maxLinks := 1
	u := links.New(httpClient, maxLinks, 3)

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	links := []entity.Link{}
	links = append(links,
		entity.Link{
			URL: "url-1",
		},
		entity.Link{
			URL: "url-2",
		},
	)

	errMsg, _, err := u.ProcessLinks(zipWriter, links, l)
	assert.ErrorIs(t, err, builderror.LinkLimitError{2, maxLinks})
	assert.Equal(t, len(errMsg), 0)

}
