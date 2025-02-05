package service

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockDirEntry struct {
	name string
}

func (m MockDirEntry) Name() string {
	return m.name
}

func (m MockDirEntry) IsDir() bool { return false }

func (m MockDirEntry) Type() os.FileMode {
	return 0
}

func (m MockDirEntry) Info() (os.FileInfo, error) { return nil, nil }

func TestOrderListHandler(t *testing.T) {
	readDirMock := func(name string) ([]os.DirEntry, error) {
		return []os.DirEntry{
			MockDirEntry{name: "file1"},
			MockDirEntry{name: "file2"},
			MockDirEntry{name: "file3"},
		}, nil
	}

	parseFilesMock := func(filenames ...string) (*template.Template, error) {
		return template.Must(template.New("template1").Parse("SomeTemplate")), nil
	}

	type args struct {
		w   http.ResponseWriter
		in1 *http.Request
	}
	tests := []struct {
		name        string
		args        args
		expectedLen int
	}{
		{
			"OrderListHandler should create order template",
			args{w: httptest.NewRecorder(), in1: httptest.NewRequest("GET", "/", nil)},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OrdersGetHandler(tt.args.w, tt.args.in1, readDirMock, parseFilesMock)
			if len(result.Orders) != tt.expectedLen {
				t.Errorf("OrdersGetHandler() got = %v, want %v", len(result.Orders), tt.expectedLen)
			}
		})
	}
}
