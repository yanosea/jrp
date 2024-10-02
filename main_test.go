// Package main is the entry point of jrp.
package main

import (
	"strings"
	"testing"

	jrprepository "github.com/yanosea/jrp/app/database/jrp/repository"
	wnjpnrepository "github.com/yanosea/jrp/app/database/wnjpn/repository"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"

	"github.com/yanosea/jrp/test/testutility"
)

const (
	TEST_OUTPUT_ANY = "ANY"
)

func Test_main(t *testing.T) {
	osProxy := osproxy.New()
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osProxy,
	)
	dbFileDirPathProvider := dbfiledirpathprovider.New(
		filepathproxy.New(),
		osProxy,
		userproxy.New(),
	)
	filepathProxy := filepathproxy.New()
	wnjpnDBFileDirPath, err := dbFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetWnjpnDBFileDirPath() : error =\n%v", err)
	}
	wnjpnDBFilePath := filepathProxy.Join(wnjpnDBFileDirPath, wnjpnrepository.WNJPN_DB_FILE_NAME)
	jrpDBFileDirPath, err := dbFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		t.Errorf("DBFileDirPathProvider.GetJrpDBFileDirPath() : error =\n%v", err)
	}
	jrpDBFilePath := filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME)
	colorProxy := colorproxy.New()
	jrpRepository := jrprepository.New(
		fmtProxy,
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)

	type fields struct {
		t        *testing.T
		fnc      func()
		capturer *testutility.Capturer
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  []string
		wantStdOut            string
		wantStdErr            string
		wantJrpCount          int
		wantFavoritedJrpCount int
		wantErr               bool
	}{
		{
			name: "initial execution (jrp)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp"},
			wantStdOut:            colorProxy.YellowString(constant.ROOT_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED) + "\n",
			wantStdErr:            "",
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
		},
		{
			name: "download wnjpn database file (jrp download)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "download"},
			wantStdOut:            colorProxy.GreenString(constant.DOWNLOAD_MESSAGE_SUCCEEDED) + "\n",
			wantStdErr:            "",
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
		},
		{
			name: "generate one jrp (jrp generate)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "generate"},
			wantStdOut:            TEST_OUTPUT_ANY,
			wantStdErr:            "",
			wantJrpCount:          1,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
		},
		{
			name: "generate 10 jrps (jrp generate 10)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "generate", "10"},
			wantStdOut:            TEST_OUTPUT_ANY,
			wantStdErr:            "",
			wantJrpCount:          11,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
		},
		{
			name: "show all histories (jrp history -a)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "history", "-a"},
			wantStdOut:            TEST_OUTPUT_ANY,
			wantStdErr:            "",
			wantJrpCount:          11,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
		},
		{
			name: "favorite first history (jrp favorite add 1)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "favorite", "add", "1"},
			wantStdOut:            colorProxy.GreenString(constant.FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY) + "\n",
			wantStdErr:            "",
			wantJrpCount:          11,
			wantFavoritedJrpCount: 1,
			wantErr:               false,
		},
		{
			name: "favorite second history (jrp favorite add 2)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "favorite", "add", "2"},
			wantStdOut:            colorProxy.GreenString(constant.FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY) + "\n",
			wantStdErr:            "",
			wantJrpCount:          11,
			wantFavoritedJrpCount: 2,
			wantErr:               false,
		},
		{
			name: "show all favorites (jrp favorite -a)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "favorite", "-a"},
			wantStdOut:            TEST_OUTPUT_ANY,
			wantStdErr:            "",
			wantJrpCount:          11,
			wantFavoritedJrpCount: 2,
			wantErr:               false,
		},
		{
			name: "remove favorited jrp (jrp history remove 1)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "history", "remove", "1"},
			wantStdOut:            colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE) + "\n",
			wantStdErr:            "",
			wantJrpCount:          11,
			wantFavoritedJrpCount: 2,
			wantErr:               false,
		},
		{
			name: "remove not favorited jrp (jrp history remove 3)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "history", "remove", "3"},
			wantStdOut:            colorProxy.GreenString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY) + "\n",
			wantStdErr:            "",
			wantJrpCount:          10,
			wantFavoritedJrpCount: 2,
			wantErr:               false,
		},
		{
			name: "remove favorited jrp (jrp favorite remove 2)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "favorite", "remove", "2"},
			wantStdOut:            colorProxy.GreenString(constant.FAVORITE_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY) + "\n",
			wantStdErr:            "",
			wantJrpCount:          10,
			wantFavoritedJrpCount: 1,
			wantErr:               false,
		},
		{
			name: "clear histories (jrp favorite clear)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "history", "clear"},
			wantStdOut:            colorProxy.GreenString(constant.HISTORY_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY) + "\n",
			wantStdErr:            "",
			wantJrpCount:          1,
			wantFavoritedJrpCount: 1,
			wantErr:               false,
		},
		{
			name: "remove histories forcely (jrp history remove -f 1)",
			fields: fields{
				t: t,
				fnc: func() {
					origOsExit := osExit
					osExit = func(code int) {
						return
					}
					defer func() {
						osExit = origOsExit
					}()
					main()
				},
				capturer: capturer,
			},
			args:                  []string{"path/to/jrp", "history", "remove", "-f", "1"},
			wantStdOut:            colorProxy.GreenString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY) + "\n",
			wantStdErr:            "",
			wantJrpCount:          0,
			wantFavoritedJrpCount: 0,
			wantErr:               false,
		},
	}
	if err := osProxy.RemoveAll(wnjpnDBFilePath); err != nil {
		t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
	}
	if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
		t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osproxy.Args = tt.args
			stdout, stderr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			stdout = removeTabAndSpaceAndLf(stdout)
			stderr = removeTabAndSpaceAndLf(stderr)
			tt.wantStdOut = removeTabAndSpaceAndLf(tt.wantStdOut)
			tt.wantStdErr = removeTabAndSpaceAndLf(tt.wantStdErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if tt.wantStdOut != TEST_OUTPUT_ANY && stdout != tt.wantStdOut {
				t.Errorf("main() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if tt.wantStdErr != TEST_OUTPUT_ANY && stderr != tt.wantStdErr {
				t.Errorf("main() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
			jrps, err := jrpRepository.GetAllHistory(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllHistory() : error =\n%v", err)
			}
			if len(jrps) != tt.wantJrpCount {
				t.Errorf("JrpRepository.GetAllHistory() : len(jrps) =\n%v, wantJrpCount =\n%v", len(jrps), tt.wantJrpCount)
			}
			favoritedJrps, err := jrpRepository.GetAllFavorite(jrpDBFilePath)
			if err != nil {
				t.Errorf("JrpRepository.GetAllFavorite() : error =\n%v", err)
			}
			if len(favoritedJrps) != tt.wantFavoritedJrpCount {
				t.Errorf("JrpRepository.GetAllFavorite() : len(favoritedJrps) =\n%v, wantFavoritedJrpCount =\n%v", len(favoritedJrps), tt.wantFavoritedJrpCount)
			}
		})
	}
	if err := osProxy.RemoveAll(jrpDBFilePath); err != nil {
		t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
	}
	if err := osProxy.RemoveAll(wnjpnDBFilePath); err != nil {
		t.Errorf("OsProxy.RemoveAll() : error =\n%v", err)
	}
}

func removeTabAndSpaceAndLf(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "\t", ""), " ", ""), "\n", "")
}
