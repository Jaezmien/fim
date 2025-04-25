package twilight

import (
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

// Parses the source string into a queue of tokens
func Parse(source string) []*token.Token {
	var t *queue.Queue[*token.Token]
	t = createPartialTokens(source)
	t = mergePartialTokens(t)

	t = createTokens(t)
	t = mergeMultiTokens(t)
	t = smartIdentifierTokens(t)
	t = mergeIdentifiers(t)
	t = cleanTokens(t)

	return t.Flatten()
}
