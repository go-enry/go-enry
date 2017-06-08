package slinguist

import (
	"math"

	"gopkg.in/src-d/simple-linguist.v1/internal/tokenizer"
)

// Classifier is the interface that contains the method Classify which is in charge to assign scores to the possibles candidates.
// The scores must order the candidates so as the highest score be the most probably language of the content. The candidates is
// a map which can be used to assign weights to languages dynamically.
type Classifier interface {
	Classify(content []byte, candidates map[string]float64) map[string]float64
}

type classifier struct {
	languagesLogProbabilities map[string]float64
	tokensLogProbabilities    map[string]map[string]float64
	tokensTotal               float64
}

func (c *classifier) Classify(content []byte, candidates map[string]float64) map[string]float64 {
	if len(content) == 0 {
		return nil
	}

	var languages map[string]float64
	if len(candidates) == 0 {
		languages = c.knownLangs()
	} else {
		languages = make(map[string]float64, len(candidates))
		for candidate, weight := range candidates {
			if lang, ok := GetLanguageByAlias(candidate); ok {
				languages[lang] = weight
			}
		}
	}

	tokens := tokenizer.Tokenize(content)
	scores := make(map[string]float64, len(languages))
	for language := range languages {
		scores[language] = c.tokensLogProbability(tokens, language) + c.languagesLogProbabilities[language]
	}

	return scores
}

func (c *classifier) knownLangs() map[string]float64 {
	langs := make(map[string]float64, len(c.languagesLogProbabilities))
	for lang := range c.languagesLogProbabilities {
		langs[lang]++
	}

	return langs
}

func (c *classifier) tokensLogProbability(tokens []string, language string) float64 {
	var sum float64
	for _, token := range tokens {
		sum += c.tokenProbability(token, language)
	}

	return sum
}

func (c *classifier) tokenProbability(token, language string) float64 {
	tokenProb, ok := c.tokensLogProbabilities[language][token]
	if !ok {
		tokenProb = math.Log(1.000000 / c.tokensTotal)
	}

	return tokenProb
}
