package _vocabulary

import "strings"

func EnglishStopWordChecker(s string) bool {
	return stopWordsEn[strings.ToLower(s)]
}

// English stop words
var stopWordsEn = map[string]bool{
	"a":       true,
	"about":   true,
	"above":   true,
	"after":   true,
	"again":   true,
	"against": true,
	"all":     true,
	"am":      true,
	"an":      true,
	"and":     true,
	"any":     true,
	"are":     true,
	"arent":   true,
	"as":      true,
	"at":      true,
	"be":      true,
	"because": true,
	"been":    true,
	"before":  true,
	"being":   true,
	"below":   true,
	"between": true,
	"both":    true,
	"but":     true,
	"by":      true,
	"can":     true,
	"cant":    true,
	"cannot":  true,
	"could":   true,
	"couldnt": true,
	"did":     true,
	"didnt":   true,
	"do":      true,
	"does":    true,
	"doesnt":  true,
	"doing":   true,
	"dont":    true,
	"down":    true,
	"during":  true,
	"each":    true,
	"few":     true,
	"for":     true,
	"from":    true,
	"further": true,
	"had":     true,
	"hadnt":   true,
	"has":     true,
	"hasnt":   true,
	"have":    true,
	"havent":  true,
	"having":  true,
	"he":      true,
	"hed":     true,
	"hell":    true,
	"hes":     true,
	"her":     true,
	"here":    true,
	"heres":   true,
	"hers":    true,
	"herself": true,
	"him":     true,
	"himself": true,
	"his":     true,
	"how":     true,
	"hows":    true,
	"i":       true,
	"id":      true,
	"ill":     true,
	"im":      true,
	"ive":     true,
	"if":      true,
	"in":      true,
	"into":    true,
	"is":      true,
	"isnt":    true,
	"it":      true,
	"its":     true,
	"itself":  true,
	"lets":    true,
	"me":      true,
	"more":    true,
	"most":    true,
	"mustnt":  true,
	"my":      true,
	"myself":  true,
	"no":      true,
	"nor":     true,
	"not":     true,
	"of":      true,
	"off":     true,
	"on":      true,
	"once":    true,
	"only":    true,
	"or":      true,
	"other":   true,
	"ought":   true,
	"our":     true,
	"ours	ourselves": true,
	"out":        true,
	"over":       true,
	"own":        true,
	"same":       true,
	"shant":      true,
	"she":        true,
	"shed":       true,
	"shell":      true,
	"shes":       true,
	"should":     true,
	"shouldnt":   true,
	"so":         true,
	"some":       true,
	"such":       true,
	"than":       true,
	"that":       true,
	"thats":      true,
	"the":        true,
	"their":      true,
	"theirs":     true,
	"them":       true,
	"themselves": true,
	"then":       true,
	"there":      true,
	"theres":     true,
	"these":      true,
	"they":       true,
	"theyd":      true,
	"theyll":     true,
	"theyre":     true,
	"theyve":     true,
	"this":       true,
	"those":      true,
	"through":    true,
	"to":         true,
	"too":        true,
	"under":      true,
	"until":      true,
	"up":         true,
	"very":       true,
	"was":        true,
	"wasnt":      true,
	"we":         true,
	"wed":        true,
	"well":       true,
	"were":       true,
	"weve":       true,
	"werent":     true,
	"what":       true,
	"whats":      true,
	"when":       true,
	"whens":      true,
	"where":      true,
	"wheres":     true,
	"which":      true,
	"while":      true,
	"who":        true,
	"whos":       true,
	"whom":       true,
	"why":        true,
	"whys":       true,
	"with":       true,
	"wont":       true,
	"would":      true,
	"wouldnt":    true,
	"you":        true,
	"youd":       true,
	"youll":      true,
	"youre":      true,
	"youve":      true,
	"your":       true,
	"yours":      true,
	"yourself":   true,
	"yourselves": true,
}
