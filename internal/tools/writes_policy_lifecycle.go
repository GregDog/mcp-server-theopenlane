package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type policyLifecycleInput struct {
	ID      string `json:"id" jsonschema:"InternalPolicy ID."`
	Confirm bool   `json:"confirm" jsonschema:"Must be true to persist. If false, returns a preview only."`
}

type policyLifecycleResult struct {
	Confirmed        bool   `json:"confirmed"`
	Error            string `json:"error,omitempty"`
	ID               string `json:"id"`
	DisplayID        string `json:"display_id,omitempty"`
	Name             string `json:"name,omitempty"`
	ApprovalRequired bool   `json:"approval_required"`
	ApproverID       string `json:"approver_id,omitempty"`
	DelegateID       string `json:"delegate_id,omitempty"`
	CurrentStatus    string `json:"current_status"`
	RequestedStatus  string `json:"requested_status"`
	ResultStatus     string `json:"result_status,omitempty"`
	Summary          string `json:"summary,omitempty"`
}

func registerPolicyLifecycle(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_policy_submit_for_approval",
		Title:       "Submit an Openlane policy for native approval",
		Description: "Set InternalPolicy status to NEEDS_APPROVAL via updateInternalPolicy. This is the native approval path (Approver/Delegate groups), not a WorkflowDefinition. Does not create a WorkflowInstance or WorkflowAssignment; use openlane_workflow_assignment_approve only when a real assignment exists. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.submitPolicyForApproval)

	addTool(server, &mcp.Tool{
		Name:        "openlane_policy_approve",
		Title:       "Approve an Openlane policy (native)",
		Description: "Approve a native InternalPolicy by setting status to APPROVED via updateInternalPolicy. Use only when status is NEEDS_APPROVAL and there is no WorkflowAssignment for this policy. If openlane_policies_awaiting_approval returns path workflow, use openlane_workflow_assignment_approve instead. Openlane HookStatusApproval allows only Approver or Delegate group members to approve. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.approvePolicy)

	addTool(server, &mcp.Tool{
		Name:        "openlane_policy_publish",
		Title:       "Publish an Openlane policy",
		Description: "Set InternalPolicy status to PUBLISHED via updateInternalPolicy. From APPROVED, or from DRAFT when approvalRequired is false. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.publishPolicy)

	addTool(server, &mcp.Tool{
		Name:        "openlane_policy_return_to_draft",
		Title:       "Return an Openlane policy to draft",
		Description: "Set InternalPolicy status to DRAFT via updateInternalPolicy from NEEDS_APPROVAL or APPROVED. This is the native way to send a policy back; there is no rejectInternalPolicy mutation. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.returnPolicyToDraft)
}

func (h *handlers) submitPolicyForApproval(ctx context.Context, _ *mcp.CallToolRequest, in policyLifecycleInput) (*mcp.CallToolResult, policyLifecycleResult, error) {
	return h.transitionPolicy(ctx, in, string(enums.DocumentNeedsApproval), "submit for approval")
}

func (h *handlers) approvePolicy(ctx context.Context, _ *mcp.CallToolRequest, in policyLifecycleInput) (*mcp.CallToolResult, policyLifecycleResult, error) {
	return h.transitionPolicy(ctx, in, string(enums.DocumentApproved), "approve")
}

func (h *handlers) publishPolicy(ctx context.Context, _ *mcp.CallToolRequest, in policyLifecycleInput) (*mcp.CallToolResult, policyLifecycleResult, error) {
	return h.transitionPolicy(ctx, in, string(enums.DocumentPublished), "publish")
}

func (h *handlers) returnPolicyToDraft(ctx context.Context, _ *mcp.CallToolRequest, in policyLifecycleInput) (*mcp.CallToolResult, policyLifecycleResult, error) {
	return h.transitionPolicy(ctx, in, string(enums.DocumentDraft), "return to draft")
}

func (h *handlers) transitionPolicy(ctx context.Context, in policyLifecycleInput, target, verb string) (*mcp.CallToolResult, policyLifecycleResult, error) {
	if in.ID == "" {
		return nil, policyLifecycleResult{}, errIDRequired
	}
	resp, err := h.api.GetInternalPolicyByID(ctx, in.ID)
	if err != nil {
		return nil, policyLifecycleResult{}, openlane.APIError(err)
	}
	p := resp.InternalPolicy
	current := openlane.Format(p.Status)
	approvalRequired := openlane.Deref(p.ApprovalRequired)
	if err := validatePolicyTransition(current, target, approvalRequired); err != nil {
		return nil, policyLifecycleResult{}, err
	}
	out := policyLifecycleResult{
		ID:               p.ID,
		DisplayID:        p.DisplayID,
		Name:             p.Name,
		ApprovalRequired: approvalRequired,
		ApproverID:       openlane.Deref(p.ApproverID),
		DelegateID:       openlane.Deref(p.DelegateID),
		CurrentStatus:    current,
		RequestedStatus:  target,
		Summary:          fmt.Sprintf("Would %s %q from %s to %s", verb, p.Name, current, target),
	}
	if !in.Confirm {
		out.Error = errConfirmationRequired
		return nil, out, nil
	}

	status := enums.DocumentStatus(target)
	if _, err := h.api.UpdateInternalPolicy(ctx, in.ID, graphclient.UpdateInternalPolicyInput{Status: &status}); err != nil {
		return nil, policyLifecycleResult{}, openlane.APIError(err)
	}
	got, err := h.api.GetInternalPolicyByID(ctx, in.ID)
	if err != nil {
		return nil, policyLifecycleResult{}, openlane.APIError(err)
	}
	result := openlane.Format(got.InternalPolicy.Status)
	if result != target {
		return nil, policyLifecycleResult{}, fmt.Errorf("policy persisted but status is %s, expected %s", result, target)
	}
	out.Confirmed = true
	out.ResultStatus = result
	out.CurrentStatus = result
	out.Summary = fmt.Sprintf("%s %q: %s → %s", strings.ToUpper(verb[:1])+verb[1:], p.Name, current, result)
	return nil, out, nil
}

func validatePolicyTransition(current, target string, approvalRequired bool) error {
	current = strings.ToUpper(strings.TrimSpace(current))
	target = strings.ToUpper(strings.TrimSpace(target))
	if current == "" {
		return fmt.Errorf("policy has no status")
	}
	if current == string(enums.DocumentArchived) {
		return fmt.Errorf("cannot change status of an ARCHIVED policy")
	}
	switch target {
	case string(enums.DocumentNeedsApproval):
		if current == string(enums.DocumentNeedsApproval) {
			return fmt.Errorf("policy is already NEEDS_APPROVAL")
		}
		if current == string(enums.DocumentArchived) {
			return fmt.Errorf("cannot submit an ARCHIVED policy for approval")
		}
		return nil
	case string(enums.DocumentApproved):
		if current != string(enums.DocumentNeedsApproval) {
			return fmt.Errorf("native policy approval requires status NEEDS_APPROVAL, current is %s", current)
		}
		return nil
	case string(enums.DocumentPublished):
		if current == string(enums.DocumentApproved) {
			return nil
		}
		if !approvalRequired && current == string(enums.DocumentDraft) {
			return nil
		}
		return fmt.Errorf("publish requires APPROVED, or DRAFT when approvalRequired is false; current is %s", current)
	case string(enums.DocumentDraft):
		if current == string(enums.DocumentNeedsApproval) || current == string(enums.DocumentApproved) {
			return nil
		}
		return fmt.Errorf("return to draft is allowed from NEEDS_APPROVAL or APPROVED, current is %s", current)
	default:
		return fmt.Errorf("unsupported status transition to %s", target)
	}
}
