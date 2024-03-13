package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/helaili/gh-ssh-cert-please/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "ssh-cert-please",
		Short:             "Generates SSH certificates for GitHub",
		RunE:              run,
		PersistentPreRunE: preRun,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		Version:           buildVersion(version, commit),
	}

	flags(cmd)
	return cmd
}

func flags(cmd *cobra.Command) {
	cmd.Flags().StringP("org", "o", "", "Organization to use as a certificate authority")
	cmd.MarkFlagRequired("org")
	cmd.Flags().StringP("repo", "r", "", "Repo to use as a certificate authority")
	cmd.MarkFlagRequired("repo")
	cmd.Flags().StringP("pubKey", "k", "", "Path to the public key file")
	cmd.MarkFlagRequired("pubKey")
	cmd.Flags().StringP("email", "m", "", "Your email address")
	cmd.MarkFlagRequired("email")

	viper.BindPFlag("org", cmd.Flags().Lookup("org"))
	viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))
	viper.BindPFlag("pubKey", cmd.Flags().Lookup("pubKey"))
	viper.BindPFlag("email", cmd.Flags().Lookup("email"))
}

func preRun(cmd *cobra.Command, args []string) error {
	// cmd.SilenceUsage = true
	return nil
}

func buildVersion(version, commit string) string {
	if commit != "" {
		version += " (" + commit + ")"
	}
	return version
}

func run(cmd *cobra.Command, args []string) (err error) {
	// Read the public key file from the disk
	publicKey, pubKeyFileReadError := utils.ReadPublicKey(viper.GetString("pubKey"))
	if pubKeyFileReadError != nil {
		return pubKeyFileReadError
	}

	// Get the current user's login
	user, userError := utils.GetUSer()
	if userError != nil {
		return userError
	}

	fmt.Println("Cleaning up previous certificates")
	// Remove the existing certificate if any
	deleteError := utils.DeleteCertificates(viper.GetString("org"), viper.GetString("repo"), *user.Login)
	if deleteError != nil {
		return deleteError
	}

	// Request a certificate from GitHub using the public key
	requestCertificateError := utils.RequestCertificateCreation(viper.GetString("org"), viper.GetString("repo"), "cert.yml", viper.GetString("email"), publicKey)
	if requestCertificateError != nil {
		return requestCertificateError
	}
	fmt.Println("Certificate requested")

	// Retrieve the certificate from GitHub every 5 seconds
	var certificates []github.Artifact
	for i := 0; i < 60 && certificates == nil; i++ {
		time.Sleep(5 * time.Second)
		certificates, err = utils.RetrieveCertificates(viper.GetString("org"), viper.GetString("repo"), *user.Login)
		if err != nil {
			return err
		}
	}

	if certificates == nil {
		return fmt.Errorf("no certificate found")
	}
	fmt.Println("Certificate generated successfully")

	// Get the path from the public key file
	certFilePath := filepath.Dir(viper.GetString("pubKey"))

	fmt.Println("Downloading certificate")
	downloadError := utils.DownloadCertificate(certificates[0], certFilePath)
	if downloadError != nil {
		return downloadError
	}
	fmt.Println("Certificate downloaded")
	return
}
