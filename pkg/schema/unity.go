package schema

type UnityMeta struct {
	App     string `json:"app"`
	Version string `json:"version"`
	Image   string `json:"image"`
	Size    struct {
		W int `json:"w"`
		H int `json:"h"`
	} `json:"size"`
}

type UnityFrame struct {
	Frame struct {
		X int `json:"x"`
		Y int `json:"y"`
		W int `json:"w"`
		H int `json:"h"`
	} `json:"frame"`
	Pivot struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"pivot"`
}

type UnityAtlasJSON struct {
	Frames map[string]UnityFrame `json:"frames"`
	Meta   UnityMeta             `json:"meta"`
}
