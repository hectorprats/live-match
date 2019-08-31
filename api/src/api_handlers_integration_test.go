package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func post(url string, contents []byte, version string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(contents))
	req.Header.Add("version", version)
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error sending request to %s", url))
	}

	return response, nil
}

func Test_NewGoal_Returns_202(t *testing.T) {
	// Arrange
	type request struct {
		team   string
		minute uint8
		player uint8
	}
	r := &request{
		team:   "MCI",
		minute: 10,
		player: 10,
	}
	contents, err := json.Marshal(r)
	if err != nil {
		t.Fatal("Unable to build request json for NewGoal.")
		return
	}

	// Act
	response, err := post("http://localhost:9090/competitions/premier-league/matches/ARS/MUT/newGoal", contents, "1.0")
	if err != nil {
		t.Fatal("Error sending post request to NewGoal endpoint.")
	}
	defer response.Body.Close()

	// Assert
	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 202, but got %d instead", response.StatusCode)
	}
}

func Test_NewPenalty_Returns_202(t *testing.T) {
	type request struct {
		team   string
		minute uint8
		player uint8
		reason string
	}
	client := &http.Client{}
	r := &request{
		team:   "MCI",
		minute: 10,
		player: 10,
		reason: "Grave foul",
	}
	contents, err := json.Marshal(r)
	if err != nil {
		t.Fatal("Unable to build request json for NewPenalty.")
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:9090/competitions/premier-league/matches/ARS/MUT/NewPenalty", bytes.NewBuffer(contents))
	req.Header.Add("version", "1.0")
	response, err := client.Do(req)
	if err != nil {
		t.Fatal("Error sending post request to NewPenalty endpoint.")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 202, but got %d instead", response.StatusCode)
	}
}

func Test_NewYellowCard_Returns_202(t *testing.T) {
	type request struct {
		team   string
		minute uint8
		player uint8
		reason string
	}
	client := &http.Client{}
	r := &request{
		team:   "MCI",
		minute: 10,
		player: 10,
		reason: "Grave foul",
	}
	contents, err := json.Marshal(r)
	if err != nil {
		t.Fatal("Unable to build request json for NewYellowCard.")
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:9090/competitions/premier-league/matches/ARS/MUT/NewYellowCard", bytes.NewBuffer(contents))
	req.Header.Add("version", "1.0")
	response, err := client.Do(req)
	if err != nil {
		t.Fatal("Error sending post request to NewYellowCard endpoint.")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 202, but got %d instead", response.StatusCode)
	}
}

func Test_NewRedCard_Returns_202(t *testing.T) {
	type request struct {
		team   string
		minute uint8
		player uint8
		reason string
	}
	client := &http.Client{}
	r := &request{
		team:   "MCI",
		minute: 10,
		player: 10,
		reason: "Grave foul",
	}
	contents, err := json.Marshal(r)
	if err != nil {
		t.Fatal("Unable to build request json for NewRedCard.")
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:9090/competitions/premier-league/matches/ARS/MUT/NewRedCard", bytes.NewBuffer(contents))
	req.Header.Add("version", "1.0")
	response, err := client.Do(req)
	if err != nil {
		t.Fatal("Error sending post request to NewRedCard endpoint.")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 202, but got %d instead", response.StatusCode)
	}
}

func Test_NewSubstitution_Returns_202(t *testing.T) {
	type request struct {
		team      string
		inPlayer  uint8
		outPlayer uint8
		minute    uint8
	}
	client := &http.Client{}
	r := &request{
		team:      "MCI",
		minute:    10,
		inPlayer:  10,
		outPlayer: 11,
	}
	contents, err := json.Marshal(r)
	if err != nil {
		t.Fatal("Unable to build request json for NewSubstitution.")
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:9090/competitions/premier-league/matches/ARS/MUT/newSubstitution", bytes.NewBuffer(contents))
	req.Header.Add("version", "1.0")
	response, err := client.Do(req)
	if err != nil {
		t.Fatal("Error sending post request to NewSubstitution endpoint.")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 202, but got %d instead", response.StatusCode)
	}
}
