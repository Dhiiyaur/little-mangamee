package entity

type IndexData struct {
	OriginalServer string `json:"original_server"`
	Id             string `json:"id"`
	Cover          string `json:"cover"`
	Title          string `json:"title"`
	LastChapter    string `json:"last_chapter"`
}

type SearchData struct {
	Cover          string `json:"cover"`
	Title          string `json:"title"`
	Id             string `json:"id"`
	OriginalServer string `json:"original_server"`
	LastChapter    string `json:"last_chapter"`
}
type Chapter struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ChapterData struct {
	Chapters       []Chapter `json:"chapters"`
	OriginalServer string    `json:"original_server"`
}

type DataChapters struct {
	ChapterName string  `json:"chapter_name"`
	Images      []Image `json:"images"`
}

type DetailData struct {
	OriginalServer string    `json:"original_server"`
	Cover          string    `json:"cover"`
	Title          string    `json:"title"`
	Summary        string    `json:"summary"`
	Chapters       []Chapter `json:"chapters"`
}

type Image struct {
	Image string `json:"image"`
}

type ImageData struct {
	OriginalServer string  `json:"original_server"`
	ChapterName    string  `json:"chapter_name"`
	Images         []Image `json:"images"`
}

type RequestLink struct {
	Url string
}
