package jrp

import (
	"errors"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

// downloadUseCase is a struct that contains the use case of the download.
type downloadUseCase struct{}

// NewDownloadUseCase returns a new instance of the DownloadUseCase struct.
func NewDownloadUseCase() *downloadUseCase {
	return &downloadUseCase{}
}

var (
	// Du is a variable that contains the DownloadUtil struct for injecting dependencies in testing.
	Du = utility.NewDownloadUtil(
		proxy.NewHttp(),
	)
	// Fu is a variable that contains the FileUtil struct for injecting dependencies in testing.
	Fu = utility.NewFileUtil(
		proxy.NewGzip(),
		proxy.NewIo(),
		proxy.NewOs(),
	)
)

// Run returns the output of the DownloadUseCase.
func (uc *downloadUseCase) Run(wnJpnDBPath string) error {
	var deferErr error
	if Fu.IsExist(wnJpnDBPath) {
		return errors.New("wnjpn.db already exists")
	}

	resp, err := Du.Download(
		"https://github.com/bond-lab/wnja/releases/download/v1.1/wnjpn.db.gz",
	)
	if err != nil {
		return err
	}
	defer func() {
		deferErr = resp.Close()
	}()

	tempFilePath, err := Fu.SaveToTempFile(resp.GetBody(), "wnjpn.db.gz")
	if err != nil {
		return err
	}

	defer func() {
		deferErr = Fu.RemoveAll(tempFilePath)
	}()

	if err := Fu.ExtractGzFile(tempFilePath, wnJpnDBPath); err != nil {
		return err
	}

	return deferErr
}
