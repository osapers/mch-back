package postgres

import (
	"testing"

	"github.com/go-pg/pg/v10"
)

func TestOpen(t *testing.T) {
	cleanup, url, _, err := StartDBInDocker()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = cleanup(); err != nil {
			t.Error(err)
		}
	}()

	connOpts, err := pg.ParseURL(url)
	if err != nil {
		t.Error(err)
	}

	type args struct {
		connOpts *pg.Options
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				connOpts: connOpts,
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				connOpts: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := open(tt.args.connOpts)
			defer func() {
				if err != nil {
					return
				}
				err = got.Close()
				if err != nil {
					t.Errorf("got.Destroy() returns error %v", err)
				}
			}()
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Error("Open() wanted error and got != nil")
			}
		})
	}
}

func Test_applyMigrations(t *testing.T) {
	cleanup, url, _, err := StartDBInDocker()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := cleanup(); err != nil {
			t.Error(err)
		}
	}()

	type args struct {
		connectionURL string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				connectionURL: url,
			},
			wantErr: false,
		},
		{
			name: "bad-url",
			args: args{
				connectionURL: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := applyMigrations(tt.args.connectionURL); (err != nil) != tt.wantErr {
				t.Errorf("Migrate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
