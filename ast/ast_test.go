package ast

import "testing"

func Test_Eval(t *testing.T) {
	type args struct {
		args map[string]interface{}
		ops  map[string]func(interface{}) bool
		rule string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				rule: `true`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				rule: `false`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
				},
				rule: `app == "media_std"`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				args: map[string]interface{}{
					"uid": 1,
				},
				rule: `uid % 10 < 3`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "5",
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
					"uid": 1,
				},
				rule: `app == "media_std" && uid % 10 < 3`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "6",
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
					"uid": 5,
				},
				rule: `app == "media_std" && uid % 10 < 3`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "7",
			args: args{
				args: map[string]interface{}{
					"app": "std_media",
					"uid": 2,
				},
				rule: `app == "media_std" && uid % 10 < 3`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "8",
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
					"uid": 5,
				},
				rule: `app == "media_std" || uid % 10 < 3`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "9",
			args: args{
				args: map[string]interface{}{
					"app": "std_media",
					"uid": 3,
				},
				rule: `app == "media_std" || uid % 10 <= 3`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "10",
			args: args{
				args: map[string]interface{}{
					"app": "std_media",
					"uid": 3,
				},
				rule: `false && (app == "media_std" || uid % 10 <= 3)`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "11",
			args: args{
				args: map[string]interface{}{
					"uid": 3,
				},
				ops: map[string]func(interface{}) bool{
					"cls_whitelist": func(uid interface{}) bool {
						return true
					},
				},
				rule: `cls_whitelist[uid]`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "12",
			args: args{
				args: map[string]interface{}{
					"uid": 3,
				},
				ops: map[string]func(interface{}) bool{
					"cls_whitelist": func(uid interface{}) bool {
						return false
					},
				},
				rule: `cls_whitelist[uid]`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "13",
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `cv[2:3]=="7"`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "14",
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `cv[8:]=="Iphone"`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "15",
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `!(cv[8:]=="Iphone")`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "16",
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `cv[8:100]=="Iphone"`,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "17",
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `contains(cv,"Iphone")`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "18",
			args: args{
				args: map[string]interface{}{
					"cv": "IK7.8.9_Iphone",
				},
				rule: `contains(cv,"Android")`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "19",
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
				},
				rule: `app == media_std`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "20",
			args: args{
				args: map[string]interface{}{
					"app": "media_std",
				},
				rule: `"media_std" == app`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "21",
			args: args{
				args: map[string]interface{}{
					"uid": 123456,
				},
				rule: `uid == 123456`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "22",
			args: args{
				args: map[string]interface{}{
					"uid": 123456,
				},
				rule: `uid == 654321`,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eval(tt.args.args, tt.args.ops, tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
