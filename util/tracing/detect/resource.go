// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package detect

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func detectResource() *resource.Resource {
	res, err := resource.New(context.Background(),
		resource.WithDetectors(serviceNameDetector{}),
		resource.WithFromEnv(),
		resource.WithDetectors(telemetrySDK{}),
	)
	if err != nil {
		otel.Handle(err)
	}
	return res
}

type (
	telemetrySDK struct{}
)

var (
	_ resource.Detector = telemetrySDK{}
)

// Detect returns a *Resource that describes the OpenTelemetry SDK used.
func (telemetrySDK) Detect(context.Context) (*resource.Resource, error) {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.TelemetrySDKName("opentelemetry"),
		semconv.TelemetrySDKLanguageGo,
		// otel.Version instead of sdk.Version: the root sdk package
		// does not exist yet at the otel version pinned by this branch,
		// and its builtin telemetrySDK detector reports otel.Version too.
		semconv.TelemetrySDKVersion(otel.Version()),
	), nil
}
