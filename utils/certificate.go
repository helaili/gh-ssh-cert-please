package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cli/go-gh/v2"
	"github.com/google/go-github/v60/github"
)

func RequestCertificateCreation(org string, repo string, workflow string, email string, key string) error {
	repoWithOwner := fmt.Sprintf("%s/%s", org, repo)
	keyParameter := fmt.Sprintf("key=%s", key)
	emailParameter := fmt.Sprintf("email=%s", email)
	_, _, err := gh.Exec("workflow", "run", workflow, "--repo", repoWithOwner, "-f", emailParameter, "-f", keyParameter)

	if err != nil {
		return err
	}
	return nil
}

func RetrieveCertificates(org string, repo string, login string) ([]github.Artifact, error) {
	artifactName := fmt.Sprintf("name=%s-%s-cert.pub", login, org)
	endpoint := fmt.Sprintf("repos/%s/%s/actions/artifacts", org, repo)

	fmt.Println("Retrieving certificates")
	artefacts, stderr, err := gh.Exec("api", endpoint, "--field", artifactName, "--method", "GET")
	if err != nil {
		fmt.Printf("error: %s, %s", err, stderr.String())
		return nil, err
	}
	fmt.Println(artefacts.String())

	// Unmarshal the response into the artifactResponse variable
	artifactResponse := struct {
		TotalCount int `json:"total_count"`
		Artifacts  []github.Artifact
	}{}

	err = json.Unmarshal([]byte(artefacts.Bytes()), &artifactResponse)
	if err != nil {
		return nil, err
	}

	// If no artifact is found, return nil
	if len(artifactResponse.Artifacts) == 0 {
		return nil, nil
	}

	return artifactResponse.Artifacts, nil
}

func DeleteCertificate(org string, repo string, id int64) error {
	endpoint := fmt.Sprintf("repos/%s/%s/actions/artifacts/%d", org, repo, id)

	_, stderr, err := gh.Exec("api", endpoint, "--method", "DELETE")
	if err != nil {
		return fmt.Errorf("error: %s, %s", err, stderr.String())
	}

	return nil
}

func DeleteCertificates(org string, repo string, login string) error {
	// Retrieve the certificate from GitHub every 5 seconds
	var certificates []github.Artifact
	var err error

	certificates, err = RetrieveCertificates(org, repo, login)
	if err != nil {
		return err
	}
	if certificates == nil {
		return nil
	}

	fmt.Printf("Deleting %d existing certificates\n", len(certificates))
	// Loop through the certificates array and remove the artifact from GitHub
	for _, certificate := range certificates {
		err := DeleteCertificate(org, repo, certificate.GetID())
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadCertificate(certificate github.Artifact) error {
	stdout, stderr, err := gh.Exec("api", *certificate.ArchiveDownloadURL, "--method", "GET", "--field", "archive_format=zip")
	if err != nil {
		return fmt.Errorf("error: %s, %s", err, stderr.String())
	}

	// Write stdout to a file
	err = os.WriteFile("cert.zip", stdout.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}
	// Unzip the file
	zipReader, err := zip.NewReader(bytes.NewReader(stdout.Bytes()), int64(stdout.Len()))
	if err != nil {
		return fmt.Errorf("error creating zip reader: %s", err)
	}

	for _, file := range zipReader.File {
		// Open each file in the zip archive
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening file in zip: %s", err)
		}
		defer fileReader.Close()

		// Create the destination file
		destPath := file.Name
		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("error creating destination file: %s", err)
		}
		defer destFile.Close()

		// Copy the contents of the file to the destination file
		_, err = io.Copy(destFile, fileReader)
		if err != nil {
			return fmt.Errorf("error copying file contents: %s", err)
		}
	}

	// Remove the zip file
	err = os.Remove("cert.zip")
	if err != nil {
		return fmt.Errorf("error removing zip file: %s", err)
	}

	return nil
}
