package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Fetch(bUrl string, auth string) []User {

	fmt.Println("Fetching users from dhis2")

	client := &http.Client{}

	fUrl := bUrl + "/api/37/users.json?fields=:owner,displayName,userGroups[id,displayName],userCredentials[id,username,lastLogin,createdBy,lastUpdatedBy,userRoles[id,displayName]]&pageSize=500"

	req, err := http.NewRequest(http.MethodGet, fUrl, http.NoBody)

	req.Header.Add("Authorization", auth)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	response, err := unmarshal(res)

	if err != nil {
		log.Fatal(err)
	}
	var users []User
	users = append(users, response.Users...)

	for response.Pager.NextPage != "" {

		req, err := http.NewRequest(http.MethodGet, response.Pager.NextPage, http.NoBody)

		req.Header.Add("Authorization", auth)

		if err != nil {
			log.Fatal(err)
		}

		res, err = client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		response, err = unmarshal(res)

		if err != nil {
			log.Fatal(err)
		}

		users = append(users, response.Users...)

	}
	fmt.Printf("%d user(s) returned\n", len(users))
	return users

}

func unmarshal(res *http.Response) (Response, error) {
	resData, err := io.ReadAll(res.Body)

	if err != nil {
		return Response{}, err
	}
	var resObj Response
	err = json.Unmarshal(resData, &resObj)
	if err != nil {

		return Response{}, err
	}
	return resObj, nil
}
