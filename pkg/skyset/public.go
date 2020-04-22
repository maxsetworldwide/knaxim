package skyset

//BuildPhrases takes a sentence as a string and returns the skyset phrases that make up that sentence
func BuildPhrases(s string) []Phrase {
	return contextualize(assemble(tokenize(s)))
}
