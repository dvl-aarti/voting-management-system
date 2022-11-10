package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var AdminTokenId = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IlJha3Vsa3VtYXJyYWowMUBnbWFpbC5jb20iLCJleHAiOjE2NjgxMTMxNzN9.Bd5kIMfwKgeIEf1pwolK0cr5UCCMBt0a_r6mbHB7IZs"
var VoterTokenId = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IkFha2FzaGt1bWFycmFqMDFAZ21haWwuY29tIiwiZXhwIjoxNjY4MTEzMjI1fQ.z9QqLT8-m9d3KbwTBsPZz4eWz48BLHn-_39turS2AeM"

func TestAddUser(t *testing.T) {
	payload := strings.NewReader(`
	{
		"role": "Admin",
		"name": "Rohit Raj",
		"email": "Rohitkumarraj01@gmail.com",
		"password": "Rohit@123!",
		"phone_number": "8524367472",
		"personal_info": {
			"name": "Rahul Raj",
			"father_name": "kishan Raj",
			"dob": "1991-11-08",
			"age": "25",
			"document_type": "Adhar Card",
			"address": {
				"street": "06",
				"city": "gurugram",
				"state": "Haryana",
				"zip_code": "120011",
				"country": "India"
			}
		},
		"uploaded_docs": {
			"document_type": "Adhar Card",
			"document_identification_no": "1232345",
			"document_path": "D:/SanDiskMemoryZone_QuickStartGuide.pdf"
		}
	}
	`)

	req := httptest.NewRequest("POST", "/api/add-user/", payload)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	addUserRecord(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdateUserDetailsById(t *testing.T) {
	payload := strings.NewReader(`
	{
		"role": "Admin",
		"name": "RakulPreet Raj",
		"email": "Rakulkumarraj01@gmail.com",
		"password": "Rakul@123!",
		"phone_number": "8524365472",
		"personal_info": {
			"name": "Rakul Raj",
			"father_name": "Rahuxcdl Raj",
			"dob": "1991-01-22",
			"age": "28",
			"document_type": "Adhar Card",
			"address": {
				"street": "03",
				"city": "gurugramdf",
				"state": "Haryanadf",
				"zip_code": "120011",
				"country": "India"
			}
		},
		"uploaded_docs": {
			"document_type": "Adharfe Card",
			"document_identification_no": "123vddf45",
			"document_path": ""
		}
	}
	`)

	req := httptest.NewRequest("PUT", "/api/update-user/635fa51ef789c99423194c53", payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("tokenid", VoterTokenId)
	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	updateUserDetailsById(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestVerifyUserDetailsById(t *testing.T) {
	payload := strings.NewReader(`
	{
		"id":"63626df1b03f76f1047b7dab",
		"email": "Aakashkumarraj01@gmail.com"
		
	}
	`)

	req := httptest.NewRequest("POST", "/api/verify-user-detail/", payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("tokenid", AdminTokenId)
	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	verifyUserDetailsById(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDeleteUser(t *testing.T) {

	req := httptest.NewRequest("DELETE", "/api/user/delete/635f9f388bbc58a8ae0eb5b6", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("tokenid", AdminTokenId)
	w := httptest.NewRecorder()
	// handler := http.HandlerFunc(controller.DeleteUser)

	deactivateUserDetails(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSearchOneUser(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/search-user-detailbyId/635fa51ef789c99423194c53", nil)
	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("tokenid", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IlJha3Vsa3VtYXJyYWowMUBnbWFpbC5jb20iLCJleHAiOjE2NjgwMDAxNDl9.Q4d79S3MTr5r85iWRvNubr5-oYVjQ7cAAwW7EaNu_aI")

	w := httptest.NewRecorder()
	// handler(w, req)
	// http.NewServeMux().HandleFunc("localhost:8080/api/user/search/63635d75a68e40fe497eac67", controller.SearchOneUser)
	// handler := http.HandlerFunc(controller.SearchOneUser)

	searchUsersDetailsById(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSearchMultipleUser(t *testing.T) {

	payload := strings.NewReader(`{
		"role": "Admin",
		"city": "",
		"state": "",
		"zip_code": "",
		"country": "",
		"is_verified": true
	}`)

	req := httptest.NewRequest("POST", "/api/search-user-detail/", payload)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	// handler := http.HandlerFunc(controller.SearchMultipleUser)

	searchUsersDetailsFilter(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// ----------------------election test cases-------------------

func TestAddElection(t *testing.T) {
	payload := strings.NewReader(`
	{
		"election_date":"2022-11-22",
		"result_date":"2022-11-25",
		"election_status":"",
		"result":"jgfkyt",
		"location":" New Delhi"
	}
	`)

	req := httptest.NewRequest("POST", "/api/add-election/", payload)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	addElectionRecord(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestFindElectionById(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/find-electionbyId/6367b19b49c2b31ed16ce479", nil)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	// handler(w, req)
	// http.NewServeMux().HandleFunc("localhost:8080/api/user/search/63635d75a68e40fe497eac67", controller.SearchOneUser)
	// handler := http.HandlerFunc(controller.SearchOneUser)

	findElectionById(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdateElectionDetailsById(t *testing.T) {
	payload := strings.NewReader(`
	{
		"election_date":"2022-11-09",
		"result_date":"2022-11-11",
		"election_status":"ready",
		"location":"Bihar"
	}
	`)

	req := httptest.NewRequest("PUT", "/api/update-election/636b80d02a455e9567c916e5", payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("tokenid", AdminTokenId)
	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	updateElectionDetailsById(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSearchElectionDetailsFilter(t *testing.T) {

	payload := strings.NewReader(`{
		"election_date":"2022-11-22",
		"result_date":"2022-11-25",
		"election_status":""
	}`)

	req := httptest.NewRequest("POST", "/api/search-election-detail/", payload)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	// handler := http.HandlerFunc(controller.SearchMultipleUser)

	searchElectionDetailsFilter(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDeleteElection(t *testing.T) {

	req := httptest.NewRequest("DELETE", "/api/election-deactivate/6363871c01e7236d16372c84", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("tokenid", AdminTokenId)
	w := httptest.NewRecorder()
	// handler := http.HandlerFunc(controller.DeleteUser)

	deactivateUserDetails(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// ----------------------candidate test cases-------------------

func TestAddCandidate(t *testing.T) {
	payload := strings.NewReader(`
	{
		"election_id":"636b80d02a455e9567c916e5",
		"user_id":"636b80d02a455e9567c916e4",
		"candidate_name":"",
		"commitments":["sdfewfd"],
		"vote_sign":"D:/SanDiskMemoryZone_QuickStartGuide.pdf",
		"nomination_status":""
	}
	`)

	req := httptest.NewRequest("POST", "/api/add-candidate/", payload)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	addCandidateRecord(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestVerifyCandidateDetails(t *testing.T) {
	payload := strings.NewReader(`
	{
		"election_id":"6367b19b49c2b31ed16ce479",
		"user_id":"635fa51ef789c99423194c53"
	}
	`)

	req := httptest.NewRequest("POST", "/api/verify-candidate/", payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("tokenid", AdminTokenId)
	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	verifyCandidateDetails(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// ---------------------------cast vote-----------------------------------
func TestCastVote(t *testing.T) {
	payload := strings.NewReader(`
	{
		"election_id":"6367b19b49c2b31ed16ce479",
		"candidate_id":"63626df1b03f76f1047b7dab"
	}
	`)

	req := httptest.NewRequest("POST", "/api/cast-vote/", payload)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	castVote(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestElectionResultById(t *testing.T) {

	req := httptest.NewRequest("PUT", "/api/election-result/6367b19b49c2b31ed16ce479", nil)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// handler := http.HandlerFunc(controller.AddUser)

	ElectionResultById(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSearchCandidateList(t *testing.T) {

	payload := strings.NewReader(`{
		"election_id":"6367b19b49c2b31ed16ce479"
	}`)

	req := httptest.NewRequest("POST", "/api/candidate-list/", payload)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	// handler := http.HandlerFunc(controller.SearchMultipleUser)

	searchCandidateList(w, req)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
