// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

// @me Discord Patch request to change Username
func (in *Instance) NameChanger(name string) (cycletls.Response, error) {
	url := "https://discord.com/api/v9/users/@me"
	data := NameChange{
		Username: name,
		Password: in.Password,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return cycletls.Response{}, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return cycletls.Response{}, fmt.Errorf("error while getting cookie %v", err)
	}
	resp, err := in.Client.Do(url, in.CycleOptions(string(bytes), in.AtMeHeaders(cookie)), "PATCH")
	if err != nil {
		return cycletls.Response{}, err
	}

	return resp, nil

}

// @me Discord Patch request to change Nickname
func (in *Instance) NickNameChanger(name string, guildid string) (cycletls.Response, error) {
	url := fmt.Sprintf("https://discord.com/api/v9/guilds/%s/members/@me", guildid)
	data := NickNameChange{
		Nickname: name,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return cycletls.Response{}, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return cycletls.Response{}, fmt.Errorf("error while getting cookie %v", err)
	}
	resp, err := in.Client.Do(url, in.CycleOptions(string(bytes), in.AtMeHeaders(cookie)), "PATCH")
	if err != nil {
		return cycletls.Response{}, err
	}
	return resp, nil
}

// @me Discord Patch request to change Avatar
func (in *Instance) AvatarChanger(avatar string) (cycletls.Response, error) {
	url := "https://discord.com/api/v9/users/@me"
	avatar = "data:image/png;base64," + avatar
	data := AvatarChange{
		Avatar: avatar,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return cycletls.Response{}, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return cycletls.Response{}, fmt.Errorf("error while getting cookie %v", err)
	}

	resp, err := in.Client.Do(url, in.CycleOptions(string(bytes), in.AtMeHeaders(cookie)), "PATCH")
	if err != nil {
		return cycletls.Response{}, err
	}
	return resp, nil
}

// Encoding images to b64
func EncodeImg(pathToImage string) (string, error) {

	image, err := os.Open(pathToImage)

	if err != nil {
		return "", err
	}

	defer image.Close()

	reader := bufio.NewReader(image)
	imagebytes, _ := ioutil.ReadAll(reader)

	extension := http.DetectContentType(imagebytes)

	switch extension {
	default:
		return "", fmt.Errorf("unsupported image type: %s", extension)
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imagebytes))
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil

	case "image/png":
		return base64.StdEncoding.EncodeToString(imagebytes), nil
	}
}

// Get all file paths in a directory
func GetFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		paths = append(paths, dir+"/"+file.Name())
	}
	return paths, nil
}

func (in *Instance) BioChanger(bios []string) error {
	chosenOne := bios[rand.Intn(len(bios))]
	site := "https://discord.com/api/v9/users/@me"
	cookie, err := in.GetCookieString()
	if err != nil {
		return fmt.Errorf("error while getting cookie: %v", err)
	}
	resp, err := in.Client.Do(site, in.CycleOptions(`{"bio": "`+chosenOne+`"}`, in.AtMeHeaders(cookie)), "PATCH")
	if err != nil {
		return fmt.Errorf("error while sending request: %v", err)
	}
	if resp.Status != 200 {
		return fmt.Errorf("error while changing bio %v %v", resp.Status, resp.Body)
	}
	return nil
}

func ValidateBios(bios []string) []string {
	var validBios []string
	for i := 0; i < len(bios); i++ {
		if len(bios[i]) > 190 {
			continue
		}
		validBios = append(validBios, bios[i])
	}
	return validBios
}

func (in *Instance) RandomHypeSquadChanger() error {
	site := "https://discord.com/api/v9/hypesquad/online"
	cookie, err := in.GetCookieString()
	if err != nil {
		return fmt.Errorf("error while getting cookie: %v", err)
	}
	resp, err := in.Client.Do(site, in.CycleOptions(fmt.Sprintf(`{"house_id": %v}`, rand.Intn(3)+1), in.AtMeHeaders(cookie)), "POST")
	if err != nil {
		return fmt.Errorf("error while sending request: %v", err)
	}
	if resp.Status != 204 {
		return fmt.Errorf("error while changing hype squad %v %v", resp.Status, resp.Body)
	}
	return nil
}

func (in *Instance) ChangeToken(newPassword string) (string, error) {
	site := "https://discord.com/api/v9/users/@me"
	payload := fmt.Sprintf(`
	{
		"password": "%v",
		"new_password": "%v"
	}
	`, in.Password, newPassword)
	cookie, err := in.GetCookieString()
	if err != nil {
		return "", fmt.Errorf("error while getting cookie: %v", err)
	}
	resp, err := in.Client.Do(site, in.CycleOptions(payload, in.AtMeHeaders(cookie)), "PATCH")
	if err != nil {
		return "", fmt.Errorf("error while sending request: %v", err)
	}
	if resp.Status != 200 {
		return "", fmt.Errorf("invalid status code %v while changing token %v", resp.Status, resp.Body)
	}
	if strings.Contains(resp.Body, "token") {
		var response map[string]interface{}
		err := json.Unmarshal([]byte(resp.Body), &response)
		if err != nil {
			return "", fmt.Errorf("error while unmarshalling response: %v", err)
		}
		return response["token"].(string), nil
	} else {
		return "", fmt.Errorf("error while changing token %v body does not contain token", resp.Body)
	}
}
