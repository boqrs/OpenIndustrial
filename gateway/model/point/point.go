package point

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/object"
)

// KindPoint is the object kind for Point.
const (
	KindPoint object.Kind = "Point"
)

// DataType represents the logical type of a point value.
type DataType uint8

const (
	TypeBool DataType = iota
	TypeInt64
	TypeFloat64
	TypeString
	TypeBytes
)

// Quality indicates whether a point value is trustworthy.
type Quality uint8

const (
	QualityGood Quality = iota
	QualityUncertain
	QualityBad
)

// PointDefinition describes a point's static information.
type PointDefinition struct {
	ID          string
	Key         string
	Name        string
	Description string
	Unit        string
	DataType    DataType
	Readable    bool
	Writable    bool
}

// PointState is the runtime value and status of a point.
type PointState struct {
	Value     any
	Quality   Quality
	Timestamp time.Time
}

// Point is the unified protocol-agnostic model exposed by the runtime.
type Point struct {
	Definition PointDefinition
	State      PointState
}

// NewPoint creates a new Point instance with an initial bad quality state.
func NewPoint(def PointDefinition) *Point {
	return &Point{
		Definition: def,
		State: PointState{
			Quality:   QualityBad,
			Timestamp: time.Now(),
		},
	}
}

// UpdateState updates the runtime state of the point.
func (p *Point) UpdateState(value any, quality Quality, timestamp time.Time) {
	p.State = PointState{
		Value:     value,
		Quality:   quality,
		Timestamp: timestamp,
	}
}

// GetID returns the unique identifier of the point.
func (p *Point) GetID() string {
	return p.Definition.ID
}

// GetKind returns the kind of the object, which is "Point".
func (p *Point) GetKind() object.Kind {
	return KindPoint
}