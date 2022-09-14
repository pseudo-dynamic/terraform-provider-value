package provider

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/uuid"

	"github.com/pseudo-dynamic/terraform-provider-value/isknown/common"
)

// PlanResourceChange function
func (s *UserProviderServer) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	var isErroneous bool
	var isWorking bool
	var diags []*tfprotov6.Diagnostic

	resp := &tfprotov6.PlanResourceChangeResponse{}
	execDiag := s.canExecute()

	if len(execDiag) > 0 {
		resp.Diagnostics = append(resp.Diagnostics, execDiag...)
		return resp, nil
	}

	resourceType := getResourceType(req.TypeName)

	proposedValueDynamic := req.ProposedNewState
	var proposedValue tftypes.Value
	var proposedValueMap map[string]tftypes.Value
	if proposedValue, proposedValueMap, diags, isErroneous = common.UnmarshalState(proposedValueDynamic, resourceType); isErroneous {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return resp, nil
	}
	_ = proposedValueMap

	configValueDynamic := req.Config
	var configValue tftypes.Value
	var configValueMap map[string]tftypes.Value
	if configValue, configValueMap, diags, isErroneous = common.UnmarshalState(configValueDynamic, resourceType); isErroneous {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return resp, nil
	}
	_ = configValue
	_ = configValueMap

	if proposedValue.IsNull() {
		// Plan to delete the resource
		resp.PlannedState = proposedValueDynamic
		return resp, nil
	}

	var seedPrefix string
	if seedPrefix, diags, isWorking = common.CanGetSeedPrefix(req.ProviderMeta); !isWorking {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return resp, nil
	}
	_ = seedPrefix

	uniqueSeedValue := proposedValueMap["unique_seed"]
	if !uniqueSeedValue.IsKnown() {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Current state of resource has a 'seed' attribute but it is not known.",
			Detail:   "The 'seed' attribute must be known during the plan phase. See attribute description for more informations.",
		})

		return resp, nil
	}
	var uniqueSeed string
	_ = uniqueSeedValue.As(&uniqueSeed) // Why it should ever fail?

	combinedSeed := seedPrefix + uniqueSeed
	isPlanPhase := !proposedValueMap["proposed_unknown"].IsKnown() // Unknown == plan

	providerTempDir := filepath.Join(os.TempDir(), "tf-provider-"+req.TypeName)
	if _, err := os.Stat(providerTempDir); os.IsNotExist(err) {
		os.MkdirAll(providerTempDir, 0700) // Create your file
	}

	deterministicFileName := uuid.DeterministicUuidFromString(combinedSeed).String()
	deterministicTempFilePath := filepath.Join(providerTempDir, deterministicFileName)
	var isValueKnown bool

	if isPlanPhase {
		// This is the plan phase
		if s.params.CheckFullyKnown {
			isValueKnown = proposedValueMap["value"].IsFullyKnown()
		} else {
			isValueKnown = proposedValueMap["value"].IsKnown()
		}

		var isValueFullyKnownByte byte
		if isValueKnown {
			isValueFullyKnownByte = 1
		}

		deterministicFile, err := os.OpenFile(deterministicTempFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)

		if err == nil {
			defer deterministicFile.Close()
			_, err = deterministicFile.Write([]byte{isValueFullyKnownByte})
		}

		if err != nil {
			resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error while working with file in the plan phase",
				Detail:   err.Error(),
			})

			return resp, nil
		}
	} else {
		readBytes, err := os.ReadFile(deterministicTempFilePath)

		if err != nil {
			resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error while working with the file in the apply phase",
				Detail:   err.Error(),
			})

			return resp, nil
		}

		readByte := readBytes[0]
		isValueKnown = readByte == 1
		os.Remove(deterministicTempFilePath)
	}

	proposedValueMap["result"] = tftypes.NewValue(tftypes.Bool, isValueKnown)
	plannedValue := tftypes.NewValue(proposedValue.Type(), proposedValueMap)
	plannedState, err := tfprotov6.NewDynamicValue(resourceType, plannedValue)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to assemble proposed state during plan",
			Detail:   err.Error(),
		})
	}

	// pass-through
	// resp.PlannedState = &newValueDynamic
	resp.PlannedState = &plannedState
	return resp, nil
}
