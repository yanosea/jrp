package testutility

import (
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strings"

	"github.com/yanosea/jrp/mock/app/proxy/os"
	"go.uber.org/mock/gomock"
)

func TestNewFileMover(t *testing.T) {
	filePathProxy := filepathproxy.New()
	osProxy := osproxy.New()
	stringsProxy := stringsproxy.New()

	type args struct {
		filePathProxy filepathproxy.FilePath
		osProxy       osproxy.Os
		stringsProxy  stringsproxy.Strings
	}
	tests := []struct {
		name string
		args args
		want *FileHider
	}{
		{
			name: "positive testing",
			args: args{
				filePathProxy: filePathProxy,
				osProxy:       osProxy,
				stringsProxy:  stringsProxy,
			},
			want: &FileHider{
				FilePathProxy: filePathProxy,
				OsProxy:       osProxy,
				StringsProxy:  stringsProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileMover(tt.args.filePathProxy, tt.args.osProxy, tt.args.stringsProxy); !reflect.DeepEqual(got.FilePathProxy, tt.want.FilePathProxy) ||
				!reflect.DeepEqual(got.OsProxy, tt.want.OsProxy) ||
				!reflect.DeepEqual(got.StringsProxy, tt.want.StringsProxy) {
				t.Errorf("NewFileMover() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestFileHider_HideFile(t *testing.T) {
	osProxy := osproxy.New()
	tempDir := osProxy.TempDir()
	filepathProxy := filepathproxy.New()

	type fields struct {
		FilePathProxy filepathproxy.FilePath
		OsProxy       osproxy.Os
		StringsProxy  stringsproxy.Strings
		HiddenFiles   []string
	}
	type args struct {
		filePathSlice []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (hide 1 file)",
			fields: fields{
				FilePathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				StringsProxy:  stringsproxy.New(),
				HiddenFiles:   []string{},
			},
			args: args{
				filePathSlice: []string{filepathProxy.Join(tempDir, "test.txt")},
			},
			want:    []int{0},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if _, err := osProxy.Create(filepathProxy.Join(tempDir, "test.txt")); err != nil {
					t.Errorf("Os.Create() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(filepathProxy.Join(tempDir, ".test.txt")); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "positive testing (hide 2 file)",
			fields: fields{
				FilePathProxy: filepathproxy.New(),
				OsProxy:       osproxy.New(),
				StringsProxy:  stringsproxy.New(),
				HiddenFiles:   []string{},
			},
			args: args{
				filePathSlice: []string{filepathProxy.Join(tempDir, "test1.txt"), filepathProxy.Join(tempDir, "test2.txt")},
			},
			want:    []int{0, 1},
			wantErr: false,
			setup: func(_ *gomock.Controller, _ *fields) {
				if _, err := osProxy.Create(filepathProxy.Join(tempDir, "test1.txt")); err != nil {
					t.Errorf("Os.Create() : error =\n%v", err)
				}
				if _, err := osProxy.Create(filepathProxy.Join(tempDir, "test2.txt")); err != nil {
					t.Errorf("Os.Create() : error =\n%v", err)
				}
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(filepathProxy.Join(tempDir, ".test1.txt")); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
				if err := osProxy.RemoveAll(filepathProxy.Join(tempDir, ".test2.txt")); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
		{
			name: "negative testing (OsProxy.Rename() failed)",
			fields: fields{
				FilePathProxy: filepathproxy.New(),
				OsProxy:       nil,
				StringsProxy:  stringsproxy.New(),
				HiddenFiles:   []string{},
			},
			args: args{
				filePathSlice: []string{filepathProxy.Join(tempDir, "test.txt")},
			},
			want:    []int{-1},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, fields *fields) {
				if _, err := osProxy.Create(filepathProxy.Join(tempDir, "test1.txt")); err != nil {
					t.Errorf("Os.Create() : error =\n%v", err)
				}
				mockOsProxy := mockosproxy.NewMockOs(mockCtrl)
				mockOsProxy.EXPECT().
			},
			cleanup: func() {
				if err := osProxy.RemoveAll(filepathProxy.Join(tempDir, "test1.txt")); err != nil {
					t.Errorf("Os.RemoveAll() : error =\n%v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				mockCtrl := gomock.NewController(t)
				tt.setup(mockCtrl, &tt.fields)
				defer mockCtrl.Finish()
			}
			f := &FileHider{
				FilePathProxy: tt.fields.FilePathProxy,
				OsProxy:       tt.fields.OsProxy,
				StringsProxy:  tt.fields.StringsProxy,
				HiddenFiles:   tt.fields.HiddenFiles,
			}
			for i, filePath := range tt.args.filePathSlice {
				got, err := f.HideFile(filePath)
				if (err != nil) != tt.wantErr {
					t.Errorf("FileHider.HideFile() : error =\n%v, wantErr=\n%v", err, tt.wantErr)
					return
				}
				if got != tt.want[i] {
					t.Errorf("FileHider.HideFile() : got =\n%v, want=\n%v", got, tt.want)
				}
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestFileHider_RestoreFile(t *testing.T) {
	type fields struct {
		FilePathProxy filepathproxy.FilePath
		OsProxy       osproxy.Os
		StringsProxy  stringsproxy.Strings
		HiddenFiles   []string
	}
	type args struct {
		index int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileHider{
				FilePathProxy: tt.fields.FilePathProxy,
				OsProxy:       tt.fields.OsProxy,
				StringsProxy:  tt.fields.StringsProxy,
				HiddenFiles:   tt.fields.HiddenFiles,
			}
			if err := f.RestoreFile(tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("FileHider.RestoreFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
