package build_constants

const (
	CustomPatrolCLIVersion = "CUSTOM_PATROL_CLI_VERSION" // Optional, using latest when empty
	TestTargetDirectory    = "TEST_TARGET_DIRECTORY"     // Required
	Platform               = "PLATFORM"                  // Required, using both as default
	BuildType              = "TEST_BUILD_TYPE"           // Required, using release as default
	Tags                   = "TAGS"                      // optional, using empty string as default
	ExcludedTags           = "EXCLUDED_TAGS"             // optional, using empty string as default
	IsVerboseMode          = "IS_VERBOSE_MODE"           // optional, using false as default
)
