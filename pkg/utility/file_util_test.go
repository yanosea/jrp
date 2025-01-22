package utility

import (
	"errors"
	"io"
	"io/fs"
	o "os"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"
	"go.uber.org/mock/gomock"
)

func TestNewFileUtil(t *testing.T) {
	gzip := proxy.NewGzip()
	io := proxy.NewIo()
	os := proxy.NewOs()

	type args struct {
		gzip proxy.Gzip
		io   proxy.Io
		os   proxy.Os
	}
	tests := []struct {
		name string
		args args
		want FileUtil
	}{
		{
			name: "positive testing",
			args: args{
				gzip: gzip,
				io:   io,
				os:   os,
			},
			want: &fileUtil{
				gzip: gzip,
				io:   io,
				os:   os,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileUtil(tt.args.gzip, tt.args.io, tt.args.os); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileUtil_ExtractGzFile(t *testing.T) {
	type fields struct {
		Gzip proxy.Gzip
		Io   proxy.Io
		Os   proxy.Os
	}
	type args struct {
		gzFilePath   string
		destFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				Gzip: nil,
				Io:   nil,
				Os:   nil,
			},
			args: args{
				gzFilePath:   "test.gz",
				destFilePath: "test",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockFile := proxy.NewMockFile(mockCtrl)
				mockFile.EXPECT().Close().Return(nil).AnyTimes()
				mockOs.EXPECT().Open("test.gz").Return(mockFile, nil)
				mockOs.EXPECT().Create("test").Return(mockFile, nil)
				mockGzip := proxy.NewMockGzip(mockCtrl)
				mockGzipReader := proxy.NewMockGzipReader(mockCtrl)
				mockGzipReader.EXPECT().Close().Return(nil)
				mockGzip.EXPECT().NewReader(mockFile).Return(mockGzipReader, nil)
				mockIo := proxy.NewMockIo(mockCtrl)
				mockIo.EXPECT().Copy(mockFile, mockGzipReader).Return(int64(0), nil)
				tt.Gzip = mockGzip
				tt.Io = mockIo
				tt.Os = mockOs
			},
		},
		{
			name: "negative testing (f.Os.Open(gzFilePath) failed)",
			fields: fields{
				Gzip: nil,
				Io:   nil,
				Os:   nil,
			},
			args: args{
				gzFilePath:   "test.gz",
				destFilePath: "test",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Open("test.gz").Return(nil, errors.New("OsProxy.Open() failed"))
				mockGzip := proxy.NewMockGzip(mockCtrl)
				mockIo := proxy.NewMockIo(mockCtrl)
				tt.Gzip = mockGzip
				tt.Io = mockIo
				tt.Os = mockOs
			},
		},
		{
			name: "negative testing (f.Gzip.NewReader(gzFile) failed)",
			fields: fields{
				Gzip: nil,
				Io:   nil,
				Os:   nil,
			},
			args: args{
				gzFilePath:   "test.gz",
				destFilePath: "test",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockFile := proxy.NewMockFile(mockCtrl)
				mockFile.EXPECT().Close().Return(nil).AnyTimes()
				mockOs.EXPECT().Open("test.gz").Return(mockFile, nil)
				mockGzip := proxy.NewMockGzip(mockCtrl)
				mockGzip.EXPECT().NewReader(mockFile).Return(nil, errors.New("GzipProxy.NewReader() failed"))
				mockIo := proxy.NewMockIo(mockCtrl)
				tt.Gzip = mockGzip
				tt.Io = mockIo
				tt.Os = mockOs
			},
		},
		{
			name: "negative testing (f.Os.Create(destFilePath) failed)",
			fields: fields{
				Gzip: nil,
				Io:   nil,
				Os:   nil,
			},
			args: args{
				gzFilePath:   "test.gz",
				destFilePath: "test",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockFile := proxy.NewMockFile(mockCtrl)
				mockFile.EXPECT().Close().Return(nil).AnyTimes()
				mockOs.EXPECT().Open("test.gz").Return(mockFile, nil)
				mockOs.EXPECT().Create("test").Return(nil, errors.New("OsProxy.Create() failed"))
				mockGzip := proxy.NewMockGzip(mockCtrl)
				mockGzipReader := proxy.NewMockGzipReader(mockCtrl)
				mockGzipReader.EXPECT().Close().Return(nil)
				mockGzip.EXPECT().NewReader(mockFile).Return(mockGzipReader, nil)
				mockIo := proxy.NewMockIo(mockCtrl)
				tt.Gzip = mockGzip
				tt.Io = mockIo
				tt.Os = mockOs
			},
		},
		{
			name: "negative testing (f.Io.Copy(destFile, gzReader) failed)",
			fields: fields{
				Gzip: nil,
				Io:   nil,
				Os:   nil,
			},
			args: args{
				gzFilePath:   "test.gz",
				destFilePath: "test",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockFile := proxy.NewMockFile(mockCtrl)
				mockFile.EXPECT().Close().Return(nil).AnyTimes()
				mockOs.EXPECT().Open("test.gz").Return(mockFile, nil)
				mockOs.EXPECT().Create("test").Return(mockFile, nil)
				mockGzip := proxy.NewMockGzip(mockCtrl)
				mockGzipReader := proxy.NewMockGzipReader(mockCtrl)
				mockGzipReader.EXPECT().Close().Return(nil)
				mockGzip.EXPECT().NewReader(mockFile).Return(mockGzipReader, nil)
				mockIo := proxy.NewMockIo(mockCtrl)
				mockIo.EXPECT().Copy(mockFile, mockGzipReader).Return(int64(0), errors.New("IoProxy.Copy() failed"))
				tt.Gzip = mockGzip
				tt.Io = mockIo
				tt.Os = mockOs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			f := &fileUtil{
				gzip: tt.fields.Gzip,
				io:   tt.fields.Io,
				os:   tt.fields.Os,
			}
			if err := f.ExtractGzFile(tt.args.gzFilePath, tt.args.destFilePath); (err != nil) != tt.wantErr {
				t.Errorf("fileUtil.ExtractGzFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileUtil_GetXDGDataHome(t *testing.T) {
	gzip := proxy.NewGzip()
	io := proxy.NewIo()
	os := proxy.NewOs()

	type fields struct {
		Gzip proxy.Gzip
		Io   proxy.Io
		Os   proxy.Os
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (XDG_DATA_HOME is set)",
			fields: fields{
				Gzip: gzip,
				Io:   io,
				Os:   os,
			},
			want:    "/home/user/.local/share",
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				if err := o.Setenv("XDG_DATA_HOME", "/home/user/.local/share"); err != nil {
					t.Errorf("os.Setenv() error = %v", err)
				}
			},
		},
		{
			name: "positive testing (XDG_DATA_HOME is not set)",
			fields: fields{
				Gzip: gzip,
				Io:   io,
				Os:   nil,
			},
			want:    "/home/user/.local/share",
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Getenv("XDG_DATA_HOME").Return("")
				mockOs.EXPECT().UserHomeDir().Return("/home/user", nil)
				tt.Os = mockOs
			},
		},
		{
			name: "negative testing (f.Os.UserHomeDir() failed)",
			fields: fields{
				Gzip: gzip,
				Io:   io,
				Os:   nil,
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Getenv("XDG_DATA_HOME").Return("")
				mockOs.EXPECT().UserHomeDir().Return("", errors.New("OsProxy.UserHomeDir() failed"))
				tt.Os = mockOs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			f := &fileUtil{
				gzip: tt.fields.Gzip,
				io:   tt.fields.Io,
				os:   tt.fields.Os,
			}
			got, err := f.GetXDGDataHome()
			if (err != nil) != tt.wantErr {
				t.Errorf("fileUtil.GetXDGDataHome() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileUtil.GetXDGDataHome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileUtil_HideFile(t *testing.T) {
	gzip := proxy.NewGzip()
	io := proxy.NewIo()

	type fields struct {
		gzip proxy.Gzip
		io   proxy.Io
		os   proxy.Os
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				gzip: gzip,
				io:   io,
				os:   nil,
			},
			args: args{
				filePath: "dir/test.txt",
			},
			want:    "dir/.test.txt",
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Rename("dir/test.txt", "dir/.test.txt").Return(nil)
				tt.os = mockOs
			},
		},
		{
			name: "negative testing (f.Os.Rename() failed)",
			fields: fields{
				gzip: gzip,
				io:   io,
				os:   nil,
			},
			args: args{
				filePath: "dir/test.txt",
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Rename("dir/test.txt", "dir/.test.txt").Return(errors.New("OsProxy.Rename() failed"))
				tt.os = mockOs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			f := &fileUtil{
				gzip: tt.fields.gzip,
				io:   tt.fields.io,
				os:   tt.fields.os,
			}
			got, err := f.HideFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileUtil.HideFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileUtil.HideFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileUtil_IsExist(t *testing.T) {
	gzip := proxy.NewGzip()
	io := proxy.NewIo()

	type fields struct {
		Gzip proxy.Gzip
		Io   proxy.Io
		Os   proxy.Os
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		setup  func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				Gzip: gzip,
				Io:   io,
				Os:   nil,
			},
			args: args{
				name: "test",
			},
			want: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Stat("test").Return(nil, fs.ErrNotExist)
				mockOs.EXPECT().IsNotExist(fs.ErrNotExist).Return(true)
				tt.Os = mockOs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			f := &fileUtil{
				gzip: tt.fields.Gzip,
				io:   tt.fields.Io,
				os:   tt.fields.Os,
			}
			if got := f.IsExist(tt.args.name); got != tt.want {
				t.Errorf("fileUtil.IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileUtil_MkdirIfNotExist(t *testing.T) {
	gzip := proxy.NewGzip()
	io := proxy.NewIo()

	type fields struct {
		Gzip proxy.Gzip
		Io   proxy.Io
		Os   proxy.Os
	}
	type args struct {
		dirPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				Gzip: gzip,
				Io:   io,
				Os:   nil,
			},
			args: args{
				dirPath: "test",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Stat("test").Return(nil, fs.ErrNotExist)
				mockOs.EXPECT().IsNotExist(fs.ErrNotExist).Return(true)
				mockOs.EXPECT().MkdirAll("test", o.FileMode(0755)).Return(nil)
				tt.Os = mockOs
			},
		},
		{
			name: "negative testing (f.Os.MkdirAll(dirPath, 0755) failed)",
			fields: fields{
				Gzip: gzip,
				Io:   io,
				Os:   nil,
			},
			args: args{
				dirPath: "test",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Stat("test").Return(nil, fs.ErrNotExist)
				mockOs.EXPECT().IsNotExist(fs.ErrNotExist).Return(true)
				mockOs.EXPECT().MkdirAll("test", o.FileMode(0755)).Return(errors.New("OsProxy.MkdirAll() failed"))
				tt.Os = mockOs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			f := &fileUtil{
				gzip: tt.fields.Gzip,
				io:   tt.fields.Io,
				os:   tt.fields.Os,
			}
			if err := f.MkdirIfNotExist(tt.args.dirPath); (err != nil) != tt.wantErr {
				t.Errorf("fileUtil.MkdirIfNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileUtil_RemoveAll(t *testing.T) {
	gzip := proxy.NewGzip()
	io := proxy.NewIo()

	type fields struct {
		Gzip proxy.Gzip
		Io   proxy.Io
		Os   proxy.Os
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				Gzip: gzip,
				Io:   io,
				Os:   nil,
			},
			args: args{
				path: "test",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().RemoveAll("test").Return(nil)
				tt.Os = mockOs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			f := &fileUtil{
				gzip: tt.fields.Gzip,
				io:   tt.fields.Io,
				os:   tt.fields.Os,
			}
			if err := f.RemoveAll(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("fileUtil.RemoveAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fileUtil_SaveToTempFile(t *testing.T) {
	type fields struct {
		Gzip proxy.Gzip
		Io   proxy.Io
		Os   proxy.Os
	}
	type args struct {
		body     io.Reader
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				Gzip: nil,
				Io:   nil,
				Os:   nil,
			},
			args: args{
				body:     nil,
				fileName: "test",
			},
			want:    "/tmp/test",
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockGzip := proxy.NewMockGzip(mockCtrl)
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().TempDir().Return("/tmp")
				mockFile := proxy.NewMockFile(mockCtrl)
				mockFile.EXPECT().Close().Return(nil)
				mockOs.EXPECT().Create("/tmp/test").Return(mockFile, nil)
				mockIo := proxy.NewMockIo(mockCtrl)
				mockIo.EXPECT().Copy(mockFile, nil).Return(int64(0), nil)
				tt.Gzip = mockGzip
				tt.Io = mockIo
				tt.Os = mockOs
			},
		},
		{
			name: "negative testing (f.Os.Create(\"/tmp/test\") failed)",
			fields: fields{
				Gzip: nil,
				Io:   nil,
				Os:   nil,
			},
			args: args{
				body:     nil,
				fileName: "test",
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockGzip := proxy.NewMockGzip(mockCtrl)
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().TempDir().Return("/tmp")
				mockOs.EXPECT().Create("/tmp/test").Return(nil, errors.New("OsProxy.Create() failed"))
				mockIo := proxy.NewMockIo(mockCtrl)
				tt.Gzip = mockGzip
				tt.Io = mockIo
				tt.Os = mockOs
			},
		},
		{
			name: "negative testing (f.Io.Copy(mockFile, nil) failed)",
			fields: fields{
				Gzip: nil,
				Io:   nil,
				Os:   nil,
			},
			args: args{
				body:     nil,
				fileName: "test",
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockGzip := proxy.NewMockGzip(mockCtrl)
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().TempDir().Return("/tmp")
				mockFile := proxy.NewMockFile(mockCtrl)
				mockFile.EXPECT().Close().Return(nil)
				mockOs.EXPECT().Create("/tmp/test").Return(mockFile, nil)
				mockIo := proxy.NewMockIo(mockCtrl)
				mockIo.EXPECT().Copy(mockFile, nil).Return(int64(0), errors.New("IoProxy.Copy() failed"))
				tt.Gzip = mockGzip
				tt.Io = mockIo
				tt.Os = mockOs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			f := &fileUtil{
				gzip: tt.fields.Gzip,
				io:   tt.fields.Io,
				os:   tt.fields.Os,
			}
			got, err := f.SaveToTempFile(tt.args.body, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileUtil.SaveToTempFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileUtil.SaveToTempFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileUtil_UnhideFile(t *testing.T) {
	gzip := proxy.NewGzip()
	io := proxy.NewIo()

	type fields struct {
		gzip proxy.Gzip
		io   proxy.Io
		os   proxy.Os
	}
	type args struct {
		hiddenFilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				gzip: gzip,
				io:   io,
				os:   nil,
			},
			args: args{
				hiddenFilePath: "dir/.test.txt",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Rename("dir/.test.txt", "dir/test.txt").Return(nil)
				tt.os = mockOs
			},
		},
		{
			name: "negative testing (f.Os.Rename() failed)",
			fields: fields{
				gzip: gzip,
				io:   io,
				os:   nil,
			},
			args: args{
				hiddenFilePath: "dir/.test.txt",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockOs := proxy.NewMockOs(mockCtrl)
				mockOs.EXPECT().Rename("dir/.test.txt", "dir/test.txt").Return(errors.New("OsProxy.Rename() failed"))
				tt.os = mockOs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			f := &fileUtil{
				gzip: tt.fields.gzip,
				io:   tt.fields.io,
				os:   tt.fields.os,
			}
			if err := f.UnhideFile(tt.args.hiddenFilePath); (err != nil) != tt.wantErr {
				t.Errorf("fileUtil.UnhideFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
