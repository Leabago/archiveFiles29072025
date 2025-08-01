package entity

type DownloadResult struct {
	Filename string
	Content  []byte
	Error    error
	FileNum  int
}
