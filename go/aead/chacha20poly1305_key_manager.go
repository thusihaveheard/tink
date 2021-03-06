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
//
////////////////////////////////////////////////////////////////////////////////

package aead

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/chacha20poly1305"
	"github.com/google/tink/go/subtle/aead"
	"github.com/google/tink/go/subtle/random"
	"github.com/google/tink/go/tink"

	cppb "github.com/google/tink/proto/chacha20_poly1305_go_proto"
	tinkpb "github.com/google/tink/proto/tink_go_proto"
)

const (
	// ChaCha20Poly1305KeyVersion is the maxmimal version of keys that this key manager supports.
	ChaCha20Poly1305KeyVersion = 0
	// ChaCha20Poly1305TypeURL is the url that this key manager supports.
	ChaCha20Poly1305TypeURL = "type.googleapis.com/google.crypto.tink.ChaCha20Poly1305Key"
)

// Common errors.
var errInvalidChaCha20Poly1305Key = fmt.Errorf("chacha20poly1305_key_manager: invalid key")
var errInvalidChaCha20Poly1305KeyFormat = fmt.Errorf("chacha20poly1305_key_manager: invalid key format")

// ChaCha20Poly1305KeyManager is an implementation of KeyManager interface.
// It generates new ChaCha20Poly1305Key keys and produces new instances of ChaCha20Poly1305 subtle.
type ChaCha20Poly1305KeyManager struct{}

// Assert that ChaCha20Poly1305KeyManager implements the KeyManager interface.
var _ tink.KeyManager = (*ChaCha20Poly1305KeyManager)(nil)

// NewChaCha20Poly1305KeyManager creates a new ChaCha20Poly1305KeyManager.
func NewChaCha20Poly1305KeyManager() *ChaCha20Poly1305KeyManager {
	return new(ChaCha20Poly1305KeyManager)
}

// GetPrimitiveFromSerializedKey creates an ChaCha20Poly1305 subtle for the given
// serialized ChaCha20Poly1305Key proto.
func (km *ChaCha20Poly1305KeyManager) GetPrimitiveFromSerializedKey(serializedKey []byte) (interface{}, error) {
	if len(serializedKey) == 0 {
		return nil, errInvalidChaCha20Poly1305Key
	}
	key := new(cppb.ChaCha20Poly1305Key)
	if err := proto.Unmarshal(serializedKey, key); err != nil {
		return nil, errInvalidChaCha20Poly1305Key
	}
	return km.GetPrimitiveFromKey(key)
}

// GetPrimitiveFromKey creates an ChaCha20Poly1305 subtle for the given ChaCha20Poly1305Key proto.
func (km *ChaCha20Poly1305KeyManager) GetPrimitiveFromKey(m proto.Message) (interface{}, error) {
	key, ok := m.(*cppb.ChaCha20Poly1305Key)
	if !ok {
		return nil, errInvalidChaCha20Poly1305Key
	}
	if err := km.validateKey(key); err != nil {
		return nil, err
	}
	ret, err := aead.NewChaCha20Poly1305(key.KeyValue)
	if err != nil {
		return nil, fmt.Errorf("chacha20poly1305_key_manager: cannot create new primitive: %s", err)
	}
	return ret, nil
}

// NewKeyFromSerializedKeyFormat is not implemented.
func (km *ChaCha20Poly1305KeyManager) NewKeyFromSerializedKeyFormat(serializedKeyFormat []byte) (proto.Message, error) {
	return km.NewChaCha20Poly1305Key(), nil
}

// NewKeyFromKeyFormat is not implemented.
func (km *ChaCha20Poly1305KeyManager) NewKeyFromKeyFormat(m proto.Message) (proto.Message, error) {
	return km.NewChaCha20Poly1305Key(), nil
}

// NewKeyData creates a new KeyData ignoring the specification in the given serialized key format
// because the key size is fixed. It should be used solely by the key management API.
func (km *ChaCha20Poly1305KeyManager) NewKeyData(serializedKeyFormat []byte) (*tinkpb.KeyData, error) {
	key := km.NewChaCha20Poly1305Key()
	serializedKey, err := proto.Marshal(key)
	if err != nil {
		return nil, err
	}
	return &tinkpb.KeyData{
		TypeUrl:         ChaCha20Poly1305TypeURL,
		Value:           serializedKey,
		KeyMaterialType: tinkpb.KeyData_SYMMETRIC,
	}, nil
}

// DoesSupport indicates if this key manager supports the given key type.
func (km *ChaCha20Poly1305KeyManager) DoesSupport(typeURL string) bool {
	return typeURL == ChaCha20Poly1305TypeURL
}

// GetKeyType returns the key type of keys managed by this key manager.
func (km *ChaCha20Poly1305KeyManager) GetKeyType() string {
	return ChaCha20Poly1305TypeURL
}

// NewChaCha20Poly1305Key returns a new ChaCha20Poly1305Key.
func (km *ChaCha20Poly1305KeyManager) NewChaCha20Poly1305Key() *cppb.ChaCha20Poly1305Key {
	keyValue := random.GetRandomBytes(chacha20poly1305.KeySize)
	return &cppb.ChaCha20Poly1305Key{
		Version:  ChaCha20Poly1305KeyVersion,
		KeyValue: keyValue,
	}
}

// validateKey validates the given ChaCha20Poly1305Key.
func (km *ChaCha20Poly1305KeyManager) validateKey(key *cppb.ChaCha20Poly1305Key) error {
	err := tink.ValidateVersion(key.Version, ChaCha20Poly1305KeyVersion)
	if err != nil {
		return fmt.Errorf("chacha20poly1305_key_manager: %s", err)
	}
	keySize := uint32(len(key.KeyValue))
	if keySize != chacha20poly1305.KeySize {
		return fmt.Errorf("chacha20poly1305_key_manager: keySize != %d", chacha20poly1305.KeySize)
	}
	return nil
}
