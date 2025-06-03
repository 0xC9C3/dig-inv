package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Tag struct {
	ent.Schema
}

func (Tag) Fields() []ent.Field {
	return withDefaults([]ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Comment("The unique identifier for the tag, which is used to track the tag in the inventory system. This is a UUID that is generated when the tag is created."),
		field.String("name").
			NotEmpty().
			Comment("The name of the tag, which is used to identify the tag in the inventory system. This is a unique identifier for the tag."),
		field.String("description").
			Optional().
			Comment("A description of the tag, which can be used to provide additional information about the tag. This is optional and can be used to provide more context about the tag's purpose or usage."),
	})
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", Item.Type).
			Comment("The items that are associated with this tag. This edge represents the many-to-many relationship between tags and items, allowing multiple items to be associated with a single tag and vice versa."),
	}
}
