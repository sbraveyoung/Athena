package ast

import (
	"testing"
)

func TestAST_Judge(t *testing.T) {
	type fields struct {
		mode MODE
	}
	type args struct {
		args map[string]interface{}
		ops  map[string]interface{}
		rule string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				rule: `true`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				rule: `false`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "3",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
				},
				rule: `${app} == "media_std"`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "4",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"uid": 1,
				},
				rule: `${uid} % 10 < 3`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "5",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
					"uid": 1,
				},
				rule: `${app} == "media_std" && ${uid} % 10 < 3`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "6",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
					"uid": 5,
				},
				rule: `${app} == "media_std" && ${uid} % 10 < 3`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "7",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "std_media",
					"uid": 2,
				},
				rule: `${app} == "media_std" && ${uid} % 10 < 3`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "8",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
					"uid": 5,
				},
				rule: `${app} == "media_std" || ${uid} % 10 < 3`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "9",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "std_media",
					"uid": 3,
				},
				rule: `${app} == "media_std" || ${uid} % 10 <= 3`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "10",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "std_media",
					"uid": 3,
				},
				rule: `false && (${app} == "media_std" || ${uid} % 10 <= 3)`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "11",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"uid": 3,
				},
				ops: map[string]interface{}{
					"cls_whitelist": func(uid interface{}) bool {
						return true
					},
				},
				rule: `cls_whitelist(${uid})`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "12",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"uid": 3,
				},
				ops: map[string]interface{}{
					"cls_whitelist": func(uid interface{}) bool {
						return false
					},
				},
				rule: `cls_whitelist(${uid})`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "13",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `${cv}[2:3]=="7"`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "14",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `${cv}[8:]=="Iphone"`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "15",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `!(${cv}[8:]=="Iphone")`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "16",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `${cv}[8:100]=="Iphone"`,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "17",
			fields: fields{
				mode: COMPATIBLE,
			},
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `contains(${cv},"Iphone")`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "18",
			fields: fields{
				mode: COMPATIBLE,
			},
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `contains(${cv},"Android")`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "19",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
				},
				rule: `${app} == media_std`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "20",
			fields: fields{
				mode: COMPATIBLE,
			},
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
				},
				rule: `${app} == media_std`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "21",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
				},
				rule: `"media_std" == ${app}`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "22",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"uid": 123456,
				},
				rule: `${uid} == 123456`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "23",
			fields: fields{
				mode: STRICT,
			},
			args: args{
				args: map[string]interface{}{
					"uid": 123456,
				},
				rule: `${uid} == 654321`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "24",
			fields: fields{
				mode: COMPATIBLE,
			},
			args: args{
				args: map[string]interface{}{
					"liveid": "123456",
				},
				rule: `mod(${liveid}, 100)==56`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "25",
			fields: fields{
				mode: COMPATIBLE,
			},
			args: args{
				rule: `true & false`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "26",
			fields: fields{
				mode: COMPATIBLE,
			},
			args: args{
				rule: `true | false`,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAST(tt.fields.mode)
			got, err := a.Judge(tt.args.args, tt.args.ops, tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("AST.Judge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AST.Judge() = %v, want %v", got, tt.want)
			}
		})
	}
}
