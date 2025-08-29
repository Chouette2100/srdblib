// Copyright © 2024-2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package srdblib

import "testing"

func TestExtractStructColumns(t *testing.T) {
	type args struct {
		model interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			"test_case_1",
			args{model: struct {
				Field1 string
				Field2 int
			}{"value1", 2}},
			"field1, field2",
		},
		{ // 構造体
			"test_case_2",
			args{
				model: Event{},
			},
			"eventid, ieventid, event_name, period, starttime, endtime, noentry, intervalmin, modmin, modsec, fromorder, toorder, resethh, resetmm, nobasis, maxdsp, cmap, target, rstatus, maxpoint, thinit, thdelta, achk, aclr",
		},
		{ // (埋め込みフィールドがある)構造体のポインタ
			"test_case_3",
			args{
				model: &Wevent{},
			},
			"eventid, ieventid, event_name, period, starttime, endtime, noentry, intervalmin, modmin, modsec, fromorder, toorder, resethh, resetmm, nobasis, maxdsp, cmap, target, rstatus, maxpoint, thinit, thdelta, achk, aclr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractStructColumns(tt.args.model); got != tt.want {
				t.Errorf("ExtractStructColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}
