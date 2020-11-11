package unicode

import (
	"testing"
)

func TestUnicodeToString(t *testing.T) {
	type args struct {
		form string
	}
	tests := []struct {
		name    string
		args    args
		wantTo  string
		wantErr bool
	}{
		{"test emoji", args{form: "\\u5f20\\ud83d\\udc37\\ud83d\\udc37"}, "张🐷🐷", false},
		{"test abc", args{form: "123"}, "123", false},
		{"test abc", args{form: "ab\\u5f20\\ud83d\\udc37\\ud83d\\udc37c"}, "ab张🐷🐷c", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTo, err := UnicodeToString(tt.args.form)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnicodeToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTo != tt.wantTo {
				t.Errorf("UnicodeToString() gotTo = %v, want %v", gotTo, tt.wantTo)
			}
		})
	}
}
