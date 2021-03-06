// VERSION MANAGER PACKAGE

package vm

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	u "github.com/zaviermiller/zen/internal/util"
)

type updateCheck struct {
	TagName            string
	UpdateSize         int
	BrowserDownloadUrl string
}

func (uc *updateCheck) Download(zenPath string) {
	t1 := time.Now()

	zenBin, err := os.Create(zenPath + ".tmp")
	check(err)

	defer zenBin.Close()

	zenUpdate, err := http.Get(uc.BrowserDownloadUrl)
	check(err)

	defer zenUpdate.Body.Close()

	if zenUpdate.StatusCode != http.StatusOK {
		check(errors.New(fmt.Sprint("HTTP Request failed with status %s", zenUpdate.StatusCode)))
	}

	updateWriter := &zenUpdateWriter{Reader: zenUpdate.Body, size: uc.UpdateSize}

	_, err = io.Copy(zenBin, updateWriter)
	check(err)

	err = os.Rename(zenPath+".tmp", zenPath)
	check(err)
	err = os.Chmod(zenPath, 0774)
	check(err)

	t2 := time.Now()
	diff := t2.Sub(t1)

	u.ZenWeirdLog(fmt.Sprintf(u.Clear+u.Bright+"Done updating in %s\n\n"+u.Normal, diff))
}

func (uc *updateCheck) ReloadBinary() {
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.Run()

}

func finishUpdate() {

}

func CheckUpdate() bool {
	GOOS := runtime.GOOS
	GOARCH := runtime.GOARCH
	if strings.Contains(os.Args[0], ".tmp") {
		finishUpdate()
		return false
	}
	resp, err := http.Get("https://api.github.com/repos/zaviermiller/zen/releases")
	check(err)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// json bullshit
	var jsonVar interface{}
	json.Unmarshal(body, &jsonVar)
	jsonVar = jsonVar.([]interface{})[0]

	// build test version
	ghVersion := ParseVersion(jsonVar.(map[string]interface{})["tag_name"].(string))

	if ghVersion.GreaterThan(VERSION) {
		assetName := fmt.Sprintf("zen%s-%s-%s", ghVersion.String(), GOOS, GOARCH)
		var asset map[string]interface{}
		jsonVar = jsonVar.(map[string]interface{})["assets"]
		for _, a := range jsonVar.([]interface{}) {
			if a.(map[string]interface{})["name"] == assetName {
				asset = a.(map[string]interface{})
			}
		}

		if asset == nil {
			return false
		}

		fmt.Println(u.Bright + u.Blue + "\n   UPDATE FOUND!" + u.Green + " Version: " + ghVersion.String() + u.Normal + " (" + binSize(int64(asset["size"].(float64))) + ")")
		u.ZenLog("Would you like to update now?" + u.Yellow + " (y/n)" + u.Normal + ": ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		check(err)

		fmt.Println("")

		switch input[0] {
		case 'y':
			dlUrl := asset["browser_download_url"].(string)

			// updateCheck object for managing the update
			update := updateCheck{TagName: ghVersion.String(), BrowserDownloadUrl: dlUrl, UpdateSize: int(asset["size"].(float64))}
			zenPath, err := exec.LookPath("zen")
			check(err)

			// download the binary update (exits on err)
			update.Download(zenPath)

			// reload the binary file if no err
			update.ReloadBinary()

			return true
		case 'n':

		default:
		}
	}

	return false

}

func check(err error) {
	if err != nil {
		fmt.Println(u.Bright + u.Red + "ZEN ERROR: " + u.Normal + err.Error())
		os.Exit(1)
	}
}

func binSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
