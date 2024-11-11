package validators

import (
	"reflect"
	"testing"
)

func Test_validateAlpha(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Positive case with valid phone number",
			args: args{
				input: "1234567890",
			},
			want: true,
		},
		{
			name: "Negative case with invalid phone number with size 9",
			args: args{
				input: "123456789",
			},
			want: false,
		},
		{
			name: "Negative case with invalid phone number with a char",
			args: args{
				input: "123456789a",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		fl := &mockFieldLevel{
			value:           reflect.ValueOf(tt.args.input),
			fieldName:       "PhoneNumber",
			structFieldName: "PhoneNumber",
			param:           "",
			tag:             "tenDigits",
		}
		// result :=
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePhoneNumber(fl); got != tt.want {
				t.Errorf("validateAlpha() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateUserName(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Positive with valid userName with .",
			args: args{
				input: "rak.ranjan",
			},
			want: true,
		},
		{
			name: "Positive with valid userName with a number",
			args: args{
				input: "rak.ranjan1",
			},
			want: true,
		},
		{
			name: "Negative with invalid userName with @",
			args: args{
				input: "rak.ranjan@",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		fl := &mockFieldLevel{
			value:           reflect.ValueOf(tt.args.input),
			fieldName:       "userName",
			structFieldName: "userName",
			param:           "",
			tag:             "alphaNumeric",
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateUserName(fl); got != tt.want {
				t.Errorf("validateUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
