package querygen

import (
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"regexp"
	"strings"
	"testing"
)

func Test_generateWhereParameter(t *testing.T) {
	type args struct {
		condition []*types.Condition
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal case",
			args: args{
				condition: []*types.Condition{
					{
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
								Operator: "AND",
								Attribute: &types.Attribute{
									Name:     "member_id",
									Operator: "=",
									Value:    "2",
									Type:     valuetype.Numeric,
								},
							},
							{
								Operator: "AND",
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
										Operator: "OR",
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
			want: `
                WHERE 
                  id = 1 
                  AND member_id = 2 
                  AND (
                    division = 'engineering' 
                    OR division = 'finance'
                  )
`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				condition: []*types.Condition{
					{
						Conditions: []*types.Condition{
							{
								Attribute: &types.Attribute{
									Name:     "id",
									Operator: "IN",
									Value:    "1,2,3,4,5",
									Type:     valuetype.Numeric,
								},
							},
							{
								Operator: "AND",
								Attribute: &types.Attribute{
									Name:     "initial",
									Operator: "IN",
									Value:    "ABC,DEF,KLM",
									Type:     valuetype.Alphanumeric,
								},
							},
						},
					},
				},
			},
			want: `
                WHERE 
                  id IN (1,2,3,4,5) 
                  AND initial IN ('ABC','DEF','KLM')
`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				condition: []*types.Condition{
					{
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
								Operator: "OR",
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
										Operator: "OR",
										Attribute: &types.Attribute{
											Name:     "division",
											Operator: "=",
											Value:    "finance",
											Type:     valuetype.Alphanumeric,
										},
									},
									{
										Operator: "OR",
										Conditions: []*types.Condition{
											{
												Attribute: &types.Attribute{
													Name:     "type",
													Operator: "=",
													Value:    "1",
													Type:     valuetype.Numeric,
												},
											},
											{
												Operator: "AND",
												Attribute: &types.Attribute{
													Name:     "type",
													Operator: "=",
													Value:    "2",
													Type:     valuetype.Numeric,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: `
                WHERE 
                  id = 1 
                  OR member_id = 2 
                  OR (
                    division = 'engineering' 
                    OR division = 'finance' 
                    OR (
                      type = 1 
                      AND type = 2
                    )
                  )
`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				condition: []*types.Condition{
					{
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
								Operator: "AND",
								Attribute: &types.Attribute{
									Name:     "member_id",
									Operator: "=",
									Value:    "2",
									Type:     valuetype.Numeric,
								},
							},
							{
								Operator: "AND",
								Attribute: &types.Attribute{
									Name:     "user_id",
									Operator: "=",
									Value:    "3",
									Type:     valuetype.Numeric,
								},
							},
							{
								Operator: "AND",
								Conditions: []*types.Condition{
									{
										Attribute: &types.Attribute{
											Name:     "province",
											Operator: "=",
											Value:    "jatim",
											Type:     valuetype.Alphanumeric,
										},
									},
									{
										Operator: "OR",
										Attribute: &types.Attribute{
											Name:     "city",
											Operator: "=",
											Value:    "mojokerto",
											Type:     valuetype.Alphanumeric,
										},
									},
									{
										Operator: "OR",
										Conditions: []*types.Condition{
											{
												Attribute: &types.Attribute{
													Name:     "warehouse_id",
													Operator: "=",
													Value:    "1",
													Type:     valuetype.Numeric,
												},
											},
											{
												Operator: "AND",
												Attribute: &types.Attribute{
													Name:     "warehouse_detail_id",
													Operator: "IN",
													Value:    "22,32,45",
													Type:     valuetype.Numeric,
												},
											},
										},
									},
								},
							},
							{
								Operator: "AND",
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
										Operator: "OR",
										Attribute: &types.Attribute{
											Name:     "division",
											Operator: "=",
											Value:    "finance",
											Type:     valuetype.Alphanumeric,
										},
									},
									{
										Operator: "OR",
										Conditions: []*types.Condition{
											{
												Attribute: &types.Attribute{
													Name:     "level",
													Operator: "=",
													Value:    "1",
													Type:     valuetype.Numeric,
												},
											},
											{
												Operator: "OR",
												Attribute: &types.Attribute{
													Name:     "level",
													Operator: "=",
													Value:    "2",
													Type:     valuetype.Numeric,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: `
                WHERE 
                  id = 1 
                  AND member_id = 2 
                  AND user_id = 3 
                  AND (
                    province = 'jatim' 
                    OR city = 'mojokerto' 
                    OR (
                      warehouse_id = 1 
                      AND warehouse_detail_id IN (22,32,45)
                    )
                  ) 
                  AND (
                    division = 'engineering' 
                    OR division = 'finance' 
                    OR (
                      level = 1 
                      OR level = 2
                    )
                  )
`,
			wantErr: false,
		},
	}
	var rgx = regexp.MustCompile(`[\s]{2,}`)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateWhereParameter(tt.args.condition)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateWhereParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			strGot := strings.TrimSpace(rgx.ReplaceAllString(got, " "))
			strWant := strings.TrimSpace(rgx.ReplaceAllString(tt.want, " "))

			if !strings.EqualFold(strGot, strWant) {
				t.Errorf("generateWhereParameter() got = %v, want %v", strGot, strWant)
			}
		})
	}
}
