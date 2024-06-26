// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitoidp

import (
	"time"
)

const (
	ResNameIdentityProvider  = "Identity Provider"
	ResNameResourceServer    = "Resource Server"
	ResNameRiskConfiguration = "Risk Configuration"
	ResNameUserPoolClient    = "User Pool Client"
	ResNameUserPoolDomain    = "User Pool Domain"
	ResNameUserPool          = "User Pool"
	ResNameUser              = "User"
)

const (
	propagationTimeout = 2 * time.Minute
)
