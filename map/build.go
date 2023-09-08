package source

/*
$ go build -gcflags=-S map.go

PS E:\gothmslee\golang\main> go build -gcflags=-S map.go
# command-line-arguments
main.main STEXT size=426 args=0x0 locals=0x1c0 funcid=0x0 align=0x0
        0x0000 00000 (E:\gothmslee\golang\main\map.go:8)        TEXT    main.main(SB), ABIInternal, $448-0
        0x0000 00000 (E:\gothmslee\golang\main\map.go:8)        LEAQ    -320(SP), R12
        0x0008 00008 (E:\gothmslee\golang\main\map.go:8)        CMPQ    R12, 16(R14)
        0x000c 00012 (E:\gothmslee\golang\main\map.go:8)        PCDATA  $0, $-2
        0x000c 00012 (E:\gothmslee\golang\main\map.go:8)        JLS     413
        0x0012 00018 (E:\gothmslee\golang\main\map.go:8)        PCDATA  $0, $-1
        0x0012 00018 (E:\gothmslee\golang\main\map.go:8)        SUBQ    $448, SP
        0x0019 00025 (E:\gothmslee\golang\main\map.go:8)        MOVQ    BP, 440(SP)
        0x0021 00033 (E:\gothmslee\golang\main\map.go:8)        LEAQ    440(SP), BP
        0x0029 00041 (E:\gothmslee\golang\main\map.go:8)        FUNCDATA        $0, gclocals·DsEZEqsV1VFgO2VgUUolkQ==(SB)
        0x0029 00041 (E:\gothmslee\golang\main\map.go:8)        FUNCDATA        $1, gclocals·HFbfPPPbONK70RRxfAk8Uw==(SB)
        0x0029 00041 (E:\gothmslee\golang\main\map.go:8)        FUNCDATA        $2, main.main.stkobj(SB)
        0x0029 00041 (E:\gothmslee\golang\main\map.go:9)        MOVUPS  X15, main..autotmp_12+88(SP)
        0x002f 00047 (E:\gothmslee\golang\main\map.go:9)        MOVUPS  X15, main..autotmp_12+104(SP)
        0x0035 00053 (E:\gothmslee\golang\main\map.go:9)        MOVUPS  X15, main..autotmp_12+120(SP)
        0x003b 00059 (E:\gothmslee\golang\main\map.go:9)        LEAQ    main..autotmp_13+232(SP), DI
        0x0043 00067 (E:\gothmslee\golang\main\map.go:9)        PCDATA  $0, $-2
        0x0043 00067 (E:\gothmslee\golang\main\map.go:9)        LEAQ    -48(DI), DI
        0x0047 00071 (E:\gothmslee\golang\main\map.go:9)        DUFFZERO        $290
        0x005a 00090 (E:\gothmslee\golang\main\map.go:9)        PCDATA  $0, $-1
        0x005a 00090 (E:\gothmslee\golang\main\map.go:9)        LEAQ    main..autotmp_13+232(SP), AX
        0x0062 00098 (E:\gothmslee\golang\main\map.go:9)        MOVQ    AX, main..autotmp_12+104(SP)
        0x0067 00103 (E:\gothmslee\golang\main\map.go:9)        PCDATA  $1, $1
        0x0067 00103 (E:\gothmslee\golang\main\map.go:9)        CALL    runtime.fastrand(SB)
        0x006c 00108 (E:\gothmslee\golang\main\map.go:9)        MOVL    AX, main..autotmp_12+100(SP)
        0x0070 00112 (E:\gothmslee\golang\main\map.go:10)       LEAQ    type:map[int]string(SB), AX
        0x0077 00119 (E:\gothmslee\golang\main\map.go:10)       LEAQ    main..autotmp_12+88(SP), BX
        0x007c 00124 (E:\gothmslee\golang\main\map.go:10)       MOVL    $1, CX
        0x0081 00129 (E:\gothmslee\golang\main\map.go:10)       CALL    runtime.mapassign_fast64(SB)
        0x0086 00134 (E:\gothmslee\golang\main\map.go:10)       MOVQ    $4, 8(AX)
        0x008e 00142 (E:\gothmslee\golang\main\map.go:10)       PCDATA  $0, $-2
        0x008e 00142 (E:\gothmslee\golang\main\map.go:10)       CMPL    runtime.writeBarrier(SB), $0
        0x0095 00149 (E:\gothmslee\golang\main\map.go:10)       JNE     163
        0x0097 00151 (E:\gothmslee\golang\main\map.go:10)       LEAQ    go:string."haha"(SB), DX
        0x009e 00158 (E:\gothmslee\golang\main\map.go:10)       MOVQ    DX, (AX)
        0x00a1 00161 (E:\gothmslee\golang\main\map.go:10)       JMP     178
        0x00a3 00163 (E:\gothmslee\golang\main\map.go:10)       MOVQ    AX, DI
        0x00a6 00166 (E:\gothmslee\golang\main\map.go:10)       LEAQ    go:string."haha"(SB), DX
        0x00ad 00173 (E:\gothmslee\golang\main\map.go:10)       CALL    runtime.gcWriteBarrierDX(SB)
        0x00b2 00178 (E:\gothmslee\golang\main\map.go:19)       PCDATA  $0, $-1
        0x00b2 00178 (E:\gothmslee\golang\main\map.go:19)       LEAQ    main..autotmp_9+136(SP), DI
        0x00ba 00186 (E:\gothmslee\golang\main\map.go:19)       PCDATA  $0, $-2
        0x00ba 00186 (E:\gothmslee\golang\main\map.go:19)       LEAQ    -32(DI), DI
        0x00be 00190 (E:\gothmslee\golang\main\map.go:19)       NOP
        0x00c0 00192 (E:\gothmslee\golang\main\map.go:19)       DUFFZERO        $331
        0x00d3 00211 (E:\gothmslee\golang\main\map.go:19)       PCDATA  $0, $-1
        0x00d3 00211 (E:\gothmslee\golang\main\map.go:19)       LEAQ    type:map[int]string(SB), AX
        0x00da 00218 (E:\gothmslee\golang\main\map.go:19)       LEAQ    main..autotmp_12+88(SP), BX
        0x00df 00223 (E:\gothmslee\golang\main\map.go:19)       LEAQ    main..autotmp_9+136(SP), CX
        0x00e7 00231 (E:\gothmslee\golang\main\map.go:19)       PCDATA  $1, $2
        0x00e7 00231 (E:\gothmslee\golang\main\map.go:19)       CALL    runtime.mapiterinit(SB)
        0x00ec 00236 (E:\gothmslee\golang\main\map.go:19)       JMP     380
        0x00f1 00241 (E:\gothmslee\golang\main\map.go:19)       MOVQ    main..autotmp_9+144(SP), DX
        0x00f9 00249 (E:\gothmslee\golang\main\map.go:19)       MOVQ    (CX), AX
        0x00fc 00252 (E:\gothmslee\golang\main\map.go:19)       MOVQ    (DX), CX
        0x00ff 00255 (E:\gothmslee\golang\main\map.go:19)       MOVQ    CX, main.v.ptr+48(SP)
        0x0104 00260 (E:\gothmslee\golang\main\map.go:19)       MOVQ    8(DX), DX
        0x0108 00264 (E:\gothmslee\golang\main\map.go:19)       MOVQ    DX, main.v.len+40(SP)
        0x010d 00269 (E:\gothmslee\golang\main\map.go:20)       MOVUPS  X15, main..autotmp_20+56(SP)
        0x0113 00275 (E:\gothmslee\golang\main\map.go:20)       MOVUPS  X15, main..autotmp_20+72(SP)
        0x0119 00281 (E:\gothmslee\golang\main\map.go:20)       PCDATA  $1, $3
        0x0119 00281 (E:\gothmslee\golang\main\map.go:20)       CALL    runtime.convT64(SB)
        0x011e 00286 (E:\gothmslee\golang\main\map.go:20)       LEAQ    type:int(SB), CX
        0x0125 00293 (E:\gothmslee\golang\main\map.go:20)       MOVQ    CX, main..autotmp_20+56(SP)
        0x012a 00298 (E:\gothmslee\golang\main\map.go:20)       MOVQ    AX, main..autotmp_20+64(SP)
        0x012f 00303 (E:\gothmslee\golang\main\map.go:20)       MOVQ    main.v.ptr+48(SP), AX
        0x0134 00308 (E:\gothmslee\golang\main\map.go:20)       MOVQ    main.v.len+40(SP), BX
        0x0139 00313 (E:\gothmslee\golang\main\map.go:20)       PCDATA  $1, $4
        0x0139 00313 (E:\gothmslee\golang\main\map.go:20)       CALL    runtime.convTstring(SB)
        0x013e 00318 (E:\gothmslee\golang\main\map.go:20)       LEAQ    type:string(SB), CX
        0x0145 00325 (E:\gothmslee\golang\main\map.go:20)       MOVQ    CX, main..autotmp_20+72(SP)
        0x014a 00330 (E:\gothmslee\golang\main\map.go:20)       MOVQ    AX, main..autotmp_20+80(SP)
        0x014f 00335 (<unknown line number>)    NOP
        0x014f 00335 (E:\Go\src\fmt\print.go:314)       MOVQ    os.Stdout(SB), BX
        0x0156 00342 (E:\Go\src\fmt\print.go:314)       LEAQ    go:itab.*os.File,io.Writer(SB), AX
        0x015d 00349 (E:\Go\src\fmt\print.go:314)       MOVL    $2, DI
        0x0162 00354 (E:\Go\src\fmt\print.go:314)       MOVQ    DI, SI
        0x0165 00357 (E:\Go\src\fmt\print.go:314)       LEAQ    main..autotmp_20+56(SP), CX
        0x016a 00362 (E:\Go\src\fmt\print.go:314)       PCDATA  $1, $2
        0x016a 00362 (E:\Go\src\fmt\print.go:314)       CALL    fmt.Fprintln(SB)
        0x016f 00367 (E:\gothmslee\golang\main\map.go:19)       LEAQ    main..autotmp_9+136(SP), AX
        0x0177 00375 (E:\gothmslee\golang\main\map.go:19)       CALL    runtime.mapiternext(SB)
        0x017c 00380 (E:\gothmslee\golang\main\map.go:19)       MOVQ    main..autotmp_9+136(SP), CX
        0x0184 00388 (E:\gothmslee\golang\main\map.go:19)       TESTQ   CX, CX
        0x0187 00391 (E:\gothmslee\golang\main\map.go:19)       JNE     241
        0x018d 00397 (E:\gothmslee\golang\main\map.go:22)       PCDATA  $1, $-1
        0x018d 00397 (E:\gothmslee\golang\main\map.go:22)       MOVQ    440(SP), BP
        0x0195 00405 (E:\gothmslee\golang\main\map.go:22)       ADDQ    $448, SP
        0x019c 00412 (E:\gothmslee\golang\main\map.go:22)       RET
        0x019d 00413 (E:\gothmslee\golang\main\map.go:22)       NOP
        0x019d 00413 (E:\gothmslee\golang\main\map.go:8)        PCDATA  $1, $-1
        0x019d 00413 (E:\gothmslee\golang\main\map.go:8)        PCDATA  $0, $-2
        0x019d 00413 (E:\gothmslee\golang\main\map.go:8)        NOP
        0x01a0 00416 (E:\gothmslee\golang\main\map.go:8)        CALL    runtime.morestack_noctxt(SB)
        0x01a5 00421 (E:\gothmslee\golang\main\map.go:8)        PCDATA  $0, $-1
        0x01a5 00421 (E:\gothmslee\golang\main\map.go:8)        JMP     0
        0x0000 4c 8d a4 24 c0 fe ff ff 4d 3b 66 10 0f 86 8b 01  L..$....M;f.....
        0x0010 00 00 48 81 ec c0 01 00 00 48 89 ac 24 b8 01 00  ..H......H..$...
        0x0020 00 48 8d ac 24 b8 01 00 00 44 0f 11 7c 24 58 44  .H..$....D..|$XD
        0x0030 0f 11 7c 24 68 44 0f 11 7c 24 78 48 8d bc 24 e8  ..|$hD..|$xH..$.
        0x0040 00 00 00 48 8d 7f d0 48 89 6c 24 f0 48 8d 6c 24  ...H...H.l$.H.l$
        0x0050 f0 e8 00 00 00 00 48 8b 6d 00 48 8d 84 24 e8 00  ......H.m.H..$..
        0x0060 00 00 48 89 44 24 68 e8 00 00 00 00 89 44 24 64  ..H.D$h......D$d
        0x0070 48 8d 05 00 00 00 00 48 8d 5c 24 58 b9 01 00 00  H......H.\$X....
        0x0080 00 e8 00 00 00 00 48 c7 40 08 04 00 00 00 83 3d  ......H.@......=
        0x0090 00 00 00 00 00 75 0c 48 8d 15 00 00 00 00 48 89  .....u.H......H.
        0x00a0 10 eb 0f 48 89 c7 48 8d 15 00 00 00 00 e8 00 00  ...H..H.........
        0x00b0 00 00 48 8d bc 24 88 00 00 00 48 8d 7f e0 66 90  ..H..$....H...f.
        0x00c0 48 89 6c 24 f0 48 8d 6c 24 f0 e8 00 00 00 00 48  H.l$.H.l$......H
        0x00d0 8b 6d 00 48 8d 05 00 00 00 00 48 8d 5c 24 58 48  .m.H......H.\$XH
        0x00e0 8d 8c 24 88 00 00 00 e8 00 00 00 00 e9 8b 00 00  ..$.............
        0x00f0 00 48 8b 94 24 90 00 00 00 48 8b 01 48 8b 0a 48  .H..$....H..H..H
        0x0100 89 4c 24 30 48 8b 52 08 48 89 54 24 28 44 0f 11  .L$0H.R.H.T$(D..
        0x0110 7c 24 38 44 0f 11 7c 24 48 e8 00 00 00 00 48 8d  |$8D..|$H.....H.
        0x0120 0d 00 00 00 00 48 89 4c 24 38 48 89 44 24 40 48  .....H.L$8H.D$@H
        0x0130 8b 44 24 30 48 8b 5c 24 28 e8 00 00 00 00 48 8d  .D$0H.\$(.....H.
        0x0140 0d 00 00 00 00 48 89 4c 24 48 48 89 44 24 50 48  .....H.L$HH.D$PH
        0x0150 8b 1d 00 00 00 00 48 8d 05 00 00 00 00 bf 02 00  ......H.........
        0x0160 00 00 48 89 fe 48 8d 4c 24 38 e8 00 00 00 00 48  ..H..H.L$8.....H
        0x0170 8d 84 24 88 00 00 00 e8 00 00 00 00 48 8b 8c 24  ..$.........H..$
        0x0180 88 00 00 00 48 85 c9 0f 85 64 ff ff ff 48 8b ac  ....H....d...H..
        0x0190 24 b8 01 00 00 48 81 c4 c0 01 00 00 c3 0f 1f 00  $....H..........
        0x01a0 e8 00 00 00 00 e9 56 fe ff ff                    ......V...
        rel 3+0 t=23 type:int+0
        rel 3+0 t=23 type:string+0
        rel 3+0 t=23 type:*os.File+0
        rel 82+4 t=7 runtime.duffzero+290
        rel 104+4 t=7 runtime.fastrand+0
        rel 115+4 t=14 type:map[int]string+0
        rel 130+4 t=7 runtime.mapassign_fast64+0
        rel 144+4 t=14 runtime.writeBarrier+-1
        rel 154+4 t=14 go:string."haha"+0
        rel 169+4 t=14 go:string."haha"+0
        rel 174+4 t=7 runtime.gcWriteBarrierDX+0
        rel 203+4 t=7 runtime.duffzero+331
        rel 214+4 t=14 type:map[int]string+0
        rel 232+4 t=7 runtime.mapiterinit+0
        rel 282+4 t=7 runtime.convT64+0
        rel 289+4 t=14 type:int+0
        rel 314+4 t=7 runtime.convTstring+0
        rel 321+4 t=14 type:string+0
        rel 338+4 t=14 os.Stdout+0
        rel 345+4 t=14 go:itab.*os.File,io.Writer+0
        rel 363+4 t=7 fmt.Fprintln+0
        rel 376+4 t=7 runtime.mapiternext+0
        rel 417+4 t=7 runtime.morestack_noctxt+0
main.read STEXT size=79 args=0x8 locals=0x28 funcid=0x0 align=0x0
        0x0000 00000 (E:\gothmslee\golang\main\map.go:24)       TEXT    main.read(SB), ABIInternal, $40-8
        0x0000 00000 (E:\gothmslee\golang\main\map.go:24)       CMPQ    SP, 16(R14)
        0x0004 00004 (E:\gothmslee\golang\main\map.go:24)       PCDATA  $0, $-2
        0x0004 00004 (E:\gothmslee\golang\main\map.go:24)       JLS     62
        0x0006 00006 (E:\gothmslee\golang\main\map.go:24)       PCDATA  $0, $-1
        0x0006 00006 (E:\gothmslee\golang\main\map.go:24)       SUBQ    $40, SP
        0x000a 00010 (E:\gothmslee\golang\main\map.go:24)       MOVQ    BP, 32(SP)
        0x000f 00015 (E:\gothmslee\golang\main\map.go:24)       LEAQ    32(SP), BP
        0x0014 00020 (E:\gothmslee\golang\main\map.go:24)       FUNCDATA        $0, gclocals·ZzMiPAiVBg7DJ6dh/CjSag==(SB)
        0x0014 00020 (E:\gothmslee\golang\main\map.go:24)       FUNCDATA        $1, gclocals·VtCL4RdUwCqwXEPeyJllRA==(SB)
        0x0014 00020 (E:\gothmslee\golang\main\map.go:24)       FUNCDATA        $5, main.read.arginfo1(SB)
        0x0014 00020 (E:\gothmslee\golang\main\map.go:24)       FUNCDATA        $6, main.read.argliveinfo(SB)
        0x0014 00020 (E:\gothmslee\golang\main\map.go:24)       PCDATA  $3, $1
        0x0014 00020 (E:\gothmslee\golang\main\map.go:24)       MOVQ    AX, main..autotmp_1+24(SP)
        0x0019 00025 (E:\gothmslee\golang\main\map.go:26)       MOVQ    AX, BX
        0x001c 00028 (E:\gothmslee\golang\main\map.go:26)       MOVL    $1, CX
        0x0021 00033 (E:\gothmslee\golang\main\map.go:26)       LEAQ    type:map[int]string(SB), AX
        0x0028 00040 (E:\gothmslee\golang\main\map.go:26)       PCDATA  $1, $1
        0x0028 00040 (E:\gothmslee\golang\main\map.go:26)       CALL    runtime.mapaccess1_fast64(SB)
        0x002d 00045 (E:\gothmslee\golang\main\map.go:27)       MOVL    $1, AX
        0x0032 00050 (E:\gothmslee\golang\main\map.go:27)       CALL    time.Sleep(SB)
        0x0037 00055 (E:\gothmslee\golang\main\map.go:26)       MOVQ    main..autotmp_1+24(SP), AX
        0x003c 00060 (E:\gothmslee\golang\main\map.go:27)       JMP     25
        0x003e 00062 (E:\gothmslee\golang\main\map.go:27)       NOP
        0x003e 00062 (E:\gothmslee\golang\main\map.go:24)       PCDATA  $1, $-1
        0x003e 00062 (E:\gothmslee\golang\main\map.go:24)       PCDATA  $0, $-2
        0x003e 00062 (E:\gothmslee\golang\main\map.go:24)       MOVQ    AX, 8(SP)
        0x0043 00067 (E:\gothmslee\golang\main\map.go:24)       CALL    runtime.morestack_noctxt(SB)
        0x0048 00072 (E:\gothmslee\golang\main\map.go:24)       MOVQ    8(SP), AX
        0x004d 00077 (E:\gothmslee\golang\main\map.go:24)       PCDATA  $0, $-1
        0x004d 00077 (E:\gothmslee\golang\main\map.go:24)       JMP     0
        0x0000 49 3b 66 10 76 38 48 83 ec 28 48 89 6c 24 20 48  I;f.v8H..(H.l$ H
        0x0010 8d 6c 24 20 48 89 44 24 18 48 89 c3 b9 01 00 00  .l$ H.D$.H......
        0x0020 00 48 8d 05 00 00 00 00 e8 00 00 00 00 b8 01 00  .H..............
        0x0030 00 00 e8 00 00 00 00 48 8b 44 24 18 eb db 48 89  .......H.D$...H.
        0x0040 44 24 08 e8 00 00 00 00 48 8b 44 24 08 eb b1     D$......H.D$...
        rel 36+4 t=14 type:map[int]string+0
        rel 41+4 t=7 runtime.mapaccess1_fast64+0
        rel 51+4 t=7 time.Sleep+0
        rel 68+4 t=7 runtime.morestack_noctxt+0
main.write STEXT size=125 args=0x8 locals=0x28 funcid=0x0 align=0x0
        0x0000 00000 (E:\gothmslee\golang\main\map.go:31)       TEXT    main.write(SB), ABIInternal, $40-8
        0x0000 00000 (E:\gothmslee\golang\main\map.go:31)       CMPQ    SP, 16(R14)
        0x0004 00004 (E:\gothmslee\golang\main\map.go:31)       PCDATA  $0, $-2
        0x0004 00004 (E:\gothmslee\golang\main\map.go:31)       JLS     108
        0x0006 00006 (E:\gothmslee\golang\main\map.go:31)       PCDATA  $0, $-1
        0x0006 00006 (E:\gothmslee\golang\main\map.go:31)       SUBQ    $40, SP
        0x000a 00010 (E:\gothmslee\golang\main\map.go:31)       MOVQ    BP, 32(SP)
        0x000f 00015 (E:\gothmslee\golang\main\map.go:31)       LEAQ    32(SP), BP
        0x0014 00020 (E:\gothmslee\golang\main\map.go:31)       FUNCDATA        $0, gclocals·ZzMiPAiVBg7DJ6dh/CjSag==(SB)
        0x0014 00020 (E:\gothmslee\golang\main\map.go:31)       FUNCDATA        $1, gclocals·VtCL4RdUwCqwXEPeyJllRA==(SB)
        0x0014 00020 (E:\gothmslee\golang\main\map.go:31)       FUNCDATA        $5, main.write.arginfo1(SB)
        0x0014 00020 (E:\gothmslee\golang\main\map.go:31)       FUNCDATA        $6, main.write.argliveinfo(SB)
        0x0014 00020 (E:\gothmslee\golang\main\map.go:31)       PCDATA  $3, $1
        0x0014 00020 (E:\gothmslee\golang\main\map.go:31)       MOVQ    AX, main..autotmp_2+24(SP)
        0x0019 00025 (E:\gothmslee\golang\main\map.go:32)       JMP     42
        0x001b 00027 (E:\gothmslee\golang\main\map.go:34)       MOVL    $1, AX
        0x0020 00032 (E:\gothmslee\golang\main\map.go:34)       PCDATA  $1, $1
        0x0020 00032 (E:\gothmslee\golang\main\map.go:34)       CALL    time.Sleep(SB)
        0x0025 00037 (E:\gothmslee\golang\main\map.go:33)       MOVQ    main..autotmp_2+24(SP), AX
        0x002a 00042 (E:\gothmslee\golang\main\map.go:33)       MOVQ    AX, BX
        0x002d 00045 (E:\gothmslee\golang\main\map.go:33)       MOVL    $1, CX
        0x0032 00050 (E:\gothmslee\golang\main\map.go:33)       LEAQ    type:map[int]string(SB), AX
        0x0039 00057 (E:\gothmslee\golang\main\map.go:33)       CALL    runtime.mapassign_fast64(SB)
        0x003e 00062 (E:\gothmslee\golang\main\map.go:33)       MOVQ    $5, 8(AX)
        0x0046 00070 (E:\gothmslee\golang\main\map.go:33)       PCDATA  $0, $-2
        0x0046 00070 (E:\gothmslee\golang\main\map.go:33)       CMPL    runtime.writeBarrier(SB), $0
        0x004d 00077 (E:\gothmslee\golang\main\map.go:33)       JNE     91
        0x004f 00079 (E:\gothmslee\golang\main\map.go:33)       LEAQ    go:string."write"(SB), CX
        0x0056 00086 (E:\gothmslee\golang\main\map.go:33)       MOVQ    CX, (AX)
        0x0059 00089 (E:\gothmslee\golang\main\map.go:33)       JMP     27
        0x005b 00091 (E:\gothmslee\golang\main\map.go:33)       MOVQ    AX, DI
        0x005e 00094 (E:\gothmslee\golang\main\map.go:33)       LEAQ    go:string."write"(SB), CX
        0x0065 00101 (E:\gothmslee\golang\main\map.go:33)       CALL    runtime.gcWriteBarrierCX(SB)
        0x006a 00106 (E:\gothmslee\golang\main\map.go:33)       JMP     27
        0x006c 00108 (E:\gothmslee\golang\main\map.go:33)       NOP
        0x006c 00108 (E:\gothmslee\golang\main\map.go:31)       PCDATA  $1, $-1
        0x006c 00108 (E:\gothmslee\golang\main\map.go:31)       PCDATA  $0, $-2
        0x006c 00108 (E:\gothmslee\golang\main\map.go:31)       MOVQ    AX, 8(SP)
        0x0071 00113 (E:\gothmslee\golang\main\map.go:31)       CALL    runtime.morestack_noctxt(SB)
        0x0076 00118 (E:\gothmslee\golang\main\map.go:31)       MOVQ    8(SP), AX
        0x007b 00123 (E:\gothmslee\golang\main\map.go:31)       PCDATA  $0, $-1
        0x007b 00123 (E:\gothmslee\golang\main\map.go:31)       JMP     0
        0x0000 49 3b 66 10 76 66 48 83 ec 28 48 89 6c 24 20 48  I;f.vfH..(H.l$ H
        0x0010 8d 6c 24 20 48 89 44 24 18 eb 0f b8 01 00 00 00  .l$ H.D$........
        0x0020 e8 00 00 00 00 48 8b 44 24 18 48 89 c3 b9 01 00  .....H.D$.H.....
        0x0030 00 00 48 8d 05 00 00 00 00 e8 00 00 00 00 48 c7  ..H...........H.
        0x0040 40 08 05 00 00 00 83 3d 00 00 00 00 00 75 0c 48  @......=.....u.H
        0x0050 8d 0d 00 00 00 00 48 89 08 eb c0 48 89 c7 48 8d  ......H....H..H.
        0x0060 0d 00 00 00 00 e8 00 00 00 00 eb af 48 89 44 24  ............H.D$
        0x0070 08 e8 00 00 00 00 48 8b 44 24 08 eb 83           ......H.D$...
        rel 33+4 t=7 time.Sleep+0
        rel 53+4 t=14 type:map[int]string+0
        rel 58+4 t=7 runtime.mapassign_fast64+0
        rel 72+4 t=14 runtime.writeBarrier+-1
        rel 82+4 t=14 go:string."write"+0
        rel 97+4 t=14 go:string."write"+0
        rel 102+4 t=7 runtime.gcWriteBarrierCX+0
        rel 114+4 t=7 runtime.morestack_noctxt+0
type:.eq.[2]interface {} STEXT dupok size=170 args=0x10 locals=0x28 funcid=0x0 align=0x0
        0x0000 00000 (<autogenerated>:1)        TEXT    type:.eq.[2]interface {}(SB), DUPOK|ABIInternal, $40-16
        0x0000 00000 (<autogenerated>:1)        CMPQ    SP, 16(R14)
        0x0004 00004 (<autogenerated>:1)        PCDATA  $0, $-2
        0x0004 00004 (<autogenerated>:1)        JLS     140
        0x000a 00010 (<autogenerated>:1)        PCDATA  $0, $-1
        0x000a 00010 (<autogenerated>:1)        SUBQ    $40, SP
        0x000e 00014 (<autogenerated>:1)        MOVQ    BP, 32(SP)
        0x0013 00019 (<autogenerated>:1)        LEAQ    32(SP), BP
        0x0018 00024 (<autogenerated>:1)        FUNCDATA        $0, gclocals·TjPuuCwdlCpTaRQGRKTrYw==(SB)
        0x0018 00024 (<autogenerated>:1)        FUNCDATA        $1, gclocals·J5F+7Qw7O7ve2QcWC7DpeQ==(SB)
        0x0018 00024 (<autogenerated>:1)        FUNCDATA        $5, type:.eq.[2]interface {}.arginfo1(SB)
        0x0018 00024 (<autogenerated>:1)        FUNCDATA        $6, type:.eq.[2]interface {}.argliveinfo(SB)
        0x0018 00024 (<autogenerated>:1)        PCDATA  $3, $1
        0x0018 00024 (<autogenerated>:1)        MOVQ    AX, main.p+48(SP)
        0x001d 00029 (<autogenerated>:1)        MOVQ    BX, main.q+56(SP)
        0x0022 00034 (<autogenerated>:1)        PCDATA  $3, $-1
        0x0022 00034 (<autogenerated>:1)        XORL    CX, CX
        0x0024 00036 (<autogenerated>:1)        JMP     56
        0x0026 00038 (<autogenerated>:1)        MOVQ    main..autotmp_6+24(SP), CX
        0x002b 00043 (<autogenerated>:1)        INCQ    CX
        0x002e 00046 (<autogenerated>:1)        MOVQ    main.q+56(SP), BX
        0x0033 00051 (<autogenerated>:1)        MOVQ    main.p+48(SP), AX
        0x0038 00056 (<autogenerated>:1)        CMPQ    CX, $2
        0x003c 00060 (<autogenerated>:1)        JGE     127
        0x003e 00062 (<autogenerated>:1)        MOVQ    CX, DX
        0x0041 00065 (<autogenerated>:1)        SHLQ    $4, CX
        0x0045 00069 (<autogenerated>:1)        MOVQ    (AX)(CX*1), SI
        0x0049 00073 (<autogenerated>:1)        MOVQ    (BX)(CX*1), DI
        0x004d 00077 (<autogenerated>:1)        MOVQ    8(CX)(AX*1), R8
        0x0052 00082 (<autogenerated>:1)        MOVQ    8(CX)(BX*1), CX
        0x0057 00087 (<autogenerated>:1)        CMPQ    DI, SI
        0x005a 00090 (<autogenerated>:1)        JNE     123
        0x005c 00092 (<autogenerated>:1)        MOVQ    DX, main..autotmp_6+24(SP)
        0x0061 00097 (<autogenerated>:1)        MOVQ    SI, AX
        0x0064 00100 (<autogenerated>:1)        MOVQ    R8, BX
        0x0067 00103 (<autogenerated>:1)        PCDATA  $1, $0
        0x0067 00103 (<autogenerated>:1)        CALL    runtime.efaceeq(SB)
        0x006c 00108 (<autogenerated>:1)        TESTB   AL, AL
        0x006e 00110 (<autogenerated>:1)        JNE     38
        0x0070 00112 (<autogenerated>:1)        MOVQ    main..autotmp_6+24(SP), CX
        0x0075 00117 (<autogenerated>:1)        CMPQ    CX, $2
        0x0079 00121 (<autogenerated>:1)        JMP     127
        0x007b 00123 (<autogenerated>:1)        CMPQ    DX, $2
        0x007f 00127 (<autogenerated>:1)        SETGE   AL
        0x0082 00130 (<autogenerated>:1)        MOVQ    32(SP), BP
        0x0087 00135 (<autogenerated>:1)        ADDQ    $40, SP
        0x008b 00139 (<autogenerated>:1)        RET
        0x008c 00140 (<autogenerated>:1)        NOP
        0x008c 00140 (<autogenerated>:1)        PCDATA  $1, $-1
        0x008c 00140 (<autogenerated>:1)        PCDATA  $0, $-2
        0x008c 00140 (<autogenerated>:1)        MOVQ    AX, 8(SP)
        0x0091 00145 (<autogenerated>:1)        MOVQ    BX, 16(SP)
        0x0096 00150 (<autogenerated>:1)        CALL    runtime.morestack_noctxt(SB)
        0x009b 00155 (<autogenerated>:1)        MOVQ    8(SP), AX
        0x00a0 00160 (<autogenerated>:1)        MOVQ    16(SP), BX
        0x00a5 00165 (<autogenerated>:1)        PCDATA  $0, $-1
        0x00a5 00165 (<autogenerated>:1)        JMP     0
        0x0000 49 3b 66 10 0f 86 82 00 00 00 48 83 ec 28 48 89  I;f.......H..(H.
        0x0010 6c 24 20 48 8d 6c 24 20 48 89 44 24 30 48 89 5c  l$ H.l$ H.D$0H.\
        0x0020 24 38 31 c9 eb 12 48 8b 4c 24 18 48 ff c1 48 8b  $81...H.L$.H..H.
        0x0030 5c 24 38 48 8b 44 24 30 48 83 f9 02 7d 41 48 89  \$8H.D$0H...}AH.
        0x0040 ca 48 c1 e1 04 48 8b 34 08 48 8b 3c 0b 4c 8b 44  .H...H.4.H.<.L.D
        0x0050 01 08 48 8b 4c 19 08 48 39 f7 75 1f 48 89 54 24  ..H.L..H9.u.H.T$
        0x0060 18 48 89 f0 4c 89 c3 e8 00 00 00 00 84 c0 75 b6  .H..L.........u.
        0x0070 48 8b 4c 24 18 48 83 f9 02 eb 04 48 83 fa 02 0f  H.L$.H.....H....
        0x0080 9d c0 48 8b 6c 24 20 48 83 c4 28 c3 48 89 44 24  ..H.l$ H..(.H.D$
        0x0090 08 48 89 5c 24 10 e8 00 00 00 00 48 8b 44 24 08  .H.\$......H.D$.
        0x00a0 48 8b 5c 24 10 e9 56 ff ff ff                    H.\$..V...
        rel 104+4 t=7 runtime.efaceeq+0
        rel 151+4 t=7 runtime.morestack_noctxt+0
go:cuinfo.producer.main SDWARFCUINFO dupok size=0
        0x0000 72 65 67 61 62 69                                regabi
go:cuinfo.packagename.main SDWARFCUINFO dupok size=0
        0x0000 6d 61 69 6e                                      main
go:info.fmt.Println$abstract SDWARFABSFCN dupok size=42
        0x0000 05 66 6d 74 2e 50 72 69 6e 74 6c 6e 00 01 01 13  .fmt.Println....
        0x0010 61 00 00 00 00 00 00 13 6e 00 01 00 00 00 00 13  a.......n.......
        0x0020 65 72 72 00 01 00 00 00 00 00                    err.......
        rel 0+0 t=22 type:[]interface {}+0
        rel 0+0 t=22 type:error+0
        rel 0+0 t=22 type:int+0
        rel 19+4 t=31 go:info.[]interface {}+0
        rel 27+4 t=31 go:info.int+0
        rel 37+4 t=31 go:info.error+0
go:itab.*os.File,io.Writer SRODATA dupok size=32
        0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0010 5a 22 ee 60 00 00 00 00 00 00 00 00 00 00 00 00  Z".`............
        rel 0+8 t=1 type:io.Writer+0
        rel 8+8 t=1 type:*os.File+0
        rel 24+8 t=-32767 os.(*File).Write+0
main..inittask SNOPTRDATA size=40
        0x0000 00 00 00 00 00 00 00 00 02 00 00 00 00 00 00 00  ................
        0x0010 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0020 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 fmt..inittask+0
        rel 32+8 t=1 time..inittask+0
go:string."haha" SRODATA dupok size=4
        0x0000 68 61 68 61                                      haha
go:string."write" SRODATA dupok size=5
        0x0000 77 72 69 74 65                                   write
type:.eqfunc.[2]interface {} SRODATA dupok size=8
        0x0000 00 00 00 00 00 00 00 00                          ........
        rel 0+8 t=1 type:.eq.[2]interface {}+0
runtime.memequal64·f SRODATA dupok size=8
        0x0000 00 00 00 00 00 00 00 00                          ........
        rel 0+8 t=1 runtime.memequal64+0
runtime.gcbits.0100000000000000 SRODATA dupok size=8
        0x0000 01 00 00 00 00 00 00 00                          ........
type:.namedata.*[2]interface {}- SRODATA dupok size=18
        0x0000 00 10 2a 5b 32 5d 69 6e 74 65 72 66 61 63 65 20  ..*[2]interface
        0x0010 7b 7d                                            {}
type:*[2]interface {} SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 53 dc 6b 00 08 08 08 36 00 00 00 00 00 00 00 00  S.k....6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*[2]interface {}-+0
        rel 48+8 t=1 type:[2]interface {}+0
runtime.gcbits.0a00000000000000 SRODATA dupok size=8
        0x0000 0a 00 00 00 00 00 00 00                          ........
type:[2]interface {} SRODATA dupok size=72
        0x0000 20 00 00 00 00 00 00 00 20 00 00 00 00 00 00 00   ....... .......
        0x0010 0a f3 b4 b4 02 08 08 11 00 00 00 00 00 00 00 00  ................
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 02 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 type:.eqfunc.[2]interface {}+0
        rel 32+8 t=1 runtime.gcbits.0a00000000000000+0
        rel 40+4 t=5 type:.namedata.*[2]interface {}-+0
        rel 44+4 t=-32763 type:*[2]interface {}+0
        rel 48+8 t=1 type:interface {}+0
        rel 56+8 t=1 type:[]interface {}+0
type:.namedata.*[8]uint8- SRODATA dupok size=11
        0x0000 00 09 2a 5b 38 5d 75 69 6e 74 38                 ..*[8]uint8
type:*[8]uint8 SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 f8 9a 95 1a 08 08 08 36 00 00 00 00 00 00 00 00  .......6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*[8]uint8-+0
        rel 48+8 t=1 type:[8]uint8+0
runtime.gcbits. SRODATA dupok size=0
type:[8]uint8 SRODATA dupok size=72
        0x0000 08 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0010 89 18 9c b4 0a 01 01 11 00 00 00 00 00 00 00 00  ................
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 08 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.+0
        rel 40+4 t=5 type:.namedata.*[8]uint8-+0
        rel 44+4 t=-32763 type:*[8]uint8+0
        rel 48+8 t=1 type:uint8+0
        rel 56+8 t=1 type:[]uint8+0
type:.namedata.*[8]int- SRODATA dupok size=9
        0x0000 00 07 2a 5b 38 5d 69 6e 74                       ..*[8]int
type:*[8]int SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 b2 24 38 0a 08 08 08 36 00 00 00 00 00 00 00 00  .$8....6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*[8]int-+0
        rel 48+8 t=1 type:noalg.[8]int+0
type:noalg.[8]int SRODATA dupok size=72
        0x0000 40 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  @...............
        0x0010 94 48 d7 e9 02 08 08 11 00 00 00 00 00 00 00 00  .H..............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 08 00 00 00 00 00 00 00                          ........
        rel 32+8 t=1 runtime.gcbits.+0
        rel 40+4 t=5 type:.namedata.*[8]int-+0
        rel 44+4 t=-32763 type:*[8]int+0
        rel 48+8 t=1 type:int+0
        rel 56+8 t=1 type:[]int+0
type:.namedata.*[8]string- SRODATA dupok size=12
        0x0000 00 0a 2a 5b 38 5d 73 74 72 69 6e 67              ..*[8]string
type:*[8]string SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 e3 bf d7 63 08 08 08 36 00 00 00 00 00 00 00 00  ...c...6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*[8]string-+0
        rel 48+8 t=1 type:noalg.[8]string+0
runtime.gcbits.5555000000000000 SRODATA dupok size=8
        0x0000 55 55 00 00 00 00 00 00                          UU......
type:noalg.[8]string SRODATA dupok size=72
        0x0000 80 00 00 00 00 00 00 00 78 00 00 00 00 00 00 00  ........x.......
        0x0010 0c 1c ff 04 02 08 08 11 00 00 00 00 00 00 00 00  ................
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 08 00 00 00 00 00 00 00                          ........
        rel 32+8 t=1 runtime.gcbits.5555000000000000+0
        rel 40+4 t=5 type:.namedata.*[8]string-+0
        rel 44+4 t=-32763 type:*[8]string+0
        rel 48+8 t=1 type:string+0
        rel 56+8 t=1 type:[]string+0
type:.namedata.*map.bucket[int]string- SRODATA dupok size=24
        0x0000 00 16 2a 6d 61 70 2e 62 75 63 6b 65 74 5b 69 6e  ..*map.bucket[in
        0x0010 74 5d 73 74 72 69 6e 67                          t]string
type:*map.bucket[int]string SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 2a 5e fd 13 08 08 08 36 00 00 00 00 00 00 00 00  *^.....6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*map.bucket[int]string-+0
        rel 48+8 t=1 type:noalg.map.bucket[int]string+0
runtime.gcbits.00aaaa0200000000 SRODATA dupok size=8
        0x0000 00 aa aa 02 00 00 00 00                          ........
type:.importpath.. SRODATA dupok size=2
        0x0000 00 00                                            ..
type:.namedata.topbits- SRODATA dupok size=9
        0x0000 00 07 74 6f 70 62 69 74 73                       ..topbits
type:.namedata.keys- SRODATA dupok size=6
        0x0000 00 04 6b 65 79 73                                ..keys
type:.namedata.elems- SRODATA dupok size=7
        0x0000 00 05 65 6c 65 6d 73                             ..elems
type:.namedata.overflow- SRODATA dupok size=10
        0x0000 00 08 6f 76 65 72 66 6c 6f 77                    ..overflow
type:noalg.map.bucket[int]string SRODATA dupok size=176
        0x0000 d0 00 00 00 00 00 00 00 d0 00 00 00 00 00 00 00  ................
        0x0010 8a 14 6e a7 02 08 08 19 00 00 00 00 00 00 00 00  ..n.............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 04 00 00 00 00 00 00 00 04 00 00 00 00 00 00 00  ................
        0x0050 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0070 00 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0080 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0090 48 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  H...............
        0x00a0 00 00 00 00 00 00 00 00 c8 00 00 00 00 00 00 00  ................
        rel 32+8 t=1 runtime.gcbits.00aaaa0200000000+0
        rel 40+4 t=5 type:.namedata.*map.bucket[int]string-+0
        rel 44+4 t=-32763 type:*map.bucket[int]string+0
        rel 48+8 t=1 type:.importpath..+0
        rel 56+8 t=1 type:noalg.map.bucket[int]string+80
        rel 80+8 t=1 type:.namedata.topbits-+0
        rel 88+8 t=1 type:[8]uint8+0
        rel 104+8 t=1 type:.namedata.keys-+0
        rel 112+8 t=1 type:noalg.[8]int+0
        rel 128+8 t=1 type:.namedata.elems-+0
        rel 136+8 t=1 type:noalg.[8]string+0
        rel 152+8 t=1 type:.namedata.overflow-+0
        rel 160+8 t=1 type:unsafe.Pointer+0
runtime.memhash64·f SRODATA dupok size=8
        0x0000 00 00 00 00 00 00 00 00                          ........
        rel 0+8 t=1 runtime.memhash64+0
type:.namedata.*map[int]string- SRODATA dupok size=17
        0x0000 00 0f 2a 6d 61 70 5b 69 6e 74 5d 73 74 72 69 6e  ..*map[int]strin
        0x0010 67                                               g
type:*map[int]string SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 fd 58 bc 94 08 08 08 36 00 00 00 00 00 00 00 00  .X.....6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*map[int]string-+0
        rel 48+8 t=1 type:map[int]string+0
type:map[int]string SRODATA dupok size=88
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 10 c9 8e 73 02 08 08 35 00 00 00 00 00 00 00 00  ...s...5........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0050 08 10 d0 00 04 00 00 00                          ........
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*map[int]string-+0
        rel 44+4 t=-32763 type:*map[int]string+0
        rel 48+8 t=1 type:int+0
        rel 56+8 t=1 type:string+0
        rel 64+8 t=1 type:noalg.map.bucket[int]string+0
        rel 72+8 t=1 runtime.memhash64·f+0
type:.namedata.*map.hdr[int]string- SRODATA dupok size=21
        0x0000 00 13 2a 6d 61 70 2e 68 64 72 5b 69 6e 74 5d 73  ..*map.hdr[int]s
        0x0010 74 72 69 6e 67                                   tring
type:*map.hdr[int]string SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 d9 34 cc 39 08 08 08 36 00 00 00 00 00 00 00 00  .4.9...6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*map.hdr[int]string-+0
        rel 48+8 t=1 type:noalg.map.hdr[int]string+0
runtime.gcbits.2c00000000000000 SRODATA dupok size=8
        0x0000 2c 00 00 00 00 00 00 00                          ,.......
type:.namedata.count- SRODATA dupok size=7
        0x0000 00 05 63 6f 75 6e 74                             ..count
type:.namedata.flags- SRODATA dupok size=7
        0x0000 00 05 66 6c 61 67 73                             ..flags
type:.namedata.B. SRODATA dupok size=3
        0x0000 01 01 42                                         ..B
type:.namedata.noverflow- SRODATA dupok size=11
        0x0000 00 09 6e 6f 76 65 72 66 6c 6f 77                 ..noverflow
type:.namedata.hash0- SRODATA dupok size=7
        0x0000 00 05 68 61 73 68 30                             ..hash0
type:.namedata.buckets- SRODATA dupok size=9
        0x0000 00 07 62 75 63 6b 65 74 73                       ..buckets
type:.namedata.oldbuckets- SRODATA dupok size=12
        0x0000 00 0a 6f 6c 64 62 75 63 6b 65 74 73              ..oldbuckets
type:.namedata.nevacuate- SRODATA dupok size=11
        0x0000 00 09 6e 65 76 61 63 75 61 74 65                 ..nevacuate
type:.namedata.extra- SRODATA dupok size=7
        0x0000 00 05 65 78 74 72 61                             ..extra
type:noalg.map.hdr[int]string SRODATA dupok size=296
        0x0000 30 00 00 00 00 00 00 00 30 00 00 00 00 00 00 00  0.......0.......
        0x0010 d0 d7 98 59 02 08 08 19 00 00 00 00 00 00 00 00  ...Y............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 09 00 00 00 00 00 00 00 09 00 00 00 00 00 00 00  ................
        0x0050 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0070 00 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0080 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0090 09 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x00a0 00 00 00 00 00 00 00 00 0a 00 00 00 00 00 00 00  ................
        0x00b0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x00c0 0c 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x00d0 00 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
        0x00e0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x00f0 18 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0100 00 00 00 00 00 00 00 00 20 00 00 00 00 00 00 00  ........ .......
        0x0110 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0120 28 00 00 00 00 00 00 00                          (.......
        rel 32+8 t=1 runtime.gcbits.2c00000000000000+0
        rel 40+4 t=5 type:.namedata.*map.hdr[int]string-+0
        rel 44+4 t=-32763 type:*map.hdr[int]string+0
        rel 48+8 t=1 type:.importpath..+0
        rel 56+8 t=1 type:noalg.map.hdr[int]string+80
        rel 80+8 t=1 type:.namedata.count-+0
        rel 88+8 t=1 type:int+0
        rel 104+8 t=1 type:.namedata.flags-+0
        rel 112+8 t=1 type:uint8+0
        rel 128+8 t=1 type:.namedata.B.+0
        rel 136+8 t=1 type:uint8+0
        rel 152+8 t=1 type:.namedata.noverflow-+0
        rel 160+8 t=1 type:uint16+0
        rel 176+8 t=1 type:.namedata.hash0-+0
        rel 184+8 t=1 type:uint32+0
        rel 200+8 t=1 type:.namedata.buckets-+0
        rel 208+8 t=1 type:*map.bucket[int]string+0
        rel 224+8 t=1 type:.namedata.oldbuckets-+0
        rel 232+8 t=1 type:*map.bucket[int]string+0
        rel 248+8 t=1 type:.namedata.nevacuate-+0
        rel 256+8 t=1 type:uintptr+0
        rel 272+8 t=1 type:.namedata.extra-+0
        rel 280+8 t=1 type:unsafe.Pointer+0
type:.namedata.*map.iter[int]string- SRODATA dupok size=22
        0x0000 00 14 2a 6d 61 70 2e 69 74 65 72 5b 69 6e 74 5d  ..*map.iter[int]
        0x0010 73 74 72 69 6e 67                                string
type:*map.iter[int]string SRODATA dupok size=56
        0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0010 22 94 16 70 08 08 08 36 00 00 00 00 00 00 00 00  "..p...6........
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00                          ........
        rel 24+8 t=1 runtime.memequal64·f+0
        rel 32+8 t=1 runtime.gcbits.0100000000000000+0
        rel 40+4 t=5 type:.namedata.*map.iter[int]string-+0
        rel 48+8 t=1 type:noalg.map.iter[int]string+0
runtime.gcbits.ff00000000000000 SRODATA dupok size=8
        0x0000 ff 00 00 00 00 00 00 00                          ........
type:.namedata.key- SRODATA dupok size=5
        0x0000 00 03 6b 65 79                                   ..key
type:.namedata.elem- SRODATA dupok size=6
        0x0000 00 04 65 6c 65 6d                                ..elem
type:.namedata.t- SRODATA dupok size=3
        0x0000 00 01 74                                         ..t
type:.namedata.h- SRODATA dupok size=3
        0x0000 00 01 68                                         ..h
type:.namedata.bptr- SRODATA dupok size=6
        0x0000 00 04 62 70 74 72                                ..bptr
type:.namedata.oldoverflow- SRODATA dupok size=13
        0x0000 00 0b 6f 6c 64 6f 76 65 72 66 6c 6f 77           ..oldoverflow
type:.namedata.startBucket- SRODATA dupok size=13
        0x0000 00 0b 73 74 61 72 74 42 75 63 6b 65 74           ..startBucket
type:.namedata.offset- SRODATA dupok size=8
        0x0000 00 06 6f 66 66 73 65 74                          ..offset
type:.namedata.wrapped- SRODATA dupok size=9
        0x0000 00 07 77 72 61 70 70 65 64                       ..wrapped
type:.namedata.i- SRODATA dupok size=3
        0x0000 00 01 69                                         ..i
type:.namedata.bucket- SRODATA dupok size=8
        0x0000 00 06 62 75 63 6b 65 74                          ..bucket
type:.namedata.checkBucket- SRODATA dupok size=13
        0x0000 00 0b 63 68 65 63 6b 42 75 63 6b 65 74           ..checkBucket
type:noalg.map.iter[int]string SRODATA dupok size=440
        0x0000 60 00 00 00 00 00 00 00 40 00 00 00 00 00 00 00  `.......@.......
        0x0010 27 fe 8b 5a 02 08 08 19 00 00 00 00 00 00 00 00  '..Z............
        0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0040 0f 00 00 00 00 00 00 00 0f 00 00 00 00 00 00 00  ................
        0x0050 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0070 00 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
        0x0080 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0090 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x00a0 00 00 00 00 00 00 00 00 18 00 00 00 00 00 00 00  ................
        0x00b0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x00c0 20 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00   ...............
        0x00d0 00 00 00 00 00 00 00 00 28 00 00 00 00 00 00 00  ........(.......
        0x00e0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x00f0 30 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  0...............
        0x0100 00 00 00 00 00 00 00 00 38 00 00 00 00 00 00 00  ........8.......
        0x0110 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0120 40 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  @...............
        0x0130 00 00 00 00 00 00 00 00 48 00 00 00 00 00 00 00  ........H.......
        0x0140 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0150 49 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  I...............
        0x0160 00 00 00 00 00 00 00 00 4a 00 00 00 00 00 00 00  ........J.......
        0x0170 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x0180 4b 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  K...............
        0x0190 00 00 00 00 00 00 00 00 50 00 00 00 00 00 00 00  ........P.......
        0x01a0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
        0x01b0 58 00 00 00 00 00 00 00                          X.......
        rel 32+8 t=1 runtime.gcbits.ff00000000000000+0
        rel 40+4 t=5 type:.namedata.*map.iter[int]string-+0
        rel 44+4 t=-32763 type:*map.iter[int]string+0
        rel 48+8 t=1 type:.importpath..+0
        rel 56+8 t=1 type:noalg.map.iter[int]string+80
        rel 80+8 t=1 type:.namedata.key-+0
        rel 88+8 t=1 type:*int+0
        rel 104+8 t=1 type:.namedata.elem-+0
        rel 112+8 t=1 type:*string+0
        rel 128+8 t=1 type:.namedata.t-+0
        rel 136+8 t=1 type:unsafe.Pointer+0
        rel 152+8 t=1 type:.namedata.h-+0
        rel 160+8 t=1 type:*map.hdr[int]string+0
        rel 176+8 t=1 type:.namedata.buckets-+0
        rel 184+8 t=1 type:*map.bucket[int]string+0
        rel 200+8 t=1 type:.namedata.bptr-+0
        rel 208+8 t=1 type:*map.bucket[int]string+0
        rel 224+8 t=1 type:.namedata.overflow-+0
        rel 232+8 t=1 type:unsafe.Pointer+0
        rel 248+8 t=1 type:.namedata.oldoverflow-+0
        rel 256+8 t=1 type:unsafe.Pointer+0
        rel 272+8 t=1 type:.namedata.startBucket-+0
        rel 280+8 t=1 type:uintptr+0
        rel 296+8 t=1 type:.namedata.offset-+0
        rel 304+8 t=1 type:uint8+0
        rel 320+8 t=1 type:.namedata.wrapped-+0
        rel 328+8 t=1 type:bool+0
        rel 344+8 t=1 type:.namedata.B.+0
        rel 352+8 t=1 type:uint8+0
        rel 368+8 t=1 type:.namedata.i-+0
        rel 376+8 t=1 type:uint8+0
        rel 392+8 t=1 type:.namedata.bucket-+0
        rel 400+8 t=1 type:uintptr+0
        rel 416+8 t=1 type:.namedata.checkBucket-+0
        rel 424+8 t=1 type:uintptr+0
type:.importpath.fmt. SRODATA dupok size=5
        0x0000 00 03 66 6d 74                                   ..fmt
type:.importpath.time. SRODATA dupok size=6
        0x0000 00 04 74 69 6d 65                                ..time
gclocals·DsEZEqsV1VFgO2VgUUolkQ== SRODATA dupok size=8
        0x0000 05 00 00 00 00 00 00 00                          ........
gclocals·HFbfPPPbONK70RRxfAk8Uw== SRODATA dupok size=43
        0x0000 05 00 00 00 31 00 00 00 00 00 00 00 00 00 00 80  ....1...........
        0x0010 05 00 00 00 00 00 00 f8 07 00 00 00 00 15 f8 07  ................
        0x0020 00 00 00 00 14 f8 07 00 00 00 00                 ...........
main.main.stkobj SRODATA static size=72
        0x0000 04 00 00 00 00 00 00 00 80 fe ff ff 20 00 00 00  ............ ...
        0x0010 20 00 00 00 00 00 00 00 a0 fe ff ff 30 00 00 00   ...........0...
        0x0020 30 00 00 00 00 00 00 00 d0 fe ff ff 60 00 00 00  0...........`...
        0x0030 40 00 00 00 00 00 00 00 30 ff ff ff d0 00 00 00  @.......0.......
        0x0040 d0 00 00 00 00 00 00 00                          ........
        rel 20+4 t=5 runtime.gcbits.0a00000000000000+0
        rel 36+4 t=5 runtime.gcbits.2c00000000000000+0
        rel 52+4 t=5 runtime.gcbits.ff00000000000000+0
        rel 68+4 t=5 runtime.gcbits.00aaaa0200000000+0
gclocals·ZzMiPAiVBg7DJ6dh/CjSag== SRODATA dupok size=11
        0x0000 03 00 00 00 01 00 00 00 01 00 00                 ...........
gclocals·VtCL4RdUwCqwXEPeyJllRA== SRODATA dupok size=11
        0x0000 03 00 00 00 01 00 00 00 00 01 00                 ...........
main.read.arginfo1 SRODATA static dupok size=3
        0x0000 00 08 ff                                         ...
main.read.argliveinfo SRODATA static dupok size=2
        0x0000 00 00                                            ..
main.write.arginfo1 SRODATA static dupok size=3
        0x0000 00 08 ff                                         ...
main.write.argliveinfo SRODATA static dupok size=2
        0x0000 00 00                                            ..
gclocals·TjPuuCwdlCpTaRQGRKTrYw== SRODATA dupok size=10
        0x0000 02 00 00 00 02 00 00 00 03 00                    ..........
gclocals·J5F+7Qw7O7ve2QcWC7DpeQ== SRODATA dupok size=8
        0x0000 02 00 00 00 00 00 00 00                          ........
type:.eq.[2]interface {}.arginfo1 SRODATA static dupok size=3
        0x0000 08 08 ff                                         ...
type:.eq.[2]interface {}.argliveinfo SRODATA static dupok size=2
        0x0000 00 00                                            ..










*/
