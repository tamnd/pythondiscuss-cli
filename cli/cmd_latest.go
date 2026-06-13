package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) latestCmd() *cobra.Command {
	var page int
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "List latest topics",
		RunE: func(cmd *cobra.Command, _ []string) error {
			a.progressf("fetching latest topics (page %d)...", page)
			topics, err := a.client.Latest(cmd.Context(), page)
			if err != nil {
				return mapFetchErr(err)
			}
			n := a.effectiveLimit(30)
			topics = applyLimit(topics, n)
			return a.renderOrEmpty(topics, len(topics))
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "page number (0-indexed)")
	return cmd
}
