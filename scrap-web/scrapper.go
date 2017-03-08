// This Project is used to get all of Indonesian synonyms from http://www.persamaankata.com/
// This Project used for collect Indonesian Synonyms Dictionary
// This Project work just for http://www.persamaankata.com/ website
// From that website there is 39046 word for search word, from index 0 - 39045

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"encoding/csv"
)

var buf bytes.Buffer

type File struct {
	Word    string
	Sinonim []string
}

func main() {
	client := &http.Client{}

	// 39045

	var ArrayFile []File

	fmt.Println("=================================== Start ========================================= ")

	for i := 0; i <= 39045; i++ {
		fmt.Println(i)
		url := fmt.Sprintf("http://www.persamaankata.com/")
		url = fmt.Sprintf("%s%d/", url, i)
		resp, err := client.Get(url)
		if err != nil {
			log.Println("Error = ", err)
			continue
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		defer resp.Body.Close()
		isi := getData(string(bytes))

		word := getWord(string(bytes))

		file := File{
			Word:    word,
			Sinonim: isi,
		}

		ArrayFile = append(ArrayFile, file)

	}

	fmt.Println("=================================================\nWrite To CSV PRocess\n=================================================")

	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)

	for i, record := range ArrayFile {
		fmt.Println(i)
		err = writer.Write(record.Sinonim)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	fmt.Println("=================================== Finish ========================================= ")
	defer writer.Flush()

}

// Tjandrayana Setiawan
// 04-03-2017
// This function is used to get search word from http://www.persamaankata.com/ website
// Parameter data is response body from that website
// Return Function is search word

func getWord(data string) string {
	text := fmt.Sprintf(`id="input_text"`)
	sizeText := len(text)
	pos := strings.Index(data, text)
	data = data[pos : len(data)-pos]
	iterAwal := strings.Index(data, ">")

	word := ""
	for i := sizeText + 9; i < iterAwal-3; i++ {
		word += string(data[i])
	}
	return word
}

// Tjandrayana Setiawan
// 04-03-2017
// This function is used to retrieve synonyms from http://www.persamaankata.com/ website
// Parameter data is response body from that website
// Return Function is array of string that return all of synonyms

func getData(data string) []string {

	word := getWord(data)
	pos := strings.Index(data, "thesaurus_group")
	cari := fmt.Sprintf(`<a id="antonim"></a><div class="thesaurus_group">`)
	antoPost := strings.Index(data, cari)

	if antoPost < 0 {
		cari =
			fmt.Sprintf(`<map name="map_synonym" id="map_synonym">`)
		antoPost = strings.Index(data, cari)
	}

	tamp := ""

	s := []string{word}
	for i := pos; i < antoPost; i++ {
		tamp += string(data[i])
	}

	length := len(tamp)
	for i := 0; i < length; i++ {
		cari2 := fmt.Sprintf(`<a href="http://www.persamaankata.com`)
		cariX := strings.Index(tamp, cari2)
		if cariX < 0 {
			break
		}
		tamp = tamp[cariX:]
		i = i + cariX
		iterAwal := strings.Index(tamp, ">")
		tamp = tamp[iterAwal:]
		i = i + iterAwal
		iterAkhir := strings.Index(tamp, "<")

		text := ""
		for j := 1; j < iterAkhir; j++ {
			text += string(tamp[j])
		}
		i = iterAwal + iterAkhir
		s = append(s, text)

	}

	iterAwal := strings.Index(tamp, ">")
	text := ""
	for i := 1; i < iterAwal; i++ {
		text += string(tamp[i])
	}
	s = append(s, text)

	return s
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}

}

func aaa() {
	fmt.Println("dfsdafasdf")
	fmt.Println("sadasd")
}

func bbb() {
	fmt.Println("masih")
}

func fery() {
	fmt.Println("lalala")
}
