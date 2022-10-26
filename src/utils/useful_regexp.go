package utils

const or = "|"

const OptionalSpace = ` \t]?`
const OrPattern = `|`
const IntNumberPattern = `\d+([ ._,]\d+)*`
const NumberPattern = IntNumberPattern + `[\.,]?+\d*`
const MoneyPattern = `[₠₢₡₣₤₥₦₧₨₩₪₫€₹₭₸₮$₯₰₷₶₱₲₳₴₵]`
const PricePattern = MoneyPattern + OptionalSpace + NumberPattern +
                        OrPattern + NumberPattern + OptionalSpace + MoneyPattern

const WordPattern = "[a-zA-Z]"

const UnknownPattern = `[^\t \n\r.?;:/\\*-+]`
const PonctuationPattern = `[.,;:?!]`
const EOSPossiblePattern = `[.?!]|$`
