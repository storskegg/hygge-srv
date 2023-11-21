package messages

import (
	"testing"

	"github.com/go-test/deep"
)

func TestParseMEF(t *testing.T) {
	type args struct {
		line string
	}
	type testCase struct {
		name    string
		args    args
		want    *MEF
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Happy Path: Hygge Payload",
			args: args{
				line: "W4PHO|1|3571c78c|4903|64.73|18.16|3.98",
			},
			want: &MEF{
				Callsign:    "W4PHO",
				MessageType: 1,
				Digest:      "3571c78c",
				Data: &Hygge{
					PacketSequence: 4903,
					Humidity:       64.73,
					Temperature:    18.16,
					Battery:        3.98,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMEF(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMEF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, diff := range deep.Equal(got, tt.want) {
				t.Errorf(diff)
			}
		})
	}
}
