//go:build !amd64 || purego
// +build !amd64 purego

package x25519

import fp "github.com/JI-0/circl-ed448/math/fp25519"

func double(x, z *fp.Elt)             { doubleGeneric(x, z) }
func diffAdd(w *[5]fp.Elt, b uint)    { diffAddGeneric(w, b) }
func ladderStep(w *[5]fp.Elt, b uint) { ladderStepGeneric(w, b) }
func mulA24(z, x *fp.Elt)             { mulA24Generic(z, x) }
