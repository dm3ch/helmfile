package cmd

import (
	"github.com/spf13/cobra"

	"github.com/helmfile/helmfile/pkg/app"
	"github.com/helmfile/helmfile/pkg/config"
	"github.com/helmfile/helmfile/pkg/runtime"
)

// NewApplyCmd returns apply subcmd
func NewApplyCmd(globalCfg *config.GlobalImpl) *cobra.Command {
	applyOptions := &config.ApplyOptions{}

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply all resources from state file only when there are changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			applyImpl := config.NewApplyImpl(globalCfg, applyOptions)

			err := config.NewCLIConfigImpl(applyImpl.GlobalImpl)
			if err != nil {
				return err
			}

			if err := applyImpl.ValidateConfig(); err != nil {
				return err
			}

			a := app.New(applyImpl)
			return toCLIError(applyImpl.GlobalImpl, a.Apply(applyImpl))
		},
	}

	f := cmd.Flags()
	f.StringArrayVar(&applyOptions.Set, "set", nil, "additional values to be merged into the helm command --set flag")
	f.StringArrayVar(&applyOptions.Values, "values", nil, "additional value files to be merged into the helm command --values flag")
	f.IntVar(&applyOptions.Concurrency, "concurrency", 0, "maximum number of concurrent helm processes to run, 0 is unlimited")
	f.BoolVar(&applyOptions.Validate, "validate", false, "validate your manifests against the Kubernetes cluster you are currently pointing at. Note that this requires access to a Kubernetes cluster to obtain information necessary for validating, like the list of available API versions")
	f.IntVar(&applyOptions.Context, "context", 0, "output NUM lines of context around changes")
	f.StringVar(&applyOptions.Output, "output", "", "output format for diff plugin")
	f.BoolVar(&applyOptions.DetailedExitcode, "detailed-exitcode", false, "return a non-zero exit code 2 instead of 0 when there were changes detected AND the changes are synced successfully")
	f.StringVar(&globalCfg.GlobalOptions.Args, "args", "", "pass args to helm exec")
	if !runtime.V1Mode {
		// TODO: Remove this function once Helmfile v0.x
		f.BoolVar(&applyOptions.RetainValuesFiles, "retain-values-files", false, "DEPRECATED: Use skip-cleanup instead")
		_ = f.MarkDeprecated("retain-values-files", "Use skip-cleanup instead")
	}

	f.BoolVar(&applyOptions.SkipCleanup, "skip-cleanup", false, "Stop cleaning up temporary values generated by helmfile and helm-secrets. Useful for debugging. Don't use in production for security")
	f.BoolVar(&applyOptions.SkipCRDs, "skip-crds", false, "if set, no CRDs will be installed on sync. By default, CRDs are installed if not already present")
	f.BoolVar(&applyOptions.SkipNeeds, "skip-needs", true, `do not automatically include releases from the target release's "needs" when --selector/-l flag is provided. Does nothing when --selector/-l flag is not provided. Defaults to true when --include-needs or --include-transitive-needs is not provided`)
	f.BoolVar(&applyOptions.IncludeNeeds, "include-needs", false, `automatically include releases from the target release's "needs" when --selector/-l flag is provided. Does nothing when --selector/-l flag is not provided`)
	f.BoolVar(&applyOptions.IncludeTransitiveNeeds, "include-transitive-needs", false, `like --include-needs, but also includes transitive needs (needs of needs). Does nothing when --selector/-l flag is not provided. Overrides exclusions of other selectors and conditions.`)
	f.BoolVar(&applyOptions.SkipDiffOnInstall, "skip-diff-on-install", false, "Skips running helm-diff on releases being newly installed on this apply. Useful when the release manifests are too huge to be reviewed, or it's too time-consuming to diff at all")
	f.BoolVar(&applyOptions.IncludeTests, "include-tests", false, "enable the diffing of the helm test hooks")
	f.StringArrayVar(&applyOptions.Suppress, "suppress", nil, "suppress specified Kubernetes objects in the diff output. Can be provided multiple times. For example: --suppress KeycloakClient --suppress VaultSecret")
	f.BoolVar(&applyOptions.SuppressSecrets, "suppress-secrets", false, "suppress secrets in the diff output. highly recommended to specify on CI/CD use-cases")
	f.BoolVar(&applyOptions.ShowSecrets, "show-secrets", false, "do not redact secret values in the diff output. should be used for debug purpose only")
	f.BoolVar(&applyOptions.NoHooks, "no-hooks", false, "do not diff changes made by hooks.")
	f.BoolVar(&applyOptions.SuppressDiff, "suppress-diff", false, "suppress diff in the output. Usable in new installs")
	f.BoolVar(&applyOptions.Wait, "wait", false, `Override helmDefaults.wait setting "helm upgrade --install --wait"`)
	f.BoolVar(&applyOptions.WaitForJobs, "wait-for-jobs", false, `Override helmDefaults.waitForJobs setting "helm upgrade --install --wait-for-jobs"`)
	f.BoolVar(&applyOptions.ReuseValues, "reuse-values", false, `Override helmDefaults.reuseValues "helm upgrade --install --reuse-values"`)
	f.BoolVar(&applyOptions.ResetValues, "reset-values", false, `Override helmDefaults.reuseValues "helm upgrade --install --reset-values"`)
	f.StringVar(&applyOptions.PostRenderer, "post-renderer", "", `pass --post-renderer to "helm template" or "helm upgrade --install"`)

	return cmd
}
