// Package hybrid defines several hybrid classical/quantum KEMs.
//
// KEMs are combined by simple concatenation of shared secrets, cipher texts,
// public keys, etc, see
//
//	https://datatracker.ietf.org/doc/draft-ietf-tls-hybrid-design/
//	https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-56Cr2.pdf
//
// Note that this is only fine if the shared secret is used in its entirety
// in a next step, such as being hashed or used as key.
//
// For deriving a KEM keypair deterministically and encapsulating
// deterministically, we expand a single seed to both using SHAKE256,
// so that a non-uniform seed (such as a shared secret generated by a hybrid
// KEM where one of the KEMs is weak) doesn't impact just one of the KEMs.
//
// Of our XOF (SHAKE256), we desire two security properties:
//
//  1. The internal state of the XOF should be big enough so that we
//     do not loose entropy.
//  2. From one of the new seeds, we shouldn't be able to derive
//     the other or the original seed.
//
// SHAKE256, and all siblings in the SHA3 family, have a 200B internal
// state, so (1) is fine if our seeds are less than 200B.
// If SHAKE256 is computationally indistinguishable from a random
// sponge, then it affords us 256b security against (2) by the
// flat sponge claim [https://keccak.team/files/SpongeFunctions.pdf].
// None of the implemented schemes claim more than 256b security
// and so SHAKE256 will do fine.
package hybrid

import (
	"errors"

	"github.com/JI-0/circl/internal/sha3"
	"github.com/JI-0/circl/kem"
	"github.com/JI-0/circl/kem/kyber/kyber1024"
	"github.com/JI-0/circl/kem/kyber/kyber512"
	"github.com/JI-0/circl/kem/kyber/kyber768"
)

var ErrUninitialized = errors.New("public or private key not initialized")

// Returns the hybrid KEM of Kyber512Draft00 and X25519.
func Kyber512X25519() kem.Scheme { return kyber512X }

// Returns the hybrid KEM of Kyber768Draft00 and X25519.
func Kyber768X25519() kem.Scheme { return kyber768X }

// Returns the hybrid KEM of Kyber768Draft00 and X448.
func Kyber768X448() kem.Scheme { return kyber768X4 }

// Returns the hybrid KEM of Kyber1024Draft00 and X448.
func Kyber1024X448() kem.Scheme { return kyber1024X }

// Returns the hybrid KEM of Kyber768Draft00 and P-256.
func P256Kyber768Draft00() kem.Scheme { return p256Kyber768Draft00 }

var p256Kyber768Draft00 kem.Scheme = &scheme{
	"P256Kyber768Draft00",
	p256Kem,
	kyber768.Scheme(),
}

var kyber512X kem.Scheme = &scheme{
	"Kyber512-X25519",
	x25519Kem,
	kyber512.Scheme(),
}

var kyber768X kem.Scheme = &scheme{
	"Kyber768-X25519",
	x25519Kem,
	kyber768.Scheme(),
}

var kyber768X4 kem.Scheme = &scheme{
	"Kyber768-X448",
	x448Kem,
	kyber768.Scheme(),
}

var kyber1024X kem.Scheme = &scheme{
	"Kyber1024-X448",
	x448Kem,
	kyber1024.Scheme(),
}

// Public key of a hybrid KEM.
type publicKey struct {
	scheme *scheme
	first  kem.PublicKey
	second kem.PublicKey
}

// Private key of a hybrid KEM.
type privateKey struct {
	scheme *scheme
	first  kem.PrivateKey
	second kem.PrivateKey
}

// Scheme for a hybrid KEM.
type scheme struct {
	name   string
	first  kem.Scheme
	second kem.Scheme
}

func (sch *scheme) Name() string { return sch.name }
func (sch *scheme) PublicKeySize() int {
	return sch.first.PublicKeySize() + sch.second.PublicKeySize()
}

func (sch *scheme) PrivateKeySize() int {
	return sch.first.PrivateKeySize() + sch.second.PrivateKeySize()
}

func (sch *scheme) SeedSize() int {
	first := sch.first.SeedSize()
	second := sch.second.SeedSize()
	ret := second
	if first > second {
		ret = first
	}
	return ret
}

func (sch *scheme) SharedKeySize() int {
	return sch.first.SharedKeySize() + sch.second.SharedKeySize()
}

func (sch *scheme) CiphertextSize() int {
	return sch.first.CiphertextSize() + sch.second.CiphertextSize()
}

func (sch *scheme) EncapsulationSeedSize() int {
	first := sch.first.EncapsulationSeedSize()
	second := sch.second.EncapsulationSeedSize()
	ret := second
	if first > second {
		ret = first
	}
	return ret
}

func (sk *privateKey) Scheme() kem.Scheme { return sk.scheme }
func (pk *publicKey) Scheme() kem.Scheme  { return pk.scheme }

func (sk *privateKey) MarshalBinary() ([]byte, error) {
	if sk.first == nil || sk.second == nil {
		return nil, ErrUninitialized
	}
	first, err := sk.first.MarshalBinary()
	if err != nil {
		return nil, err
	}
	second, err := sk.second.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(first, second...), nil
}

func (sk *privateKey) Equal(other kem.PrivateKey) bool {
	oth, ok := other.(*privateKey)
	if !ok {
		return false
	}
	if sk.first == nil && sk.second == nil && oth.first == nil && oth.second == nil {
		return true
	}
	if sk.first == nil || sk.second == nil || oth.first == nil || oth.second == nil {
		return false
	}
	return sk.first.Equal(oth.first) && sk.second.Equal(oth.second)
}

func (sk *privateKey) Public() kem.PublicKey {
	return &publicKey{sk.scheme, sk.first.Public(), sk.second.Public()}
}

func (pk *publicKey) Equal(other kem.PublicKey) bool {
	oth, ok := other.(*publicKey)
	if !ok {
		return false
	}
	if pk.first == nil && pk.second == nil && oth.first == nil && oth.second == nil {
		return true
	}
	if pk.first == nil || pk.second == nil || oth.first == nil || oth.second == nil {
		return false
	}
	return pk.first.Equal(oth.first) && pk.second.Equal(oth.second)
}

func (pk *publicKey) MarshalBinary() ([]byte, error) {
	if pk.first == nil || pk.second == nil {
		return nil, ErrUninitialized
	}
	first, err := pk.first.MarshalBinary()
	if err != nil {
		return nil, err
	}
	second, err := pk.second.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(first, second...), nil
}

func (sch *scheme) GenerateKeyPair() (kem.PublicKey, kem.PrivateKey, error) {
	pk1, sk1, err := sch.first.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	pk2, sk2, err := sch.second.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	return &publicKey{sch, pk1, pk2}, &privateKey{sch, sk1, sk2}, nil
}

func (sch *scheme) DeriveKeyPair(seed []byte) (kem.PublicKey, kem.PrivateKey) {
	if len(seed) != sch.SeedSize() {
		panic(kem.ErrSeedSize)
	}
	h := sha3.NewShake256()
	_, _ = h.Write(seed)
	first := make([]byte, sch.first.SeedSize())
	second := make([]byte, sch.second.SeedSize())
	_, _ = h.Read(first)
	_, _ = h.Read(second)

	pk1, sk1 := sch.first.DeriveKeyPair(first)
	pk2, sk2 := sch.second.DeriveKeyPair(second)

	return &publicKey{sch, pk1, pk2}, &privateKey{sch, sk1, sk2}
}

func (sch *scheme) Encapsulate(pk kem.PublicKey) (ct, ss []byte, err error) {
	pub, ok := pk.(*publicKey)
	if !ok {
		return nil, nil, kem.ErrTypeMismatch
	}

	ct1, ss1, err := sch.first.Encapsulate(pub.first)
	if err != nil {
		return nil, nil, err
	}

	ct2, ss2, err := sch.second.Encapsulate(pub.second)
	if err != nil {
		return nil, nil, err
	}

	return append(ct1, ct2...), append(ss1, ss2...), nil
}

func (sch *scheme) EncapsulateDeterministically(
	pk kem.PublicKey, seed []byte,
) (ct, ss []byte, err error) {
	if len(seed) != sch.EncapsulationSeedSize() {
		return nil, nil, kem.ErrSeedSize
	}

	h := sha3.NewShake256()
	_, _ = h.Write(seed)
	first := make([]byte, sch.first.EncapsulationSeedSize())
	second := make([]byte, sch.second.EncapsulationSeedSize())
	_, _ = h.Read(first)
	_, _ = h.Read(second)

	pub, ok := pk.(*publicKey)
	if !ok {
		return nil, nil, kem.ErrTypeMismatch
	}

	ct1, ss1, err := sch.first.EncapsulateDeterministically(pub.first, first)
	if err != nil {
		return nil, nil, err
	}
	ct2, ss2, err := sch.second.EncapsulateDeterministically(pub.second, second)
	if err != nil {
		return nil, nil, err
	}
	return append(ct1, ct2...), append(ss1, ss2...), nil
}

func (sch *scheme) Decapsulate(sk kem.PrivateKey, ct []byte) ([]byte, error) {
	if len(ct) != sch.CiphertextSize() {
		return nil, kem.ErrCiphertextSize
	}

	priv, ok := sk.(*privateKey)
	if !ok {
		return nil, kem.ErrTypeMismatch
	}

	firstSize := sch.first.CiphertextSize()
	ss1, err := sch.first.Decapsulate(priv.first, ct[:firstSize])
	if err != nil {
		return nil, err
	}
	ss2, err := sch.second.Decapsulate(priv.second, ct[firstSize:])
	if err != nil {
		return nil, err
	}
	return append(ss1, ss2...), nil
}

func (sch *scheme) UnmarshalBinaryPublicKey(buf []byte) (kem.PublicKey, error) {
	if len(buf) != sch.PublicKeySize() {
		return nil, kem.ErrPubKeySize
	}
	firstSize := sch.first.PublicKeySize()
	pk1, err := sch.first.UnmarshalBinaryPublicKey(buf[:firstSize])
	if err != nil {
		return nil, err
	}
	pk2, err := sch.second.UnmarshalBinaryPublicKey(buf[firstSize:])
	if err != nil {
		return nil, err
	}
	return &publicKey{sch, pk1, pk2}, nil
}

func (sch *scheme) UnmarshalBinaryPrivateKey(buf []byte) (kem.PrivateKey, error) {
	if len(buf) != sch.PrivateKeySize() {
		return nil, kem.ErrPrivKeySize
	}
	firstSize := sch.first.PrivateKeySize()
	sk1, err := sch.first.UnmarshalBinaryPrivateKey(buf[:firstSize])
	if err != nil {
		return nil, err
	}
	sk2, err := sch.second.UnmarshalBinaryPrivateKey(buf[firstSize:])
	if err != nil {
		return nil, err
	}
	return &privateKey{sch, sk1, sk2}, nil
}
