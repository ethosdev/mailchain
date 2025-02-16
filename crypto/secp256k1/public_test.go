// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secp256k1

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicKeyFromBytes(t *testing.T) {
	type args struct {
		keyBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"success-65",
			args{
				func() []byte {
					pub := make([]byte, 65)
					pub[0] = byte(4)
					copy(pub[1:], ecdsaPrivateKeyAlice().X.Bytes())
					copy(pub[33:], ecdsaPrivateKeyAlice().Y.Bytes())
					return pub
				}(),
			},
			alicePublicKeyBytes,
			false,
		},
		{
			"success-65-bob",
			args{
				func() []byte {
					pub := make([]byte, 65)
					pub[0] = byte(4)
					copy(pub[1:], ecdsaPrivateKeyBob().X.Bytes())
					copy(pub[33:], ecdsaPrivateKeyBob().Y.Bytes())
					return pub
				}(),
			},
			bobPublicKeyBytes,
			false,
		},
		{
			"err-65",
			args{
				[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			nil,
			true,
		},
		{
			"success-64",
			args{
				func() []byte {
					pub := make([]byte, 64)
					copy(pub, ecdsaPrivateKeyAlice().X.Bytes())
					copy(pub[32:], ecdsaPrivateKeyAlice().Y.Bytes())
					return pub
				}(),
			},
			alicePublicKeyBytes,
			false,
		},
		{
			"err-64",
			args{
				[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			nil,
			true,
		},
		{
			"success-33",
			args{
				bobPublicKey.Bytes(),
			},
			bobPublicKeyBytes,
			false,
		},
		{
			"err-63",
			args{
				[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			nil,
			true,
		},
		{
			"err",
			args{
				[]byte{0x3, 0xbd, 0xf6, 0xfb, 0x97},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromBytes(tt.args.keyBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotBytes []byte
			if got != nil {
				gotBytes = got.Bytes()
			}

			if !assert.Equal(t, tt.want, gotBytes) {
				t.Errorf("PublicKeyFromBytes() = %v, want %v", gotBytes, tt.want)
			}
		})
	}
}

func TestPublicKey_Bytes(t *testing.T) {
	tests := []struct {
		name   string
		pubKey PublicKey
		want   []byte
	}{
		{
			"bob",
			bobPublicKey,
			bobPublicKeyBytes,
		},
		{
			"alice",
			alicePublicKey,
			alicePublicKeyBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pubKey.Bytes(); !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKey_Kind(t *testing.T) {
	type fields struct {
		ecdsa ecdsa.PublicKey
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"bob",
			fields{
				ecdsaPublicKeyBob(),
			},
			"secp256k1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PublicKey{
				ecdsa: tt.fields.ecdsa,
			}
			if got := pk.Kind(); !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKey.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKey_Verify(t *testing.T) {
	tests := []struct {
		name    string
		pk      PublicKey
		message []byte
		sig     []byte
		want    bool
	}{
		{
			"success-bob-with-recid",
			bobPublicKey,
			[]byte("message"),
			[]byte{0x9d, 0xf7, 0x76, 0xab, 0xde, 0x8c, 0x20, 0x55, 0xc3, 0x4, 0x68, 0x37, 0xa8, 0x66, 0xf8, 0x89, 0x95, 0xf9, 0x82, 0xf0, 0x4b, 0xb8, 0x23, 0x40, 0xf0, 0x3, 0x8, 0x6a, 0x32, 0xa7, 0xac, 0xef, 0x5f, 0xa, 0xea, 0xda, 0x60, 0xbf, 0x9, 0xd5, 0xc3, 0x27, 0x61, 0xa, 0xc5, 0xc8, 0x33, 0xe3, 0xa0, 0x79, 0xdf, 0x6d, 0xe1, 0x9c, 0xa8, 0xcc, 0x33, 0xea, 0x1d, 0xe6, 0x3, 0x34, 0xb1, 0xa1, 0x0},
			true,
		},
		{
			"success-bob-without-recid",
			bobPublicKey,
			[]byte("message"),
			[]byte{0x9d, 0xf7, 0x76, 0xab, 0xde, 0x8c, 0x20, 0x55, 0xc3, 0x4, 0x68, 0x37, 0xa8, 0x66, 0xf8, 0x89, 0x95, 0xf9, 0x82, 0xf0, 0x4b, 0xb8, 0x23, 0x40, 0xf0, 0x3, 0x8, 0x6a, 0x32, 0xa7, 0xac, 0xef, 0x5f, 0xa, 0xea, 0xda, 0x60, 0xbf, 0x9, 0xd5, 0xc3, 0x27, 0x61, 0xa, 0xc5, 0xc8, 0x33, 0xe3, 0xa0, 0x79, 0xdf, 0x6d, 0xe1, 0x9c, 0xa8, 0xcc, 0x33, 0xea, 0x1d, 0xe6, 0x3, 0x34, 0xb1, 0xa1},
			true,
		},
		{
			"success-alice-with-recid",
			alicePublicKey,
			[]byte("egassem"),
			[]byte{0xe9, 0x33, 0xe, 0x4a, 0xe3, 0x5, 0x19, 0xea, 0x36, 0x37, 0x19, 0xdd, 0xbc, 0x91, 0xfd, 0x4f, 0xd3, 0x64, 0x9b, 0xdc, 0xf0, 0x74, 0x36, 0x16, 0xc9, 0x81, 0xfc, 0x6d, 0x3c, 0x7e, 0xb0, 0xd0, 0x6e, 0xdd, 0x4, 0x13, 0xfd, 0x15, 0xe5, 0xec, 0x64, 0x6e, 0x63, 0xe0, 0x84, 0xdb, 0xb2, 0xd7, 0xcf, 0x18, 0x3d, 0x81, 0x1e, 0x31, 0x36, 0x77, 0x39, 0x86, 0x4b, 0x58, 0xb8, 0x23, 0xed, 0xc, 0x1},
			true,
		},
		{
			"success-alice-without-recid",
			alicePublicKey,
			[]byte("egassem"),
			[]byte{0xe9, 0x33, 0xe, 0x4a, 0xe3, 0x5, 0x19, 0xea, 0x36, 0x37, 0x19, 0xdd, 0xbc, 0x91, 0xfd, 0x4f, 0xd3, 0x64, 0x9b, 0xdc, 0xf0, 0x74, 0x36, 0x16, 0xc9, 0x81, 0xfc, 0x6d, 0x3c, 0x7e, 0xb0, 0xd0, 0x6e, 0xdd, 0x4, 0x13, 0xfd, 0x15, 0xe5, 0xec, 0x64, 0x6e, 0x63, 0xe0, 0x84, 0xdb, 0xb2, 0xd7, 0xcf, 0x18, 0x3d, 0x81, 0x1e, 0x31, 0x36, 0x77, 0x39, 0x86, 0x4b, 0x58, 0xb8, 0x23, 0xed, 0xc},
			true,
		},
		{
			"err-invalid-signature-bob",
			bobPublicKey,
			[]byte("message"),
			[]byte{0xdb, 0x12, 0x7f, 0x98, 0x10, 0x91, 0xc6, 0xb0, 0x52, 0x9c, 0x17, 0x0, 0x3, 0xfe, 0xef, 0xfc, 0x7b, 0xd2, 0x97, 0x62, 0x95, 0x90, 0x1f, 0xd9, 0x3f, 0x21, 0x26, 0xe2, 0xbc, 0x5d, 0x11, 0xd9, 0x37, 0xd2, 0x5c, 0x63, 0xd8, 0xa3, 0x78, 0x4f, 0xb0, 0x42, 0x80, 0x3d, 0xf8, 0x30, 0x55, 0x79, 0x9d, 0x5f, 0x8f, 0x35, 0x85, 0xca, 0xf0, 0xad, 0x60, 0xdc, 0x4a, 0x74, 0x2d, 0xa4, 0xb0, 0x9d, 0x1},
			false,
		},
		{
			"err-invalid-signature-alice",
			alicePublicKey,
			[]byte("egassem"),
			[]byte{0xcf, 0x44, 0x7f, 0xc1, 0xd4, 0x6f, 0xd, 0xe2, 0xd4, 0xfc, 0x5c, 0x18, 0x5a, 0x7e, 0x89, 0xba, 0x5e, 0xca, 0x16, 0x68, 0xca, 0x73, 0xa3, 0x4d, 0x43, 0xf1, 0x8d, 0xa2, 0x45, 0xf1, 0xd1, 0xd7, 0x6f, 0x8f, 0x6e, 0xd1, 0x7c, 0x8c, 0xdf, 0x95, 0xd8, 0x46, 0xea, 0x8d, 0x5e, 0xa1, 0x50, 0x8f, 0x97, 0x18, 0xd7, 0xfe, 0x1a, 0x99, 0x69, 0x3b, 0x50, 0xe8, 0x9a, 0x30, 0x9c, 0x41, 0x2c, 0xb0, 0x1},
			false,
		},
		{
			"err-invalid-signature",
			bobPublicKey,
			[]byte("message"),
			[]byte{0xd},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pk.Verify(tt.message, tt.sig)
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKey.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKey_ECDSA(t *testing.T) {
	type fields struct {
		ecdsa ecdsa.PublicKey
	}
	tests := []struct {
		name   string
		fields fields
		want   *ecdsa.PublicKey
	}{
		{
			"success-bob",
			fields{
				ecdsaPublicKeyBob(),
			},
			func() *ecdsa.PublicKey {
				k := ecdsaPublicKeyBob()
				return &k
			}(),
		},
		{
			"success-alice",
			fields{
				ecdsaPublicKeyAlice(),
			},
			func() *ecdsa.PublicKey {
				k := ecdsaPublicKeyAlice()
				return &k
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PublicKey{
				ecdsa: tt.fields.ecdsa,
			}
			if got := pk.ECDSA(); !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKey.ECDSA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPublicKey_UncompressedBytes(t *testing.T) {
	type fields struct {
		ecdsa ecdsa.PublicKey
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"alice",
			fields{
				ecdsaPublicKeyAlice(),
			},
			[]byte{0x69, 0xd9, 0x8, 0x51, 0xe, 0x35, 0x5b, 0xeb, 0x1d, 0x5b, 0xf2, 0xdf, 0x81, 0x29, 0xe5, 0xb6, 0x40, 0x1e, 0x19, 0x69, 0x89, 0x1e, 0x80, 0x16, 0xa0, 0xb2, 0x30, 0x7, 0x39, 0xbb, 0xb0, 0x6, 0x87, 0x5, 0x5e, 0x59, 0x24, 0xa2, 0xfd, 0x8d, 0xd3, 0x5f, 0x6, 0x9d, 0xc1, 0x4d, 0x81, 0x47, 0xaa, 0x11, 0xc1, 0xf7, 0xe2, 0xf2, 0x71, 0x57, 0x34, 0x87, 0xe1, 0xbe, 0xeb, 0x2b, 0xe9, 0xd0},
		},
		{
			"bob",
			fields{
				ecdsaPublicKeyBob(),
			},
			[]byte{0xbd, 0xf6, 0xfb, 0x97, 0xc9, 0x7c, 0x12, 0x6b, 0x49, 0x21, 0x86, 0xa4, 0xd5, 0xb2, 0x8f, 0x34, 0xf0, 0x67, 0x1a, 0x5a, 0xac, 0xc9, 0x74, 0xda, 0x3b, 0xde, 0xb, 0xe9, 0x3e, 0x45, 0xa1, 0xc5, 0xf, 0x89, 0xce, 0xff, 0x72, 0xbd, 0x4, 0xac, 0x9e, 0x25, 0xa0, 0x4a, 0x1a, 0x6c, 0xb0, 0x10, 0xae, 0xda, 0xf6, 0x5f, 0x91, 0xce, 0xc8, 0xeb, 0xe7, 0x59, 0x1, 0xc4, 0x9b, 0x63, 0x35, 0x5d},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PublicKey{
				ecdsa: tt.fields.ecdsa,
			}
			if got := pk.UncompressedBytes(); !assert.Equal(t, tt.want, got) {
				t.Errorf("PublicKey.UncompressedBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
