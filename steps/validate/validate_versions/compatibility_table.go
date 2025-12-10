package validate_versions

import (
	v "github.com/Masterminds/semver/v3"
)

/*
This struct contains the compatibility table for Patrol CLI and Patrol
The compatibility table is used to check if the Patrol CLI version is compatible
*/
type VersionRange struct {
	Min *v.Version
	Max *v.Version
}

// CompatibilityEntry links Patrol CLI, Patrol, and Flutter version requirements.
type CompatibilityEntry struct {
	PatrolCLIRange    VersionRange
	PatrolRange       VersionRange
	MinFlutterVersion *v.Version
}

var CompatibilityTable = []CompatibilityEntry{
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("4.0.0"), Max: v.MustParse("4.0.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("4.0.0"), Max: v.MustParse("4.0.0")},
		MinFlutterVersion: v.MustParse("3.32.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.11.0"), Max: v.MustParse("3.11.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.20.0"), Max: v.MustParse("3.20.0")},
		MinFlutterVersion: v.MustParse("3.32.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.9.0"), Max: v.MustParse("3.10.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.18.0"), Max: v.MustParse("3.19.0")},
		MinFlutterVersion: v.MustParse("3.32.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.7.0"), Max: v.MustParse("3.8.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.16.0"), Max: v.MustParse("3.17.0")},
		MinFlutterVersion: v.MustParse("3.32.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.5.0"), Max: v.MustParse("3.6.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.14.0"), Max: v.MustParse("3.15.2")},
		MinFlutterVersion: v.MustParse("3.24.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.4.1"), Max: v.MustParse("3.4.1")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.13.1"), Max: v.MustParse("3.13.2")},
		MinFlutterVersion: v.MustParse("3.24.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.4.0"), Max: v.MustParse("3.4.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.13.0"), Max: v.MustParse("3.13.0")},
		MinFlutterVersion: v.MustParse("3.24.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.3.0"), Max: v.MustParse("3.3.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.12.0"), Max: v.MustParse("3.12.0")},
		MinFlutterVersion: v.MustParse("3.24.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.2.1"), Max: v.MustParse("3.2.1")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.11.2"), Max: v.MustParse("3.11.2")},
		MinFlutterVersion: v.MustParse("3.24.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.2.0"), Max: v.MustParse("3.2.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.11.0"), Max: v.MustParse("3.11.1")},
		MinFlutterVersion: v.MustParse("3.22.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("3.1.0"), Max: v.MustParse("3.1.1")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.10.0"), Max: v.MustParse("3.10.0")},
		MinFlutterVersion: v.MustParse("3.22.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("2.6.5"), Max: v.MustParse("3.0.1")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.6.0"), Max: v.MustParse("3.10.0")},
		MinFlutterVersion: v.MustParse("3.16.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("2.6.0"), Max: v.MustParse("2.6.4")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.4.0"), Max: v.MustParse("3.5.2")},
		MinFlutterVersion: v.MustParse("3.16.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("2.3.0"), Max: v.MustParse("2.5.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("3.0.0"), Max: v.MustParse("3.3.0")},
		MinFlutterVersion: v.MustParse("3.16.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("2.2.0"), Max: v.MustParse("2.2.2")},
		PatrolRange:       VersionRange{Min: v.MustParse("2.3.0"), Max: v.MustParse("2.3.2")},
		MinFlutterVersion: v.MustParse("3.3.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("2.0.1"), Max: v.MustParse("2.1.5")},
		PatrolRange:       VersionRange{Min: v.MustParse("2.0.1"), Max: v.MustParse("2.2.5")},
		MinFlutterVersion: v.MustParse("3.3.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("2.0.0"), Max: v.MustParse("2.0.0")},
		PatrolRange:       VersionRange{Min: v.MustParse("2.0.0"), Max: v.MustParse("2.0.0")},
		MinFlutterVersion: v.MustParse("3.3.0"),
	},
	{
		PatrolCLIRange:    VersionRange{Min: v.MustParse("1.1.4"), Max: v.MustParse("1.1.11")},
		PatrolRange:       VersionRange{Min: v.MustParse("1.0.9"), Max: v.MustParse("1.1.11")},
		MinFlutterVersion: v.MustParse("3.3.0"),
	},
}
