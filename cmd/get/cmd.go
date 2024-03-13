package get

import (
	"fmt"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/helaili/gh-ssh-cert-please/utils"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get a SSH Certificate from GitHub",
		RunE:  run,
	}
}

func run(cmd *cobra.Command, args []string) (err error) {
	// Read the public key file from the disk
	publicKey, pubKeyFileReadError := utils.ReadPublicKey("/Users/helaili/.ssh/github-ssh-authority/helaili.pub")
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
	deleteError := utils.DeleteCertificates("github-ssh-authority", "cert-broker", *user.Login)
	if deleteError != nil {
		return deleteError
	}

	// Request a certificate from GitHub using the public key
	requestCertificateError := utils.RequestCertificateCreation("github-ssh-authority", "cert-broker", "cert.yml", "helaili@github.com", publicKey)
	if requestCertificateError != nil {
		return requestCertificateError
	}
	fmt.Println("Certificate requested")

	// Retrieve the certificate from GitHub every 5 seconds
	var certificates []github.Artifact
	for i := 0; i < 60 && certificates == nil; i++ {
		time.Sleep(5 * time.Second)
		certificates, err = utils.RetrieveCertificates("github-ssh-authority", "cert-broker", *user.Login)
		if err != nil {
			return err
		}
	}

	if certificates == nil {
		return fmt.Errorf("no certificate found")
	}
	fmt.Println("Certificate generated successfully")

	downloadError := utils.DownloadCertificate(certificates[0])
	if downloadError != nil {
		return downloadError
	}
	fmt.Println("Certificate downloaded")
	return
}
