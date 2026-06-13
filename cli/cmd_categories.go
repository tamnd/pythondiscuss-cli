package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) categoriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "List forum categories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			a.progressf("fetching categories...")
			cats, err := a.client.Categories(cmd.Context())
			if err != nil {
				return mapFetchErr(err)
			}
			n := a.effectiveLimit(0)
			cats = applyLimit(cats, n)
			return a.renderOrEmpty(cats, len(cats))
		},
	}
}
