package elastic

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Ivanhahanov/GoLibrary/models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	pdf "github.com/pdfcpu/pdfcpu/pkg/api"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)


func Put(book *models.Book, index string) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	tmp, err := ioutil.TempDir("", book.Slug)
	if err != nil {
		os.Exit(1)
	}


	// Create single page PDFs for in.pdf in outDir using the default configuration.
	pdf.SplitFile(book.Path, tmp, 1, nil)

	files, err := ioutil.ReadDir(tmp)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := tmp + "/"+ file.Name()
		fmt.Println(filePath)
		splitedFilename := strings.Split(file.Name(), "_")
		pageNum := strings.Split(splitedFilename[len(splitedFilename)-1], ".")[0]
		f, _ := os.Open(filePath)
		content, _ := ioutil.ReadAll(bufio.NewReader(f))
		str := base64.StdEncoding.EncodeToString(content)
		intPageNum, _ := strconv.Atoi(pageNum)
		book.Data = append(book.Data, map[string]interface{}{
			"data": str,
			"page": intPageNum,
		})
	}
	//os.RemoveAll(tmp)


	//if err != nil {
	//	log.Fatal(err)
	//}

	req := esapi.IndexRequest{
		Index:    index,
		DocumentID: book.Slug,
		Body:     esutil.NewJSONReader(&book),
		Pipeline: "multiple_attachment",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil{
		log.Println(err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

