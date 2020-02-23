package url

import "testing"

func Test_parseScheme(t *testing.T) {
	type args struct {
		rawurl string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{
			name: "gs://...",
			args: args{
				rawurl: "gs://bucket/object",
			},
			wantResult: "gs",
			wantErr:    false,
		},
		{
			name: "no scheme",
			args: args{
				rawurl: "bad_url",
			},
			wantErr: true,
		},
		{
			name: "bigquery://...",
			args: args{
				rawurl: "bigquery://project.dataset.table",
			},
			wantResult: "bigquery",
			wantErr: false,
		},
		{
			name: "case shouldn't matter",
			args: args{
				rawurl: "BigQuery://project.dataset.table",
			},
			wantResult: "bigquery",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := parseScheme(tt.args.rawurl)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseScheme() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("parseScheme() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
