// Package ischema contains the types for schema 'information_schema'.
package ischema

import "github.com/jennal/xo/examples/pgcatalog/pgtypes"

// Code generated by xo. DO NOT EDIT.

// TriggeredUpdateColumn represents a row from 'information_schema.triggered_update_columns'.
type TriggeredUpdateColumn struct {
	TriggerCatalog     pgtypes.SQLIdentifier `json:"trigger_catalog"`      // trigger_catalog
	TriggerSchema      pgtypes.SQLIdentifier `json:"trigger_schema"`       // trigger_schema
	TriggerName        pgtypes.SQLIdentifier `json:"trigger_name"`         // trigger_name
	EventObjectCatalog pgtypes.SQLIdentifier `json:"event_object_catalog"` // event_object_catalog
	EventObjectSchema  pgtypes.SQLIdentifier `json:"event_object_schema"`  // event_object_schema
	EventObjectTable   pgtypes.SQLIdentifier `json:"event_object_table"`   // event_object_table
	EventObjectColumn  pgtypes.SQLIdentifier `json:"event_object_column"`  // event_object_column
}
