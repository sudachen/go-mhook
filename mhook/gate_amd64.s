// +build amd64,windows,!gccgo

TEXT _setupWin64Gate_1(SB), 4, $0
        MOVQ   	CX,  8(SP)
        MOVQ   	DX, 16(SP)
        MOVQ   	R8, 24(SP)
        MOVQ  	R9, 32(SP)
        LEAQ	8(SP), DX
        XORQ 	AX, AX
        PUSHQ 	AX
        MOVQ	SP, CX
        SUBQ	$32, SP
        CALL	*R10
        ADDQ	$32, SP
        XCHGQ	(SP), AX
        CMPQ	AX, 0
        JNE     L
        POPQ	AX
        RET
L:
	POPQ	CX
	MOVQ   	8(SP),  CX
	MOVQ   	16(SP), DX
	MOVQ   	24(SP), R8
	MOVQ   	32(SP), R9
	JMP	*AX

TEXT _setupWin64Gate_2(SB), 4, $0
	MOVQ 	$0xf000000000000001, R10
	MOVQ	$0xf000000000000001, AX
	JMP	*AX
	XCHGL 	AX, AX
	XCHGL 	AX ,AX

// func setupWin64Gate(mem uintptr,cb uintptr)
TEXT ·setupWin64Gate(SB), 0, $0
	MOVQ 	_setupWin64Gate_2(SB), SI
	MOVQ 	mem+0(FP), DI
	MOVQ 	(SI), AX
	MOVQ 	AX, (DI)
	MOVQ 	8(SI), AX
	MOVQ 	AX, 8(DI)
	MOVQ 	16(SI), AX
	MOVQ 	AX, 16(DI)
	MOVQ 	cb+8(FP),AX
	MOVQ 	AX, 2(DI)
	MOVQ 	_setupWin64Gate_1(SB), AX
	MOVQ 	AX, 12(DI)
	RET


