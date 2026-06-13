package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validPeriods = map[string]bool{
	"all": true, "yearly": true, "monthly": true, "weekly": true, "daily": true,
}

func (a *App) topCmd() *cobra.Command {
	var period string
	cmd := &cobra.Command{
		Use:   "top",
		Short: "List top topics by period",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !validPeriods[period] {
				return codeError(exitUsage, fmt.Errorf("invalid period %q: must be one of all, yearly, monthly, weekly, daily", period))
			}
			a.progressf("fetching top topics (%s)...", period)
			topics, err := a.client.Top(cmd.Context(), period, 0)
			if err != nil {
				return mapFetchErr(err)
			}
			n := a.effectiveLimit(30)
			topics = applyLimit(topics, n)
			return a.renderOrEmpty(topics, len(topics))
		},
	}
	cmd.Flags().StringVar(&period, "period", "monthly", "time period: all|yearly|monthly|weekly|daily")
	return cmd
}
