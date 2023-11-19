package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LogEntry struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Level      string             `bson:"level"`
	Message    string             `bson:"message"`
	ResourceId string             `bson:"resourceId"`
	Timestamp  time.Time          `bson:"timestamp"`
	TraceId    string             `bson:"traceId"`
	SpanId     string             `bson:"spanId"`
	Commit     string             `bson:"commit"`
	Metadata   Metadata           `bson:"metadata"`
}

type Metadata struct {
	ParentResourceId string `bson:"parentResourceId"`
}
