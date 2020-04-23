package multigenerator

import (
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"testing"
)

//BENCHMARK GenerateCondition
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  301020	      3341 ns/op
//  372768	      2793 ns/op (now)
//------------------------------------
func BenchmarkGenerateCondition(b *testing.B) {
	query := `id=1 && ( division = engineering || division = finance )`
	for n := 0; n < b.N; n++ {
		GenerateCondition(query)
	}
}

//BENCHMARK Validate
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  484363	      2316 ns/op (now)
//------------------------------------
func BenchmarkValidate(b *testing.B) {
	object := struct {
		ID       string `json:"id"`
		MemberID string `json:"member_id"`
		Division string `json:"division"`
	}{
		ID:       "1",
		MemberID: "2",
		Division: "finance",
	}

	query := "(id=1 && (member_id=12||member_id=2))  &&   (division=engineering || division=finance)"
	condition, _ := GenerateCondition(query)
	for n := 0; n < b.N; n++ {
		Validate(condition, object)
	}
}

//BENCHMARK ValidateObjects
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  359617	      2829 ns/op (now)
//------------------------------------
func BenchmarkValidateObjects(b *testing.B) {
	object := struct {
		ID       string `json:"id"`
		MemberID string `json:"member_id"`
		Division string `json:"division"`
	}{
		ID:       "1",
		MemberID: "2",
		Division: "finance",
	}

	query := "(id=1 && (member_id=12||member_id=2))  &&   (division=engineering || division=finance)"
	condition, _ := GenerateCondition(query)
	for n := 0; n < b.N; n++ {
		ValidateObjects(condition, object)
	}
}

//BENCHMARK ValidateCondition
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  1869136	       625 ns/op (now)
//------------------------------------
func BenchmarkValidateCondition(b *testing.B) {
	referenceQuery := `id=1 && ( division = engineering || division = finance )`
	input := `id=1 && division = engineering`
	referenceCondition, _ := GenerateCondition(referenceQuery)
	inputCondition, _ := GenerateCondition(input)

	for n := 0; n < b.N; n++ {
		ValidateCondition(referenceCondition, inputCondition)
	}
}

//BENCHMARK GenerateQuery
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  293460	      3800 ns/op
//  542161        1850 ns/op (now)
//------------------------------------
func BenchmarkGenerateQuery(b *testing.B) {
	type args struct {
		mainQuery     string
		baseCondition types.BaseCondition
	}
	req := struct {
		args args
	}{
		args: args{
			mainQuery: "select * from data_member",
			baseCondition: types.BaseCondition{
				Fields: []string{"memberId", "name", "phoneNumber"},
				Conditions: []*types.Condition{
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
		},
	}

	for n := 0; n < b.N; n++ {
		GenerateQuery(req.args.mainQuery, req.args.baseCondition)
	}
}
