// Package ischema contains the types for schema 'information_schema'.
package ischema

import "github.com/jennal/xo/examples/pgcatalog/pgtypes"

// Code generated by xo. DO NOT EDIT.

// SQLImplementationInfo represents a row from 'information_schema.sql_implementation_info'.
type SQLImplementationInfo struct {
	Tableoid               pgtypes.Oid            `json:"tableoid"`                 // tableoid
	Cmax                   pgtypes.Cid            `json:"cmax"`                     // cmax
	Xmax                   pgtypes.Xid            `json:"xmax"`                     // xmax
	Cmin                   pgtypes.Cid            `json:"cmin"`                     // cmin
	Xmin                   pgtypes.Xid            `json:"xmin"`                     // xmin
	Ctid                   pgtypes.Tid            `json:"ctid"`                     // ctid
	ImplementationInfoID   pgtypes.CharacterData  `json:"implementation_info_id"`   // implementation_info_id
	ImplementationInfoName pgtypes.CharacterData  `json:"implementation_info_name"` // implementation_info_name
	IntegerValue           pgtypes.CardinalNumber `json:"integer_value"`            // integer_value
	CharacterValue         pgtypes.CharacterData  `json:"character_value"`          // character_value
	Comments               pgtypes.CharacterData  `json:"comments"`                 // comments
}
