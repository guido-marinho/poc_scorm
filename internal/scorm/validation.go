package scorm

type AssessmentQuestion struct {
	Text         string `json:"text" validate:"required"`
	Type         string `json:"type" validate:"required,oneof=MULTI SINGLE"`
	UUID         string `json:"uuid" validate:"required"`
	Points       int    `json:"points" validate:"gte=0"`
	Alternatives []struct {
		Text     string `json:"text" validate:"required"`
		UUID     string `json:"uuid" validate:"required"`
		Correct  bool   `json:"correct"`
		Feedback string `json:"feedback"`
	} `json:"alternatives,omitempty" validate:"omitempty,dive"`
}

type Assessment struct {
	UUID      string               `json:"uuid" validate:"required"`
	Questions []AssessmentQuestion `json:"questions" validate:"omitempty,dive"`
}

type Topic struct {
	Name                  string      `json:"name" validate:"required"`
	Type                  string      `json:"type" validate:"required,oneof=LECTURE ASSESSMENT"`
	UUID                  string      `json:"uuid" validate:"required"`
	Order                 int         `json:"order" validate:"gte=0"`
	ScrimbaUrl            string      `json:"scrimbaUrl" validate:"omitempty,url"`
	Description           string      `json:"description"`
	VideoLength           *int        `json:"videoLength,omitempty" validate:"omitempty,gte=0"`
	Observations          string      `json:"observations"`
	MuxPlaybackId         string      `json:"muxPlaybackId,omitempty"`
	ExternalVideoUrl      string      `json:"externalVideoUrl,omitempty"`
	DigitalCourseId       string      `json:"digitalCourseId" validate:"required"`
	DigitalCourseModuleId string      `json:"digitalCourseModuleId" validate:"required"`
	Assessment            *Assessment `json:"assessment,omitempty" validate:"omitempty"`
}

type Module struct {
	Name   string  `json:"name" validate:"required"`
	UUID   string  `json:"uuid" validate:"required"`
	Order  int     `json:"order" validate:"gte=0"`
	Topics []Topic `json:"topics" validate:"omitempty,dive"`
}

type DigitalCourse struct {
	Logo         string   `json:"logo"`
	Name         string   `json:"name" validate:"required"`
	UUID         string   `json:"uuid" validate:"required"`
	Modules      []Module `json:"modules" validate:"omitempty,dive"`
	Description  string   `json:"description,omitempty"`
	ThumbnailUrl string   `json:"thumbnailUrl,omitempty" validate:"omitempty,url"`
	CourseType   string   `json:"courseType" validate:"required"`
}
