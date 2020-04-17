package process

import (
	"context"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
)

func ProcessContent(ctx context.Context, fs *types.FileStore, db database.Database) {

	// Tika
	//   Split Sentences
	//     > ContentLines
	//     skyset
	//       aggregate data > NLP Tags
	//   Split Words > ContentTags
	// Gotenburg > View
}
