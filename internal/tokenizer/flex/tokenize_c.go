package flex

// #include <stdlib.h>
// #include "linguist.h"
// #include "lex.linguist_yy.h"
// int linguist_yywrap(yyscan_t yyscanner) {
// 	return 1;
// }
import "C"
import "unsafe"

// TokenizeC is only calling a C-flex based tokenizer from linguist
func TokenizeC(content []byte) []string {
	cs := C.CBytes(content)
	defer C.free(unsafe.Pointer(cs))
	// C.tokenizer_extract_tokens((*C.char)(cs))
	return nil
}

const maxTokenLen = 32


// TokenizeFlex implements tokenizer by calling Flex generated code from linguist in C
func TokenizeFlex(content []byte) []string {
	var buf C.YY_BUFFER_STATE
	var scanner C.yyscan_t
	var extra C.struct_tokenizer_extra
	// var scanner *C.yyscan_t = (*C.yyscan_t)(C.malloc(C.sizeof_yyscan_t))
	// var extra *C.struct_tokenizer_extra = (*C.struct_tokenizer_extra)(C.malloc(C.sizeof_struct_tokenizer_extra))
	var _len C.ulong
	var r C.int

	_len = C.ulong(len(content))
	cs := C.CBytes(content)
	defer C.free(unsafe.Pointer(cs))

	C.linguist_yylex_init_extra(&extra, &scanner)
	buf = C.linguist_yy_scan_bytes((*C.char)(cs), _len, scanner)


	ary := []string{}
	for {
		extra._type = C.NO_ACTION
		extra.token = nil
		r = C.linguist_yylex(scanner)
		switch (extra._type) {
		case C.NO_ACTION:
			break
		case C.REGULAR_TOKEN:
			_len = C.strlen(extra.token)
			if (_len <= maxTokenLen) {
				ary = append(ary, C.GoStringN(extra.token, (C.int)(_len)))
				//rb_ary_push(ary, rb_str_new(extra.token, len))
			}
			C.free(unsafe.Pointer(extra.token))
			break
		case C.SHEBANG_TOKEN:
			_len = C.strlen(extra.token)
			if (_len <= maxTokenLen) {
				s := "SHEBANG#!" + C.GoStringN(extra.token, (C.int)(_len))
				ary = append(ary, s)
				//s = rb_str_new2("SHEBANG#!");
				//rb_str_cat(s, extra.token, len);
				//rb_ary_push(ary, s);
			}
			C.free(unsafe.Pointer(extra.token))
			break
		case C.SGML_TOKEN:
			_len = C.strlen(extra.token)
			if (_len <= maxTokenLen) {
				s := C.GoStringN(extra.token, (C.int)(_len)) + ">"
				ary = append(ary, s)
				//s = rb_str_new(extra.token, len);
				//rb_str_cat2(s, ">");
				//rb_ary_push(ary, s);
			}
			C.free(unsafe.Pointer(extra.token))
			break
		}
		if r == 0 {
			break
		}
	}

	C.linguist_yy_delete_buffer(buf, scanner)
	C.linguist_yylex_destroy(scanner)
	// C.free(unsafe.Pointer(extra))
	// C.free(unsafe.Pointer(scanner))

	return ary
}
