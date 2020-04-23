package querygen

import (
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"testing"
)

//BENCHMARK GenerateWhereParameter
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  945962	      1258 ns/op
//  1117884	       991 ns/op (now)
//------------------------------------
func BenchmarkGenerateWhereParameter(b *testing.B) {
	type args struct {
		condition []*types.Condition
	}
	req := struct {
		args args
	}{
		args: args{
			condition: []*types.Condition{
				{
					Operator: "AND",
					Conditions: []*types.Condition{
						{
							Attribute: &types.Attribute{
								Name:     "id",
								Operator: "=",
								Value:    "1",
								Type:     valuetype.Numeric,
							},
						},
						{
							Attribute: &types.Attribute{
								Name:     "member_id",
								Operator: "=",
								Value:    "2",
								Type:     valuetype.Numeric,
							},
						},
						{
							Operator: "OR",
							Conditions: []*types.Condition{
								{
									Attribute: &types.Attribute{
										Name:     "division",
										Operator: "=",
										Value:    "engineering",
										Type:     valuetype.Alphanumeric,
									},
								},
								{
									Attribute: &types.Attribute{
										Name:     "division",
										Operator: "=",
										Value:    "finance",
										Type:     valuetype.Alphanumeric,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for n := 0; n < b.N; n++ {
		generateWhereParameter(req.args.condition)
	}
}
