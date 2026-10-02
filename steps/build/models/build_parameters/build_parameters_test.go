package build_parameters

import (
	"reflect"
	"testing"
)

func TestNewBuildParameters_Validation(t *testing.T) {
	tests := []struct {
		name    string
		envMap  map[string]string
		wantErr bool
	}{
		{
			name: "missing platform",
			envMap: map[string]string{
				"buildType": "release",
			},
			wantErr: true,
		},
		{
			name: "empty platform",
			envMap: map[string]string{
				"platform":  "   ",
				"buildType": "release",
			},
			wantErr: true,
		},
		{
			name: "invalid platform",
			envMap: map[string]string{
				"platform":  "windows",
				"buildType": "release",
			},
			wantErr: true,
		},
		{
			name: "missing buildType",
			envMap: map[string]string{
				"platform": "android",
			},
			wantErr: true,
		},
		{
			name: "invalid buildType",
			envMap: map[string]string{
				"platform":  "android",
				"buildType": "profile",
			},
			wantErr: true,
		},
		{
			name: "valid android release",
			envMap: map[string]string{
				"platform":  "android",
				"buildType": "release",
			},
			wantErr: false,
		},
		{
			name: "valid ios debug",
			envMap: map[string]string{
				"platform":  "ios",
				"buildType": "debug",
			},
			wantErr: false,
		},
		{
			name: "valid both with all options",
			envMap: map[string]string{
				"platform":     "both",
				"buildType":    "debug",
				"target":       "integration_test/app_test.dart",
				"tags":         "smoke, e2e",
				"excludedTags": "flaky",
				"verbose":      "true",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bp, err := NewBuildParameters(tt.envMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBuildParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && bp == nil {
				t.Error("expected non-nil BuildParameters")
			}
		})
	}
}

func TestBuildParameters_Command(t *testing.T) {
	tests := []struct {
		name     string
		params   BuildParameters
		expected []string
	}{
		{
			name: "android release",
			params: BuildParameters{
				Platform:  "android",
				BuildType: "release",
			},
			expected: []string{"patrol build android --release"},
		},
		{
			name: "android debug",
			params: BuildParameters{
				Platform:  "android",
				BuildType: "debug",
			},
			expected: []string{"patrol build android --debug"},
		},
		{
			name: "ios release",
			params: BuildParameters{
				Platform:  "ios",
				BuildType: "release",
			},
			expected: []string{"patrol build ios --release"},
		},
		{
			name: "ios debug (adds --simulator)",
			params: BuildParameters{
				Platform:  "ios",
				BuildType: "debug",
			},
			expected: []string{"patrol build ios --debug --simulator"},
		},
		{
			name: "both release",
			params: BuildParameters{
				Platform:  "both",
				BuildType: "release",
			},
			expected: []string{
				"patrol build android --release",
				"patrol build ios --release",
			},
		},
		{
			name: "both debug",
			params: BuildParameters{
				Platform:  "both",
				BuildType: "debug",
			},
			expected: []string{
				"patrol build android --debug",
				"patrol build ios --debug --simulator",
			},
		},
		{
			name: "android with target and tags and verbose",
			params: BuildParameters{
				Platform:     "android",
				BuildType:    "release",
				Target:       "integration_test/login_test.dart",
				Tags:         "'( smoke && e2e )'",
				ExcludedTags: "'( flaky )'",
				IsVerbose:    "--verbose",
			},
			expected: []string{
				"patrol build android --release --target integration_test/login_test.dart --tags '( smoke && e2e )' --excludedTags '( flaky )' --verbose",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.params.Command()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Command() = %q, want %q", got, tt.expected)
			}
		})
	}
}
