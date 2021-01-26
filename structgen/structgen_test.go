package structgen

import (
	"bytes"
	"encoding/json"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"reflect"
	"strings"
	"testing"
)

func TestGenerateConditionQueryStructure(t *testing.T) {
	type args struct {
		query string
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
				query: `
                      (id = 1 
                      && member_id = 2 )
                      || (
                        division = engineering 
                        || division = finance
                      )
`,
			},
			want:    `{"conditions":[{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance","type":"alphanumeric"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                      (id = 1 
                      && member_id = 2 )
                      || (
                        division = engineering 
                        || division = finance
                      ) && user_id = 43
`,
			},
			want:    `{"conditions":[{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance","type":"alphanumeric"}}]},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"43","type":"numeric"}}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                      id = 1 
                      && member_id = 2 
                      && (
                        division = engineering 
                        || division = finance
                      )
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}},{"operator":"AND","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance","type":"alphanumeric"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `id=1 &&  member_id=2   &&   (division=engineering || division=finance)`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}},{"operator":"AND","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance","type":"alphanumeric"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                  id = 1 
                  && member_id = 2 
                  && user_id = 3 
                  && (
                    province = jatim 
                    || city = mojokerto
                    || (
                      warehouse_id = 1 
                      && warehouse_detail_id = 2
                    )
                  )
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"3","type":"numeric"}},{"operator":"AND","conditions":[{"attribute":{"name":"province","operator":"=","value":"jatim","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"city","operator":"=","value":"mojokerto","type":"alphanumeric"}},{"operator":"OR","conditions":[{"attribute":{"name":"warehouse_id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"warehouse_detail_id","operator":"=","value":"2","type":"numeric"}}]}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                  id = 1
                  && member_id = 2
                  && user_id = 3
                  && (
                    province = jatim
                    || city = mojokerto
                    || (
                      warehouse_id = 1
                      && warehouse_detail_id = 2
                    )
                  )
				  && data_id = 54
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2","type":"numeric"}},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"3","type":"numeric"}},{"operator":"AND","conditions":[{"attribute":{"name":"province","operator":"=","value":"jatim","type":"alphanumeric"}},{"operator":"OR","attribute":{"name":"city","operator":"=","value":"mojokerto","type":"alphanumeric"}},{"operator":"OR","conditions":[{"attribute":{"name":"warehouse_id","operator":"=","value":"1","type":"numeric"}},{"operator":"AND","attribute":{"name":"warehouse_detail_id","operator":"=","value":"2","type":"numeric"}}]}]},{"operator":"AND","attribute":{"name":"data_id","operator":"=","value":"54","type":"numeric"}}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: "((date<=2019-09-09 && date > 2019-08-08) || (p_date>=2019-01-01 && p_date<2019-02-02)) && (member_type=1||member_type=2)",
			},
			want:    `{"conditions":[{"conditions":[{"conditions":[{"attribute":{"name":"date","operator":"\u003c=","value":"2019-09-09","type":"alphanumeric"}},{"operator":"AND","attribute":{"name":"date","operator":"\u003e","value":"2019-08-08","type":"alphanumeric"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"p_date","operator":"\u003e=","value":"2019-01-01","type":"alphanumeric"}},{"operator":"AND","attribute":{"name":"p_date","operator":"\u003c","value":"2019-02-02","type":"alphanumeric"}}]}]},{"operator":"AND","conditions":[{"attribute":{"name":"member_type","operator":"=","value":"1","type":"numeric"}},{"operator":"OR","attribute":{"name":"member_type","operator":"=","value":"2","type":"numeric"}}]}]}`,
			wantErr: false,
		},
	}
	s := StructGen{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GenerateCondition(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			byteBuf, _ := json.Marshal(got)
			if !strings.EqualFold(string(byteBuf), tt.want) {
				t.Errorf("GenerateCondition() = %v, want %v", string(byteBuf), tt.want)
			}
		})
	}
}

func Test_getToken(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want []*types.TokenAttribute
	}{
		{
			name: "Normal case",
			args: args{
				value: "id=1 &&  member_id=2   &&   (division=engineering      || division=finance)",
			},
			want: []*types.TokenAttribute{
				{
					Value: "id",
				},
				{
					Value: "=",
				},
				{
					Value: "1",
				},
				{
					Value: "&&",
				},
				{
					Value: "member_id",
				},
				{
					Value: "=",
				},
				{
					Value: "2",
				},
				{
					Value: "&&",
				},
				{
					Value: "(",
				},
				{
					Value: "division",
				},
				{
					Value: "=",
				},
				{
					Value: "engineering",
				},
				{
					Value: "||",
				},
				{
					Value: "division",
				},
				{
					Value: "=",
				},
				{
					Value: "finance",
				},
				{
					Value: ")",
				},
			},
		},
		{
			name: "Normal case",
			args: args{
				value: "id>1 &&  member_id>=2 && (test_id<10 || pr_id<=28)",
			},
			want: []*types.TokenAttribute{
				{
					Value: "id",
				},
				{
					Value: ">",
				},
				{
					Value: "1",
				},
				{
					Value: "&&",
				},
				{
					Value: "member_id",
				},
				{
					Value: ">=",
				},
				{
					Value: "2",
				},
				{
					Value: "&&",
				},
				{
					Value: "(",
				},
				{
					Value: "test_id",
				},
				{
					Value: "<",
				},
				{
					Value: "10",
				},
				{
					Value: "||",
				},
				{
					Value: "pr_id",
				},
				{
					Value: "<=",
				},
				{
					Value: "28",
				},
				{
					Value: ")",
				},
			},
		},
		{
			name: "Normal case",
			args: args{
				value: "date>2019-09-01 && date<=2019-10-10 && (segment_id=12||segment_id=13)",
			},
			want: []*types.TokenAttribute{
				{
					Value: "date",
				},
				{
					Value: ">",
				},
				{
					Value: "2019-09-01",
				},
				{
					Value: "&&",
				},
				{
					Value: "date",
				},
				{
					Value: "<=",
				},
				{
					Value: "2019-10-10",
				},
				{
					Value: "&&",
				},
				{
					Value: "(",
				},
				{
					Value: "segment_id",
				},
				{
					Value: "=",
				},
				{
					Value: "12",
				},
				{
					Value: "||",
				},
				{
					Value: "segment_id",
				},
				{
					Value: "=",
				},
				{
					Value: "13",
				},
				{
					Value: ")",
				},
			},
		},
		{
			name: "Normal case",
			args: args{
				value: `date>"2019-09-01 00:10:00" && date<=2019-10-10 && (segment_id="12"||segment_id=13)`,
			},
			want: []*types.TokenAttribute{
				{
					Value: "date",
				},
				{
					Value: ">",
				},
				{
					Value:          "2019-09-01 00:10:00",
					IsAlphanumeric: true,
				},
				{
					Value: "&&",
				},
				{
					Value: "date",
				},
				{
					Value: "<=",
				},
				{
					Value: "2019-10-10",
				},
				{
					Value: "&&",
				},
				{
					Value: "(",
				},
				{
					Value: "segment_id",
				},
				{
					Value: "=",
				},
				{
					Value:          "12",
					IsAlphanumeric: true,
				},
				{
					Value: "||",
				},
				{
					Value: "segment_id",
				},
				{
					Value: "=",
				},
				{
					Value: "13",
				},
				{
					Value: ")",
				},
			},
		},
		{
			name: "Nil case",
			args: args{
				value: ` `,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTokenAttributes(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				strbGot := bytes.Buffer{}
				for _, g := range got {
					strbGot.WriteString("\"" + g.Value + "\" ")
				}
				strbWant := bytes.Buffer{}
				for _, g := range tt.want {
					strbWant.WriteString("\"" + g.Value + "\" ")
				}
				t.Errorf("getTokenAttributes() = %v, want %v", strbGot.String(), strbWant.String())
			}
		})
	}
}
