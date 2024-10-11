package jrpwriter

import (
	"reflect"
	"testing"

	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/proxy/buffer"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/tablewriter"
	"github.com/yanosea/jrp/app/proxy/time"

	"github.com/yanosea/jrp/test/testutility"
)

func TestNew(t *testing.T) {
	strconvProxy := strconvproxy.New()
	tableWriterProxy := tablewriterproxy.New()

	type args struct {
		strconvProxy     strconvproxy.Strconv
		tableWriterProxy tablewriterproxy.TableWriter
	}
	tests := []struct {
		name string
		args args
		want *JrpWriter
	}{
		{
			name: "positive testing",
			args: args{
				strconvProxy:     strconvProxy,
				tableWriterProxy: tableWriterProxy,
			},
			want: &JrpWriter{
				StrconvProxy:     strconvProxy,
				TableWriterProxy: tableWriterProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.strconvProxy, tt.args.tableWriterProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestJrpWriter_WriteGenerateResultAsTable(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
	jrpWriter := New(
		strconvproxy.New(),
		tablewriterproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()

	type fields struct {
		t        *testing.T
		fnc      func()
		capturer *testutility.Capturer
	}
	tests := []struct {
		name       string
		fields     fields
		wantStdOut string
		wantStdErr string
		wantErr    bool
	}{
		{
			name: "positive testing (jrps are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					jrpWriter.WriteGenerateResultAsTable(osproxy.Stdout, nil, false)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					jrpWriter.WriteGenerateResultAsTable(osproxy.Stdout, []*model.Jrp{}, false)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are one, does not show id)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString("prefix"),
							Suffix:    sqlProxy.StringToNullString("suffix"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					}
					jrpWriter.WriteGenerateResultAsTable(osproxy.Stdout, jrps, false)
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\tPREFIX\tSUFFIX\tCREATED AT\ntest\tprefix\tsuffix\t9999-12-31 00:00:00\n\t\t\t\nTOTAL : 1\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are one, show id)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							Phrase:    "test",
							Prefix:    sqlProxy.StringToNullString("prefix"),
							Suffix:    sqlProxy.StringToNullString("suffix"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					}
					jrpWriter.WriteGenerateResultAsTable(osproxy.Stdout, jrps, true)
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tCREATED AT\n0\ttest\tprefix\tsuffix\t9999-12-31 00:00:00\n\t\t\t\nTOTAL : 1\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are two, does not show id)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							Phrase:    "test1",
							Prefix:    sqlProxy.StringToNullString("prefix1"),
							Suffix:    sqlProxy.StringToNullString("suffix1"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test2",
							Prefix:    sqlProxy.StringToNullString("prefix2"),
							Suffix:    sqlProxy.StringToNullString("suffix2"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					}
					jrpWriter.WriteGenerateResultAsTable(osproxy.Stdout, jrps, false)
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\tPREFIX\tSUFFIX\tCREATED AT\ntest1\tprefix1\tsuffix1\t9999-12-31 00:00:00\ntest2\tprefix2\tsuffix2\t9999-12-31 00:00:00\n\t\t\t\nTOTAL : 2\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are two, show id)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							Phrase:    "test1",
							Prefix:    sqlProxy.StringToNullString("prefix1"),
							Suffix:    sqlProxy.StringToNullString("suffix1"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							Phrase:    "test2",
							Prefix:    sqlProxy.StringToNullString("prefix2"),
							Suffix:    sqlProxy.StringToNullString("suffix2"),
							CreatedAt: timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					}
					jrpWriter.WriteGenerateResultAsTable(osproxy.Stdout, jrps, true)
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tCREATED AT\n0\ttest1\tprefix1\tsuffix1\t9999-12-31 00:00:00\n0\ttest2\tprefix2\tsuffix2\t9999-12-31 00:00:00\n\t\t\t\nTOTAL : 2\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			stdout = testutility.RemoveTabAndSpaceAndLf(stdout)
			stderr = testutility.RemoveTabAndSpaceAndLf(stderr)
			tt.wantStdOut = testutility.RemoveTabAndSpaceAndLf(tt.wantStdOut)
			tt.wantStdErr = testutility.RemoveTabAndSpaceAndLf(tt.wantStdErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if stdout != tt.wantStdOut {
				t.Errorf("JrpWriter.WriteGenerateResultAsTable() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("JrpWriter.WriteGenerateResultAsTable() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func TestJrpWriter_WriteAsTable(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
	jrpWriter := New(
		strconvproxy.New(),
		tablewriterproxy.New(),
	)
	sqlProxy := sqlproxy.New()
	timeProxy := timeproxy.New()

	type fields struct {
		t        *testing.T
		fnc      func()
		capturer *testutility.Capturer
	}
	tests := []struct {
		name       string
		fields     fields
		wantStdOut string
		wantStdErr string
		wantErr    bool
	}{
		{
			name: "positive testing (jrps are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					jrpWriter.WriteAsTable(osproxy.Stdout, nil)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					jrpWriter.WriteAsTable(osproxy.Stdout, []*model.Jrp{})
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are one)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							ID:          1,
							Phrase:      "test",
							Prefix:      sqlProxy.StringToNullString("prefix"),
							Suffix:      sqlProxy.StringToNullString("suffix"),
							IsFavorited: 0,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					}
					jrpWriter.WriteAsTable(osproxy.Stdout, jrps)
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest\tprefix\tsuffix\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\t\t\t\nTOTAL : 1\t\t\t\t\t\t\n",
			wantStdErr: "",
		},
		{
			name: "positive testing (jrps are two)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							ID:          1,
							Phrase:      "test1",
							Prefix:      sqlProxy.StringToNullString("prefix1"),
							Suffix:      sqlProxy.StringToNullString("suffix1"),
							IsFavorited: 0,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
						{
							ID:          1,
							Phrase:      "test2",
							Prefix:      sqlProxy.StringToNullString("prefix2"),
							Suffix:      sqlProxy.StringToNullString("suffix2"),
							IsFavorited: 1,
							CreatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
							UpdatedAt:   timeProxy.Date(9999, 12, 31, 0, 0, 0, 0, &timeproxy.UTC),
						},
					}
					jrpWriter.WriteAsTable(osproxy.Stdout, jrps)
				},
				capturer: capturer,
			},
			wantStdOut: "ID\tPHRASE\tPREFIX\tSUFFIX\tIS FAVORITED\tCREATED AT\tUPDATED AT\n1\ttest1\tprefix1\tsuffix1\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n1\ttest2\tprefix2\tsuffix2\t○\t9999-12-31 00:00:00\t9999-12-31 00:00:00\n\t\t\t\t\t\t\nTOTAL : 2\t\t\t\t\t\t\n",
			wantStdErr: "",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			stdout = testutility.RemoveTabAndSpaceAndLf(stdout)
			stderr = testutility.RemoveTabAndSpaceAndLf(stderr)
			tt.wantStdOut = testutility.RemoveTabAndSpaceAndLf(tt.wantStdOut)
			tt.wantStdErr = testutility.RemoveTabAndSpaceAndLf(tt.wantStdErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if stdout != tt.wantStdOut {
				t.Errorf("JrpWriter.WriteAsTable() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("JrpWriter.WriteAsTable() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func TestJrpWriter_writeTable(t *testing.T) {
	capturer := testutility.NewCapturer(
		bufferproxy.New(),
		bufferproxy.New(),
		osproxy.New(),
	)
	jrpWriter := New(
		strconvproxy.New(),
		tablewriterproxy.New(),
	)
	headers := []string{"phrase"}
	rowFunc := func(jrp *model.Jrp) []string {
		return []string{
			jrp.Phrase,
		}
	}

	type fields struct {
		t        *testing.T
		fnc      func()
		capturer *testutility.Capturer
	}
	tests := []struct {
		name       string
		fields     fields
		wantStdOut string
		wantStdErr string
		wantErr    bool
	}{
		{
			name: "positive testing (jrps are nil)",
			fields: fields{
				t: t,
				fnc: func() {
					jrpWriter.writeTable(osproxy.Stdout, nil, headers, rowFunc, false)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are empty)",
			fields: fields{
				t: t,
				fnc: func() {
					jrpWriter.writeTable(osproxy.Stdout, []*model.Jrp{}, headers, rowFunc, false)
				},
				capturer: capturer,
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are one, show total)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							Phrase: "test",
						},
					}
					jrpWriter.writeTable(osproxy.Stdout, jrps, headers, rowFunc, true)
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\ntest\n\t\nTOTAL : 1\n",
			wantStdErr: "",
		},
		{
			name: "positive testing (jrps are two, show total)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							Phrase: "test1",
						}, {
							Phrase: "test2",
						},
					}
					jrpWriter.writeTable(osproxy.Stdout, jrps, headers, rowFunc, true)
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\ntest1\ntest2\n\t\nTOTAL : 2\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (jrps are one, do not show total)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							Phrase: "test",
						},
					}
					jrpWriter.writeTable(osproxy.Stdout, jrps, headers, rowFunc, false)
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\ntest\n",
			wantStdErr: "",
		}, {
			name: "positive testing (jrps are two, do not show total)",
			fields: fields{
				t: t,
				fnc: func() {
					jrps := []*model.Jrp{
						{
							Phrase: "test1",
						}, {
							Phrase: "test2",
						},
					}
					jrpWriter.writeTable(osproxy.Stdout, jrps, headers, rowFunc, false)
				},
				capturer: capturer,
			},
			wantStdOut: "PHRASE\ntest1\ntest2\n",
			wantStdErr: "",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := tt.fields.capturer.CaptureOutput(
				tt.fields.t,
				tt.fields.fnc,
			)
			stdout = testutility.RemoveTabAndSpaceAndLf(stdout)
			stderr = testutility.RemoveTabAndSpaceAndLf(stderr)
			tt.wantStdOut = testutility.RemoveTabAndSpaceAndLf(tt.wantStdOut)
			tt.wantStdErr = testutility.RemoveTabAndSpaceAndLf(tt.wantStdErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
			}
			if stdout != tt.wantStdOut {
				t.Errorf("JrpWriter.writeTable() : stdout =\n%v, wantStdOut =\n%v", stdout, tt.wantStdOut)
			}
			if stderr != tt.wantStdErr {
				t.Errorf("JrpWriter.writeTable() : stderr =\n%v, wantStdErr =\n%v", stderr, tt.wantStdErr)
			}
		})
	}
}

func TestJrpWriter_getDefaultTableWriter(t *testing.T) {
	tablewriterProxy := tablewriterproxy.New()
	table := tablewriterProxy.NewTable(osproxy.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriterproxy.ALIGN_LEFT)
	table.SetAlignment(tablewriterproxy.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	type fields struct {
		StrconvProxy     strconvproxy.Strconv
		TableWriterProxy tablewriterproxy.TableWriter
	}
	type args struct {
		o ioproxy.WriterInstanceInterface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   tablewriterproxy.TableInstanceInterface
	}{
		{
			name: "positive testing",
			fields: fields{
				StrconvProxy:     strconvproxy.New(),
				TableWriterProxy: tablewriterproxy.New(),
			},
			args: args{
				o: osproxy.Stdout,
			},
			want: table,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := New(
				tt.fields.StrconvProxy,
				tt.fields.TableWriterProxy,
			)
			if got := j.getDefaultTableWriter(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JrpWriter.getDefaultTableWriter() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}
