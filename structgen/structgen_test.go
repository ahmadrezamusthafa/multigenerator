package structgen

import (
	"bytes"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"reflect"
	"testing"
)

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
					Value: "2019-09-01 00:10:00",
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
