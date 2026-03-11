package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DeviceRule holds the schema definition for the DeviceRule entity.
type DeviceRule struct {
	ent.Schema
}

// Fields of the DeviceRule.
func (DeviceRule) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("tenant_id", uuid.UUID{}),
		field.String("name").
			NotEmpty(),
		field.String("rule_type").
			NotEmpty(),
		field.JSON("trigger_conditions", map[string]any{}).
			Optional(),
		field.JSON("action", map[string]any{}).
			Optional(),
		field.Bool("is_active").
			Default(true),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the DeviceRule.
func (DeviceRule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("executions", RuleExecution.Type),
	}
}

