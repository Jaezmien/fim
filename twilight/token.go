package twilight

import "git.jaezmien.com/Jaezmien/fim/twilight/token"
import "git.jaezmien.com/Jaezmien/fim/twilight/queue"

func Parse(source string) *queue.Queue[*token.Token] {
	var t *queue.Queue[*token.Token]
	t = createPartialTokens(source)
	t = mergePartialTokens(t)

	t = createTokens(t)
	t = mergeMultitokens(t)
	t = mergeLiterals(t)
	t = cleanTokens(t)

	return t
}
