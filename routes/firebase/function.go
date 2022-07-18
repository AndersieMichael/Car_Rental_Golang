package firebase

import (
	"bytes"
	"context"
	"fmt"
	"golangSecond/utilities/firebaseconfig"
	"io"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

func UploadExcel(fileInput []byte, Filename string, Mime_Type string) (string, error) {
	id := uuid.New()
	var url string
	client := firebaseconfig.StorageClient()

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Println(err)
	}

	// upload_id := strconv.Itoa(Periode_Incentive_ID)
	ctx := context.Background()
	object := bucket.Object("test/" + Filename + "_"  + fmt.Sprint(time.Now().Unix()))

	writer := object.NewWriter(ctx)
	writer.ContentType = (Mime_Type)
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}

	defer writer.Close()

	if _, err := io.Copy(writer, bytes.NewReader(fileInput)); err != nil {
		return url, err
	}

	bucket_name := "project1-9c741.appspot.com"
	filename := ("test/" + Filename + "_" + fmt.Sprint(time.Now().Unix()))
	expired := time.Date(2200, time.November, 1, 23, 0, 0, 0, time.Local)


	url, err = storage.SignedURL(bucket_name, filename, &storage.SignedURLOptions{
		GoogleAccessID: firebaseconfig.Credential["client_email"],
		PrivateKey:     []byte(firebaseconfig.Credential["private_key"]),
		Method:         "GET",
		Expires:        expired,
	})
	
	return url, err
}