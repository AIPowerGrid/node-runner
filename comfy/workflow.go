package comfy

type Workflow struct {
	Num3 struct {
		Inputs struct {
			Seed        int64  `json:"seed"`
			Steps       int    `json:"steps"`
			Cfg         int    `json:"cfg"`
			SamplerName string `json:"sampler_name"`
			Scheduler   string `json:"scheduler"`
			Denoise     int    `json:"denoise"`
			Model       []any  `json:"model"`
			Positive    []any  `json:"positive"`
			Negative    []any  `json:"negative"`
			LatentImage []any  `json:"latent_image"`
		} `json:"inputs"`
		ClassType string `json:"class_type"`
		Meta      struct {
			Title string `json:"title"`
		} `json:"_meta"`
	} `json:"3"`
	Num4 struct {
		Inputs struct {
			CkptName string `json:"ckpt_name"`
		} `json:"inputs"`
		ClassType string `json:"class_type"`
		Meta      struct {
			Title string `json:"title"`
		} `json:"_meta"`
	} `json:"4"`
	Num5 struct {
		Inputs struct {
			Width     int `json:"width"`
			Height    int `json:"height"`
			BatchSize int `json:"batch_size"`
		} `json:"inputs"`
		ClassType string `json:"class_type"`
		Meta      struct {
			Title string `json:"title"`
		} `json:"_meta"`
	} `json:"5"`
	Num6 struct {
		Inputs struct {
			Text string `json:"text"`
			Clip []any  `json:"clip"`
		} `json:"inputs"`
		ClassType string `json:"class_type"`
		Meta      struct {
			Title string `json:"title"`
		} `json:"_meta"`
	} `json:"6"`
	Num7 struct {
		Inputs struct {
			Text string `json:"text"`
			Clip []any  `json:"clip"`
		} `json:"inputs"`
		ClassType string `json:"class_type"`
		Meta      struct {
			Title string `json:"title"`
		} `json:"_meta"`
	} `json:"7"`
	Num8 struct {
		Inputs struct {
			Samples []any `json:"samples"`
			Vae     []any `json:"vae"`
		} `json:"inputs"`
		ClassType string `json:"class_type"`
		Meta      struct {
			Title string `json:"title"`
		} `json:"_meta"`
	} `json:"8"`
	Num9 struct {
		Inputs struct {
			FilenamePrefix string `json:"filename_prefix"`
			Images         []any  `json:"images"`
		} `json:"inputs"`
		ClassType string `json:"class_type"`
		Meta      struct {
			Title string `json:"title"`
		} `json:"_meta"`
	} `json:"9"`
}
