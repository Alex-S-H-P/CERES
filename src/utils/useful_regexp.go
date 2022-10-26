package utils


const OptionalSpace = `[ \t]?`

const OrPattern = `|`

const IntNumberPattern = `\d+([ ._,]\d+)*`
const NumberPattern = IntNumberPattern + `[\.,]?\d*`
const CurrencyPattern = `[₠₢₡₣₤₥₦₧₨₩₪₫€₹₭₸₮$₯₰₷₶₱₲₳₴₵]`
const PricePattern = CurrencyPattern + OptionalSpace + NumberPattern +
	OrPattern + NumberPattern + OptionalSpace + CurrencyPattern

const WordPattern = `[a-zà-üA-ZÀ-Ü]+\'?`

const UnknownPattern = `[^\t \n\r.?;:/\\*-+]`
const PonctuationPattern = `[.,;:?!]`
const EOSPossiblePattern = `[.?!]|$`

const TokenPattern = WordPattern + OrPattern + NumberPattern +
	   OrPattern + PricePattern + OrPattern + UnknownPattern
