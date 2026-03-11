package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TelemetryData holds the schema definition for the TelemetryData entity.
type TelemetryData struct {
	ent.Schema
}

// Fields of the TelemetryData.
func (TelemetryData) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("tenant_id", uuid.UUID{}),
		field.UUID("device_id", uuid.UUID{}),
		field.String("metric_name").
			NotEmpty(),
		field.Float("metric_value"),
		field.String("unit").
			Optional(),
		field.Time("timestamp").
			Default(time.Now),
		field.String("geo_point").
			Optional(),
		field.JSON("metadata", map[string]any{}).
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the TelemetryData.
func (TelemetryData) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("device", Device.Type).
			Ref("telemetry_data").
			Field("device_id").
			Unique().
			Required(),
	}
}

