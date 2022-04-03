package domain

type (
	FileStatus int
)

type File struct {
	ContentType     string             `json:"contentType"`
	Name            string             `json:"name"`
	Size            int64              `json:"size"`
	URL             string             `json:"url"`
}