// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package flex

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ResourcePrefixCtxKey string

const (
	ResourcePrefix        ResourcePrefixCtxKey = "RESOURCE_PREFIX"
	ResourcePrefixRecurse ResourcePrefixCtxKey = "RESOURCE_PREFIX_RECURSE"
	MapBlockKey                                = "MapBlockKey"
)

// Expand  = TF -->  AWS
// Flatten = AWS --> TF

// autoFlexer is the interface implemented by an auto-flattener or expander.
type autoFlexer interface {
	convert(context.Context, reflect.Value, reflect.Value) diag.Diagnostics
}

// AutoFlexOptionsFunc is a type alias for an autoFlexer functional option.
type AutoFlexOptionsFunc func(autoFlexer)

// autoFlexConvert converts `from` to `to` using the specified auto-flexer.
func autoFlexConvert(ctx context.Context, from, to any, flexer autoFlexer) diag.Diagnostics {
	var diags diag.Diagnostics

	valFrom, valTo, d := autoFlexValues(ctx, from, to)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	// Top-level struct to struct conversion.
	if valFrom.IsValid() && valTo.IsValid() {
		if typFrom, typTo := valFrom.Type(), valTo.Type(); typFrom.Kind() == reflect.Struct && typTo.Kind() == reflect.Struct {
			diags.Append(autoFlexConvertStruct(ctx, from, to, flexer)...)
			return diags
		}
	}

	// Anything else.
	diags.Append(flexer.convert(ctx, valFrom, valTo)...)
	return diags
}

// autoFlexValues returns the underlying `reflect.Value`s of `from` and `to`.
func autoFlexValues(_ context.Context, from, to any) (reflect.Value, reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	valFrom, valTo := reflect.ValueOf(from), reflect.ValueOf(to)
	if kind := valFrom.Kind(); kind == reflect.Ptr {
		valFrom = valFrom.Elem()
	}
	if kind := valTo.Kind(); kind != reflect.Ptr {
		diags.AddError("AutoFlEx", fmt.Sprintf("target (%T): %s, want pointer", to, kind))
		return reflect.Value{}, reflect.Value{}, diags
	}
	valTo = valTo.Elem()

	return valFrom, valTo, diags
}

var (
	plural = pluralize.NewClient()
)

// autoFlexConvertStruct traverses struct `from` calling `flexer` for each exported field.
func autoFlexConvertStruct(ctx context.Context, from any, to any, flexer autoFlexer) diag.Diagnostics {
	var diags diag.Diagnostics

	valFrom, valTo, d := autoFlexValues(ctx, from, to)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	for i, typFrom := 0, valFrom.Type(); i < typFrom.NumField(); i++ {
		field := typFrom.Field(i)
		if field.PkgPath != "" {
			continue // Skip unexported fields.
		}
		fieldName := field.Name
		if fieldName == "Tags" {
			continue // Resource tags are handled separately.
		}
		if fieldName == MapBlockKey {
			continue
		}

		toFieldVal := findFieldFuzzy(ctx, fieldName, valTo, valFrom)
		if !toFieldVal.IsValid() {
			continue // Corresponding field not found in to.
		}
		if !toFieldVal.CanSet() {
			continue // Corresponding field value can't be changed.
		}

		diags.Append(flexer.convert(ctx, valFrom.Field(i), toFieldVal)...)
		if diags.HasError() {
			diags.AddError("AutoFlEx", fmt.Sprintf("convert (%s)", fieldName))
			return diags
		}
	}

	return diags
}

func findFieldFuzzy(ctx context.Context, fieldNameFrom string, valTo, valFrom reflect.Value) reflect.Value {
	// first precedence is exact match (case sensitive)
	if v := valTo.FieldByName(fieldNameFrom); v.IsValid() {
		return v
	}

	// If a "from" field fuzzy matches a "to" field, we are certain the fuzzy match
	// is NOT correct if "from" also contains a field by the fuzzy matched name.
	// For example, if "from" has "Value" and "Values", "Values" should *never*
	// fuzzy match "Value" in "to" since "from" also has "Value". We check "from"
	// to make sure fuzzy matches are not in "from".

	// second precedence is exact match (case insensitive)
	for i, typTo := 0, valTo.Type(); i < typTo.NumField(); i++ {
		field := typTo.Field(i)
		if field.PkgPath != "" {
			continue // Skip unexported fields.
		}
		fieldNameTo := field.Name
		if fieldNameTo == "Tags" {
			continue // Resource tags are handled separately.
		}
		if v := valTo.FieldByName(fieldNameTo); v.IsValid() && strings.EqualFold(fieldNameFrom, fieldNameTo) && !fieldExistsInStruct(fieldNameTo, valFrom) {
			// probably could assume validity here since reflect gave the field name
			return v
		}
	}

	// third precedence is singular/plural
	if plural.IsSingular(fieldNameFrom) && !fieldExistsInStruct(plural.Plural(fieldNameFrom), valFrom) {
		if v := valTo.FieldByName(plural.Plural(fieldNameFrom)); v.IsValid() {
			return v
		}
	}

	if plural.IsPlural(fieldNameFrom) && !fieldExistsInStruct(plural.Singular(fieldNameFrom), valFrom) {
		if v := valTo.FieldByName(plural.Singular(fieldNameFrom)); v.IsValid() {
			return v
		}
	}

	// fourth precedence is using resource prefix
	if v, ok := ctx.Value(ResourcePrefix).(string); ok && v != "" {
		v = strings.ReplaceAll(v, " ", "")
		if ctx.Value(ResourcePrefixRecurse) == nil {
			// so it will only recurse once
			ctx = context.WithValue(ctx, ResourcePrefixRecurse, true)
			if strings.HasPrefix(fieldNameFrom, v) {
				return findFieldFuzzy(ctx, strings.TrimPrefix(fieldNameFrom, v), valTo, valFrom)
			}
			return findFieldFuzzy(ctx, v+fieldNameFrom, valTo, valFrom)
		}
	}

	// no finds, fuzzy or otherwise - return zero value
	return valTo.FieldByName(fieldNameFrom)
}

func fieldExistsInStruct(field string, str reflect.Value) bool {
	if v := str.FieldByName(field); v.IsValid() {
		return true
	}

	return false
}
