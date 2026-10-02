package get_patrol_version

import (
	"testing"

	"patrol_build_flutter/commands"

	v "github.com/Masterminds/semver/v3"
)

func Test_GetPatrolVersionFromLog(t *testing.T) {
	tests := []struct {
		name    string
		log     string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid patrol line",
			log:     "- patrol 3.15.1 [boolean_selector equatable flutter flutter_test http json_annotation meta patrol_finders patrol_log shelf test_api]",
			want:    "3.15.1",
			wantErr: false,
		}, {
			name:    "Invalid patrol name with -",
			log:     "- patrol-abc 3.15.1",
			want:    "",
			wantErr: true,
		}, {
			name:    "Invalid patrol name with underscore",
			log:     "- patrol_bad_scenario 3.15.1",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Non installed patrol",
			log:     "- animations 2.1.0 [flutter]",
			want:    "",
			wantErr: true,
		},
		{
			name: "Multiple patrol lines",
			log: `
			- patrol_finders 2.7.2 [flutter flutter_test meta patrol_log]
            - patrol_log 0.3.0 [dispose_scope equatable json_annotation]
            - patrol 3.15.1 [boolean_selector equatable flutter flutter_test http json_annotation meta patrol_finders patrol_log shelf test_api]
			`,
			want:    "3.15.1",
			wantErr: false,
		},
		{
			name:    "Invalid patrol version format",
			log:     "- patrol b3.v15.1+1 [boolean_selector equatable flutter flutter_test http json_annotation meta patrol_finders patrol_log shelf test_api]",
			want:    "3.15.1+1",
			wantErr: true,
		},
		{
			name:    "Valid Semantic version format",
			log:     "- patrol v3.15.1+1 [boolean_selector equatable flutter flutter_test http json_annotation meta patrol_finders patrol_log shelf test_api]",
			want:    "3.15.1+1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPatrolVersionFromLog(tt.log)
			t.Logf("📝 %s\n  Log: %s\n  Want: %s\n  WantErr: %t\n  Got: %v\n",
				tt.name, tt.log, tt.want, tt.wantErr, got)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != nil && tt.want != "" {
				wantVer, _ := v.NewVersion(tt.want)
				if !got.Equal(wantVer) {
					t.Errorf("got = %v, want %v", got, wantVer)
				}
			} else if got != nil && tt.want == "" {
				t.Errorf("expected nil, got %v", got)
			}
		})
	}
}

func Test_GetPatrolVersion(t *testing.T) {
	t.Run("wrong command returns error", func(t *testing.T) {
		wrongCmd := commands.Command{Name: "echo", Args: []string{"hello"}}
		_, err := GetPatrolVersion(wrongCmd)
		if err == nil {
			t.Error("expected error for wrong command, got nil")
		}
	})
}
