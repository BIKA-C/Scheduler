package course

import (
	"scheduler/util"
	"time"
)

// Class is an instance of a section
type Class struct {
	Meta
	SectionID    util.ID       `json:"section"`
	Index        int           `json:"index"`
	Title        string        `json:"title"`
	Time         time.Time     `json:"time"`
	Duration     time.Duration `json:"duration"`
	UnitPrice    int           `json:"unitPrice"` // Unit Price in $/minute
	Remark       string        `json:"remark,omitempty"`
	InstructorID util.UUID     `json:"instructorID,omitempty"`
}
