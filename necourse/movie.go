package necourse

type SubTitleData struct {
	Name string `json:"subName"`
	Size int    `json:"subSize"`
	Url  string `json:"subUrl"`
}

type VideoData struct {
	ImgPath               string         `json:"imgpath"`
	MId                   string         `json:"mid"`
	MLength               int            `json:"mlength"`
	Mp4Size               int            `json:"mp4size"`
	Mp4SizeOrigin         int            `json:"mp4sizeOrigin"`
	Repovideourl          string         `json:"repovideourl"`
	RepovideourlOrigin    string         `json:"repovideourlOrigin"`
	Repovideourlmp4       string         `json:"repovideourlmp4"`
	Repovideourlmp4Origin string         `json:"repovideourlmp4Origin"`
	SubList               []SubTitleData `json:"subList"`
	Subtitle              string         `json:"subtitle"`
	SubtitleLanguage      string         `json:"subtitle_language"`
	Title                 string         `json:"title"`
	Weburl                string         `json:"weburl"`
}

type MovieData struct {
	CCPic       string      `json:"ccPic"`
	CCUrl       string      `json:"ccUrl"`
	Description string      `json:"description"`
	Director    string      `json:"director"`
	Hits        int         `json:"hits"`
	ImgPath     string      `json:"imgpath"`
	LargeImgUrl string      `json:"largeimgurl"`
	LTime       int         `json:"ltime"`
	School      string      `json:"school"`
	Source      string      `json:"source"`
	SubTitle    string      `json:"subtitle"`
	Tags        string      `json:"tags"`
	Title       string      `json:"title"`
	VideoList   []VideoData `json:"videoList"`
}

type Subtitle struct {
	data SubTitleData
}

func (s *Subtitle) Name() string {
	return s.getData().Name
}

func (s *Subtitle) Size() int {
	return s.getData().Size
}

func (s *Subtitle) Url() string {
	return s.getData().Url
}

func (s *Subtitle) getData() *SubTitleData {
	return &s.data
}

type Video struct {
	data VideoData
}

func (v *Video) Id() string {
	return v.getData().MId
}

func (v *Video) ImgUrl() string {
	return v.getData().ImgPath
}

func (v *Video) Title() string {
	return v.getData().Title
}

func (v *Video) Length() int {
	return v.getData().MLength
}

func (v *Video) Size() int {
	return v.getData().Mp4Size
}

func (v *Video) Subtitles() []Subtitle {
	dataList := v.getData().SubList
	subtitles := make([]Subtitle, len(dataList))

	for index, data := range dataList {
		subtitles[index] = Subtitle{data}
	}

	return subtitles
}

func (v *Video) Url(quality Quality) string {
	if quality == SD {
		return v.SDUrl()
	}
	return v.HDUrl()
}

func (v *Video) SDUrl() string {
	return v.getData().Repovideourlmp4Origin
}

func (v *Video) HDUrl() string {
	return v.getData().Repovideourl
}

func (v *Video) getData() *VideoData {
	return &v.data
}
