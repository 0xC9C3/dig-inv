package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// An asset class represents a category or type of asset that can be tracked in the inventory system.
// Asset classes are used to group similar assets together, allowing for easier management and reporting.

type AssetClass struct {
	ent.Schema
}

func (AssetClass) Fields() []ent.Field {
	return withDefaults([]ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Comment("The unique identifier for the asset class."),
		field.String("name").
			NotEmpty().
			Comment("The name of the asset class, which is used to identify it in the inventory system."),
		field.String("description").
			Optional().
			Comment("A description of the asset class, which can be used to provide additional information about the asset class."),
		field.String("icon").
			Optional().
			Comment("An icon representing the asset class, which can be used in the user interface to visually distinguish different asset classes."),
		field.String("color").
			Optional().
			Comment("A color associated with the asset class, which can be used in the user interface to visually distinguish different asset classes."),
		field.String("provider").
			Optional().
			Comment("The provider of the asset class, which is used to identify the provider of the asset class. This is used to determine which provider's table to use for storing additional information about the asset class."),
	})
}

func (AssetClass) Edges() []ent.Edge {
	return []ent.Edge{}
}
