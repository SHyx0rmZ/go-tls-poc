#include "textflag.h"

// This has been automatically generated by running: go run tls_gen.go
TEXT ·goid(SB),NOSPLIT,$-8
	MOVQ TLS, AX
	MOVQ 0(AX)(TLS*1), AX
	MOVQ 152(AX), AX
	MOVQ AX, ret+0(FP)
	RET