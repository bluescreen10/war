package text

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// https://webassembly.github.io/spec/core/text/values.html#text-digit
const digit = "0123456789"
const letter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const symbol = "!#$%&'*+-./:<=>?@\\^_`|~"
const sign = "+-"
const keyword = letter + digit + "_.:="

// https://webassembly.github.io/spec/core/text/values.html#text-idchar
var idChar = digit + letter + symbol

// https://webassembly.github.io/spec/core/text/values.html#text-hexdigit
const hexDigit = digit + "abcdefABCDEF"

var ErrEOF error = errors.New("EOF")

type tokenKind int

const (
	tokenError tokenKind = iota
	tokenEOF
	tokenLParen
	tokenRParen
	tokenIdent
	tokenNumber
	tokenString
	tokenNumtype
	tokenVectype
	tokenReftype
	tokenOffsetEq
	tokenAlignEq
	tokenKeyword
	tokenModule
	tokenType
	tokenFunc
	tokenParam
	tokenResult
	tokenBlock
	tokenLoop
	tokenIf
	tokenThen
	tokenElse
	tokenEnd
	tokenUnreachable
	tokenNop
	tokenBr
	tokenBrIf
	tokenBrTable
	tokenReturn
	tokenCall
	tokenCallIndirect
	tokenDrop
	tokenSelect
	tokenLocalGet
	tokenLocalTee
	tokenLocalSet
	tokenGlobalGet
	tokenGlobalSet
	tokenTableGet
	tokenTableSet
	tokenTableSize
	tokenTableGrow
	tokenTableFill
	tokenTableCopy
	tokenTableInit
	tokenElemDrop
	tokenMemorySize
	tokenMemoryGrow
	tokenMemoryFill
	tokenMemoryCopy
	tokenMemoryInit
	tokenDataDrop
	tokenI32Load
	tokenI64Load
	tokenF32Load
	tokenF64Load
	tokenI32Store
	tokenI64Store
	tokenF32Store
	tokenF64Store
	tokenI32Load8U
	tokenI32Load8S
	tokenI32Load16U
	tokenI32Load16S
	tokenI64Load8U
	tokenI64Load8S
	tokenI64Load16U
	tokenI64Load16S
	tokenI64Load32U
	tokenI64Load32S
	tokenI32Store8
	tokenI32Store16
	tokenI64Store8
	tokenI64Store16
	tokenI64Store32
	tokenV128Load
	tokenV128Store
	tokenV128Load8x8U
	tokenV128Load8x8S
	tokenV128Load16x4U
	tokenV128Load16x4S
	tokenV128Load32x2U
	tokenV128Load32x2S
	tokenV128Load8Splat
	tokenV128Load16Splat
	tokenV128Load32Splat
	tokenV128Load64Splat
	tokenV128Load32Zero
	tokenV128Load64Zero
	tokenV128Load8Lane
	tokenV128Load16Lane
	tokenV128Load32Lane
	tokenV128Load64Lane
	tokenV128Store8Lane
	tokenV128Store16Lane
	tokenV128Store32Lane
	tokenV128Store64Lane
	tokenI32Const
	tokenI64Const
	tokenF32Const
	tokenF64Const
	tokenV128Const
	tokenRefNull
	tokenRefFunc
	tokenRefExtern
	tokenRefIsNull
	tokenI32Clz
	tokenI32Ctz
	tokenI32Popcnt
	tokenI32Extend8S
	tokenI32Extend16S
	tokenI64Clz
	tokenI64Ctz
	tokenI64Popcnt
	tokenI64Extend8S
	tokenI64Extend16S
	tokenI64Extend32S
	tokenF32Neg
	tokenF32Abs
	tokenF32Sqrt
	tokenF32Ceil
	tokenF32Floor
	tokenF32Trunc
	tokenF32Nearest
	tokenF64Neg
	tokenF64Abs
	tokenF64Sqrt
	tokenF64Ceil
	tokenF64Floor
	tokenF64Trunc
	tokenF64Nearest
	tokenI32Add
	tokenI32Sub
	tokenI32Mul
	tokenI32DivU
	tokenI32DivS
	tokenI32RemU
	tokenI32RemS
	tokenI32And
	tokenI32Or
	tokenI32Xor
	tokenI32Shl
	tokenI32ShrU
	tokenI32ShrS
	tokenI32Rotl
	tokenI32Rotr
	tokenI64Add
	tokenI64Sub
	tokenI64Mul
	tokenI64DivU
	tokenI64DivS
	tokenI64RemU
	tokenI64RemS
	tokenI64And
	tokenI64Or
	tokenI64Xor
	tokenI64Shl
	tokenI64ShrU
	tokenI64ShrS
	tokenI64Rotl
	tokenI64Rotr
	tokenF32Add
	tokenF32Sub
	tokenF32Mul
	tokenF32Div
	tokenF32Min
	tokenF32Max
	tokenF32Copysign
	tokenF64Add
	tokenF64Sub
	tokenF64Mul
	tokenF64Div
	tokenF64Min
	tokenF64Max
	tokenF64Copysign
	tokenI32Eqz
	tokenI64Eqz
	tokenI32Eq
	tokenI32Ne
	tokenI32LtU
	tokenI32LtS
	tokenI32LeU
	tokenI32LeS
	tokenI32GtU
	tokenI32GtS
	tokenI32GeU
	tokenI32GeS
	tokenI64Eq
	tokenI64Ne
	tokenI64LtU
	tokenI64LtS
	tokenI64LeU
	tokenI64LeS
	tokenI64GtU
	tokenI64GtS
	tokenI64GeU
	tokenI64GeS
	tokenF32Eq
	tokenF32Ne
	tokenF32Lt
	tokenF32Le
	tokenF32Gt
	tokenF32Ge
	tokenF64Eq
	tokenF64Ne
	tokenF64Lt
	tokenF64Le
	tokenF64Gt
	tokenF64Ge
	tokenI32WrapI64
	tokenI64ExtendI32S
	tokenI64ExtendI32U
	tokenF32DemoteF64
	tokenF64PromoteF32
	tokenI32TruncF32U
	tokenI32TruncF32S
	tokenI64TruncF32U
	tokenI64TruncF32S
	tokenI32TruncF64U
	tokenI32TruncF64S
	tokenI64TruncF64U
	tokenI64TruncF64S
	tokenI32TruncSatF32U
	tokenI32TruncSatF32S
	tokenI64TruncSatF32U
	tokenI64TruncSatF32S
	tokenI32TruncSatF64U
	tokenI32TruncSatF64S
	tokenI64TruncSatF64U
	tokenI64TruncSatF64S
	tokenF32ConvertI32U
	tokenF32ConvertI32S
	tokenF64ConvertI32U
	tokenF64ConvertI32S
	tokenF32ConvertI64U
	tokenF32ConvertI64S
	tokenF64ConvertI64U
	tokenF64ConvertI64S
	tokenF32ReinterpretI32
	tokenF64ReinterpretI64
	tokenI32ReinterpretF32
	tokenI64ReinterpretF64
	tokenV128Not
	tokenV128And
	tokenV128Andnot
	tokenV128Or
	tokenV128Xor
	tokenV128Bitselect
	tokenV128AnyTrue
	tokenI8x16Neg
	tokenI16x8Neg
	tokenI32x4Neg
	tokenI64x2Neg
	tokenI8x16Abs
	tokenI16x8Abs
	tokenI32x4Abs
	tokenI64x2Abs
	tokenI8x16Popcnt
	tokenI8x16AvgrU
	tokenI16x8AvgrU
	tokenF32x4Neg
	tokenF64x2Neg
	tokenF32x4Abs
	tokenF64x2Abs
	tokenF32x4Sqrt
	tokenF64x2Sqrt
	tokenF32x4Ceil
	tokenF64x2Ceil
	tokenF32x4Floor
	tokenF64x2Floor
	tokenF32x4Trunc
	tokenF64x2Trunc
	tokenF32x4Nearest
	tokenF64x2Nearest
	tokenI32x4TruncSatF32x4U
	tokenI32x4TruncSatF32x4S
	tokenI32x4TruncSatF64x2UZero
	tokenI32x4TruncSatF64x2SZero
	tokenF64x2PromoteLowF32x4
	tokenF32x4DemoteF64x2Zero
	tokenF32x4ConvertI32x4U
	tokenF32x4ConvertI32x4S
	tokenF64x2ConvertLowI32x4U
	tokenF64x2ConvertLowI32x4S
	tokenI16x8ExtaddPairwiseI8x16U
	tokenI16x8ExtaddPairwiseI8x16S
	tokenI32x4ExtaddPairwiseI16x8U
	tokenI32x4ExtaddPairwiseI16x8S
	tokenI8x16Eq
	tokenI16x8Eq
	tokenI32x4Eq
	tokenI64x2Eq
	tokenI8x16Ne
	tokenI16x8Ne
	tokenI32x4Ne
	tokenI64x2Ne
	tokenI8x16LtU
	tokenI8x16LtS
	tokenI16x8LtU
	tokenI16x8LtS
	tokenI32x4LtU
	tokenI32x4LtS
	tokenI64x2LtS
	tokenI8x16LeU
	tokenI8x16LeS
	tokenI16x8LeU
	tokenI16x8LeS
	tokenI32x4LeU
	tokenI32x4LeS
	tokenI64x2LeS
	tokenI8x16GtU
	tokenI8x16GtS
	tokenI16x8GtU
	tokenI16x8GtS
	tokenI32x4GtU
	tokenI32x4GtS
	tokenI64x2GtS
	tokenI8x16GeU
	tokenI8x16GeS
	tokenI16x8GeU
	tokenI16x8GeS
	tokenI32x4GeU
	tokenI32x4GeS
	tokenI64x2GeS
	tokenF32x4Eq
	tokenF64x2Eq
	tokenF32x4Ne
	tokenF64x2Ne
	tokenF32x4Lt
	tokenF64x2Lt
	tokenF32x4Le
	tokenF64x2Le
	tokenF32x4Gt
	tokenF64x2Gt
	tokenF32x4Ge
	tokenF64x2Ge
	tokenI8x16Swizzle
	tokenI8x16Add
	tokenI16x8Add
	tokenI32x4Add
	tokenI64x2Add
	tokenI8x16Sub
	tokenI16x8Sub
	tokenI32x4Sub
	tokenI64x2Sub
	tokenI16x8Mul
	tokenI32x4Mul
	tokenI64x2Mul
	tokenI8x16AddSatU
	tokenI8x16AddSatS
	tokenI16x8AddSatU
	tokenI16x8AddSatS
	tokenI8x16SubSatU
	tokenI8x16SubSatS
	tokenI16x8SubSatU
	tokenI16x8SubSatS
	tokenI32x4DotI16x8S
	tokenI8x16MinU
	tokenI16x8MinU
	tokenI32x4MinU
	tokenI8x16MinS
	tokenI16x8MinS
	tokenI32x4MinS
	tokenI8x16MaxU
	tokenI16x8MaxU
	tokenI32x4MaxU
	tokenI8x16MaxS
	tokenI16x8MaxS
	tokenI32x4MaxS
	tokenF32x4Add
	tokenF64x2Add
	tokenF32x4Sub
	tokenF64x2Sub
	tokenF32x4Mul
	tokenF64x2Mul
	tokenF32x4Div
	tokenF64x2Div
	tokenF32x4Min
	tokenF64x2Min
	tokenF32x4Max
	tokenF64x2Max
	tokenF32x4Pmin
	tokenF64x2Pmin
	tokenF32x4Pmax
	tokenF64x2Pmax
	tokenI16x8Q15mulrSatS
	tokenI8x16NarrowI16x8U
	tokenI8x16NarrowI16x8S
	tokenI16x8NarrowI32x4U
	tokenI16x8NarrowI32x4S
	tokenI16x8ExtendLowI8x16U
	tokenI16x8ExtendLowI8x16S
	tokenI16x8ExtendHighI8x16U
	tokenI16x8ExtendHighI8x16S
	tokenI32x4ExtendLowI16x8U
	tokenI32x4ExtendLowI16x8S
	tokenI32x4ExtendHighI16x8U
	tokenI32x4ExtendHighI16x8S
	tokenI64x2ExtendLowI32x4U
	tokenI64x2ExtendLowI32x4S
	tokenI64x2ExtendHighI32x4U
	tokenI64x2ExtendHighI32x4S
	tokenI16x8ExtmulLowI8x16U
	tokenI16x8ExtmulLowI8x16S
	tokenI16x8ExtmulHighI8x16U
	tokenI16x8ExtmulHighI8x16S
	tokenI32x4ExtmulLowI16x8U
	tokenI32x4ExtmulLowI16x8S
	tokenI32x4ExtmulHighI16x8U
	tokenI32x4ExtmulHighI16x8S
	tokenI64x2ExtmulLowI32x4U
	tokenI64x2ExtmulLowI32x4S
	tokenI64x2ExtmulHighI32x4U
	tokenI64x2ExtmulHighI32x4S
	tokenI8x16AllTrue
	tokenI16x8AllTrue
	tokenI32x4AllTrue
	tokenI64x2AllTrue
	tokenI8x16Bitmask
	tokenI16x8Bitmask
	tokenI32x4Bitmask
	tokenI64x2Bitmask
	tokenI8x16Shl
	tokenI16x8Shl
	tokenI32x4Shl
	tokenI64x2Shl
	tokenI8x16ShrU
	tokenI8x16ShrS
	tokenI16x8ShrU
	tokenI16x8ShrS
	tokenI32x4ShrU
	tokenI32x4ShrS
	tokenI64x2ShrU
	tokenI64x2ShrS
	tokenI8x16Shuffle
	tokenI8x16Splat
	tokenI16x8Splat
	tokenI32x4Splat
	tokenI64x2Splat
	tokenF32x4Splat
	tokenF64x2Splat
	tokenI8x16ExtractLaneU
	tokenI8x16ExtractLaneS
	tokenI16x8ExtractLaneU
	tokenI16x8ExtractLaneS
	tokenI32x4ExtractLane
	tokenI64x2ExtractLane
	tokenF32x4ExtractLane
	tokenF64x2ExtractLane
	tokenI8x16ReplaceLane
	tokenI16x8ReplaceLane
	tokenI32x4ReplaceLane
	tokenI64x2ReplaceLane
	tokenF32x4ReplaceLane
	tokenF64x2ReplaceLane
	tokenStart
	tokenLocal
	tokenGlobal
	tokenTable
	tokenMemory
	tokenElem
	tokenData
	tokenDeclare
	tokenOffset
	tokenItem
	tokenImport
	tokenExport
	tokenBin
	tokenQuote
	tokenScript
	tokenRegister
	tokenInvoke
	tokenGet
	tokenAssertMalformed
	tokenAssertInvalid
	tokenAssertUnlinkable
	tokenAssertReturn
	tokenAssertTrap
	tokenAssertExhaustion
	tokenNanCanonical
	tokenNanArithmetic
	tokenInput
	tokenOuput
	tokenExtern
	tokenExternRef
	tokenFuncRef
	tokenMut
)

const (
	eof = 0
)

type token struct {
	kind tokenKind
	val  []byte
}

var key = map[string]tokenKind{
	// numtypes: https://webassembly.github.io/spec/core/text/types.html#number-types
	"i32": tokenNumtype,
	"i64": tokenNumtype,
	"f32": tokenNumtype,
	"f64": tokenNumtype,

	// vectype: https://webassembly.github.io/spec/core/text/types.html#vector-types
	"v128": tokenVectype,

	// instructions: https://webassembly.github.io/spec/core/text/instructions.html#text-plaininstr
	"unreachable":   tokenUnreachable,
	"nop":           tokenNop,
	"br":            tokenBr,
	"br_if":         tokenBrIf,
	"br_table":      tokenBrTable,
	"return":        tokenReturn,
	"call":          tokenCall,
	"call_indirect": tokenCallIndirect,
	"drop":          tokenDrop,

	// control
	"module": tokenModule,
	"type":   tokenType,
	"func":   tokenFunc,
	"param":  tokenParam,
	"result": tokenResult,
	"block":  tokenBlock,
	"loop":   tokenLoop,
	"if":     tokenIf,
	"else":   tokenElse,
	"end":    tokenEnd,
	"select": tokenSelect,

	// memory
	"offset=":     tokenOffsetEq,
	"align=":      tokenAlignEq,
	"memory.size": tokenMemorySize,
	"memory.grow": tokenMemoryGrow,
	"memory.fill": tokenMemoryFill,
	"memory.copy": tokenMemoryCopy,
	"memory.init": tokenMemoryInit,
	"data.drop":   tokenDataDrop,

	// variables
	"local.get":  tokenLocalGet,
	"local.set":  tokenLocalSet,
	"local.tee":  tokenLocalTee,
	"global.get": tokenGlobalGet,
	"global.set": tokenGlobalSet,

	// table
	"table.get":  tokenTableGet,
	"table.set":  tokenTableSet,
	"table.size": tokenTableSize,
	"table.grow": tokenTableGrow,
	"table.fill": tokenTableFill,
	"table.copy": tokenTableCopy,
	"table.init": tokenTableInit,
	"elem.drop":  tokenElemDrop,

	// more instructions
	"start":             tokenStart,
	"local":             tokenLocal,
	"global":            tokenGlobal,
	"table":             tokenTable,
	"memory":            tokenMemory,
	"elem":              tokenElem,
	"data":              tokenData,
	"declare":           tokenDeclare,
	"offset":            tokenOffset,
	"item":              tokenItem,
	"import":            tokenImport,
	"export":            tokenExport,
	"binary":            tokenBin,
	"quote":             tokenQuote,
	"script":            tokenScript,
	"register":          tokenRegister,
	"invoke":            tokenInvoke,
	"get":               tokenGet,
	"assert_malformed":  tokenAssertMalformed,
	"assert_invalid":    tokenAssertInvalid,
	"assert_unlinkable": tokenAssertUnlinkable,
	"assert_return":     tokenAssertReturn,
	"assert_trap":       tokenAssertTrap,
	"assert_exhaustion": tokenAssertExhaustion,
	"nan:canonical":     tokenNanCanonical,
	"nan:arithmetic":    tokenNanArithmetic,
	"input":             tokenInput,
	"output":            tokenOuput,

	// load  & store
	"i32.load":          tokenI32Load,
	"i64.load":          tokenI64Load,
	"f32.load":          tokenF32Load,
	"f64.load":          tokenF64Load,
	"i32.store":         tokenI32Store,
	"i64.store":         tokenI64Store,
	"f32.store":         tokenF32Store,
	"f64.store":         tokenF64Store,
	"i32.load8_u":       tokenI32Load8U,
	"i32.load8_s":       tokenI32Load8S,
	"i32.load16_u":      tokenI32Load16U,
	"i32.load16_s":      tokenI32Load16S,
	"i64.load8_u":       tokenI64Load8U,
	"i64.load8_s":       tokenI64Load8S,
	"i64.load16_u":      tokenI64Load16U,
	"i64.load16_s":      tokenI64Load16S,
	"i64.load32_u":      tokenI64Load32U,
	"i64.load32_s":      tokenI64Load32S,
	"i32.store8":        tokenI32Store8,
	"i32.store16":       tokenI32Store16,
	"i64.store8":        tokenI64Store8,
	"i64.store16":       tokenI64Store16,
	"i64.store32":       tokenI64Store32,
	"v128.load":         tokenV128Load,
	"v128.store":        tokenV128Store,
	"v128.load8x8_u":    tokenV128Load8x8U,
	"v128.load8x8_s":    tokenV128Load8x8S,
	"v128.load16x4_u":   tokenV128Load16x4U,
	"v128.load16x4_s":   tokenV128Load16x4S,
	"v128.load32x2_u":   tokenV128Load32x2U,
	"v128.load32x2_s":   tokenV128Load32x2S,
	"v128.load8_splat":  tokenV128Load8Splat,
	"v128.load16_splat": tokenV128Load16Splat,
	"v128.load32_splat": tokenV128Load32Splat,
	"v128.load64_splat": tokenV128Load64Splat,
	"v128.load32_zero":  tokenV128Load32Zero,
	"v128.load64_zero":  tokenV128Load64Zero,
	"v128.load8_lane":   tokenV128Load8Lane,
	"v128.load16_lane":  tokenV128Load16Lane,
	"v128.load32_lane":  tokenV128Load32Lane,
	"v128.load64_lane":  tokenV128Load64Lane,
	"v128.store8_lane":  tokenV128Store8Lane,
	"v128.store16_lane": tokenV128Store16Lane,
	"v128.store32_lane": tokenV128Store32Lane,
	"v128.store64_lane": tokenV128Store64Lane,

	// constants
	"i32.const":  tokenI32Const,
	"i64.const":  tokenI64Const,
	"f32.const":  tokenF32Const,
	"f64.const":  tokenF64Const,
	"v128.const": tokenV128Const,

	// refs
	"extern":    tokenExtern,
	"externref": tokenExternRef,
	"funcref":   tokenFuncRef,
	"mut":       tokenMut,

	// references
	"ref.null":    tokenRefNull,
	"ref.func":    tokenRefFunc,
	"ref.extern":  tokenRefExtern,
	"ref.is_null": tokenRefIsNull,

	// i32 ops
	"i32.clz":        tokenI32Clz,
	"i32.ctz":        tokenI32Ctz,
	"i32.popcnt":     tokenI32Popcnt,
	"i32.extend8_s":  tokenI32Extend8S,
	"i32.extend16_s": tokenI32Extend16S,

	// i64 ops
	"i64.clz":        tokenI64Clz,
	"i64.ctz":        tokenI64Ctz,
	"i64.popcnt":     tokenI64Popcnt,
	"i64.extend8_s":  tokenI64Extend8S,
	"i64.extend16_s": tokenI64Extend16S,
	"i64.extend32_s": tokenI64Extend32S,

	// f32 ops
	"f32.neg":     tokenF32Neg,
	"f32.abs":     tokenF32Abs,
	"f32.sqrt":    tokenF32Sqrt,
	"f32.ceil":    tokenF32Ceil,
	"f32.floor":   tokenF32Floor,
	"f32.trunc":   tokenF32Trunc,
	"f32.nearest": tokenF32Nearest,

	// f64 ops
	"f64.neg":     tokenF64Neg,
	"f64.abs":     tokenF64Abs,
	"f64.sqrt":    tokenF64Sqrt,
	"f64.ceil":    tokenF64Ceil,
	"f64.floor":   tokenF64Floor,
	"f64.trunc":   tokenF64Trunc,
	"f64.nearest": tokenF64Nearest,

	// i32 ops
	"i32.add":   tokenI32Add,
	"i32.sub":   tokenI32Sub,
	"i32.mul":   tokenI32Mul,
	"i32.div_u": tokenI32DivU,
	"i32.div_s": tokenI32DivS,
	"i32.rem_u": tokenI32RemU,
	"i32.rem_s": tokenI32RemS,
	"i32.and":   tokenI32And,
	"i32.or":    tokenI32Or,
	"i32.xor":   tokenI32Xor,
	"i32.shl":   tokenI32Shl,
	"i32.shr_u": tokenI32ShrU,
	"i32.shr_s": tokenI32ShrS,
	"i32.rotl":  tokenI32Rotl,
	"i32.rotr":  tokenI32Rotr,

	// i64 ops
	"i64.add":   tokenI64Add,
	"i64.sub":   tokenI64Sub,
	"i64.mul":   tokenI64Mul,
	"i64.div_u": tokenI64DivU,
	"i64.div_s": tokenI64DivS,
	"i64.rem_u": tokenI64RemU,
	"i64.rem_s": tokenI64RemS,
	"i64.and":   tokenI64And,
	"i64.or":    tokenI64Or,
	"i64.xor":   tokenI64Xor,
	"i64.shl":   tokenI64Shl,
	"i64.shr_u": tokenI64ShrU,
	"i64.shr_s": tokenI64ShrS,
	"i64.rotl":  tokenI64Rotl,
	"i64.rotr":  tokenI64Rotr,

	// f32 ops
	"f32.add":      tokenF32Add,
	"f32.sub":      tokenF32Sub,
	"f32.mul":      tokenF32Mul,
	"f32.div":      tokenF32Div,
	"f32.min":      tokenF32Min,
	"f32.max":      tokenF32Max,
	"f32.copysign": tokenF32Copysign,

	// f64 ops
	"f64.add":      tokenF64Add,
	"f64.sub":      tokenF64Sub,
	"f64.mul":      tokenF64Mul,
	"f64.div":      tokenF64Div,
	"f64.min":      tokenF64Min,
	"f64.max":      tokenF64Max,
	"f64.copysign": tokenF64Copysign,
	"i32.eqz":      tokenI32Eqz,
	"i64.eqz":      tokenI64Eqz,

	// i32 tests
	"i32.eq":   tokenI32Eq,
	"i32.ne":   tokenI32Ne,
	"i32.lt_u": tokenI32LtU,
	"i32.lt_s": tokenI32LtS,
	"i32.le_u": tokenI32LeU,
	"i32.le_s": tokenI32LeS,
	"i32.gt_u": tokenI32GtU,
	"i32.gt_s": tokenI32GtS,
	"i32.ge_u": tokenI32GeU,
	"i32.ge_s": tokenI32GeS,

	// i64 tests
	"i64.eq":   tokenI64Eq,
	"i64.ne":   tokenI64Ne,
	"i64.lt_u": tokenI64LtU,
	"i64.lt_s": tokenI64LtS,
	"i64.le_u": tokenI64LeU,
	"i64.le_s": tokenI64LeS,
	"i64.gt_u": tokenI64GtU,
	"i64.gt_s": tokenI64GtS,
	"i64.ge_u": tokenI64GeU,
	"i64.ge_s": tokenI64GeS,

	// f32 tests
	"f32.eq": tokenF32Eq,
	"f32.ne": tokenF32Ne,
	"f32.lt": tokenF32Lt,
	"f32.le": tokenF32Le,
	"f32.gt": tokenF32Gt,
	"f32.ge": tokenF32Ge,

	// f64 tests
	"f64.eq": tokenF64Eq,
	"f64.ne": tokenF64Ne,
	"f64.lt": tokenF64Lt,
	"f64.le": tokenF64Le,
	"f64.gt": tokenF64Gt,
	"f64.ge": tokenF64Ge,

	// i32
	"i32.wrap_i64": tokenI32WrapI64,

	// i64
	"i64.extend_i32_s": tokenI64ExtendI32S,
	"i64.extend_i32_u": tokenI64ExtendI32U,

	// f32
	"f32.demote_f64": tokenF32DemoteF64,

	// f64
	"f64.promote_f32": tokenF64PromoteF32,

	// i32
	"i32.trunc_f32_u": tokenI32TruncF32U,
	"i32.trunc_f32_s": tokenI32TruncF32S,

	// i64
	"i64.trunc_f32_u": tokenI64TruncF32U,
	"i64.trunc_f32_s": tokenI64TruncF32S,

	// i32
	"i32.trunc_f64_u": tokenI32TruncF64U,
	"i32.trunc_f64_s": tokenI32TruncF64S,

	// i64
	"i64.trunc_f64_u":     tokenI64TruncF64U,
	"i64.trunc_f64_s":     tokenI64TruncF64S,
	"i32.trunc_sat_f32_u": tokenI32TruncSatF32U,
	"i32.trunc_sat_f32_s": tokenI32TruncSatF32S,
	"i64.trunc_sat_f32_u": tokenI64TruncSatF32U,
	"i64.trunc_sat_f32_s": tokenI64TruncSatF32S,
	"i32.trunc_sat_f64_u": tokenI32TruncSatF64U,
	"i32.trunc_sat_f64_s": tokenI32TruncSatF64S,
	"i64.trunc_sat_f64_u": tokenI64TruncSatF64U,
	"i64.trunc_sat_f64_s": tokenI64TruncSatF64S,

	// convert
	"f32.convert_i32_u":             tokenF32ConvertI32U,
	"f32.convert_i32_s":             tokenF32ConvertI32S,
	"f64.convert_i32_u":             tokenF64ConvertI32U,
	"f64.convert_i32_s":             tokenF64ConvertI32S,
	"f32.convert_i64_u":             tokenF32ConvertI64U,
	"f32.convert_i64_s":             tokenF32ConvertI64S,
	"f64.convert_i64_u":             tokenF64ConvertI64U,
	"f64.convert_i64_s":             tokenF64ConvertI64S,
	"f32.reinterpret_i32":           tokenF32ReinterpretI32,
	"f64.reinterpret_i64":           tokenF64ReinterpretI64,
	"i32.reinterpret_f32":           tokenI32ReinterpretF32,
	"i64.reinterpret_f64":           tokenI64ReinterpretF64,
	"v128.not":                      tokenV128Not,
	"v128.and":                      tokenV128And,
	"v128.andnot":                   tokenV128Andnot,
	"v128.or":                       tokenV128Or,
	"v128.xor":                      tokenV128Xor,
	"v128.bitselect":                tokenV128Bitselect,
	"v128.any_true":                 tokenV128AnyTrue,
	"i8x16.neg":                     tokenI8x16Neg,
	"i16x8.neg":                     tokenI16x8Neg,
	"i32x4.neg":                     tokenI32x4Neg,
	"i64x2.neg":                     tokenI64x2Neg,
	"i8x16.abs":                     tokenI8x16Abs,
	"i16x8.abs":                     tokenI16x8Abs,
	"i32x4.abs":                     tokenI32x4Abs,
	"i64x2.abs":                     tokenI64x2Abs,
	"i8x16.popcnt":                  tokenI8x16Popcnt,
	"i8x16.avgr_u":                  tokenI8x16AvgrU,
	"i16x8.avgr_u":                  tokenI16x8AvgrU,
	"f32x4.neg":                     tokenF32x4Neg,
	"f64x2.neg":                     tokenF64x2Neg,
	"f32x4.abs":                     tokenF32x4Abs,
	"f64x2.abs":                     tokenF64x2Abs,
	"f32x4.sqrt":                    tokenF32x4Sqrt,
	"f64x2.sqrt":                    tokenF64x2Sqrt,
	"f32x4.ceil":                    tokenF32x4Ceil,
	"f64x2.ceil":                    tokenF64x2Ceil,
	"f32x4.floor":                   tokenF32x4Floor,
	"f64x2.floor":                   tokenF64x2Floor,
	"f32x4.trunc":                   tokenF32x4Trunc,
	"f64x2.trunc":                   tokenF64x2Trunc,
	"f32x4.nearest":                 tokenF32x4Nearest,
	"f64x2.nearest":                 tokenF64x2Nearest,
	"i32x4.trunc_sat_f32x4_u":       tokenI32x4TruncSatF32x4U,
	"i32x4.trunc_sat_f32x4_s":       tokenI32x4TruncSatF32x4S,
	"i32x4.trunc_sat_f64x2_u_zero":  tokenI32x4TruncSatF64x2UZero,
	"i32x4.trunc_sat_f64x2_s_zero":  tokenI32x4TruncSatF64x2SZero,
	"f64x2.promote_low_f32x4":       tokenF64x2PromoteLowF32x4,
	"f32x4.demote_f64x2_zero":       tokenF32x4DemoteF64x2Zero,
	"f32x4.convert_i32x4_u":         tokenF32x4ConvertI32x4U,
	"f32x4.convert_i32x4_s":         tokenF32x4ConvertI32x4S,
	"f64x2.convert_low_i32x4_u":     tokenF64x2ConvertLowI32x4U,
	"f64x2.convert_low_i32x4_s":     tokenF64x2ConvertLowI32x4S,
	"i16x8.extadd_pairwise_i8x16_u": tokenI16x8ExtaddPairwiseI8x16U,
	"i16x8.extadd_pairwise_i8x16_s": tokenI16x8ExtaddPairwiseI8x16S,
	"i32x4.extadd_pairwise_i16x8_u": tokenI32x4ExtaddPairwiseI16x8U,
	"i32x4.extadd_pairwise_i16x8_s": tokenI32x4ExtaddPairwiseI16x8S,
	"i8x16.eq":                      tokenI8x16Eq,
	"i16x8.eq":                      tokenI16x8Eq,
	"i32x4.eq":                      tokenI32x4Eq,
	"i64x2.eq":                      tokenI64x2Eq,
	"i8x16.ne":                      tokenI8x16Ne,
	"i16x8.ne":                      tokenI16x8Ne,
	"i32x4.ne":                      tokenI32x4Ne,
	"i64x2.ne":                      tokenI64x2Ne,
	"i8x16.lt_u":                    tokenI8x16LtU,
	"i8x16.lt_s":                    tokenI8x16LtS,
	"i16x8.lt_u":                    tokenI16x8LtU,
	"i16x8.lt_s":                    tokenI16x8LtS,
	"i32x4.lt_u":                    tokenI32x4LtU,
	"i32x4.lt_s":                    tokenI32x4LtS,
	"i64x2.lt_s":                    tokenI64x2LtS,
	"i8x16.le_u":                    tokenI8x16LeU,
	"i8x16.le_s":                    tokenI8x16LeS,
	"i16x8.le_u":                    tokenI16x8LeU,
	"i16x8.le_s":                    tokenI16x8LeS,
	"i32x4.le_u":                    tokenI32x4LeU,
	"i32x4.le_s":                    tokenI32x4LeS,
	"i64x2.le_s":                    tokenI64x2LeS,
	"i8x16.gt_u":                    tokenI8x16GtU,
	"i8x16.gt_s":                    tokenI8x16GtS,
	"i16x8.gt_u":                    tokenI16x8GtU,
	"i16x8.gt_s":                    tokenI16x8GtS,
	"i32x4.gt_u":                    tokenI32x4GtU,
	"i32x4.gt_s":                    tokenI32x4GtS,
	"i64x2.gt_s":                    tokenI64x2GtS,
	"i8x16.ge_u":                    tokenI8x16GeU,
	"i8x16.ge_s":                    tokenI8x16GeS,
	"i16x8.ge_u":                    tokenI16x8GeU,
	"i16x8.ge_s":                    tokenI16x8GeS,
	"i32x4.ge_u":                    tokenI32x4GeU,
	"i32x4.ge_s":                    tokenI32x4GeS,
	"i64x2.ge_s":                    tokenI64x2GeS,
	"f32x4.eq":                      tokenF32x4Eq,
	"f64x2.eq":                      tokenF64x2Eq,
	"f32x4.ne":                      tokenF32x4Ne,
	"f64x2.ne":                      tokenF64x2Ne,
	"f32x4.lt":                      tokenF32x4Lt,
	"f64x2.lt":                      tokenF64x2Lt,
	"f32x4.le":                      tokenF32x4Le,
	"f64x2.le":                      tokenF64x2Le,
	"f32x4.gt":                      tokenF32x4Gt,
	"f64x2.gt":                      tokenF64x2Gt,
	"f32x4.ge":                      tokenF32x4Ge,
	"f64x2.ge":                      tokenF64x2Ge,
	"i8x16.swizzle":                 tokenI8x16Swizzle,
	"i8x16.add":                     tokenI8x16Add,
	"i16x8.add":                     tokenI16x8Add,
	"i32x4.add":                     tokenI32x4Add,
	"i64x2.add":                     tokenI64x2Add,
	"i8x16.sub":                     tokenI8x16Sub,
	"i16x8.sub":                     tokenI16x8Sub,
	"i32x4.sub":                     tokenI32x4Sub,
	"i64x2.sub":                     tokenI64x2Sub,
	"i16x8.mul":                     tokenI16x8Mul,
	"i32x4.mul":                     tokenI32x4Mul,
	"i64x2.mul":                     tokenI64x2Mul,
	"i8x16.add_sat_u":               tokenI8x16AddSatU,
	"i8x16.add_sat_s":               tokenI8x16AddSatS,
	"i16x8.add_sat_u":               tokenI16x8AddSatU,
	"i16x8.add_sat_s":               tokenI16x8AddSatS,
	"i8x16.sub_sat_u":               tokenI8x16SubSatU,
	"i8x16.sub_sat_s":               tokenI8x16SubSatS,
	"i16x8.sub_sat_u":               tokenI16x8SubSatU,
	"i16x8.sub_sat_s":               tokenI16x8SubSatS,
	"i32x4.dot_i16x8_s":             tokenI32x4DotI16x8S,
	"i8x16.min_u":                   tokenI8x16MinU,
	"i16x8.min_u":                   tokenI16x8MinU,
	"i32x4.min_u":                   tokenI32x4MinU,
	"i8x16.min_s":                   tokenI8x16MinS,
	"i16x8.min_s":                   tokenI16x8MinS,
	"i32x4.min_s":                   tokenI32x4MinS,
	"i8x16.max_u":                   tokenI8x16MaxU,
	"i16x8.max_u":                   tokenI16x8MaxU,
	"i32x4.max_u":                   tokenI32x4MaxU,
	"i8x16.max_s":                   tokenI8x16MaxS,
	"i16x8.max_s":                   tokenI16x8MaxS,
	"i32x4.max_s":                   tokenI32x4MaxS,
	"f32x4.add":                     tokenF32x4Add,
	"f64x2.add":                     tokenF64x2Add,
	"f32x4.sub":                     tokenF32x4Sub,
	"f64x2.sub":                     tokenF64x2Sub,
	"f32x4.mul":                     tokenF32x4Mul,
	"f64x2.mul":                     tokenF64x2Mul,
	"f32x4.div":                     tokenF32x4Div,
	"f64x2.div":                     tokenF64x2Div,
	"f32x4.min":                     tokenF32x4Min,
	"f64x2.min":                     tokenF64x2Min,
	"f32x4.max":                     tokenF32x4Max,
	"f64x2.max":                     tokenF64x2Max,
	"f32x4.pmin":                    tokenF32x4Pmin,
	"f64x2.pmin":                    tokenF64x2Pmin,
	"f32x4.pmax":                    tokenF32x4Pmax,
	"f64x2.pmax":                    tokenF64x2Pmax,
	"i16x8.q15mulr_sat_s":           tokenI16x8Q15mulrSatS,
	"i8x16.narrow_i16x8_u":          tokenI8x16NarrowI16x8U,
	"i8x16.narrow_i16x8_s":          tokenI8x16NarrowI16x8S,
	"i16x8.narrow_i32x4_u":          tokenI16x8NarrowI32x4U,
	"i16x8.narrow_i32x4_s":          tokenI16x8NarrowI32x4S,
	"i16x8.extend_low_i8x16_u":      tokenI16x8ExtendLowI8x16U,
	"i16x8.extend_low_i8x16_s":      tokenI16x8ExtendLowI8x16S,
	"i16x8.extend_high_i8x16_u":     tokenI16x8ExtendHighI8x16U,
	"i16x8.extend_high_i8x16_s":     tokenI16x8ExtendHighI8x16S,
	"i32x4.extend_low_i16x8_u":      tokenI32x4ExtendLowI16x8U,
	"i32x4.extend_low_i16x8_s":      tokenI32x4ExtendLowI16x8S,
	"i32x4.extend_high_i16x8_u":     tokenI32x4ExtendHighI16x8U,
	"i32x4.extend_high_i16x8_s":     tokenI32x4ExtendHighI16x8S,
	"i64x2.extend_low_i32x4_u":      tokenI64x2ExtendLowI32x4U,
	"i64x2.extend_low_i32x4_s":      tokenI64x2ExtendLowI32x4S,
	"i64x2.extend_high_i32x4_u":     tokenI64x2ExtendHighI32x4U,
	"i64x2.extend_high_i32x4_s":     tokenI64x2ExtendHighI32x4S,
	"i16x8.extmul_low_i8x16_u":      tokenI16x8ExtmulLowI8x16U,
	"i16x8.extmul_low_i8x16_s":      tokenI16x8ExtmulLowI8x16S,
	"i16x8.extmul_high_i8x16_u":     tokenI16x8ExtmulHighI8x16U,
	"i16x8.extmul_high_i8x16_s":     tokenI16x8ExtmulHighI8x16S,
	"i32x4.extmul_low_i16x8_u":      tokenI32x4ExtmulLowI16x8U,
	"i32x4.extmul_low_i16x8_s":      tokenI32x4ExtmulLowI16x8S,
	"i32x4.extmul_high_i16x8_u":     tokenI32x4ExtmulHighI16x8U,
	"i32x4.extmul_high_i16x8_s":     tokenI32x4ExtmulHighI16x8S,
	"i64x2.extmul_low_i32x4_u":      tokenI64x2ExtmulLowI32x4U,
	"i64x2.extmul_low_i32x4_s":      tokenI64x2ExtmulLowI32x4S,
	"i64x2.extmul_high_i32x4_u":     tokenI64x2ExtmulHighI32x4U,
	"i64x2.extmul_high_i32x4_s":     tokenI64x2ExtmulHighI32x4S,
	"i8x16.all_true":                tokenI8x16AllTrue,
	"i16x8.all_true":                tokenI16x8AllTrue,
	"i32x4.all_true":                tokenI32x4AllTrue,
	"i64x2.all_true":                tokenI64x2AllTrue,
	"i8x16.bitmask":                 tokenI8x16Bitmask,
	"i16x8.bitmask":                 tokenI16x8Bitmask,
	"i32x4.bitmask":                 tokenI32x4Bitmask,
	"i64x2.bitmask":                 tokenI64x2Bitmask,
	"i8x16.shl":                     tokenI8x16Shl,
	"i16x8.shl":                     tokenI16x8Shl,
	"i32x4.shl":                     tokenI32x4Shl,
	"i64x2.shl":                     tokenI64x2Shl,
	"i8x16.shr_u":                   tokenI8x16ShrU,
	"i8x16.shr_s":                   tokenI8x16ShrS,
	"i16x8.shr_u":                   tokenI16x8ShrU,
	"i16x8.shr_s":                   tokenI16x8ShrS,
	"i32x4.shr_u":                   tokenI32x4ShrU,
	"i32x4.shr_s":                   tokenI32x4ShrS,
	"i64x2.shr_u":                   tokenI64x2ShrU,
	"i64x2.shr_s":                   tokenI64x2ShrS,
	"i8x16.shuffle":                 tokenI8x16Shuffle,
	"i8x16.splat":                   tokenI8x16Splat,
	"i16x8.splat":                   tokenI16x8Splat,
	"i32x4.splat":                   tokenI32x4Splat,
	"i64x2.splat":                   tokenI64x2Splat,
	"f32x4.splat":                   tokenF32x4Splat,
	"f64x2.splat":                   tokenF64x2Splat,
	"i8x16.extract_lane_u":          tokenI8x16ExtractLaneU,
	"i8x16.extract_lane_s":          tokenI8x16ExtractLaneS,
	"i16x8.extract_lane_u":          tokenI16x8ExtractLaneU,
	"i16x8.extract_lane_s":          tokenI16x8ExtractLaneS,
	"i32x4.extract_lane":            tokenI32x4ExtractLane,
	"i64x2.extract_lane":            tokenI64x2ExtractLane,
	"f32x4.extract_lane":            tokenF32x4ExtractLane,
	"f64x2.extract_lane":            tokenF64x2ExtractLane,
	"i8x16.replace_lane":            tokenI8x16ReplaceLane,
	"i16x8.replace_lane":            tokenI16x8ReplaceLane,
	"i32x4.replace_lane":            tokenI32x4ReplaceLane,
	"i64x2.replace_lane":            tokenI64x2ReplaceLane,
	"f32x4.replace_lane":            tokenF32x4ReplaceLane,
	"f64x2.replace_lane":            tokenF64x2ReplaceLane,
}

func (t token) String() string {
	switch {
	case t.kind == tokenEOF:
		return "EOF"
	case t.kind == tokenError:
		return string(t.val)
	case t.kind == tokenKeyword:
		return fmt.Sprintf("<%s>", t.val)
	case len(t.val) > 10:
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

type stateFn func(l *lexer) stateFn

type lexer struct {
	input  []byte
	pos    int
	start  int
	width  int
	state  stateFn
	tokens chan token
}

func (l *lexer) nextToken() token {
	for {
		select {
		case token := <-l.tokens:
			return token
		default:
			if l.state == nil {
				close(l.tokens)
				return token{kind: tokenEOF}
			}
			l.state = l.state(l)
		}
	}
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	var r rune
	r, l.width = utf8.DecodeRune(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	w := l.width
	r := l.next()
	l.backup()
	l.width = w
	return r
}

func (l *lexer) lexeme() []byte {
	return l.input[l.start:l.pos]
}

func (l *lexer) emit(kind tokenKind) {
	l.tokens <- token{kind, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) emitWithData(kind tokenKind, data []byte) {
	l.tokens <- token{kind, data}
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) int {
	count := 0
	for strings.IndexRune(valid, l.next()) >= 0 {
		count++
	}
	l.backup()
	return count
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.tokens <- token{tokenError, []byte(fmt.Sprintf(format, args...))}
	return nil
}

func lexDefault(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(tokenEOF)
			return nil
		case isSpace(r):
			l.ignore()
		case r == ';':
			return lexComment
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case r == '$':
			l.backup()
			return lexIdentifier
		case r == '"':
			return lexString
		case r == '(':
			if l.peek() == ';' {
				return lexBlockComment
			} else {
				l.emit(tokenLParen)
			}
			return lexDefault
		case r == ')':
			l.emit(tokenRParen)
			return lexDefault
		case isLowercaseLetter(r):
			l.backup()
			return lexKeyword
		default:
			return l.errorf("unknown token: %q", r)
		}
	}
}

func lexBlockComment(l *lexer) stateFn {
	level := 1
	l.accept(";")
	for level > 0 {
		switch r := l.next(); {
		case r == '(':
			if l.peek() == ';' {
				level++
			}
		case r == ';':
			if l.peek() == ')' {
				level--
			}
		}
	}
	return lexDefault
}

func lexNumber(l *lexer) stateFn {
	l.accept(sign)
	// is it hex?
	valid := digit
	if l.accept("0") && l.accept("xX") {
		valid = hexDigit
	}

	valid += "_"

	l.acceptRun(valid)
	if l.accept(".") {
		l.acceptRun(valid)
	}
	if l.accept("eEpP") {
		l.accept(sign)
		l.acceptRun(digit + "_")
	}

	l.emit(tokenNumber)
	return lexDefault
}

func lexKeyword(l *lexer) stateFn {
	l.acceptRun(keyword)
	if kind, ok := key[string(l.lexeme())]; ok {
		l.emit(kind)
	} else {
		l.emit(tokenKeyword)
	}
	return lexDefault
}

func lexComment(l *lexer) stateFn {
	if !l.accept(";") {
		return l.errorf("expected ';' but got: %q", l.next())
	}

	for {
		if r := l.next(); r == eof || r == '\n' {
			break
		}
	}
	return lexDefault
}

// https://webassembly.github.io/spec/core/text/values.html#text-id
func lexIdentifier(l *lexer) stateFn {
	if l.acceptRun(idChar) < 2 {
		return l.errorf("empty identifier")
	}
	l.emit(tokenIdent)
	return lexDefault
}

// https://webassembly.github.io/spec/core/text/values.html#strings
func lexString(l *lexer) stateFn {
	s := strings.Builder{}
	for {
		switch r := l.next(); {
		case r == '\\':
			s1, err := escapeSeq(l)
			s.WriteString(s1)
			if err != nil {
				return l.errorf("invalid escape sequence: %q", l.input[l.start:l.pos])
			}
		case r == '"':
			l.emitWithData(tokenString, []byte(s.String()))
			return lexDefault
		case r == eof || r == '\n' || r == '\r':
			return l.errorf("unclosed quote: %q", l.input[l.start:l.pos])
		default:
			s.WriteRune(r)
		}
	}
}

func escapeSeq(l *lexer) (string, error) {
	switch r := l.next(); {
	case r == 't':
		return "\t", nil
	case r == 'r':
		return "\r", nil
	case r == '"':
		return "\"", nil
	case r == '\'':
		return "'", nil
	case r == '\\':
		return "\\", nil
	case isHexDigit(r):
		if r2 := l.next(); r2 != eof && isHexDigit(r2) {
			v, _ := strconv.ParseUint(string(r)+string(r2), 16, 16)
			return string(rune(v)), nil
		} else {
			return "", fmt.Errorf("invalid escape sequence: %q%q", r, r2)
		}
	case r == 'u':
		if !l.accept("{") {
			return "", fmt.Errorf("invalid unicode")
		}

		var s string
		for i := 0; i < 7; i++ {
			d := l.next()
			if d == '}' {
				l.backup()
				break
			}
			if d == eof {
				return "", fmt.Errorf("unterminated unicode")
			}

			if !isHexDigit(d) {
				return "", fmt.Errorf("invalid hex digit: %q", d)
			}

			s += string(d)
		}

		v, _ := strconv.ParseUint(s, 16, 16)

		if !l.accept("}") {
			return "", fmt.Errorf("invalid unicode")
		}
		return string(rune(v)), nil
	default:
		return "", fmt.Errorf("unexpected escape sequence: %q", r)
	}
}

func isLowercaseLetter(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isHexDigit(r rune) bool {
	return strings.IndexRune(hexDigit, r) >= 0
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t' || r == '\r'
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func NewLexer(input []byte) *lexer {
	return &lexer{
		input:  input,
		state:  lexDefault,
		tokens: make(chan token, 3),
	}
}
