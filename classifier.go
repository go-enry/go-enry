package slinguist

import (
	"math"

	"gopkg.in/src-d/simple-linguist.v1/internal/tokenizer"
)

func getLanguageByClassifier(content []byte, candidates []string, classifier Classifier) string {
	if classifier == nil {
		classifier = DefaultClassifier
	}

	scores := classifier.Classify(content, candidates)
	if len(scores) == 0 {
		return OtherLanguage
	}

	return getLangugeHigherScore(scores)
}

func getLangugeHigherScore(scores map[string]float64) string {
	var language string
	higher := -math.MaxFloat64
	for lang, score := range scores {
		if higher < score {
			language = lang
			higher = score
		}
	}

	return language
}

// Classifier is the interface that contains the method Classify which is in charge to assign scores to the possibles candidates.
// The scores must order the candidates so as the highest score be the most probably language of the content.
type Classifier interface {
	Classify(content []byte, candidates []string) map[string]float64
}

type classifier struct {
	languagesLogProbabilities map[string]float64
	tokensLogProbabilities    map[string]map[string]float64
	tokensTotal               float64
}

func (c *classifier) Classify(content []byte, candidates []string) map[string]float64 {
	if len(content) == 0 {
		return nil
	}

	var languages []string
	if len(candidates) == 0 {
		languages = c.knownLangs()
	} else {
		languages = make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			if lang, ok := GetLanguageByAlias(candidate); ok {
				languages = append(languages, lang)
			}
		}
	}

	tokens := tokenizer.Tokenize(content)
	scores := make(map[string]float64, len(languages))
	for _, language := range languages {
		scores[language] = c.tokensLogProbability(tokens, language) + c.languagesLogProbabilities[language]
	}

	return scores
}

func (c *classifier) knownLangs() []string {
	langs := make([]string, 0, len(c.languagesLogProbabilities))
	for lang := range c.languagesLogProbabilities {
		langs = append(langs, lang)
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
