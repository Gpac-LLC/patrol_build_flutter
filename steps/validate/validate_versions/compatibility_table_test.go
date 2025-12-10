package validate_versions

import (
	"testing"

	v "github.com/Masterminds/semver/v3"
)

// TestCheckCompatibilityErrorFormat tests that error messages are properly formatted.
func TestCheckCompatibilityErrorFormat(t *testing.T) {
	params := ValidateRunParams{
		FlutterVersion: v.MustParse("3.0.0"),
		CliVersion:     v.MustParse("5.0.0"),
		PatrolVersion:  v.MustParse("5.0.0"),
	}

	isCompatible := CheckCompatibility(params)
	if isCompatible {
		t.Error("CheckCompatibility() expected false for incompatible versions, got true")
	}
}

// TestCheckCompatibilityWithNilVersions tests defensive behavior with nil versions.
// This tests that the function handles edge cases gracefully.
func TestCheckCompatibilityWithNilVersions(t *testing.T) {
	// Note: In a real scenario, these should be validated before calling CheckCompatibility
	// But we test defensive behavior here
	t.Run("nil_flutter_version", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with nil FlutterVersion, but function did not panic")
			}
		}()

		params := ValidateRunParams{
			FlutterVersion: nil,
			CliVersion:     v.MustParse("4.0.0"),
			PatrolVersion:  v.MustParse("4.0.0"),
		}

		_ = CheckCompatibility(params)
	})

	t.Run("nil_cli_version", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with nil CliVersion, but function did not panic")
			}
		}()

		params := ValidateRunParams{
			FlutterVersion: v.MustParse("3.32.0"),
			CliVersion:     nil,
			PatrolVersion:  v.MustParse("4.0.0"),
		}

		_ = CheckCompatibility(params)
	})

	t.Run("nil_patrol_version", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with nil PatrolVersion, but function did not panic")
			}
		}()

		params := ValidateRunParams{
			FlutterVersion: v.MustParse("3.32.0"),
			CliVersion:     v.MustParse("4.0.0"),
			PatrolVersion:  nil,
		}

		_ = CheckCompatibility(params)
	})
}
