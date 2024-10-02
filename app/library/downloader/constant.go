package downloader

// DownloadStatus is a type for download status
type DownloadStatus int

const (
	// DownloadedSuccessfully is a status for downloaded successfully
	DownloadedSuccessfully DownloadStatus = iota
	// DownloadedFailed is a status for downloaded failed
	DownloadedFailed
	// DownloadedAlready is a status for downloaded already
	DownloadedAlready
)

const (
	// WordNet Japan database archive file URL
	WNJPN_DB_ARCHIVE_FILE_URL = "https://github.com/bond-lab/wnja/releases/download/v1.1/wnjpn.db.gz"
	// WordNet Japan database archive file name
	WNJPN_DB_ARCHIVE_FILE_NAME = "wnjpn.db.gz"
	// WordNet Japan sqlite3 database file name
	WNJPN_DB_FILE_NAME = "wnjpn.db"
)
