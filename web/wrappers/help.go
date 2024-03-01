package wrappers

type (
	MapDetailReplace struct { //替换
		Video map[string]string `json:"video"`
		Img   map[string]string `json:"img"`
		//File  map[string]string `json:"file"`
	}
)
