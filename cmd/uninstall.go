package cmd

import (
	"fmt"
	"path"

        "github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/whalebrew/whalebrew/hooks"
	"github.com/whalebrew/whalebrew/packages"
)

var forceUninstall bool

func init() {
	uninstallCommand.Flags().BoolVarP(&assumeYes, "assume-yes", "y", false, "Assume 'yes' as answer to all prompts and run non-interactively. Defaults to false.")

	RootCmd.AddCommand(uninstallCommand)
}

var uninstallCommand = &cobra.Command{
	Use:   "uninstall PACKAGENAME",
	Short: "Uninstall a package",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}
		if len(args) > 1 {
			return fmt.Errorf("Only one image can be uninstalled at a time")
		}

		pm := packages.NewPackageManager(viper.GetString("install_path"))
		packageName := args[0]

		path := path.Join(pm.InstallPath, packageName)

		if err := hooks.Run("pre-uninstall", packageName); err != nil {
			return fmt.Errorf("pre-uninstall install script failed: %s", err.Error())
		}
		
                if !assumeYes {
                	if !prompter.YN(fmt.Sprintf("This will permanently delete '%s'. Are you sure?", path), false) {
				return nil
			}
		}

		err := pm.Uninstall(packageName)
		if err != nil {
			return err
		}

		if err := hooks.Run("post-uninstall", packageName); err != nil {
			return fmt.Errorf("post-uninstall install script failed: %s", err.Error())
		}
		fmt.Printf("🚽  Uninstalled %s\n", path)

		return nil
	},
}
