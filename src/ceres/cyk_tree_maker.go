package ceres

import "regexp"
import "fmt"

var re_rule = regexp.MustCompile(`(?m)(.*) -> ([a-zA-Zà-ùÀ-Ù0-9|]*) ?([a-zA-Zà-ùÀ-Ù0-9|]*)`)

type group_info struct {
    group group
    k int
    assignment_idx int_pair
}

type group struct {
    name string
    // instanceSolver is the lowest EntityType such that all ideas
    // represented in that group are children of that EntityType
    instanceSolver *EntityType
}

func MakeGroup(s string, g*grammar)group {
    return g.find(s)
}

func (g group)String() string{
    return g.name
}

type int_pair struct {
    a, b int
}

type cyk_case struct {
    assignments []group
    assignment_source []int
    source_assignment []int_pair // the first one for the left, the second one for down

    debug_pointed string
}

func (cc*cyk_case)join(assignment group_info){
    cc.assignments = append(cc.assignments, assignment.group)
    cc.assignment_source = append(cc.assignment_source, assignment.k)
    cc.source_assignment = append(cc.source_assignment, assignment.assignment_idx)
}

type cyk_table [][]cyk_case

func (cykt cyk_table) rslt() *cyk_case{
    if len(cykt) == 0 {return nil}
    return &(cykt[0][len(cykt)-1])
}

func MakeCYKTable(n int) cyk_table{
    table := make([][]cyk_case, n)
    for i := 0; i < n; i++ {
        table[i] = make([]cyk_case, n)
        for j := i; j < n; j++ {
            table[i][j].assignments = make([]group, 0, 3)
            table[i][j].assignment_source = make([]int, 0, 3)
        }
    }

    return cyk_table(table)
}


func CYK_PARSE(sentence []RecognizedEntity, grammar grammar) cyk_table {
    table := MakeCYKTable(len(sentence))
    table.display()

    for j:=0; j<len(sentence);j++{
        for _, assignment := range grammar.singleAssign(sentence[j]) {
            table[j][j].join(group_info{group:assignment, k:0})
            //fmt.Printf("joined \"%s\" onto \"%v\"\n", assignment, table[j][j].assignments)
        }
        for i := j-1; i >=0; i-- {
            for k := 1; k <= j-i; k ++ {
                fmt.Println("looking", k, "cases to the left")
                fmt.Println(i, j-k, j-k-1, j)
                left, down := &table[i][j-k], &table[j-k+1][j]

                /*table[i][j].debug_pointed = "+"
                left.debug_pointed = "<"
                down.debug_pointed = "v" //*/
                table.display()
                for _, assignment := range grammar.assign(*left, *down){
                    assignment.k = k
                    table[i][j].join(assignment)
                }
                /*
                table[i][j].debug_pointed = ""
                left.debug_pointed = ""
                down.debug_pointed = ""
                fmt.Println("========================")//*/
            }
        }
        table.display()
    }


    return table
}

func (table cyk_table) display() {
    for line := range table{
        fmt.Println(table[line])
    }
}

func (table cyk_table) display_tree(i int) error{
    if len(table.rslt().assignments) == 0 {
        return fmt.Errorf("No solution found")
    }
    if i > len(table.rslt().assignments) || i < 0{
        return fmt.Errorf("Cannot use solution n°%v. Please only use between n°0 and n°%v",
                    i, len(table.rslt().assignments)-1)
    }

    _, lines := table.display_node(0, len(table)-1, i)

    for _, line := range lines {
        fmt.Println(line)
    }


    return nil
}

func (table cyk_table) display_node(row, col, i int) (int, []string){
    var loffsetter string

    node := &(table[row][col])
    k := node.assignment_source[i]

    if k != 0 {

        left, down := col-k, col-k+1
        left_i := node.source_assignment[i].a
        down_i := node.source_assignment[i].b

        width, llines := table.display_node(row, left, left_i)
        sided, rlines := table.display_node(down, col, down_i)

        for i := 0; i < width + 1; i++ {
            loffsetter += " "
        }

        nbLines := len(llines)
        if nbLines < len(rlines) {nbLines = len(rlines)}
        const nbHeaderLines int = 2
        var returnedLines = make([]string, nbLines+nbHeaderLines)
        var lines []string = returnedLines[nbHeaderLines:]
        var header = returnedLines[:nbHeaderLines]

        for i := range lines {
            switch {
            case i >= len(llines):
                lines[i] = loffsetter + rlines[i]
            case i >= len(rlines):
                lines[i] = llines[i]
            default:
                lines[i] = llines[i] + loffsetter[len(llines[i]):] + rlines[i]
            }
        }
        header[0] = fmt.Sprintf("(%v)", node.assignments[i])
        if len(header[0]) < width {
            header[0] = loffsetter[:width-len(header[0])/2 + 1] + header[0]
            header[1] = loffsetter[:width-len(header[0])/2] + "/" +
                        loffsetter[:len(header)] + "\\"
        } else {
            header[1] = "|" + loffsetter[:width-1] + "\\"
        }
        return width + sided + 1, returnedLines
    }
    lines := []string{fmt.Sprintf("(%v)", node.assignments[i])}
    return  len(lines[0]), lines

}
