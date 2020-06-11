package get

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) == 1 {

		// p : The argument passed in by user
		p := args[0]

		// Construct the URL to the archive zip
		urlTokens := []string{"https://", p, "/archive/master.zip"}
		url := strings.Join(urlTokens, "")

		fmt.Printf("Attempt to download %v ...\n", url)

		// Download the archive zip to current folder
		err := DownloadFile("master.zip", url)
		if err != nil {
			panic(err)
		} else {

			fmt.Println("Download success!")
			fmt.Println("Unzipping files...")

			// Unpack the zip file to current folder
			files, err := Unzip("master.zip", ".")
			if err != nil {
				log.Fatal(err)
			} else {

				// Print out the list of unzipped files
				fmt.Println(strings.Join(files, "\n") + "\n")

				// Tokenize
				packageTokens := strings.Split(p, "/")
				mGithub := packageTokens[0]
				mAuthor := packageTokens[1]
				mPackage := packageTokens[2]

				// Set up folder in $GOPATH/src
				gopath := os.Getenv("GOPATH")
				packageSrcTokens := []string{gopath, "src", mGithub, mAuthor}
				packageSrcPath := strings.Join(packageSrcTokens, "/")
				fmt.Println("Creating package folder in $GOPATH/src ...")
				os.MkdirAll(packageSrcPath, 0755)

				fmt.Println("Moving package files into the src folder ...")
				os.Rename(mPackage+"-master", packageSrcPath+"/"+mPackage)

				fmt.Println("Cleaning up ...")
				os.RemoveAll(mPackage + "-master")
				os.Remove("master.zip")

				fmt.Println("Now you need to run the following command to install the package:\n")
				fmt.Println("cd " + packageSrcPath + "/" + mPackage + " && go install && cd -")
			}
		}
	} else {
		panic("USAGE: get github.com/<AUTHOR>/<PACKAGE>")
	}
}

// Credits to: https://golangcode.com/download-a-file-from-a-url/
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Credits to: https://golangcode.com/unzip-files-in-go/
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}
