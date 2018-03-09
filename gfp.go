package bn256

import (
	"fmt"
	"math/big"
)

type gfP [4]uint64

func newGFp(x int64) (out *gfP) {
	if x >= 0 {
		out = &gfP{uint64(x)}
	} else {
		out = &gfP{uint64(-x)}
		gfpNeg(out, out)
	}

	montEncode(out, out)
	return out
}

func (e *gfP) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x", e[3], e[2], e[1], e[0])
}

func (e *gfP) Set(f *gfP) {
	e[0] = f[0]
	e[1] = f[1]
	e[2] = f[2]
	e[3] = f[3]
}

func (e *gfP) Invert(f *gfP) {
	bits := [4]uint64{0x185cac6c5e089665, 0xee5b88d120b5b59e, 0xaa6fecb86184dc21, 0x8fb501e34aa387f9}

	sum, power := &gfP{}, &gfP{}
	sum.Set(rN1)
	power.Set(f)

	for word := 0; word < 4; word++ {
		for bit := uint(0); bit < 64; bit++ {
			if (bits[word]>>bit)&1 == 1 {
				gfpMul(sum, sum, power)
			}
			gfpMul(power, power, power)
		}
	}

	gfpMul(sum, sum, r3)
	e.Set(sum)
}

func (e *gfP) Exp(base *gfP, power *big.Int) {
	sum := &gfP{1}
	t := &gfP{0}
	for i := power.BitLen() - 1; i >= 0; i-- {
		gfpMul(t, sum, sum)
		if power.Bit(i) != 0 {
			gfpMul(sum, t, base)
		} else {
			sum.Set(t)
		}
	}
	e.Set(sum)
}

func isQuadraticResidue(number *gfP) bool {
	exp := new(big.Int).Set(p)
	exp.Sub(exp, new(big.Int).SetUint64(1))
	exp.Div(exp, new(big.Int).SetUint64(2))
	e := &gfP{0}
	e.Exp(number, exp)
	fmt.Println(e)
	if e[0] == 1 {
		return true
	}
	return false
}

func (e *gfP) calcQuadraticResidue(number *gfP) {
	k := new(big.Int).Add(p, new(big.Int).SetInt64(1))
	k.Div(k, new(big.Int).SetInt64(4))
	e.Exp(number, k)
}

func (e *gfP) Marshal(out []byte) {
	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			out[8*w+b] = byte(e[3-w] >> (56 - 8*b))
		}
	}
}

func (e *gfP) Unmarshal(in []byte) {
	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			e[3-w] += uint64(in[8*w+b]) << (56 - 8*b)
		}
	}
}

func montEncode(c, a *gfP) { gfpMul(c, a, r2) }
func montDecode(c, a *gfP) { gfpMul(c, a, &gfP{1}) }

// go:noescape
func gfpNeg(c, a *gfP)

//go:noescape
func gfpAdd(c, a, b *gfP)

//go:noescape
func gfpSub(c, a, b *gfP)

//go:noescape
func gfpMul(c, a, b *gfP)
