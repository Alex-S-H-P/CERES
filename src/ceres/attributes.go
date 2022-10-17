package ceres


type AttributeTypeList struct {
    attrs []*EntityType
}

type AttributeInstanceList struct {
    attrs []*EntityType
    values map[EntityType]Word
}




func (atl *AttributeTypeList) parentType(atl2 *AttributeTypeList) {
    COPY_ATTR_LOOP:
    for _, attr := range atl.attrs {
        // look if one element matches.
        for _, attr2 := range atl2.attrs {
            if IsTypeOf(attr2, attr) {
                continue COPY_ATTR_LOOP
            }
        }
        // No element match
        atl2.attrs = append(atl2.attrs, attr)
    }
}



func (atl *AttributeTypeList) parentInstance(atl2 *AttributeInstanceList) {
    COPY_ATTR_LOOP:
    for _, attr := range atl.attrs {
        // look if one element matches.
        for _, attr2 := range atl2.attrs {
            if IsTypeOf(attr2, attr) {
                continue COPY_ATTR_LOOP
            }
        }
        // No element match
        atl2.attrs = append(atl2.attrs, attr)
    }
}
