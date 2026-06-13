package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) searchCmd() *cobra.Command {
	var page int
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search topics and posts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			a.progressf("searching for %q (page %d)...", query, page)
			topics, err := a.client.Search(cmd.Context(), query, page)
			if err != nil {
				return mapFetchErr(err)
			}
			n := a.effectiveLimit(20)
			topics = applyLimit(topics, n)
			return a.renderOrEmpty(topics, len(topics))
		},
	}
	cmd.Flags().IntVar(&page, "page", 1, "page number (1-indexed)")
	return cmd
}
