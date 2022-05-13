package path

import (
	"reflect"
	"testing"

	"github.com/dualm/ethernet-ip/types"
)

func TestPortBuild(t *testing.T) {
	type args struct {
		link   []byte
		portID uint16
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				link:   []byte("130.151.132.1"),
				portID: 2,
			},
			want:    []byte{0x12, 0x0d, 0x31, 0x33, 0x30, 0x2e, 0x31, 0x35, 0x31, 0x2e, 0x31, 0x33, 0x32, 0x2e, 0x31, 0x00},
			wantErr: false,
		},

		{
			name: "2",
			args: args{
				link:   []byte("plc.controlnet.org"),
				portID: 3,
			},
			want:    []byte{0x13, 0x12, 0x70, 0x6c, 0x63, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6e, 0x65, 0x74, 0x2e, 0x6f, 0x72, 0x67},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				link:   []byte("130.151.132.55:0x3210"),
				portID: 6,
			},
			want:    []byte{0x16, 0x15, 0x31, 0x33, 0x30, 0x2e, 0x31, 0x35, 0x31, 0x2e, 0x31, 0x33, 0x32, 0x2e, 0x35, 0x35, 0x3a, 0x30, 0x78, 0x33, 0x32, 0x31, 0x30, 0x00},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				link:   []byte("plc.controlnet.org:9876"),
				portID: 5,
			},
			want:    []byte{0x15, 0x17, 0x70, 0x6c, 0x63, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6e, 0x65, 0x74, 0x2e, 0x6f, 0x72, 0x67, 0x3a, 0x39, 0x38, 0x37, 0x36, 0x00},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PortBuild(tt.args.link, tt.args.portID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortBuild() = % x, want % x", got, tt.want)
			}
		})
	}
}

func TestLogicalBuild(t *testing.T) {
	type args struct {
		logicalType LogicalType
		value       types.UDINT
		format      uint8
		padded      bool
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				logicalType: LogicalClassID,
				value:       6,
				format:      0,
				padded:      true,
			},
			want:    []byte{0x20, 0x06},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				logicalType: LogicalInstaceID,
				value:       2,
				format:      0,
				padded:      false,
			},
			want:    []byte{0x24, 0x02},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				logicalType: LogicalClassID,
				value:       5,
				format:      1,
				padded:      true,
			},
			want:    []byte{0x21, 0x00, 0x05, 0x00},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				logicalType: LogicalClassID,
				value:       5,
				format:      1,
				padded:      false,
			},
			want:    []byte{0x21, 0x05, 0x00},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LogicalBuild(tt.args.logicalType, tt.args.value, tt.args.format, tt.args.padded)
			if (err != nil) != tt.wantErr {
				t.Errorf("LogicalBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogicalBuild() = % x, want % x", got, tt.want)
			}
		})
	}
}

func TestDataBuild(t *testing.T) {
	type args struct {
		datatype DataSegmentSubType
		raw      []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				datatype: SimpleDataSegment,
				raw:      []byte{0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00},
			},
			want:    []byte{0x80, 0x07, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				datatype: SymbolSegment,
				raw:      []byte("starter"),
			},
			want:    []byte{0x91, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x00},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DataBuild(tt.args.datatype, tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataBuild() = % x, want % x", got, tt.want)
			}
		})
	}
}
