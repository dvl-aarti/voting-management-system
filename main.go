package main

import (
	"encoding/json"
	"errors"
	"onlineVoting/auth"
	"onlineVoting/config"
	"onlineVoting/pojo"
	"onlineVoting/services"
	"strings"

	"fmt"
	"log"
	"net/http"
)

var con = services.Connection{}
var responseData pojo.Response

func init() {
	con.Server = "mongodb://localhost:27017"
	//  cla.Server = "mongodb+srv://m001-student:m001-mongodb-basics@sandbox.7zffz3a.mongodb.net/?retryWrites=true&w=majority"
	// con.Database = "onlineVotingSystem"
	con.Database = config.APP_CONFIG.Database
	if config.APP_CONFIG.Environment == "tests" {
		con.Database = config.APP_CONFIG.TestDatabase
	}
	con.Collection = "user"
	con.Collection2 = "login"
	con.Collection3 = "election"

	con.Connect()
}
func main() {
	// http.HandleFunc("/add-blood-group-data/", addBloodGroupData)
	http.HandleFunc("/api/add-user/", addUserRecord)
	http.HandleFunc("/api/update-user/", updateUserDetailsById)
	http.HandleFunc("/api/verify-user-detail/", verifyUserDetailsById)
	http.HandleFunc("/api/deactivate-user-detail/", deactivateUserDetails)
	http.HandleFunc("/api/search-user-detailbyId/", searchUsersDetailsById)
	http.HandleFunc("/api/search-user-detail/", searchUsersDetailsFilter)
	http.HandleFunc("/api/add-election/", addElectionRecord)
	http.HandleFunc("/api/find-electionbyId/", findElectionById)
	http.HandleFunc("/api/update-election/", updateElectionDetailsById)
	http.HandleFunc("/api/search-election-detail/", searchElectionDetailsFilter)
	http.HandleFunc("/api/election-deactivate/", DeactivateElection)
	http.HandleFunc("/api/add-candidate/", addCandidateRecord)
	http.HandleFunc("/api/verify-candidate/", verifyCandidateDetails)
	http.HandleFunc("/api/cast-vote/", castVote)
	http.HandleFunc("/api/election-result/", ElectionResultById)
	http.HandleFunc("/api/candidate-list/", searchCandidateList)
	http.HandleFunc("/api/login/", login)
	fmt.Println("Excecuted Main Method")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addUserRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method", "")
		return
	}

	var data pojo.UserDetailsRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "")
		return
	}

	if data.Email == "" || data.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter mailId or Password", "")
		return
	}
	if data.Role == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter role field value Admin or Voter", "")
		return
	}
	if result, msg, err := con.SaveUserDetails(data); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}, err, msg string) {

	if err == "error" {
		responseData.Success = "false"
	} else {
		responseData.Success = "true"
	}
	responseData.SucessMsg = msg
	responseData.SucessCode = fmt.Sprintf("%v", code)
	responseData.Response = payload
	response, _ := json.Marshal(responseData)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string, msg2 string) {
	respondWithJson(w, code, map[string]string{"error": msg}, "error", msg2)
}

func updateUserDetailsById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "PUT" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	_, role, err := validateToken(token)
	if role != "Voter" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]

	var dataBody pojo.UserDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "Error Occurred")
		return
	}

	if result, msg, err := con.UpdateUserDetailsById(dataBody, id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func verifyUserDetailsById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	mail, role, err := validateToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	if role != "Admin" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}

	var dataBody pojo.VerifyUser
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := con.VerifyUserDetails(dataBody, mail); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func deactivateUserDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "DELETE" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	_, role, err := validateToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	if role != "Admin" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}

	segment := strings.Split(r.URL.Path, "/")
	id := segment[len(segment)-1]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide Id for Search", "Error Occurred")
	}

	if result, msg, err := con.DeactivateUserDetails(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func searchUsersDetailsById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	segment := strings.Split(r.URL.Path, "/")
	id := segment[len(segment)-1]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide Id for Search", "Error Occurred")
	}

	if result, msg, err := con.SearchUsersDetailsById(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func searchUsersDetailsFilter(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var dataBody pojo.SearchUser
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := con.SearchFilterUserDetails(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

//===========================election details=====================================

func addElectionRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var dataBody pojo.ElectionDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := con.SaveElectionDetails(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}
func findElectionById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	segment := strings.Split(r.URL.Path, "/")
	id := segment[len(segment)-1]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide Id for Search", "Error Occurred")
	}

	if result, msg, err := con.FindElectionById(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func updateElectionDetailsById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "PUT" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	_, role, err := validateToken(token)
	if role != "Admin" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]

	var dataBody pojo.ElectionDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "Error Occurred")
		return
	}

	if result, msg, err := con.UpdateElectionDetailsById(dataBody, id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func searchElectionDetailsFilter(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var dataBody pojo.SearchFilterElectionReq
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := con.SearchFilterElectionDetails(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func DeactivateElection(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "DELETE" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	_, role, err := validateToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	if role != "Admin" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}

	segment := strings.Split(r.URL.Path, "/")
	id := segment[len(segment)-1]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide Id for Search", "Error Occurred")
	}

	if result, msg, err := con.DeactivateElection(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

//===================================candidate details======================================

func addCandidateRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method", "")
		return
	}

	var data pojo.CandidatesRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "")
		return
	}

	if result, msg, err := con.SaveCandidateDetails(data); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func verifyCandidateDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	mail, role, err := validateToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	if role != "Admin" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}

	var dataBody pojo.VerifyCandidates
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := con.VerifyCandidate(dataBody, mail); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

// =======================login===========================================================

func login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "")
		return
	}

	var loginDetails pojo.SignInInput
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "")
	}

	if result, msg, err := con.GenerateToken(loginDetails); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

// ============================valid token=======================================
func validateToken(token string) (string, string, error) {
	if token == "" {
		return "", "", errors.New("Please Enter Token")
	}
	mail, err := auth.ValidateToken(token)
	if err != nil {
		return "", "", errors.New("Either Token Is Invalid Or Expired")
	}
	role := con.FetchRole(mail)
	fmt.Println("MailId:", mail)
	fmt.Println("Role:", role)
	return mail, role, err
}

// ================================vote Result=======================================
func castVote(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method", "")
		return
	}

	token := r.Header.Get("tokenid")
	mail, _, err := validateToken(token)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	var data pojo.CastVote

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "")
		return
	}

	if result, msg, err := con.CastVote(data, mail); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func ElectionResultById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	segment := strings.Split(r.URL.Path, "/")
	id := segment[len(segment)-1]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide Id for Search", "Error Occurred")
	}

	if result, msg, err := con.ElectionResultById(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func searchCandidateList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var dataBody pojo.SearchElectionResult
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := con.SearchCandidateList(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}
