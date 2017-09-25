// Package ischema contains the types for schema 'information_schema'.
package ischema

import "github.com/jennal/xo/examples/pgcatalog/pgtypes"

// Code generated by xo. DO NOT EDIT.

// InformationSchemaCatalogName represents a row from 'information_schema.information_schema_catalog_name'.
type InformationSchemaCatalogName struct {
	CatalogName pgtypes.SQLIdentifier `json:"catalog_name"` // catalog_name
}
