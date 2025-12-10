package validate_versions

import (
	"testing"

	v "github.com/Masterminds/semver/v3"
)

// TestIsVersionInRangeFunction tests the isVersionInRange function with various scenarios.
func TestIsVersionInRangeFunction(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		rangeMin  string
		rangeMax  string
		isInRange bool
	}{
		{
			name:      "version_equals_min_and_less_than_max",
			version:   "3.0.0",
			rangeMin:  "3.0.0",
			rangeMax:  "3.1.0",
			isInRange: true,
		},
		{
			name:      "version_between_min_and_max",
			version:   "3.5.0",
			rangeMin:  "3.0.0",
			rangeMax:  "4.0.0",
			isInRange: true,
		},
		{
			name:      "version_equals_max",
			version:   "3.1.0",
			rangeMin:  "3.0.0",
			rangeMax:  "3.1.0",
			isInRange: true,
		},
		{
			name:      "version_equals_both_min_and_max",
			version:   "4.0.0",
			rangeMin:  "4.0.0",
			rangeMax:  "4.0.0",
			isInRange: true,
		},
		{
			name:      "version_less_than_min",
			version:   "2.9.0",
			rangeMin:  "3.0.0",
			rangeMax:  "3.1.0",
			isInRange: false,
		},
		{
			name:      "version_greater_than_max",
			version:   "3.2.0",
			rangeMin:  "3.0.0",
			rangeMax:  "3.1.0",
			isInRange: false,
		},
		{
			name:      "version_greater_than_min_less_than_max",
			version:   "3.0.5",
			rangeMin:  "3.0.0",
			rangeMax:  "3.1.0",
			isInRange: true,
		},
		{
			name:      "version_just_below_min",
			version:   "2.8.9",
			rangeMin:  "2.9.0",
			rangeMax:  "3.0.0",
			isInRange: false,
		},
		{
			name:      "version_just_above_max",
			version:   "3.1.1",
			rangeMin:  "3.0.0",
			rangeMax:  "3.1.0",
			isInRange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := v.MustParse(tt.version)
			versionRange := VersionRange{
				Min: v.MustParse(tt.rangeMin),
				Max: v.MustParse(tt.rangeMax),
			}

			got := isVersionInRange(version, versionRange)
			if got != tt.isInRange {
				t.Errorf("isVersionInRange(%s, Range[%s-%s]) = %v, want %v",
					tt.version, tt.rangeMin, tt.rangeMax, got, tt.isInRange)
			}
		})
	}
}

// TestTableCheckCompatibility
// This test checks the CheckCompatibility function against the real CompatibilityTable.
// Each test case includes a comment referencing the relevant compatibility entry or why it fails.
func TestTableCheckCompatibility(t *testing.T) {
	tests := []struct {
		name             string
		flutterVersion   string
		patrolCLIVersion string
		patrolVersion    string
		areCompatible    bool
		context          string
	}{
		{
			name:             "compatible_latest_max_versions", // ✅ Matches: CLI [4.0.0-4.0.1], Patrol [4.0.0-4.0.0], Flutter 3.32.0
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "4.0.1",
			patrolVersion:    "4.0.0",
			areCompatible:    true,
			context:          "✅ Matches table entry: CLI [4.0.0-4.0.1], Patrol [4.0.0], Flutter 3.32.0",
		},
		{
			name:             "compatible_min_version_match", // ✅ Matches: CLI [3.7.0-3.8.0], Patrol [3.16.0-3.17.0], Flutter 3.32.0
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "3.7.0",
			patrolVersion:    "3.16.0",
			areCompatible:    true,
			context:          "✅ Matches table entry: CLI [3.7.0-3.8.0], Patrol [3.16.0-3.17.0], Flutter 3.32.0",
		},

		{
			name:             "compatible_patrol_version_in_between_range", // ✅ Matches: CLI [3.9.0-3.10.0], Patrol [3.18.0-3.19.0], Flutter 3.32.0
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "3.9.0",
			patrolVersion:    "3.18.5",
			areCompatible:    true,
			context:          "✅ Matches table entry: CLI [3.9.0-3.10.0], Patrol [3.18.0-3.19.0], Flutter 3.32.0",
		},
		{
			name:             "compatible_patrol_cli_version_in_between_range", // ✅ Matches: CLI [3.7.0-3.8.0], Patrol [3.16.0-3.17.0], Flutter 3.32.0
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "3.7.5",
			patrolVersion:    "3.17.0",
			areCompatible:    true,
			context:          "✅ Matches table entry: CLI [3.7.0-3.8.0], Patrol [3.16.0-3.17.0], Flutter 3.32.0",
		},

		{
			name:             "compatible_cli_2.6.5_to_3.7.5", // ✅ Matches: CLI [2.6.5-3.0.1], Patrol [3.6.0-3.10.0], Flutter 3.16.0
			flutterVersion:   "3.16.0",
			patrolCLIVersion: "2.6.5",
			patrolVersion:    "3.7.5",
			areCompatible:    true,
			context:          "✅ Matches table entry: CLI [2.6.5-3.0.1], Patrol [3.6.0-3.10.0], Flutter 3.16.0",
		},
		{
			name:             "compatible_2_6_5_to_3_0_1_at_max", // ✅ Matches: CLI [2.6.5-3.0.1], Patrol [3.6.0-3.10.0], Flutter 3.16.0
			flutterVersion:   "3.16.0",
			patrolCLIVersion: "3.0.1",
			patrolVersion:    "3.10.0",
			areCompatible:    true,
			context:          "✅ Matches table entry: CLI [2.6.5-3.0.1], Patrol [3.6.0-3.10.0], Flutter 3.16.0",
		},
		{
			name:             "compatible_oldest_entry", // ✅ Matches: CLI [1.1.4-1.1.11], Patrol [1.0.9-1.1.11], Flutter 3.3.0
			flutterVersion:   "3.3.0",
			patrolCLIVersion: "1.1.4",
			patrolVersion:    "1.0.9",
			areCompatible:    true,
			context:          "✅ Matches table entry: CLI [1.1.4-1.1.11], Patrol [1.0.9-1.1.11], Flutter 3.3.0",
		},
		{
			name:             "compatible_oldest_at_max", // ✅ Matches: CLI [1.1.4-1.1.11], Patrol [1.0.9-1.1.11], Flutter 3.3.0
			flutterVersion:   "3.3.0",
			patrolCLIVersion: "1.1.11",
			patrolVersion:    "1.1.11",
			areCompatible:    true,
			context:          "✅ Matches table entry: CLI [1.1.4-1.1.11], Patrol [1.0.9-1.1.11], Flutter 3.3.0",
		},
		{
			name:             "nonCompatible_flutter_above_required_version", // ❌ Flutter version too high, not in table
			flutterVersion:   "3.38.1",
			patrolCLIVersion: "4.0.1",
			patrolVersion:    "4.0.0",
			areCompatible:    false,
			context:          "❌ Flutter version 3.38.1 not in table (only 3.32.0 allowed for this entry)",
		},
		{
			name:             "nonCompatible_flutter_version_lower", // ❌ Flutter version 3.22.0 only allowed for CLI 3.2.0, Patrol 3.11.0-3.11.1
			flutterVersion:   "3.22.0",
			patrolCLIVersion: "3.2.1",
			patrolVersion:    "3.11.2",
			areCompatible:    false,
			context:          "❌ Flutter 3.22.0 with CLI 3.2.1, Patrol 3.11.2 not in table (only CLI 3.2.0, Patrol 3.11.0-3.11.1 allowed)",
		},
		{
			name:             "nonCompatible_patrol_cli_and_patrol_versions_dont_match_any_entry", // ❌ Versions too high, not in table
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "5.0.0",
			patrolVersion:    "5.0.0",
			areCompatible:    false,
			context:          "❌ CLI 5.0.0 and Patrol 5.0.0 not in any table entry",
		},
		{
			name:             "nonCompatible_patrol_cli_version_above_range_max", // ❌ CLI version above max for any entry
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "4.1.0",
			patrolVersion:    "4.0.0",
			areCompatible:    false,
			context:          "❌ CLI 4.1.0 above max (4.0.1) for Flutter 3.32.0",
		},
		{
			name:             "nonCompatible_patrol_version_above_range_max", // ❌ Patrol version above max for any entry
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "4.0.0",
			patrolVersion:    "4.1.0",
			areCompatible:    false,
			context:          "❌ Patrol 4.1.0 above max (4.0.0) for Flutter 3.32.0",
		},
		{
			name:             "nonCompatible_patrol_cli_version_below_range_min", // ❌ CLI version below min for any entry
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "3.8.0",
			patrolVersion:    "3.18.0",
			areCompatible:    false,
			context:          "❌ CLI 3.8.0 not in any valid range for Flutter 3.32.0",
		},
		{
			name:             "nonCompatible_patrol_version_above_range_max_for_range", // ❌ Patrol version above max for entry
			flutterVersion:   "3.32.0",
			patrolCLIVersion: "3.9.0",
			patrolVersion:    "3.20.0",
			areCompatible:    false,
			context:          "❌ Patrol 3.20.0 above max (3.19.0) for CLI 3.9.0, Flutter 3.32.0",
		},
		{
			name:             "nonCompatible_cli_and_patrol_version_above_range", // ❌ CLI & Patrol version above max for entry
			flutterVersion:   "3.24.0",
			patrolCLIVersion: "3.6.1",
			patrolVersion:    "3.15.3",
			areCompatible:    false,
			context:          "❌ CLI and Patrol right above max for Flutter 3.24.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ValidateRunParams{
				FlutterVersion: v.MustParse(tt.flutterVersion),
				CliVersion:     v.MustParse(tt.patrolCLIVersion),
				PatrolVersion:  v.MustParse(tt.patrolVersion),
			}

			got := CheckCompatibility(params)

			t.Logf("📝 %s\n  Flutter: %s\n  Patrol CLI: %s\n  Patrol: %s\n  Expected: %v\n  Got: %v\n  Context: %s\n",
				tt.name, tt.flutterVersion, tt.patrolCLIVersion, tt.patrolVersion, tt.areCompatible, got, tt.context)

			if got != tt.areCompatible {
				t.Errorf("❌ CheckCompatibility() failed for test '%s'\n  Expected: %v\n  Got:      %v\n  Context:  %s",
					tt.name, tt.areCompatible, got, tt.context)
			}
		})
	}
}
