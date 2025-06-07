package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// An Item holds a single asset, which can be a physical item, a digital asset, or any other entity that can be tracked.
// providers are meant to save asset specific information, such as the location, condition, or any other relevant metadata in
// their own table.
// The Item schema is the core of the inventory system, representing the items that are being tracked and holds
// the basic information about each item as well as the type (provider) of the item.

type Item struct {
	ent.Schema
}

func (Item) Fields() []ent.Field {
	return withDefaults([]ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Comment("The unique identifier for the item, which is used to track the item in the inventory system. This is a UUID that is generated when the item is created."),
		field.String("name").
			NotEmpty().
			Comment("The name of the item, which is used to identify it in the inventory system."),
		field.String("description").
			Optional().
			Comment("A description of the item, which can be used to provide additional information about the item."),
	})
}

func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tags", Tag.Type).
			Comment("The tags that are associated with this item. This edge represents the many-to-many relationship between items and tags, allowing multiple tags to be associated with a single item and vice versa."),
		edge.To("user_groups", UserGroup.Type).
			Comment("The user groups that are associated with this item. This edge represents the many-to-many relationship between items and user groups, allowing multiple user groups to be associated with a single item and vice versa."),
		edge.To("asset_class", AssetClass.Type).
			Required().
			Comment("The asset class that this item belongs to. This edge represents the many-to-one relationship between items and asset classes, allowing multiple items to be associated with a single asset class. The asset class is defined in the AssetClass schema."),
	}
}
