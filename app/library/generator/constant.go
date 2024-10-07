package generator

// GenerateMode is a type for generate mode.
type GenerateMode int

const (
	// WithNoPrefixOrSuffix is a mode for generate without prefix or suffix.
	WithNoPrefixOrSuffix GenerateMode = iota
	// WithPrefix is a mode for generate with prefix.
	WithPrefix
	// WithSuffix is a mode for generate with suffix.
	WithSuffix
)

// GenerateResult is a type for generate result.
type GenerateResult int

const (
	// GeneratedSuccessfully is a result for generated successfully.
	GeneratedSuccessfully GenerateResult = iota
	// GeneratedFailed is a result for generated failed.
	GeneratedFailed
	// DBFileNotFound is a result for generated failed because db file is not found.
	DBFileNotFound
)
