// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdkv2

import (
	"testing"
)

func TestSimpleSchemaSetFunc(t *testing.T) {
	t.Parallel()

	v1 := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": true,
	}
	v2 := map[string]interface{}{
		"key1": "value1",
		"key2": "value2-new",
		"key3": 3,
		"key4": true,
	}
	v3 := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": 4,
		"key4": true,
	}
	f := SimpleSchemaSetFunc("key1", "key3", "key4")

	if f(v1) != f(v2) {
		t.Errorf("expected equal")
	}
	if f(v1) == f(v3) {
		t.Errorf("expected not equal")
	}
}
