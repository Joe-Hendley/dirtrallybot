package emoji

type Classification int

const (
	NotSet Classification = iota
	Negative
	Positive
)

func Classify(emoji string) Classification {
	switch emoji {
	case "👍":
		return Positive
	case "👎":
		return Negative
	}

	return NotSet
}
