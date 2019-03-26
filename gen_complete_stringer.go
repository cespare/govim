// Code generated by "stringer -type=Complete -linecomment -output gen_complete_stringer.go"; DO NOT EDIT.

package govim

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CompleteArglist-0]
	_ = x[CompleteAugroup-1]
	_ = x[CompleteBuffer-2]
	_ = x[CompleteBehave-3]
	_ = x[CompleteColor-4]
	_ = x[CompleteCommand-5]
	_ = x[CompleteCompiler-6]
	_ = x[CompleteCscope-7]
	_ = x[CompleteDir-8]
	_ = x[CompleteEnvironment-9]
	_ = x[CompleteEvent-10]
	_ = x[CompleteExpression-11]
	_ = x[CompleteFile-12]
	_ = x[CompleteFileInPath-13]
	_ = x[CompleteFiletype-14]
	_ = x[CompleteFunction-15]
	_ = x[CompleteHelp-16]
	_ = x[CompleteHighlight-17]
	_ = x[CompleteHistory-18]
	_ = x[CompleteLocale-19]
	_ = x[CompleteMapclear-20]
	_ = x[CompleteMapping-21]
	_ = x[CompleteMenu-22]
	_ = x[CompleteMessages-23]
	_ = x[CompleteOption-24]
	_ = x[CompletePackadd-25]
	_ = x[CompleteShellCmd-26]
	_ = x[CompleteSign-27]
	_ = x[CompleteSyntax-28]
	_ = x[CompleteSyntime-29]
	_ = x[CompleteTag-30]
	_ = x[CompleteTagListFiles-31]
	_ = x[CompleteUser-32]
	_ = x[CompleteVar-33]
}

const _Complete_name = "-complete=arglist-complete=augroup-complete=buffer-complete=behave-complete=color-complete=command-complete=compiler-complete=cscope-complete=dir-complete=environment-complete=event-complete=expression-complete=file-complete=file_in_path-complete=filetype-complete=function-complete=help-complete=highlight-complete=history-complete=locale-complete=mapclear-complete=mapping-complete=menu-complete=messages-complete=option-complete=packadd-complete=shellcmd-complete=sign-complete=syntax-complete=syntime-complete=tag-complete=tag_listfiles-complete=user-complete=var"

var _Complete_index = [...]uint16{0, 17, 34, 50, 66, 81, 98, 116, 132, 145, 166, 181, 201, 215, 237, 255, 273, 287, 306, 323, 339, 357, 374, 388, 406, 422, 439, 457, 471, 487, 504, 517, 540, 554, 567}

func (i Complete) String() string {
	if i >= Complete(len(_Complete_index)-1) {
		return "Complete(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Complete_name[_Complete_index[i]:_Complete_index[i+1]]
}
