package ceres


import (
    "regexp"
    "fmt"
)

type ruleString string

func (r ruleString) matches(arg1, arg2 string) (string, bool) {
    var cond1, cond2 string
    match := re_rule.FindStringSubmatch(string(r))
    if len(match) < 4 {
        panic("this rule cannot be processed")
    }

    cond1, cond2 = match[2], match[3]

    match1 := regexp.MustCompile(cond1).FindString(arg1)
    match2 := regexp.MustCompile(cond2).FindString(arg2)
    if len(match1) == 0 || (len(match2) == 0 && !(len(arg2) == 0 && len(cond2) == 0)){
        return "", false
    }
    fmt.Printf("ruleString :\"%s\", \"%s\", \"%s\" <- \"%s\". Matched \"%s\" & \"%s\"\n",
        match[1], match[2], match[3], match[0], match1, match2)

    return match[1], true
}

func (r ruleString) String() string {return string(r)}

type rule interface {
    matches(string, string) (string, bool)
    String() string
}

type grammar struct {
    rules []rule

    groups map[string]group
}


func (g*grammar) singleAssign(element RecognizedEntity) []group {
    return g.internalassign(MakeGroup(element.s, g), group{})
}

func (g*grammar)assign(left, down cyk_case) []group_info{
    slice := make([]group_info, 0, len(left.assignments)*len(down.assignments))
    for i, left_group := range left.assignments {
        for j, down_group := range down.assignments {
            idx := int_pair{a:i, b:j}
            for _, assignment := range g.internalassign(left_group, down_group){
                ginfo := group_info{group:assignment, assignment_idx:idx}
                slice = append(slice, ginfo)
            }
        }
    }
    return slice
}

func (g*grammar) find(name string) group {
    if gr, ok := g.groups[name]; ok {
        return gr
    } else {
        gr := group{name:name}
        g.groups[name] = gr
        return gr
    }
}

func (g*grammar)internalassign(a, b group) []group {
    slice := make([]group, 0, 3)
    for _, rule := range g.rules {
        if group, ok := rule.matches(a.String(), b.String()); ok && len(group) > 0{
            slice = append(slice, g.find(group))
        }
    }
    return slice
}

func (g*grammar) RefreshAllGroups(ics*ICS) {
    gen := ics.allEntities()
    for {
        e,done := gen.Next()
        if done {return}

        if et, ok := e.(*EntityType);ok {
            if et.grammar_group.instanceSolver == nil {
                et.grammar_group.instanceSolver = et
            } else if IsTypeOf(e, et.grammar_group.instanceSolver) {
                ancestor := ClosestAncestor(e, et.grammar_group.instanceSolver)
                et.grammar_group.instanceSolver = ancestor
            }
        }
    }

}
