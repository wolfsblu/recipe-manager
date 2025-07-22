package security

import (
	"strings"
	"testing"
)

func TestCreateHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		params   *HashParams
		wantErr  bool
	}{
		{
			name:     "Valid password with default params",
			password: "test-password",
			params:   DefaultHashParams,
			wantErr:  false,
		},
		{
			name:     "Empty password with default params",
			password: "",
			params:   DefaultHashParams,
			wantErr:  false,
		},
		{
			name:     "Valid password with custom params",
			password: "test-password",
			params: &HashParams{
				Memory:      16 * 1024,
				Iterations:  2,
				Parallelism: 1,
				SaltLength:  8,
				KeyLength:   16,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := CreateHash(tc.password, tc.params)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateHash() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err == nil {
				// Verify hash format
				if !strings.HasPrefix(hash, "$argon2id$v=") {
					t.Errorf("CreateHash() hash format invalid: %v", hash)
				}

				// Verify parts count
				parts := strings.Split(hash, "$")
				if len(parts) != 6 {
					t.Errorf("CreateHash() hash parts count invalid: got %d, want 6", len(parts))
				}
			}
		})
	}
}

func TestComparePasswordAndHash(t *testing.T) {
	// Create a hash for testing
	password := "test-password"
	hash, err := CreateHash(password, DefaultHashParams)
	if err != nil {
		t.Fatalf("Failed to create hash for testing: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
		wantErr  bool
	}{
		{
			name:     "Matching password and hash",
			password: password,
			hash:     hash,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "Non-matching password and hash",
			password: "wrong-password",
			hash:     hash,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "Empty password with valid hash",
			password: "",
			hash:     hash,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "Valid password with invalid hash format",
			password: password,
			hash:     "invalid-hash-format",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "Valid password with empty hash",
			password: password,
			hash:     "",
			want:     false,
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ComparePasswordAndHash(tc.password, tc.hash)
			if (err != nil) != tc.wantErr {
				t.Errorf("ComparePasswordAndHash() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("ComparePasswordAndHash() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDecodeHash(t *testing.T) {
	// Create a valid hash
	password := "test-password"
	validHash, err := CreateHash(password, DefaultHashParams)
	if err != nil {
		t.Fatalf("Failed to create hash for testing: %v", err)
	}

	tests := []struct {
		name    string
		hash    string
		wantErr bool
	}{
		{
			name:    "Valid hash",
			hash:    validHash,
			wantErr: false,
		},
		{
			name:    "Empty hash",
			hash:    "",
			wantErr: true,
		},
		{
			name:    "Invalid hash format",
			hash:    "invalid-format",
			wantErr: true,
		},
		{
			name:    "Hash with wrong number of segments",
			hash:    "$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ",
			wantErr: true,
		},
		{
			name:    "Hash with incompatible variant",
			hash:    "$argon2i$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params, salt, key, err := DecodeHash(tc.hash)
			if (err != nil) != tc.wantErr {
				t.Errorf("DecodeHash() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if err == nil {
				// Check that params and data were extracted correctly
				if params == nil {
					t.Error("DecodeHash() params is nil for valid hash")
				}
				if salt == nil {
					t.Error("DecodeHash() salt is nil for valid hash")
				}
				if key == nil {
					t.Error("DecodeHash() key is nil for valid hash")
				}

				// Check that we have the expected parameters
				if params.Memory != DefaultHashParams.Memory {
					t.Errorf("DecodeHash() params.Memory = %v, want %v", params.Memory, DefaultHashParams.Memory)
				}
				if params.Iterations != DefaultHashParams.Iterations {
					t.Errorf("DecodeHash() params.Iterations = %v, want %v", params.Iterations, DefaultHashParams.Iterations)
				}
				if params.Parallelism != DefaultHashParams.Parallelism {
					t.Errorf("DecodeHash() params.Parallelism = %v, want %v", params.Parallelism, DefaultHashParams.Parallelism)
				}
				if params.SaltLength != DefaultHashParams.SaltLength {
					t.Errorf("DecodeHash() params.SaltLength = %v, want %v", params.SaltLength, DefaultHashParams.SaltLength)
				}
				if params.KeyLength != DefaultHashParams.KeyLength {
					t.Errorf("DecodeHash() params.KeyLength = %v, want %v", params.KeyLength, DefaultHashParams.KeyLength)
				}
			}
		})
	}
}
