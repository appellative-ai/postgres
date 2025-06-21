package exec

import (
	"github.com/behavioral-ai/core/json"
	"github.com/jackc/pgx/v5/pgconn"
)

func newCmdTag(tag pgconn.CommandTag) CommandTag {
	return CommandTag{
		Sql:          tag.String(),
		RowsAffected: tag.RowsAffected(),
		Insert:       tag.Insert(),
		Update:       tag.Update(),
		Delete:       tag.Delete(),
		Select:       tag.Select(),
	}
}

func NewCommandTag(url string) CommandTag {
	tag, status := json.New[CommandTag](url, nil)
	if status != nil {
		return CommandTag{}
	}
	return tag
}
