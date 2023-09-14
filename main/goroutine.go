package main

import "fmt"

func main() {
	//gomaxprocs := runtime.GOMAXPROCS(0) // 8
	//fmt.Println(gomaxprocs)
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		fmt.Println(i)
	//	}()
	//}

	ch := make(chan int) // runtime.makechan(SB)
	ch <- 1              // runtime.chansend1(SB)
	select {
	case v := <-ch:
		fmt.Println(v) // runtime.selectnbrecv(SB)
	default:
	}
}

/*
select

PS E:\gothmslee\golang\main> go build -gcflags=-S goroutine.go
# command-line-arguments
main.main STEXT size=159 args=0x0 locals=0x50 funcid=0x0 align=0x0
        0x0000 00000 (E:\gothmslee\golang\main\goroutine.go:5)  TEXT    main.main(SB), ABIInternal, $80-0
        0x0000 00000 (E:\gothmslee\golang\main\goroutine.go:5)  CMPQ    SP, 16(R14)
        0x0004 00004 (E:\gothmslee\golang\main\goroutine.go:5)  PCDATA  $0, $-2
        0x0004 00004 (E:\gothmslee\golang\main\goroutine.go:5)  JLS     149
        0x000a 00010 (E:\gothmslee\golang\main\goroutine.go:5)  PCDATA  $0, $-1
        0x000a 00010 (E:\gothmslee\golang\main\goroutine.go:5)  SUBQ    $80, SP
        0x000e 00014 (E:\gothmslee\golang\main\goroutine.go:5)  MOVQ    BP, 72(SP)
        0x0013 00019 (E:\gothmslee\golang\main\goroutine.go:5)  LEAQ    72(SP), BP
        0x0018 00024 (E:\gothmslee\golang\main\goroutine.go:5)  FUNCDATA        $0, gclocals·ykHN0vawYuq1dUW4zEe2gA==(SB)
        0x0018 00024 (E:\gothmslee\golang\main\goroutine.go:5)  FUNCDATA        $1, gclocals·SAgskpdBM8mlWzn6XNUFrQ==(SB)
        0x0018 00024 (E:\gothmslee\golang\main\goroutine.go:5)  FUNCDATA        $2, main.main.stkobj(SB)
        0x0018 00024 (E:\gothmslee\golang\main\goroutine.go:14) LEAQ    type:chan int(SB), AX
        0x001f 00031 (E:\gothmslee\golang\main\goroutine.go:14) XORL    BX, BX
        0x0021 00033 (E:\gothmslee\golang\main\goroutine.go:14) PCDATA  $1, $0
        0x0021 00033 (E:\gothmslee\golang\main\goroutine.go:14) CALL    runtime.makechan(SB)
        0x0026 00038 (E:\gothmslee\golang\main\goroutine.go:14) MOVQ    AX, main.ch+48(SP)
        0x002b 00043 (E:\gothmslee\golang\main\goroutine.go:15) LEAQ    main..stmp_0(SB), BX
        0x0032 00050 (E:\gothmslee\golang\main\goroutine.go:15) PCDATA  $1, $1
        0x0032 00050 (E:\gothmslee\golang\main\goroutine.go:15) CALL    runtime.chansend1(SB)
        0x0037 00055 (E:\gothmslee\golang\main\goroutine.go:17) LEAQ    main..autotmp_8+40(SP), AX
        0x003c 00060 (E:\gothmslee\golang\main\goroutine.go:17) MOVQ    main.ch+48(SP), BX
        0x0041 00065 (E:\gothmslee\golang\main\goroutine.go:17) PCDATA  $1, $0
        0x0041 00065 (E:\gothmslee\golang\main\goroutine.go:17) CALL    runtime.selectnbrecv(SB)
        0x0046 00070 (E:\gothmslee\golang\main\goroutine.go:17) TESTB   AL, AL
        0x0048 00072 (E:\gothmslee\golang\main\goroutine.go:17) JEQ     139
        0x004a 00074 (E:\gothmslee\golang\main\goroutine.go:17) MOVQ    main..autotmp_8+40(SP), AX
        0x004f 00079 (E:\gothmslee\golang\main\goroutine.go:18) MOVUPS  X15, main..autotmp_13+56(SP)
        0x0055 00085 (E:\gothmslee\golang\main\goroutine.go:18) PCDATA  $1, $2
        0x0055 00085 (E:\gothmslee\golang\main\goroutine.go:18) CALL    runtime.convT64(SB)
        0x005a 00090 (E:\gothmslee\golang\main\goroutine.go:18) LEAQ    type:int(SB), CX
        0x0061 00097 (E:\gothmslee\golang\main\goroutine.go:18) MOVQ    CX, main..autotmp_13+56(SP)
        0x0066 00102 (E:\gothmslee\golang\main\goroutine.go:18) MOVQ    AX, main..autotmp_13+64(SP)
        0x006b 00107 (<unknown line number>)    NOP
        0x006b 00107 (E:\Go\src\fmt\print.go:314)       MOVQ    os.Stdout(SB), BX
        0x0072 00114 (E:\Go\src\fmt\print.go:314)       LEAQ    go:itab.*os.File,io.Writer(SB), AX
        0x0079 00121 (E:\Go\src\fmt\print.go:314)       LEAQ    main..autotmp_13+56(SP), CX
        0x007e 00126 (E:\Go\src\fmt\print.go:314)       MOVL    $1, DI
        0x0083 00131 (E:\Go\src\fmt\print.go:314)       MOVQ    DI, SI
        0x0086 00134 (E:\Go\src\fmt\print.go:314)       PCDATA  $1, $0
        0x0086 00134 (E:\Go\src\fmt\print.go:314)       CALL    fmt.Fprintln(SB)
        0x008b 00139 (E:\gothmslee\golang\main\goroutine.go:21) PCDATA  $1, $-1
        0x008b 00139 (E:\gothmslee\golang\main\goroutine.go:21) MOVQ    72(SP), BP
        0x0090 00144 (E:\gothmslee\golang\main\goroutine.go:21) ADDQ    $80, SP
        0x0094 00148 (E:\gothmslee\golang\main\goroutine.go:21) RET
        0x0095 00149 (E:\gothmslee\golang\main\goroutine.go:21) NOP
        0x0095 00149 (E:\gothmslee\golang\main\goroutine.go:5)  PCDATA  $1, $-1
        0x0095 00149 (E:\gothmslee\golang\main\goroutine.go:5)  PCDATA  $0, $-2
        0x0095 00149 (E:\gothmslee\golang\main\goroutine.go:5)  CALL    runtime.morestack_noctxt(SB)
        0x009a 00154 (E:\gothmslee\golang\main\goroutine.go:5)  PCDATA  $0, $-1
        0x009a 00154 (E:\gothmslee\golang\main\goroutine.go:5)  JMP     0
*/
