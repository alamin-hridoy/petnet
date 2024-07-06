package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/spf13/viper"
)

func main() {
	Execute()
}

var n int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pgunk",
	Short: "Generate gunk packages in parallel",
	Long:  `Generate gRPC from gunk packages in parallel.`,
	RunE:  pexec,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// No configuation available.
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().
	// StringVar(&cfgFile, "config", "", "config file (default is $PWD/.gunkconfig)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().IntVarP(&n, "num-threads", "n", 4, "Number of threads to use in generation.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find current directory.
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Search config in home directory with name ".pgunk" (without extension).
	viper.AddConfigPath(pwd)
	viper.SetConfigName(".gunkconfig")
	viper.SetConfigType("ini")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func pexec(c *cobra.Command, args []string) error {
	if n <= 0 {
		n = 4
	}
	ch := make(chan *exec.Cmd, len(args))
	ctx := c.Context()
	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < n; i++ {
		eg.Go(func() error {
			for cmd := range ch {
				fmt.Println("Generating", cmd.Dir)
				if out, err := cmd.CombinedOutput(); err != nil {
					fmt.Println(err, string(out))
					return err
				}
				fmt.Println("Done", cmd.Dir)
			}
			return nil
		})
	}
	for _, d := range args {
		cmd := exec.CommandContext(ctx, "gunk", "generate")
		cmd.Dir = d
		ch <- cmd
	}
	close(ch)
	return eg.Wait()
}
