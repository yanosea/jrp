package repository

// SaveStatus is a type for save status
type SaveStatus int

const (
	// SavedSuccessfully is a status for saved successfully
	SavedSuccessfully SaveStatus = iota
	// SavedFailed is a status for saved failed
	SavedFailed
	// SavedNone is a status for saved none
	SavedNone
	// SavedNotAll is a status for saved not all
	SavedNotAll
)

// RemoveStatus is a type for remove status
type RemoveStatus int

const (
	// RemovedSuccessfully is a status for removed successfully
	RemovedSuccessfully RemoveStatus = iota
	// RemovedFailed is a status for removed failed
	RemovedFailed
	// RemovedNone is a status for removed none
	RemovedNone
	// RemovedNotAll is a status for removed not all
	RemovedNotAll
)

// AddStatus is a type for add status
type AddStatus int

const (
	// AddedSuccessfully is a status for added successfully
	AddedSuccessfully AddStatus = iota
	// AddedFailed is a status for added failed
	AddedFailed
	// AddedNone is a status for added none
	AddedNone
	// AddedNotAll is a status for added not all
	AddedNotAll
)

const (
	// jrp sqlite database file name
	JRP_DB_FILE_NAME = "jrp.db"
)
