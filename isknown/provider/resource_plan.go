package provider

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/schema"
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
	if proposedValue, proposedValueMap, diags, isErroneous = schema.UnmarshalState(proposedValueDynamic, resourceType); isErroneous {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return resp, nil
	}

	// configValueDynamic := req.Config
	// var configValue tftypes.Value
	// var configValueMap map[string]tftypes.Value
	// if configValue, configValueMap, diags, isErroneous = schema.UnmarshalState(configValueDynamic, resourceType); isErroneous {
	// 	resp.Diagnostics = append(resp.Diagnostics, diags...)
	// 	return resp, nil
	// }
	// _ = configValue
	// _ = configValueMap

	if proposedValue.IsNull() {
		// Plan to delete the resource
		resp.PlannedState = proposedValueDynamic
		return resp, nil
	}

	var providerMetaSeedAddition string
	if providerMetaSeedAddition, diags, isWorking = common.TryExtractProviderMetaGuidSeedAddition(req.ProviderMeta); !isWorking {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return resp, nil
	}

	guidSeedValue := proposedValueMap["guid_seed"]
	if !guidSeedValue.IsKnown() {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Current resource state has a 'guid_seed' attribute but it is not known.",
			Detail:   "The 'guid_seed' attribute must be known during the plan phase. See attribute description for more informations.",
		})

		return resp, nil
	}
	var guidSeed string
	_ = guidSeedValue.As(&guidSeed) // Why it should ever fail?

	combinedSeed := s.ProviderConfigSeedAddition + "|" +
		providerMetaSeedAddition + "|" +
		req.TypeName + "|" +
		guidSeed

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
		deterministicFile, err := os.Open(deterministicTempFilePath)
		var readBytes []byte

		if errors.Is(err, os.ErrNotExist) {
			resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "File does not exist",
				Detail: `The file does not exist. This can mean
1. the file got deleted before apply-phase or
2. this plan method got called the third time`,
			})
		} else if err == nil {
			defer deterministicFile.Close() // ignore error intentionally
			readBytes = make([]byte, 1)
			deterministicFile.Read(readBytes)
		}

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
