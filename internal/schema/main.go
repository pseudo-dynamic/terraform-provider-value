package schema

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// GetObjectTypeFromSchema returns a tftypes.Type that can wholy represent the schema input
// TODO: Outsource this and reduce redudancy across terraform-plugin-go providers
func GetObjectTypeFromSchema(schema *tfprotov6.Schema) tftypes.Type {
	bm := map[string]tftypes.Type{}

	for _, att := range schema.Block.Attributes {
		bm[att.Name] = att.Type
	}

	for _, b := range schema.Block.BlockTypes {
		a := map[string]tftypes.Type{}
		for _, att := range b.Block.Attributes {
			a[att.Name] = att.Type
		}
		bm[b.TypeName] = tftypes.List{
			ElementType: tftypes.Object{AttributeTypes: a},
		}

		// FIXME we can make this function recursive to handle
		// n levels of nested blocks
		for _, bb := range b.Block.BlockTypes {
			aa := map[string]tftypes.Type{}
			for _, att := range bb.Block.Attributes {
				aa[att.Name] = att.Type
			}
			a[bb.TypeName] = tftypes.List{
				ElementType: tftypes.Object{AttributeTypes: aa},
			}
		}
	}

	return tftypes.Object{AttributeTypes: bm}
}

func UnmarshalState(state *tfprotov6.DynamicValue, stateType tftypes.Type) (tftypes.Value, map[string]tftypes.Value, []*tfprotov6.Diagnostic, bool) {
	diags := []*tfprotov6.Diagnostic{}
	valueMap := make(map[string]tftypes.Value)
	value, err := state.Unmarshal(stateType)

	if err != nil {
		diags = append(diags, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to unmarshal state",
			Detail:   err.Error(),
		})

		goto End
	}

	err = value.As(&valueMap)

	if err != nil {
		diags = append(diags, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to extract state from tftypes.Value",
			Detail:   err.Error(),
		})

		goto End
	}

End:
	return value, valueMap, diags, len(diags) != 0
}
