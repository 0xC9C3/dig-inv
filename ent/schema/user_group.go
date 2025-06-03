package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// A UserGroup represents a group of users in the inventory system. Items can be assigned to user groups,
// allowing for better organization and management of items based on user roles or departments. Groups are assigned
// using OIDC (OpenID Connect) scopes, which allows for integration with external identity providers.

type UserGroup struct {
	ent.Schema
}

func (UserGroup) Fields() []ent.Field {
	return withDefaults([]ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable().
			Comment("The unique identifier for the user group, which is used to track the group in the inventory system. This is a UUID that is generated when the group is created."),
		field.String("name").
			NotEmpty().
			Comment("The name of the user group, which is used to identify the group in the inventory system. This is a unique identifier for the group."),
		field.String("description").
			Optional().
			Comment("A description of the user group, which can be used to provide additional information about the group. This is optional and can be used to provide more context about the group's purpose or usage."),
		field.String("oidc_scope").
			NotEmpty().
			Comment("The OIDC scope associated with the user group, which is used to integrate with external identity providers. This allows for the assignment of items to user groups based on OIDC scopes, enabling better organization and management of items based on user roles or departments."),
	})
}

func (UserGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", Item.Type).
			Comment("The items that are assigned to this user group. This edge represents the many-to-many relationship between user groups and items, allowing multiple items to be associated with a single user group and vice versa."),
	}
}
