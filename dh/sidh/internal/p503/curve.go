// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package p503

import (
	"errors"
	. "github.com/JI-0/circl/dh/sidh/internal/common"
	"math"
)

// Stores isogeny 3 curve constants
type isogeny3 struct {
	K1 Fp2
	K2 Fp2
}

// Stores isogeny 4 curve constants
type isogeny4 struct {
	isogeny3
	K3 Fp2
}

// Computes j-invariant for a curve y2=x3+A/Cx+x with A,C in F_(p^2). Result
// is returned in jBytes buffer, encoded in little-endian format. Caller
// provided jBytes buffer has to be big enough to j-invariant value. In case
// of SIDH, buffer size must be at least size of shared secret.
// Implementation corresponds to Algorithm 9 from SIKE.
func Jinvariant(cparams *ProjectiveCurveParameters, j *Fp2) {
	var t0, t1 Fp2

	sqr(j, &cparams.A)   // j  = A^2
	sqr(&t1, &cparams.C) // t1 = C^2
	add(&t0, &t1, &t1)   // t0 = t1 + t1
	sub(&t0, j, &t0)     // t0 = j - t0
	sub(&t0, &t0, &t1)   // t0 = t0 - t1
	sub(j, &t0, &t1)     // t0 = t0 - t1
	sqr(&t1, &t1)        // t1 = t1^2
	mul(j, j, &t1)       // j = j * t1
	add(&t0, &t0, &t0)   // t0 = t0 + t0
	add(&t0, &t0, &t0)   // t0 = t0 + t0
	sqr(&t1, &t0)        // t1 = t0^2
	mul(&t0, &t0, &t1)   // t0 = t0 * t1
	add(&t0, &t0, &t0)   // t0 = t0 + t0
	add(&t0, &t0, &t0)   // t0 = t0 + t0
	inv(j, j)            // j  = 1/j
	mul(j, &t0, j)       // j  = t0 * j
}

// Given affine points x(P), x(Q) and x(Q-P) in a extension field F_{p^2}, function
// recovers projective coordinate A of a curve. This is Algorithm 10 from SIKE.
func RecoverCoordinateA(curve *ProjectiveCurveParameters, xp, xq, xr *Fp2) {
	var t0, t1 Fp2

	add(&t1, xp, xq)                        // t1 = Xp + Xq
	mul(&t0, xp, xq)                        // t0 = Xp * Xq
	mul(&curve.A, xr, &t1)                  // A  = X(q-p) * t1
	add(&curve.A, &curve.A, &t0)            // A  = A + t0
	mul(&t0, &t0, xr)                       // t0 = t0 * X(q-p)
	sub(&curve.A, &curve.A, &params.OneFp2) // A  = A - 1
	add(&t0, &t0, &t0)                      // t0 = t0 + t0
	add(&t1, &t1, xr)                       // t1 = t1 + X(q-p)
	add(&t0, &t0, &t0)                      // t0 = t0 + t0
	sqr(&curve.A, &curve.A)                 // A  = A^2
	inv(&t0, &t0)                           // t0 = 1/t0
	mul(&curve.A, &curve.A, &t0)            // A  = A * t0
	sub(&curve.A, &curve.A, &t1)            // A  = A - t1
}

// Computes equivalence (A:C) ~ (A+2C : A-2C)
func CalcCurveParamsEquiv3(cparams *ProjectiveCurveParameters) CurveCoefficientsEquiv {
	var coef CurveCoefficientsEquiv
	var c2 Fp2

	add(&c2, &cparams.C, &cparams.C)
	// A24p = A+2*C
	add(&coef.A, &cparams.A, &c2)
	// A24m = A-2*C
	sub(&coef.C, &cparams.A, &c2)
	return coef
}

// Computes equivalence (A:C) ~ (A+2C : 4C)
func CalcCurveParamsEquiv4(cparams *ProjectiveCurveParameters) CurveCoefficientsEquiv {
	var coefEq CurveCoefficientsEquiv

	add(&coefEq.C, &cparams.C, &cparams.C)
	// A24p = A+2C
	add(&coefEq.A, &cparams.A, &coefEq.C)
	// C24 = 4*C
	add(&coefEq.C, &coefEq.C, &coefEq.C)
	return coefEq
}

// Helper function for RightToLeftLadder(). Returns A+2C / 4.
func CalcAplus2Over4(cparams *ProjectiveCurveParameters) (ret Fp2) {
	var tmp Fp2

	// 2C
	add(&tmp, &cparams.C, &cparams.C)
	// A+2C
	add(&ret, &cparams.A, &tmp)
	// 1/4C
	add(&tmp, &tmp, &tmp)
	inv(&tmp, &tmp)
	// A+2C/4C
	mul(&ret, &ret, &tmp)
	return
}

// Recovers (A:C) curve parameters from projectively equivalent (A+2C:A-2C).
func RecoverCurveCoefficients3(cparams *ProjectiveCurveParameters, coefEq *CurveCoefficientsEquiv) {
	add(&cparams.A, &coefEq.A, &coefEq.C)
	// cparams.A = 2*(A+2C+A-2C) = 4A
	add(&cparams.A, &cparams.A, &cparams.A)
	// cparams.C = (A+2C-A+2C) = 4C
	sub(&cparams.C, &coefEq.A, &coefEq.C)
	return
}

// Recovers (A:C) curve parameters from projectively equivalent (A+2C:4C).
func RecoverCurveCoefficients4(cparams *ProjectiveCurveParameters, coefEq *CurveCoefficientsEquiv) {
	// cparams.C = (4C)*1/2=2C
	mul(&cparams.C, &coefEq.C, &params.HalfFp2)
	// cparams.A = A+2C - 2C = A
	sub(&cparams.A, &coefEq.A, &cparams.C)
	// cparams.C = 2C * 1/2 = C
	mul(&cparams.C, &cparams.C, &params.HalfFp2)
}

// Combined coordinate doubling and differential addition. Takes projective points
// P,Q,Q-P and (A+2C)/4C curve E coefficient. Returns 2*P and P+Q calculated on E.
// Function is used only by RightToLeftLadder. Corresponds to Algorithm 5 of SIKE
func xDbladd(P, Q, QmP *ProjectivePoint, a24 *Fp2) (dblP, PaQ ProjectivePoint) {
	var t0, t1, t2 Fp2

	xQmP, zQmP := &QmP.X, &QmP.Z
	xPaQ, zPaQ := &PaQ.X, &PaQ.Z
	x2P, z2P := &dblP.X, &dblP.Z
	xP, zP := &P.X, &P.Z
	xQ, zQ := &Q.X, &Q.Z

	add(&t0, xP, zP)      // t0   = Xp+Zp
	sub(&t1, xP, zP)      // t1   = Xp-Zp
	sqr(x2P, &t0)         // 2P.X = t0^2
	sub(&t2, xQ, zQ)      // t2   = Xq-Zq
	add(xPaQ, xQ, zQ)     // Xp+q = Xq+Zq
	mul(&t0, &t0, &t2)    // t0   = t0 * t2
	mul(z2P, &t1, &t1)    // 2P.Z = t1 * t1
	mul(&t1, &t1, xPaQ)   // t1   = t1 * Xp+q
	sub(&t2, x2P, z2P)    // t2   = 2P.X - 2P.Z
	mul(x2P, x2P, z2P)    // 2P.X = 2P.X * 2P.Z
	mul(xPaQ, a24, &t2)   // Xp+q = A24 * t2
	sub(zPaQ, &t0, &t1)   // Zp+q = t0 - t1
	add(z2P, xPaQ, z2P)   // 2P.Z = Xp+q + 2P.Z
	add(xPaQ, &t0, &t1)   // Xp+q = t0 + t1
	mul(z2P, z2P, &t2)    // 2P.Z = 2P.Z * t2
	sqr(zPaQ, zPaQ)       // Zp+q = Zp+q ^ 2
	sqr(xPaQ, xPaQ)       // Xp+q = Xp+q ^ 2
	mul(zPaQ, xQmP, zPaQ) // Zp+q = Xq-p * Zp+q
	mul(xPaQ, zQmP, xPaQ) // Xp+q = Zq-p * Xp+q
	return
}

// Given the curve parameters, xP = x(P), computes xP = x([2^k]P)
// Safe to overlap xP, x2P.
func Pow2k(xP *ProjectivePoint, params *CurveCoefficientsEquiv, k uint32) {
	var t0, t1 Fp2

	x, z := &xP.X, &xP.Z
	for i := uint32(0); i < k; i++ {
		sub(&t0, x, z)           // t0  = Xp - Zp
		add(&t1, x, z)           // t1  = Xp + Zp
		sqr(&t0, &t0)            // t0  = t0 ^ 2
		sqr(&t1, &t1)            // t1  = t1 ^ 2
		mul(z, &params.C, &t0)   // Z2p = C24 * t0
		mul(x, z, &t1)           // X2p = Z2p * t1
		sub(&t1, &t1, &t0)       // t1  = t1 - t0
		mul(&t0, &params.A, &t1) // t0  = A24+ * t1
		add(z, z, &t0)           // Z2p = Z2p + t0
		mul(z, z, &t1)           // Zp  = Z2p * t1
	}
}

// Given the curve parameters, xP = x(P), and k >= 0, compute xP = x([3^k]P).
//
// Safe to overlap xP, xR.
func Pow3k(xP *ProjectivePoint, params *CurveCoefficientsEquiv, k uint32) {
	var t0, t1, t2, t3, t4, t5, t6 Fp2

	x, z := &xP.X, &xP.Z
	for i := uint32(0); i < k; i++ {
		sub(&t0, x, z)           // t0  = Xp - Zp
		sqr(&t2, &t0)            // t2  = t0^2
		add(&t1, x, z)           // t1  = Xp + Zp
		sqr(&t3, &t1)            // t3  = t1^2
		add(&t4, &t1, &t0)       // t4  = t1 + t0
		sub(&t0, &t1, &t0)       // t0  = t1 - t0
		sqr(&t1, &t4)            // t1  = t4^2
		sub(&t1, &t1, &t3)       // t1  = t1 - t3
		sub(&t1, &t1, &t2)       // t1  = t1 - t2
		mul(&t5, &t3, &params.A) // t5  = t3 * A24+
		mul(&t3, &t3, &t5)       // t3  = t5 * t3
		mul(&t6, &t2, &params.C) // t6  = t2 * A24-
		mul(&t2, &t2, &t6)       // t2  = t2 * t6
		sub(&t3, &t2, &t3)       // t3  = t2 - t3
		sub(&t2, &t5, &t6)       // t2  = t5 - t6
		mul(&t1, &t2, &t1)       // t1  = t2 * t1
		add(&t2, &t3, &t1)       // t2  = t3 + t1
		sqr(&t2, &t2)            // t2  = t2^2
		mul(x, &t2, &t4)         // X3p = t2 * t4
		sub(&t1, &t3, &t1)       // t1  = t3 - t1
		sqr(&t1, &t1)            // t1  = t1^2
		mul(z, &t1, &t0)         // Z3p = t1 * t0
	}
}

// Set (y1, y2, y3)  = (1/x1, 1/x2, 1/x3).
//
// All xi, yi must be distinct.
func Fp2Batch3Inv(x1, x2, x3, y1, y2, y3 *Fp2) {
	var x1x2, t Fp2

	mul(&x1x2, x1, x2) // x1*x2
	mul(&t, &x1x2, x3) // 1/(x1*x2*x3)
	inv(&t, &t)
	mul(y1, &t, x2) // 1/x1
	mul(y1, y1, x3)
	mul(y2, &t, x1) // 1/x2
	mul(y2, y2, x3)
	mul(y3, &t, &x1x2) // 1/x3
}

// Scalarmul3Pt is a right-to-left point multiplication that given the
// x-coordinate of P, Q and P-Q calculates the x-coordinate of R=Q+[scalar]P.
// nbits must be smaller or equal to len(scalar).
func ScalarMul3Pt(cparams *ProjectiveCurveParameters, P, Q, PmQ *ProjectivePoint, nbits uint, scalar []uint8) ProjectivePoint {
	var R0, R2, R1 ProjectivePoint
	aPlus2Over4 := CalcAplus2Over4(cparams)
	R1 = *P
	R2 = *PmQ
	R0 = *Q

	// Iterate over the bits of the scalar, bottom to top
	prevBit := uint8(0)
	for i := uint(0); i < nbits; i++ {
		bit := (scalar[i>>3] >> (i & 7) & 1)
		swap := prevBit ^ bit
		prevBit = bit
		cswap(&R1.X, &R1.Z, &R2.X, &R2.Z, swap)
		R0, R2 = xDbladd(&R0, &R2, &R1, &aPlus2Over4)
	}
	cswap(&R1.X, &R1.Z, &R2.X, &R2.Z, prevBit)
	return R1
}

// Given a three-torsion point p = x(PB) on the curve E_(A:C), construct the
// three-isogeny phi : E_(A:C) -> E_(A:C)/<P_3> = E_(A':C').
//
// Input: (XP_3: ZP_3), where P_3 has exact order 3 on E_A/C
// Output:
//   - Curve coordinates (A' + 2C', A' - 2C') corresponding to E_A'/C' = A_E/C/<P3>
//   - Isogeny phi with constants in F_p^2
func (phi *isogeny3) GenerateCurve(p *ProjectivePoint) CurveCoefficientsEquiv {
	var t0, t1, t2, t3, t4 Fp2
	var coefEq CurveCoefficientsEquiv
	K1, K2 := &phi.K1, &phi.K2

	sub(K1, &p.X, &p.Z)            // K1 = XP3 - ZP3
	sqr(&t0, K1)                   // t0 = K1^2
	add(K2, &p.X, &p.Z)            // K2 = XP3 + ZP3
	sqr(&t1, K2)                   // t1 = K2^2
	add(&t2, &t0, &t1)             // t2 = t0 + t1
	add(&t3, K1, K2)               // t3 = K1 + K2
	sqr(&t3, &t3)                  // t3 = t3^2
	sub(&t3, &t3, &t2)             // t3 = t3 - t2
	add(&t2, &t1, &t3)             // t2 = t1 + t3
	add(&t3, &t3, &t0)             // t3 = t3 + t0
	add(&t4, &t3, &t0)             // t4 = t3 + t0
	add(&t4, &t4, &t4)             // t4 = t4 + t4
	add(&t4, &t1, &t4)             // t4 = t1 + t4
	mul(&coefEq.C, &t2, &t4)       // A24m = t2 * t4
	add(&t4, &t1, &t2)             // t4 = t1 + t2
	add(&t4, &t4, &t4)             // t4 = t4 + t4
	add(&t4, &t0, &t4)             // t4 = t0 + t4
	mul(&t4, &t3, &t4)             // t4 = t3 * t4
	sub(&t0, &t4, &coefEq.C)       // t0 = t4 - A24m
	add(&coefEq.A, &coefEq.C, &t0) // A24p = A24m + t0
	return coefEq
}

// Given a 3-isogeny phi and a point pB = x(PB), compute x(QB), the x-coordinate
// of the image QB = phi(PB) of PB under phi : E_(A:C) -> E_(A':C').
//
// The output xQ = x(Q) is then a point on the curve E_(A':C'); the curve
// parameters are returned by the GenerateCurve function used to construct phi.
func (phi *isogeny3) EvaluatePoint(p *ProjectivePoint) {
	var t0, t1, t2 Fp2
	K1, K2 := &phi.K1, &phi.K2
	px, pz := &p.X, &p.Z

	add(&t0, px, pz)   // t0 = XQ + ZQ
	sub(&t1, px, pz)   // t1 = XQ - ZQ
	mul(&t0, K1, &t0)  // t2 = K1 * t0
	mul(&t1, K2, &t1)  // t1 = K2 * t1
	add(&t2, &t0, &t1) // t2 = t0 + t1
	sub(&t0, &t1, &t0) // t0 = t1 - t0
	sqr(&t2, &t2)      // t2 = t2 ^ 2
	sqr(&t0, &t0)      // t0 = t0 ^ 2
	mul(px, px, &t2)   // XQ'= XQ * t2
	mul(pz, pz, &t0)   // ZQ'= ZQ * t0
}

// Given a four-torsion point p = x(PB) on the curve E_(A:C), construct the
// four-isogeny phi : E_(A:C) -> E_(A:C)/<P_4> = E_(A':C').
//
// Input: (XP_4: ZP_4), where P_4 has exact order 4 on E_A/C
// Output:
//   - Curve coordinates (A' + 2C', 4C') corresponding to E_A'/C' = A_E/C/<P4>
//   - Isogeny phi with constants in F_p^2
func (phi *isogeny4) GenerateCurve(p *ProjectivePoint) CurveCoefficientsEquiv {
	var coefEq CurveCoefficientsEquiv
	xp4, zp4 := &p.X, &p.Z
	K1, K2, K3 := &phi.K1, &phi.K2, &phi.K3

	sub(K2, xp4, zp4)
	add(K3, xp4, zp4)
	sqr(K1, zp4)
	add(K1, K1, K1)
	sqr(&coefEq.C, K1)
	add(K1, K1, K1)
	sqr(&coefEq.A, xp4)
	add(&coefEq.A, &coefEq.A, &coefEq.A)
	sqr(&coefEq.A, &coefEq.A)
	return coefEq
}

// Given a 4-isogeny phi and a point xP = x(P), compute x(Q), the x-coordinate
// of the image Q = phi(P) of P under phi : E_(A:C) -> E_(A':C').
//
// Input: Isogeny returned by GenerateCurve and point q=(Qx,Qz) from E0_A/C
// Output: Corresponding point q from E1_A'/C', where E1 is 4-isogenous to E0
func (phi *isogeny4) EvaluatePoint(p *ProjectivePoint) {
	var t0, t1 Fp2
	xq, zq := &p.X, &p.Z
	K1, K2, K3 := &phi.K1, &phi.K2, &phi.K3

	add(&t0, xq, zq)
	sub(&t1, xq, zq)
	mul(xq, &t0, K2)
	mul(zq, &t1, K3)
	mul(&t0, &t0, &t1)
	mul(&t0, &t0, K1)
	add(&t1, xq, zq)
	sub(zq, xq, zq)
	sqr(&t1, &t1)
	sqr(zq, zq)
	add(xq, &t0, &t1)
	sub(&t0, zq, &t0)
	mul(xq, xq, &t1)
	mul(zq, zq, &t0)
}

// PublicKeyValidation preforms public key/ciphertext validation using the CLN test.
// CLN test: Check that P and Q are both of order 3^e3 and they generate the torsion E_A[3^e3]
// A countermeasure for remote timing attacks on SIKE; suggested by https://eprint.iacr.org/2022/054.pdf
// Any curve E_A (SIKE 434, 503, 751) that passes CLN test is supersingular.
// Input: The public key / ciphertext P, Q, PmQ. The projective coordinate A of the curve defined by (P, Q, PmQ)
// Outputs: Whether (P,Q,PmQ) follows the CLN test
func PublicKeyValidation(cparams *ProjectiveCurveParameters, P, Q, PmQ *ProjectivePoint, nbits uint) error {

	var PmQX, PmQZ Fp2
	FromMontgomery(&PmQX, &PmQ.X)
	FromMontgomery(&PmQZ, &PmQ.Z)

	// PmQ is not point T or O
	if (isZero(&PmQX) == 1) || (isZero(&PmQZ) == 1) {
		return errors.New("curve: PmQ is invalid")
	}

	cparam := CalcCurveParamsEquiv3(cparams)

	// Compute e_3 = log3(2^(nbits+1))
	var e3 uint32
	e3_float := float64(int(nbits)+1) / math.Log2(3)
	e3 = uint32(e3_float)

	// Verify that P and Q generate E_A[3^e_3] by checking: [3^(e_3-1)]P != [+-3^(e_3-1)]Q
	var test_P, test_Q ProjectivePoint
	test_P = *P
	test_Q = *Q

	Pow3k(&test_P, &cparam, e3-1)
	Pow3k(&test_Q, &cparam, e3-1)

	var PZ, QZ Fp2
	FromMontgomery(&PZ, &test_P.Z)
	FromMontgomery(&QZ, &test_Q.Z)

	// P, Q are not of full order 3^e_3
	if (isZero(&PZ) == 1) || (isZero(&QZ) == 1) {
		return errors.New("curve: ciphertext/public key are not of full order 3^e3")
	}

	// PX/PZ = affine(PX)
	// QX/QZ = affine(QX)
	// If PX/PZ = QX/QZ, we have P=+-Q
	var PXQZ_PZQX_fromMont, PXQZ_PZQX, PXQZ, PZQX Fp2
	mul(&PXQZ, &test_P.X, &test_Q.Z)
	mul(&PZQX, &test_P.Z, &test_Q.X)
	sub(&PXQZ_PZQX, &PXQZ, &PZQX)
	FromMontgomery(&PXQZ_PZQX_fromMont, &PXQZ_PZQX)

	// [3^(e_3-1)]P == [+-3^(e_3-1)]Q
	if isZero(&PXQZ_PZQX_fromMont) == 1 {
		return errors.New("curve: ciphertext/public key are not linearly independent")
	}

	// Check that Ord(P) = Ord(Q) = 3^(e_3)
	Pow3k(&test_P, &cparam, 1)
	Pow3k(&test_Q, &cparam, 1)

	FromMontgomery(&PZ, &test_P.Z)
	FromMontgomery(&QZ, &test_Q.Z)

	// P, Q are not of correct order 3^e_3
	if (isZero(&PZ) == 0) || (isZero(&QZ) == 0) {
		return errors.New("curve: ciphertext/public key are not of correct order 3^e3")
	}
	return nil
}
