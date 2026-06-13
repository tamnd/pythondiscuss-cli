package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func (a *App) topicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topic <id>",
		Short: "Fetch a topic and its posts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return codeError(exitUsage, fmt.Errorf("topic id must be an integer: %w", err))
			}
			a.progressf("fetching topic %d...", id)
			posts, err := a.client.GetTopic(cmd.Context(), id)
			if err != nil {
				return mapFetchErr(err)
			}
			n := a.effectiveLimit(20)
			posts = applyLimit(posts, n)
			return a.renderOrEmpty(posts, len(posts))
		},
	}
	return cmd
}
