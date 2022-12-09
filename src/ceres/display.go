package ceres

import (
	"fmt"
)

func (erf*EntangledRecognitionForest) Display() {
	fmt.Printf("Displaying forest @%p\n", erf)
	for tree_idx, tree := range *erf {
		fmt.Printf("tree n°%v (out of %v)\n", tree_idx + 1, len(*erf))
		tree.Display()
	}
}

func (rt*RecognitionTree) Display() {
	t_lines := rt.display()
	for _, line := range t_lines {
		fmt.Println(line)
	}
}

type formattedTreeBranch []string

func mergeLines(a formattedTreeBranch, a_width int,
                b formattedTreeBranch, b_width int,
				sep string) (formattedTreeBranch, int){
    if a_width < 0 {a_width=a.maxWidth()+1}
    if b_width < 0 {b_width=b.maxWidth()+1}

	var offseter string = ""
	for i := 0; i < a_width; i++ {
		offseter += " "
	}

	var max_len int = len(a)
	var width int = a_width
	if len(b) > max_len {max_len = len(b)}
	//fmt.Println("we have a resulting line with length", max_len, "from", len(a), "and", len(b))

	var answer = make([]string, max_len)
	for i := 0; i < max_len; i++ {
		switch {
		case i >= len(a):
			answer[i] = offseter + sep + b[i]
		case i >= len(b):
			answer[i] = a[i]
			if sep != "" {
				answer[i] += offseter[len(a[i]):] + sep
			}
			continue
		default: // i < len(a) && i < len(b)
			answer[i] = a[i] + offseter[len(a):] + sep + b[i]
		}
		if width < len(answer[i])  {
			width = len(answer[i])
		}
	}

	return formattedTreeBranch(answer), width
}

func (ftb *formattedTreeBranch) maxWidth() int {
	var maxWidth int = 0

	for _, a := range *ftb {
		if len(a) > maxWidth {
			maxWidth = len(a)
		}
	}

	return maxWidth
}

func (rt *RecognitionTree) display() []string {
	var lines []string
	var offset int = 0

ROOT_DISPLAY_LOOP:
	for _, root := range rt.roots {
		sublines, sublines_maxwidth := root.display()
		for i, l := range lines { // update pre-existing lines
			if i >= len(sublines) {
				continue ROOT_DISPLAY_LOOP
			}
			missing_offset := offset - len(l)
			var offseter string = ""
			for j := 0; j < missing_offset; j++ {
				offseter += " "
			}

			lines[i] = lines[i] + offseter + sublines[i]
		}
		offseter := ""
		for j := 0; j < offset; j++ { // making the space on the left
			offseter += " "
		}
		for i := len(lines); i < len(sublines); i++ { // append new lines
			lines = append(lines, offseter+sublines[i])
		}
		offset += sublines_maxwidth + 1
	}

	return lines
}

func (ern*entangledRecognitionNode) display() (formattedTreeBranch, int) {
    var ftb formattedTreeBranch
	var offset int
    var offseter, header string = "", "("+ern.Content.s+")"

    ftb = append(ftb, header)
    for _, node := range ern.nodes() {
        lines, width := node.display()
        for i, l := range lines {
            if len(ftb) <= i + 1 { // ftb[i+1] <- lines[i]
                ftb = append(ftb, "")
            }
            ftb[i+1] = fmt.Sprintf("%s%s[%s", ftb[1+i],
				offseter[:len(ftb[1+i])], l)
        }
        offset += width + 1
		for len(offseter) < offset {
			offseter += " "
		}
    }
    if offset > len(ftb[0]) {
        ftb[0] = string(offseter[:(offset-len(ftb[0]))/2]) + ftb[0]
    }
	if offset < len(header) {offset = len(header)}
	/*fmt.Println("Preparing to display", ern.Content.s, "(width:",offset,
		",", len(ern.nodes()), "nodes treated)")//*/
	for i, line := range ftb {
		if i == 0 {continue}
		if len(line) < offset {
			ftb[i] = line + offseter[len(line):] + "]"
		} else {
			ftb[i] = line + "]"
		}
	}

    return ftb, offset + 1
}

func (rn*recognitionNode) display() (formattedTreeBranch, int) {
	if rn == nil {
		return nil, 0
	}

	var ftb formattedTreeBranch

    var loffseter, roffseter,
        header string = "", "",
        fmt.Sprintf("{%s}%v", rn.Surround.String(), len(rn.ChildMap))

	/*fmt.Println("Displaying node with", len(rn.ChildMap),
		"direct children. header :", header,
		"searching from", rn.Surround.minPos,
		"to", rn.Surround.maxPos)*/

    var demi_header_width int = len(header) / 2

	var left_subtree, right_subtree formattedTreeBranch

    ftb = append(ftb, header)
    for pos := rn.Surround.minPos; pos <= rn.Surround.maxPos; pos ++ {
        child, ok := rn.ChildMap[pos]
        if !ok {continue}

        lines, width := child.display()

		var l_width, r_width int

        switch {
        case pos == 0:
            continue
        case pos < 0:
			left_subtree, l_width = mergeLines(left_subtree, len(loffseter),
												lines, width, " ")
			/*fmt.Println("adding left", left_subtree.maxWidth(), l_width, "on", child.Content.s, ":", lines)*/
            for len(loffseter) < l_width {
				loffseter += " "
			}
        default:
			right_subtree, r_width = mergeLines(right_subtree, len(roffseter),
												 lines, width, " ")
            for len(roffseter) < r_width {
				roffseter += " "
			}
        }
    }

	MERGE_IN, total_width := mergeLines(left_subtree, len(loffseter),
										right_subtree, len(roffseter), "|")
	ftb = append(ftb, MERGE_IN...)

	if demi_header_width < total_width / 2 {
		ftb[0] = (loffseter+roffseter)[:total_width/2-demi_header_width] + ftb[0]
	} else {
		total_width = len(header)
	}

    return ftb, total_width
}
