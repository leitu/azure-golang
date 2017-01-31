// This is the sample code with get sasToken from Azure blob storage
// and inject the key to template
// it's pretty ugly code, but can be improved when you really need it.

package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/Azure/azure-sdk-for-go/storage"
)

//define account
const (
	AzureAccount   = "CHANGEME"
	AzureAccessKey = "CHANGEME"
	containerName  = "CHANGEME"
	blockBlobName  = "CHANGME"
)

func main() {

	client, err := storage.NewBasicClient(AzureAccount, AzureAccessKey)
	if err != nil {
		fmt.Print(err)
	}

	blobClient := client.GetBlobService()

	if _, err := blobClient.CreateContainerIfNotExists(containerName, storage.ContainerAccessTypePrivate); err != nil {
		fmt.Print(err)
	}

	t := time.Now().AddDate(0, 0, 1)
	fmt.Println("Get sasTokenUri")
	url, err := blobClient.GetBlobSASURI(containerName, blockBlobName, t, "rwd")
	if err != nil {
		fmt.Print(err)
	}
	bloburl := "https://" + AzureAccount + ".blob.core.windows.net/" + containerName + "/" + blockBlobName
	sasToken := strings.Trim(url, bloburl)
	fmt.Printf("this is '%v' \n", sasToken)

	fmt.Print(sasToken)

	generateFile(sasToken)

}

func generateFile(sasToken string) {
	file := "parameters.json.tmpl"
	if _, err := os.Stat(file); err != nil {
		fmt.Print(err)
	}

	config := map[string]string{
		"sasToken": sasToken,
	}

	newfile := "parameters.json"

	templateFile, err := template.ParseFiles(file)
	if err != nil {
		fmt.Print(err)
	}

	f, err := os.Create(newfile)
	if err != nil {
		fmt.Println("create file: ", err)
	}
	err = templateFile.Execute(f, config)
	if err != nil {
		fmt.Print(err)
	}
	f.Close()

	fmt.Printf("Successful created '%v' \n", newfile)
}
