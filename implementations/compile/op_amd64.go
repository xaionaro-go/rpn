package rpn

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/nelhage/gojit"
	asm "github.com/twitchyliquid64/golang-asm"
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/obj/x86"
	"github.com/xaionaro-go/rpn/types"
)

func noop(builder *asm.Builder) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.ANOPL
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = x86.REG_AX
	return prog
}

func fAddDPImmediate(builder *asm.Builder, in float64) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AFADDDP
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = x86.REG_AL
	prog.From.Type = obj.TYPE_CONST
	prog.From.Offset = int64(math.Float64bits(in))
	return prog
}

func fMovDPImmediate(builder *asm.Builder, reg int16, in float64) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AFMOVDP
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = reg
	prog.From.Type = obj.TYPE_CONST
	prog.From.Offset = int64(math.Float64bits(in))
	return prog
}

func addQImmediateConst(builder *asm.Builder, reg int16, in int64) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AADDQ
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = reg
	prog.From.Type = obj.TYPE_CONST
	prog.From.Offset = in
	return prog
}

func movQImmediateConst(builder *asm.Builder, reg int16, in int64) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AMOVQ
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = reg
	prog.From.Type = obj.TYPE_CONST
	prog.From.Offset = in
	return prog
}

func subQImmediateConst(builder *asm.Builder, reg int16, in int64) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.ASUBQ
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = reg
	prog.From.Type = obj.TYPE_CONST
	prog.From.Offset = in
	return prog
}

func movQImmediate(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AMOVQ
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func movQ(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AMOVQ
	prog.To.Type = obj.TYPE_MEM
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_MEM
	prog.From.Reg = regFrom
	return prog
}

func load(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AMOVQ
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_MEM
	prog.From.Reg = regFrom
	return prog
}

func loadSDOffset(builder *asm.Builder, regTo, regFrom int16, offset int64) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AMOVSD
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_MEM
	prog.From.Reg = regFrom
	prog.From.Offset = offset
	return prog
}

func store(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AMOVQ
	prog.To.Type = obj.TYPE_MEM
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func storeSDOffset(builder *asm.Builder, regTo, regFrom int16, offset int64) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AMOVSD
	prog.To.Type = obj.TYPE_MEM
	prog.To.Reg = regTo
	prog.To.Offset = offset
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func addQ(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AADDQ
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func fAddDP(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AFADDDP
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func fSubDP(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AFSUBDP
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func fMulDP(builder *asm.Builder, regFrom, regTo int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AFMULDP
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func fDivDP(builder *asm.Builder, regFrom, regTo int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AFDIVDP
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func addSD(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AADDSD
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func subSD(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.ASUBSD
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func mulSD(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.AMULSD
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func divSD(builder *asm.Builder, regTo, regFrom int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.ADIVSD
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = regTo
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = regFrom
	return prog
}

func pushQ(builder *asm.Builder, reg int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.APUSHQ
	prog.From.Type = obj.TYPE_REG
	prog.From.Reg = reg
	return prog
}

func popQ(builder *asm.Builder, reg int16) *obj.Prog {
	prog := builder.NewProg()
	prog.As = x86.APOPQ
	prog.To.Type = obj.TYPE_REG
	prog.To.Reg = reg
	return prog
}

func ret(builder *asm.Builder) *obj.Prog {
	prog := builder.NewProg()
	prog.As = obj.ARET
	return prog
}

// Compile converts ops to a native code which could be executed by calling
// function `eval`. It will always read incoming values from the pointer
// stored in slice `valuesRaw`.
func (ops Ops) Compile(stackRaw []float64, valuesRaw []float64) (eval func() float64, cleanup func()) {
	// See also: http://staffwww.fullcoll.edu/aclifton/cs241/lecture-floating-point-simd.html

	stackPtr := uint64((*reflect.SliceHeader)(unsafe.Pointer(&stackRaw)).Data)
	valuesPtr := uint64((*reflect.SliceHeader)(unsafe.Pointer(&valuesRaw)).Data)

	builder, _ := asm.NewBuilder("amd64", 64)

	tempReg := int16(x86.REG_AX)
	stackPtrReg := int16(x86.REG_DI)
	valuesPtrReg := int16(x86.REG_SI)

	itemSize := int64(unsafe.Sizeof(float64(0)))
	builder.AddInstruction(pushQ(builder, x86.REG_BP))
	builder.AddInstruction(movQImmediate(builder, x86.REG_BP, x86.REG_SP))
	builder.AddInstruction(subQImmediateConst(builder, x86.REG_SP, 32))
	builder.AddInstruction(storeSDOffset(builder, x86.REG_SP, x86.REG_X0, 0))
	builder.AddInstruction(storeSDOffset(builder, x86.REG_SP, x86.REG_X1, 16))
	builder.AddInstruction(pushQ(builder, tempReg))
	builder.AddInstruction(pushQ(builder, stackPtrReg))
	builder.AddInstruction(pushQ(builder, valuesPtrReg))

	builder.AddInstruction(movQImmediateConst(builder, stackPtrReg, int64(stackPtr)))
	builder.AddInstruction(movQImmediateConst(builder, valuesPtrReg, int64(valuesPtr)))

	for _, op := range ops {
		if op == types.OpFetch {
			builder.AddInstruction(load(builder, tempReg, valuesPtrReg))
			builder.AddInstruction(addQImmediateConst(builder, valuesPtrReg, itemSize))
			builder.AddInstruction(store(builder, stackPtrReg, tempReg))
			builder.AddInstruction(addQImmediateConst(builder, stackPtrReg, itemSize))
			continue
		}

		builder.AddInstruction(addQImmediateConst(builder, stackPtrReg, -itemSize))
		builder.AddInstruction(load(builder, tempReg, stackPtrReg))
		builder.AddInstruction(movQImmediate(builder, x86.REG_X1, tempReg))
		builder.AddInstruction(addQImmediateConst(builder, stackPtrReg, -itemSize))
		builder.AddInstruction(load(builder, tempReg, stackPtrReg))
		builder.AddInstruction(movQImmediate(builder, x86.REG_X0, tempReg))
		switch op {
		case types.OpPlus:
			builder.AddInstruction(addSD(builder, x86.REG_X0, x86.REG_X1))
		case types.OpMinus:
			builder.AddInstruction(subSD(builder, x86.REG_X0, x86.REG_X1))
		case types.OpMultiply:
			builder.AddInstruction(mulSD(builder, x86.REG_X0, x86.REG_X1))
		case types.OpDivide:
			builder.AddInstruction(divSD(builder, x86.REG_X0, x86.REG_X1))
		case types.OpPower:
			panic("not implemented")
		case types.OpIf:
			panic("not implemented")
		}
		builder.AddInstruction(movQImmediate(builder, tempReg, x86.REG_X0))
		builder.AddInstruction(store(builder, stackPtrReg, tempReg))
		builder.AddInstruction(addQImmediateConst(builder, stackPtrReg, itemSize))
	}
	builder.AddInstruction(popQ(builder, valuesPtrReg))
	builder.AddInstruction(popQ(builder, stackPtrReg))
	builder.AddInstruction(popQ(builder, tempReg))
	builder.AddInstruction(loadSDOffset(builder, x86.REG_X1, x86.REG_SP, 16))
	builder.AddInstruction(loadSDOffset(builder, x86.REG_X0, x86.REG_SP, 0))
	builder.AddInstruction(movQImmediate(builder, x86.REG_SP, x86.REG_BP))
	builder.AddInstruction(popQ(builder, x86.REG_BP))

	builder.AddInstruction(ret(builder))

	code := builder.Assemble()
	if len(code) > gojit.PageSize {
		panic("too large code")
	}
	b, e := gojit.Alloc(gojit.PageSize)
	if e != nil {
		panic(e)
	}
	copy(b, code)

	var fn func()
	gojit.BuildTo(b, &fn)
	eval = func() float64 {
		fn()
		return stackRaw[0]
	}
	cleanup = func() {
		err := gojit.Release(b)
		if err != nil {
			panic(err)
		}
	}
	return
}
