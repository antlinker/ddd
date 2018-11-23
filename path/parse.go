package path

const (
	dstart  = '['
	dend    = ']'
	pathSpe = '.'
	modePre = iota
	modeDstart
	modeNameStart
	modeNameRuning
	modeNameEnd
	whitespace = (1 << '\t') | (1 << '\n') | (1 << '\r') | (1 << ' ')
)

func parseNamespaces(pathstr, startflag string, aliaslen int) (out []string, remainder string, ok bool) {
	l := len(pathstr)
	dl := len(startflag)
	mode := modePre

	dspos := 0
	var ch byte
	runflag := false
one:
	for i := 0; i < l; i++ {
		ch = pathstr[i]
		switch mode {
		case modePre:
			if i < dl && ch == startflag[i] {
				// 匹配$domain
				continue
			} else if i == aliaslen {
				// 匹配$d
				if ch == dstart {
					mode = modeNameStart
				} else if isWhite(ch) {
					mode = modeDstart
					continue
				} else {
					break one
				}
				continue
			} else if i == dl {
				if ch == dstart {
					mode = modeNameStart
					continue
				}
				mode = modeDstart
				continue
			} else if i < dl {

				// 匹配失败找不到领域名称
				break one

			} else {
				// 不应该运行到这里
				panic("算法有错误")
			}
		case modeDstart:
			if isWhite(ch) {
				continue
			}
			if ch == dstart {
				mode = modeNameStart
				continue
			}
			// 运行到这里解析失败
			return
		case modeNameStart:

			if isWhite(ch) {
				continue
			}
			if ch == dend {
				// 没有指定域可以认为是当前域
				mode = modeNameEnd
				continue
			}
			mode = modeNameRuning
			dspos = i
		case modeNameRuning:
			if runflag {
				if isWhite(ch) {
					dspos = i + 1
					continue
				}
				if ch == pathSpe {
					runflag = false
					dspos = i + 1
					continue
				} else if ch == dend {
					mode = modeNameEnd
					continue
				} else {
					// 分析出错
					return
				}
			}
			if isWhite(ch) {
				if dspos == i {
					dspos = i + 1
					continue
				}
				o := pathstr[dspos:i]
				// if isWhite(o[0]) {
				// 	dspos = i
				// 	continue
				// }
				out = append(out, o)
				runflag = true
				dspos = i
			} else if ch == pathSpe {
				o := pathstr[dspos:i]
				if !isWhite(o[0]) {
					out = append(out, o)
				}
				dspos = i + 1
			} else if ch == dend {
				o := pathstr[dspos:i]
				if !isWhite(o[0]) {
					out = append(out, o)
				}
				mode = modeNameEnd
			}

		case modeNameEnd:
			if isWhite(ch) {
				continue
			}
			if ch == pathSpe {
				// 获取剩余字符串
				remainder = pathstr[i+1:]
				ok = true
				break one
			}
		default:
			// 此处不进行任何处理不应该进入这个分支
			panic("进入这里算法有错误")
		}

	}
	if mode == modeNameEnd {
		ok = true
	} else if mode == modePre {
		ok = true
		remainder = pathstr
	}

	switch mode {
	case modeNameEnd:
		ok = true
	case modePre:
		ok = true
		remainder = pathstr
	case modeDstart:
		return
	}
	return
}
func isWhite(ch uint8) bool {
	return (whitespace & (1 << ch)) != 0
}
