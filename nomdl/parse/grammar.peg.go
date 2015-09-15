package parse

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type alias struct {
	Name   string
	Target string
}

type namespaceIdent struct {
	Namespace string
	ID        string
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Package",
			pos:  position{line: 15, col: 1, offset: 133},
			expr: &actionExpr{
				pos: position{line: 15, col: 12, offset: 144},
				run: (*parser).callonPackage1,
				expr: &seqExpr{
					pos: position{line: 15, col: 12, offset: 144},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 15, col: 12, offset: 144},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 15, col: 14, offset: 146},
							label: "dd",
							expr: &oneOrMoreExpr{
								pos: position{line: 15, col: 17, offset: 149},
								expr: &ruleRefExpr{
									pos:  position{line: 15, col: 17, offset: 149},
									name: "Definition",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 15, col: 29, offset: 161},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 15, col: 31, offset: 163},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Definition",
			pos:  position{line: 50, col: 1, offset: 1101},
			expr: &choiceExpr{
				pos: position{line: 50, col: 15, offset: 1115},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 50, col: 15, offset: 1115},
						name: "Struct",
					},
					&ruleRefExpr{
						pos:  position{line: 50, col: 24, offset: 1124},
						name: "Using",
					},
					&ruleRefExpr{
						pos:  position{line: 50, col: 32, offset: 1132},
						name: "Alias",
					},
					&ruleRefExpr{
						pos:  position{line: 50, col: 40, offset: 1140},
						name: "Enum",
					},
				},
			},
		},
		{
			name: "Alias",
			pos:  position{line: 52, col: 1, offset: 1146},
			expr: &actionExpr{
				pos: position{line: 52, col: 10, offset: 1155},
				run: (*parser).callonAlias1,
				expr: &seqExpr{
					pos: position{line: 52, col: 10, offset: 1155},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 52, col: 10, offset: 1155},
							val:        "alias",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 52, col: 18, offset: 1163},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 52, col: 20, offset: 1165},
							label: "i",
							expr: &ruleRefExpr{
								pos:  position{line: 52, col: 22, offset: 1167},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 52, col: 28, offset: 1173},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 52, col: 30, offset: 1175},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 52, col: 34, offset: 1179},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 52, col: 36, offset: 1181},
							val:        "import",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 52, col: 45, offset: 1190},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 52, col: 47, offset: 1192},
							label: "q",
							expr: &ruleRefExpr{
								pos:  position{line: 52, col: 49, offset: 1194},
								name: "QuotedString",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 52, col: 62, offset: 1207},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Enum",
			pos:  position{line: 56, col: 1, offset: 1257},
			expr: &actionExpr{
				pos: position{line: 56, col: 9, offset: 1265},
				run: (*parser).callonEnum1,
				expr: &seqExpr{
					pos: position{line: 56, col: 9, offset: 1265},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 56, col: 9, offset: 1265},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 16, offset: 1272},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 56, col: 18, offset: 1274},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 56, col: 21, offset: 1277},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 27, offset: 1283},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 56, col: 29, offset: 1285},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 33, offset: 1289},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 56, col: 35, offset: 1291},
							label: "l",
							expr: &oneOrMoreExpr{
								pos: position{line: 56, col: 37, offset: 1293},
								expr: &ruleRefExpr{
									pos:  position{line: 56, col: 37, offset: 1293},
									name: "EnumEntry",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 48, offset: 1304},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 56, col: 50, offset: 1306},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 56, col: 54, offset: 1310},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "EnumEntry",
			pos:  position{line: 65, col: 1, offset: 1485},
			expr: &actionExpr{
				pos: position{line: 65, col: 14, offset: 1498},
				run: (*parser).callonEnumEntry1,
				expr: &seqExpr{
					pos: position{line: 65, col: 14, offset: 1498},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 65, col: 14, offset: 1498},
							label: "i",
							expr: &ruleRefExpr{
								pos:  position{line: 65, col: 16, offset: 1500},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 65, col: 22, offset: 1506},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Using",
			pos:  position{line: 70, col: 1, offset: 1538},
			expr: &actionExpr{
				pos: position{line: 70, col: 10, offset: 1547},
				run: (*parser).callonUsing1,
				expr: &seqExpr{
					pos: position{line: 70, col: 10, offset: 1547},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 70, col: 10, offset: 1547},
							val:        "using",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 70, col: 18, offset: 1555},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 70, col: 20, offset: 1557},
							label: "ct",
							expr: &ruleRefExpr{
								pos:  position{line: 70, col: 23, offset: 1560},
								name: "CompoundType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 70, col: 36, offset: 1573},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Struct",
			pos:  position{line: 75, col: 1, offset: 1597},
			expr: &actionExpr{
				pos: position{line: 75, col: 11, offset: 1607},
				run: (*parser).callonStruct1,
				expr: &seqExpr{
					pos: position{line: 75, col: 11, offset: 1607},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 75, col: 11, offset: 1607},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 20, offset: 1616},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 75, col: 22, offset: 1618},
							label: "i",
							expr: &ruleRefExpr{
								pos:  position{line: 75, col: 24, offset: 1620},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 30, offset: 1626},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 75, col: 32, offset: 1628},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 36, offset: 1632},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 75, col: 38, offset: 1634},
							label: "l",
							expr: &oneOrMoreExpr{
								pos: position{line: 75, col: 40, offset: 1636},
								expr: &ruleRefExpr{
									pos:  position{line: 75, col: 40, offset: 1636},
									name: "StructEntry",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 53, offset: 1649},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 75, col: 55, offset: 1651},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 59, offset: 1655},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "StructEntry",
			pos:  position{line: 100, col: 1, offset: 2314},
			expr: &choiceExpr{
				pos: position{line: 100, col: 16, offset: 2329},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 100, col: 16, offset: 2329},
						name: "Union",
					},
					&ruleRefExpr{
						pos:  position{line: 100, col: 24, offset: 2337},
						name: "Field",
					},
				},
			},
		},
		{
			name: "Union",
			pos:  position{line: 103, col: 1, offset: 2345},
			expr: &actionExpr{
				pos: position{line: 103, col: 10, offset: 2354},
				run: (*parser).callonUnion1,
				expr: &seqExpr{
					pos: position{line: 103, col: 10, offset: 2354},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 103, col: 10, offset: 2354},
							val:        "union",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 103, col: 18, offset: 2362},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 103, col: 20, offset: 2364},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 103, col: 24, offset: 2368},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 103, col: 26, offset: 2370},
							label: "u",
							expr: &oneOrMoreExpr{
								pos: position{line: 103, col: 28, offset: 2372},
								expr: &ruleRefExpr{
									pos:  position{line: 103, col: 28, offset: 2372},
									name: "Field",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 103, col: 35, offset: 2379},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 103, col: 37, offset: 2381},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 103, col: 41, offset: 2385},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 118, col: 1, offset: 2743},
			expr: &actionExpr{
				pos: position{line: 118, col: 10, offset: 2752},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 118, col: 10, offset: 2752},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 118, col: 10, offset: 2752},
							label: "i",
							expr: &ruleRefExpr{
								pos:  position{line: 118, col: 12, offset: 2754},
								name: "Ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 18, offset: 2760},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 118, col: 20, offset: 2762},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 24, offset: 2766},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 118, col: 26, offset: 2768},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 118, col: 28, offset: 2770},
								name: "Type",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 118, col: 33, offset: 2775},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Type",
			pos:  position{line: 122, col: 1, offset: 2826},
			expr: &actionExpr{
				pos: position{line: 122, col: 9, offset: 2834},
				run: (*parser).callonType1,
				expr: &labeledExpr{
					pos:   position{line: 122, col: 9, offset: 2834},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 122, col: 12, offset: 2837},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 122, col: 12, offset: 2837},
								name: "PrimitiveType",
							},
							&ruleRefExpr{
								pos:  position{line: 122, col: 28, offset: 2853},
								name: "CompoundType",
							},
							&ruleRefExpr{
								pos:  position{line: 122, col: 43, offset: 2868},
								name: "Union",
							},
							&ruleRefExpr{
								pos:  position{line: 122, col: 51, offset: 2876},
								name: "NamespaceIdent",
							},
						},
					},
				},
			},
		},
		{
			name: "CompoundType",
			pos:  position{line: 135, col: 1, offset: 3183},
			expr: &choiceExpr{
				pos: position{line: 135, col: 17, offset: 3199},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 135, col: 17, offset: 3199},
						run: (*parser).callonCompoundType2,
						expr: &seqExpr{
							pos: position{line: 135, col: 17, offset: 3199},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 135, col: 17, offset: 3199},
									val:        "List",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 24, offset: 3206},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 135, col: 26, offset: 3208},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 30, offset: 3212},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 135, col: 32, offset: 3214},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 135, col: 34, offset: 3216},
										name: "Type",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 39, offset: 3221},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 135, col: 41, offset: 3223},
									val:        ")",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 135, col: 45, offset: 3227},
									name: "_",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 137, col: 5, offset: 3302},
						run: (*parser).callonCompoundType13,
						expr: &seqExpr{
							pos: position{line: 137, col: 5, offset: 3302},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 137, col: 5, offset: 3302},
									val:        "Map",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 137, col: 11, offset: 3308},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 137, col: 13, offset: 3310},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 137, col: 17, offset: 3314},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 137, col: 19, offset: 3316},
									label: "k",
									expr: &ruleRefExpr{
										pos:  position{line: 137, col: 21, offset: 3318},
										name: "Type",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 137, col: 26, offset: 3323},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 137, col: 28, offset: 3325},
									val:        ",",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 137, col: 32, offset: 3329},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 137, col: 34, offset: 3331},
									label: "v",
									expr: &ruleRefExpr{
										pos:  position{line: 137, col: 36, offset: 3333},
										name: "Type",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 137, col: 41, offset: 3338},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 137, col: 43, offset: 3340},
									val:        ")",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 137, col: 47, offset: 3344},
									name: "_",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 139, col: 5, offset: 3431},
						run: (*parser).callonCompoundType29,
						expr: &seqExpr{
							pos: position{line: 139, col: 5, offset: 3431},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 139, col: 5, offset: 3431},
									val:        "Set",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 139, col: 11, offset: 3437},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 139, col: 13, offset: 3439},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 139, col: 17, offset: 3443},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 139, col: 19, offset: 3445},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 139, col: 21, offset: 3447},
										name: "Type",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 139, col: 26, offset: 3452},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 139, col: 28, offset: 3454},
									val:        ")",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 139, col: 32, offset: 3458},
									name: "_",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 141, col: 5, offset: 3532},
						run: (*parser).callonCompoundType40,
						expr: &seqExpr{
							pos: position{line: 141, col: 5, offset: 3532},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 141, col: 5, offset: 3532},
									val:        "Ref",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 141, col: 11, offset: 3538},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 141, col: 13, offset: 3540},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 141, col: 17, offset: 3544},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 141, col: 19, offset: 3546},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 141, col: 21, offset: 3548},
										name: "Type",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 141, col: 26, offset: 3553},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 141, col: 28, offset: 3555},
									val:        ")",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 141, col: 32, offset: 3559},
									name: "_",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PrimitiveType",
			pos:  position{line: 145, col: 1, offset: 3632},
			expr: &actionExpr{
				pos: position{line: 145, col: 18, offset: 3649},
				run: (*parser).callonPrimitiveType1,
				expr: &labeledExpr{
					pos:   position{line: 145, col: 18, offset: 3649},
					label: "p",
					expr: &choiceExpr{
						pos: position{line: 145, col: 21, offset: 3652},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 145, col: 21, offset: 3652},
								val:        "UInt64",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 32, offset: 3663},
								val:        "UInt32",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 43, offset: 3674},
								val:        "UInt16",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 54, offset: 3685},
								val:        "UInt8",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 64, offset: 3695},
								val:        "Int64",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 74, offset: 3705},
								val:        "Int32",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 84, offset: 3715},
								val:        "Int16",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 94, offset: 3725},
								val:        "Int8",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 103, offset: 3734},
								val:        "Float64",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 115, offset: 3746},
								val:        "Float32",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 127, offset: 3758},
								val:        "Bool",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 136, offset: 3767},
								val:        "String",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 147, offset: 3778},
								val:        "Blob",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 145, col: 156, offset: 3787},
								val:        "Value",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "QuotedString",
			pos:  position{line: 149, col: 1, offset: 3856},
			expr: &actionExpr{
				pos: position{line: 149, col: 17, offset: 3872},
				run: (*parser).callonQuotedString1,
				expr: &seqExpr{
					pos: position{line: 149, col: 17, offset: 3872},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 149, col: 17, offset: 3872},
							val:        "\"",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 149, col: 21, offset: 3876},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 149, col: 23, offset: 3878},
								name: "String",
							},
						},
						&litMatcher{
							pos:        position{line: 149, col: 30, offset: 3885},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 153, col: 1, offset: 3918},
			expr: &actionExpr{
				pos: position{line: 153, col: 11, offset: 3928},
				run: (*parser).callonString1,
				expr: &choiceExpr{
					pos: position{line: 153, col: 12, offset: 3929},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 153, col: 12, offset: 3929},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 153, col: 12, offset: 3929},
									name: "StringPiece",
								},
								&litMatcher{
									pos:        position{line: 153, col: 24, offset: 3941},
									val:        "\\\"",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 153, col: 29, offset: 3946},
									name: "StringPiece",
								},
								&litMatcher{
									pos:        position{line: 153, col: 41, offset: 3958},
									val:        "\\\"",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 153, col: 46, offset: 3963},
									name: "StringPiece",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 153, col: 60, offset: 3977},
							name: "StringPiece",
						},
					},
				},
			},
		},
		{
			name: "StringPiece",
			pos:  position{line: 157, col: 1, offset: 4023},
			expr: &zeroOrMoreExpr{
				pos: position{line: 157, col: 16, offset: 4038},
				expr: &choiceExpr{
					pos: position{line: 157, col: 17, offset: 4039},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 157, col: 17, offset: 4039},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 157, col: 17, offset: 4039},
									val:        "\\",
									ignoreCase: false,
								},
								&notExpr{
									pos: position{line: 157, col: 21, offset: 4043},
									expr: &litMatcher{
										pos:        position{line: 157, col: 22, offset: 4044},
										val:        "\"",
										ignoreCase: false,
									},
								},
							},
						},
						&charClassMatcher{
							pos:        position{line: 157, col: 28, offset: 4050},
							val:        "[^\"\\\\]",
							chars:      []rune{'"', '\\'},
							ignoreCase: false,
							inverted:   true,
						},
					},
				},
			},
		},
		{
			name: "NamespaceIdent",
			pos:  position{line: 159, col: 1, offset: 4060},
			expr: &actionExpr{
				pos: position{line: 159, col: 19, offset: 4078},
				run: (*parser).callonNamespaceIdent1,
				expr: &seqExpr{
					pos: position{line: 159, col: 19, offset: 4078},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 159, col: 19, offset: 4078},
							label: "n",
							expr: &zeroOrMoreExpr{
								pos: position{line: 159, col: 21, offset: 4080},
								expr: &seqExpr{
									pos: position{line: 159, col: 22, offset: 4081},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 159, col: 22, offset: 4081},
											name: "Ident",
										},
										&litMatcher{
											pos:        position{line: 159, col: 28, offset: 4087},
											val:        ".",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 159, col: 34, offset: 4093},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 159, col: 37, offset: 4096},
								name: "Ident",
							},
						},
					},
				},
			},
		},
		{
			name: "Ident",
			pos:  position{line: 168, col: 1, offset: 4294},
			expr: &actionExpr{
				pos: position{line: 168, col: 10, offset: 4303},
				run: (*parser).callonIdent1,
				expr: &seqExpr{
					pos: position{line: 168, col: 10, offset: 4303},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 168, col: 10, offset: 4303},
							val:        "[\\pL_]",
							chars:      []rune{'_'},
							classes:    []*unicode.RangeTable{rangeTable("L")},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 168, col: 17, offset: 4310},
							expr: &charClassMatcher{
								pos:        position{line: 168, col: 17, offset: 4310},
								val:        "[\\pL\\pN_]",
								chars:      []rune{'_'},
								classes:    []*unicode.RangeTable{rangeTable("L"), rangeTable("N")},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"optional whitespace\"",
			pos:         position{line: 172, col: 1, offset: 4354},
			expr: &actionExpr{
				pos: position{line: 172, col: 28, offset: 4381},
				run: (*parser).callon_1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 172, col: 28, offset: 4381},
					expr: &charClassMatcher{
						pos:        position{line: 172, col: 28, offset: 4381},
						val:        "[\\r\\n\\t\\pZ]",
						chars:      []rune{'\r', '\n', '\t'},
						classes:    []*unicode.RangeTable{rangeTable("Z")},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 176, col: 1, offset: 4416},
			expr: &seqExpr{
				pos: position{line: 176, col: 8, offset: 4423},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 176, col: 8, offset: 4423},
						name: "_",
					},
					&notExpr{
						pos: position{line: 176, col: 10, offset: 4425},
						expr: &anyMatcher{
							line: 176, col: 11, offset: 4426,
						},
					},
				},
			},
		},
	},
}

func (c *current) onPackage1(dd interface{}) (interface{}, error) {
	aliases := map[string]string{}
	usings := []TypeRef{}
	named := map[string]TypeRef{}
	for _, d := range dd.([]interface{}) {
		switch d := d.(type) {
		default:
			return nil, fmt.Errorf("Unknown definition: %v", d)
		case alias:
			if _, present := aliases[d.Name]; present {
				return nil, fmt.Errorf("Redefinition of " + d.Name)
			}
			aliases[d.Name] = d.Target
		case TypeRef:
			switch d.Desc.Kind() {
			default:
				return nil, fmt.Errorf("%v can't be defined at the top-level", d)
			case ListKind, MapKind, RefKind, SetKind:
				for _, u := range usings {
					if u.Equals(d) {
						return nil, fmt.Errorf("%v is a duplicate using declaration", d)
					}
				}
				usings = append(usings, d)
			case EnumKind, StructKind:
				if _, present := named[d.Name]; present {
					return nil, fmt.Errorf("Redefinition of enum " + d.Name)
				}
				named[d.Name] = d
			}
		}
	}
	return Package{aliases, usings, named}, nil
}

func (p *parser) callonPackage1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPackage1(stack["dd"])
}

func (c *current) onAlias1(i, q interface{}) (interface{}, error) {
	return alias{i.(string), q.(string)}, nil
}

func (p *parser) callonAlias1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAlias1(stack["i"], stack["q"])
}

func (c *current) onEnum1(id, l interface{}) (interface{}, error) {
	entries := l.([]interface{})
	ids := make([]string, len(entries))
	for i, e := range entries {
		ids[i] = e.(string)
	}
	return makeEnumTypeRef(id.(string), ids), nil
}

func (p *parser) callonEnum1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnum1(stack["id"], stack["l"])
}

func (c *current) onEnumEntry1(i interface{}) (interface{}, error) {
	return i.(string), nil
}

func (p *parser) callonEnumEntry1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumEntry1(stack["i"])
}

func (c *current) onUsing1(ct interface{}) (interface{}, error) {
	return ct, nil
}

func (p *parser) callonUsing1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUsing1(stack["ct"])
}

func (c *current) onStruct1(i, l interface{}) (interface{}, error) {
	ll := l.([]interface{})
	var u *UnionDesc
	fieldNames := map[string]bool{}
	fields := make([]Field, 0, len(ll))
	for _, e := range ll {
		switch e := e.(type) {
		case UnionDesc:
			if u != nil {
				return nil, fmt.Errorf("Only one anonymous union per struct.")
			}
			u = &e
		case Field:
			if fieldNames[e.Name] {
				return nil, fmt.Errorf("Redefinition of field %s in struct %s", e.Name, i.(string))
			}
			fieldNames[e.Name] = true
			fields = append(fields, e)
		default:
			return nil, fmt.Errorf("Structs must be made up of field declarations and at most one anonymous union.")
		}
	}
	return makeStructTypeRef(i.(string), fields, u), nil
}

func (p *parser) callonStruct1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStruct1(stack["i"], stack["l"])
}

func (c *current) onUnion1(u interface{}) (interface{}, error) {
	uu := u.([]interface{})
	choiceNames := map[string]bool{}
	desc := UnionDesc{Choices: make([]Field, 0, len(uu))}
	for _, f := range uu {
		ff := f.(Field)
		if choiceNames[ff.Name] {
			return nil, fmt.Errorf("Redefinition of union choice %s", ff.Name)
		}
		choiceNames[ff.Name] = true
		desc.Choices = append(desc.Choices, ff)
	}
	return desc, nil
}

func (p *parser) callonUnion1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnion1(stack["u"])
}

func (c *current) onField1(i, t interface{}) (interface{}, error) {
	return Field{i.(string), t.(TypeRef)}, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["i"], stack["t"])
}

func (c *current) onType1(t interface{}) (interface{}, error) {
	switch t := t.(type) {
	case TypeRef:
		return t, nil
	case UnionDesc:
		return makeStructTypeRef("", nil, &t), nil
	case namespaceIdent:
		return TypeRef{PkgRef: t.Namespace, Name: t.ID}, nil
	default:
		return nil, fmt.Errorf("%v is %T, not something that satisfies TypeRef", t, t)
	}
}

func (p *parser) callonType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onType1(stack["t"])
}

func (c *current) onCompoundType2(t interface{}) (interface{}, error) {
	return makeCompoundTypeRef(ListKind, []TypeRef{t.(TypeRef)}), nil
}

func (p *parser) callonCompoundType2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCompoundType2(stack["t"])
}

func (c *current) onCompoundType13(k, v interface{}) (interface{}, error) {
	return makeCompoundTypeRef(MapKind, []TypeRef{k.(TypeRef), v.(TypeRef)}), nil
}

func (p *parser) callonCompoundType13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCompoundType13(stack["k"], stack["v"])
}

func (c *current) onCompoundType29(t interface{}) (interface{}, error) {
	return makeCompoundTypeRef(SetKind, []TypeRef{t.(TypeRef)}), nil
}

func (p *parser) callonCompoundType29() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCompoundType29(stack["t"])
}

func (c *current) onCompoundType40(t interface{}) (interface{}, error) {
	return makeCompoundTypeRef(RefKind, []TypeRef{t.(TypeRef)}), nil
}

func (p *parser) callonCompoundType40() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCompoundType40(stack["t"])
}

func (c *current) onPrimitiveType1(p interface{}) (interface{}, error) {
	return makePrimitiveTypeRef(string(p.([]uint8))), nil
}

func (p *parser) callonPrimitiveType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPrimitiveType1(stack["p"])
}

func (c *current) onQuotedString1(n interface{}) (interface{}, error) {
	return n.(string), nil
}

func (p *parser) callonQuotedString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuotedString1(stack["n"])
}

func (c *current) onString1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

func (c *current) onNamespaceIdent1(n, id interface{}) (interface{}, error) {
	nn := n.([]interface{})
	ns := make([]string, len(nn))
	for i, e := range nn {
		ns[i] = e.([]interface{})[0].(string)
	}
	return namespaceIdent{strings.Join(ns, "."), id.(string)}, nil
}

func (p *parser) callonNamespaceIdent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamespaceIdent1(stack["n"], stack["id"])
}

func (c *current) onIdent1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdent1()
}

func (c *current) on_1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callon_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.on_1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}