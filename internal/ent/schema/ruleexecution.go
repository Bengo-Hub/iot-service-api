package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// RuleExecution holds the schema definition for the RuleExecution entity.
type RuleExecution struct {
	ent.Schema
}

// Fields of the RuleExecution.
func (RuleExecution) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("rule_id", uuid.UUID{}),
		field.UUID("device_id", uuid.UUID{}).
			Optional(),
		field.Time("triggered_at").
			Default(time.Now),
		field.String("execution_status").
			Default("pending"),
		field.JSON("result", map[string]any{}).
			Optional(),
		field.String("error_message").
			Optional(),
	}
}

// Edges of the RuleExecution.
func (RuleExecution) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rule", DeviceRule.Type).
			Ref("executions").
			Field("rule_id").
			Unique().
			Required(),
	}
}

