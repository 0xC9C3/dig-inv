package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"slices"
	"time"
)

func withDefaults(pre []ent.Field) []ent.Field {
	return slices.Concat(pre, []ent.Field{
		field.String("created_by").
			Comment("The user who created the resource in the inventory system. This is used for auditing purposes and to track who added the resource."),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When the resource was created in the inventory system."),
		field.String("updated_by").
			Comment("The user who last updated the resource in the inventory system. This is used for auditing purposes and to track who modified the resource."),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The time when the resource was last updated in the inventory system."),
		field.String("deleted_by").
			Optional().
			Comment("The user who deleted the resource from the inventory system. This is used for auditing purposes and to track who removed the resource."),
		field.Time("deleted_at").
			Optional().
			Nillable().
			Comment("The time when the resource was deleted from the inventory system."),
	})
}
