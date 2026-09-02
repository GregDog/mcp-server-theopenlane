package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/enums"
)

func writeAnnotations() *mcp.ToolAnnotations {
	destructive := false
	openWorld := true
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    false,
		DestructiveHint: &destructive,
		OpenWorldHint:   &openWorld,
		IdempotentHint:  false,
	}
}

func enumPtr[T ~string](value string) *T {
	if value == "" {
		return nil
	}
	v := T(value)
	return &v
}

func controlStatus(value string) *enums.ControlStatus {
	return enumPtr[enums.ControlStatus](value)
}

func documentStatus(value string) *enums.DocumentStatus {
	return enumPtr[enums.DocumentStatus](value)
}

func riskStatus(value string) *enums.RiskStatus {
	return enumPtr[enums.RiskStatus](value)
}

func riskImpact(value string) *enums.RiskImpact {
	return enumPtr[enums.RiskImpact](value)
}

func riskLikelihood(value string) *enums.RiskLikelihood {
	return enumPtr[enums.RiskLikelihood](value)
}

func deleteAnnotations() *mcp.ToolAnnotations {
	destructive := true
	openWorld := true
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    false,
		DestructiveHint: &destructive,
		OpenWorldHint:   &openWorld,
		IdempotentHint:  true,
	}
}

func taskStatus(value string) *enums.TaskStatus {
	return enumPtr[enums.TaskStatus](value)
}
