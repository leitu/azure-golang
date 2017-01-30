// This is the summarize how to interactive with Azure blob storage
// Create container/upload file/get the uri
// it's pretty ugly code, but can be improved when you really need it.

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/storage"
)

//define account
const (
	AzureAccount   = "CHANGEME"
	AzureAccessKey = "CHANGEME"
	containerName  = "CHANGEME"
)

func main() {

	fmt.Println("Create storage client...")
	client, err := storage.NewBasicClient(AzureAccount, AzureAccessKey)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("Create blob client...")
	blobClient := client.GetBlobService()

	fmt.Println("Create container with private access type...")
	if _, err := blobClient.CreateContainerIfNotExists(containerName, storage.ContainerAccessTypePrivate); err != nil {
		fmt.Print(err)
	}

	fileName := "samplefile"
	if _, err := os.Stat(fileName); err == nil {
		fmt.Print(err)
	}
	bytesRead, err := ioutil.ReadFile(fileName)
	blockBlobName := "filesample"

	fmt.Println("Create an empty block blob...")
	if err := blobClient.CreateBlockBlob(containerName, blockBlobName); err != nil {
		fmt.Print(err)
	}

	fmt.Println("Upload a block...")
	blockID := base64.StdEncoding.EncodeToString([]byte("random"))
	blockData := bytesRead
	if err := blobClient.PutBlock(containerName, blockBlobName, blockID, blockData); err != nil {
		fmt.Print(err)
	}

	fmt.Println("Build uncommitted blocks list...")
	blocksList, err := blobClient.GetBlockList(containerName, blockBlobName, storage.BlockListTypeUncommitted)
	if err != nil {
		fmt.Print(err)
	}
	uncommittedBlocksList := make([]storage.Block, len(blocksList.UncommittedBlocks))
	for i := range blocksList.UncommittedBlocks {
		uncommittedBlocksList[i].ID = blocksList.UncommittedBlocks[i].Name
		uncommittedBlocksList[i].Status = storage.BlockStatusUncommitted
	}
	fmt.Println("Commit blocks...")
	if err = blobClient.PutBlockList(containerName, blockBlobName, uncommittedBlocksList); err != nil {
		fmt.Print(err)
	}

	fmt.Println("Get block list...")
	list, err := blobClient.GetBlockList(containerName, blockBlobName, storage.BlockListTypeAll)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Block blob '%v' block list\n", blockBlobName)
	fmt.Println("\tCommitted Blocks' IDs")
	for _, b := range list.CommittedBlocks {
		fmt.Printf("\t\t%v\n", b.Name)
	}
	fmt.Println("\tUncommited Blocks' IDs")
	for _, b := range list.UncommittedBlocks {
		fmt.Printf("\t\t%v\n", b.Name)
	}

	t := time.Now()
	fmt.Println("Get sasTokenUri")
	url, err := blobClient.GetBlobSASURI(containerName, blockBlobName, t, "read")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("this is '%v \n", url)
}

func randomData(strLen int) []byte {
	ran := 'z' - '0'
	text := make([]byte, strLen)
	for i := range text {
		char := rand.Int()
		char %= int(ran)
		char += '0'
		text[i] = byte(char)
	}
	return text
}
