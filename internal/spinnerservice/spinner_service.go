package spinnerservice

import (
	"time"

	"github.com/briandowns/spinner"
)

type SpinnerService interface {
	Start()
	Stop()
	SetColor(colors ...string) error
	SetSuffix(suffix string)
}

type RealSpinnerService struct {
	sp *spinner.Spinner
}

func NewRealSpinnerService() *RealSpinnerService {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Reverse()
	return &RealSpinnerService{sp: s}
}

func (rss *RealSpinnerService) Start() {
	rss.sp.Start()
}

func (rss *RealSpinnerService) Stop() {
	rss.sp.Stop()
}

func (rss *RealSpinnerService) SetColor(colors ...string) error {
	return rss.sp.Color(colors...)
}

func (rss *RealSpinnerService) SetSuffix(suffix string) {
	rss.sp.Suffix = suffix
}
