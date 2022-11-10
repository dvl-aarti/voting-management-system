package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"onlineVoting/auth"
	"onlineVoting/pojo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	Server      string
	Database    string
	Collection  string
	Collection2 string
	Collection3 string
}

var CollectionUser *mongo.Collection
var CollectionLogin *mongo.Collection
var CollectionElection *mongo.Collection
var ctx = context.TODO()
var insertDocs int

const uploadPath = "upload/"

func (c *Connection) Connect() {
	clientOptions := options.Client().ApplyURI(c.Server)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	CollectionUser = client.Database(c.Database).Collection(c.Collection)
	CollectionLogin = client.Database(c.Database).Collection(c.Collection2)
	CollectionElection = client.Database(c.Database).Collection(c.Collection3)

}

// ===================================userDetails============================================
func (e *Connection) SaveUserDetails(reqBody pojo.UserDetailsRequest) ([]*pojo.UserDetails, string, error) {
	var data []*pojo.UserDetails
	bool, err := ValidateByNameAndDob(reqBody)
	if err != nil {
		return data, "", err
	}
	if !bool {
		return data, "", errors.New("User already present")
	}
	if err != nil {
		log.Println(err)
		return data, "", err
	}

	boolean, err := validateByEmail(reqBody)
	if err != nil {
		return data, "", err
	}
	if !boolean {
		return data, "", errors.New("email already present")
	}
	if err != nil {
		log.Println(err)
		return data, "", err
	}
	msg, err := UploadFile(reqBody.UploadedDocs.DocumentPath)
	if err != nil {
		log.Println(err)
		return data, "", errors.New("Unable to upload file")
	}
	fmt.Println("Upload file:", msg)
	reqBody.IsVerified = false
	finalData, err := CollectionUser.InsertOne(ctx, reqBody)
	if err != nil {
		log.Println(err)
		return data, "", errors.New("Unable to store data")
	}
	result, err := CollectionUser.Find(ctx, bson.D{primitive.E{Key: "_id", Value: finalData.InsertedID}})
	if err != nil {
		log.Println(err)
		return data, "", err
	}
	data, err = convertDbResultIntoUserStruct(result)
	if err != nil {
		log.Println(err)
		return data, "", err
	}
	return data, "Saved Successfully", nil
}

func ValidateByNameAndDob(reqbody pojo.UserDetailsRequest) (bool, error) {
	dobStr := reqbody.PersonalInfo.DOB
	fmt.Println(dobStr)
	var result []*pojo.UserDetails
	data, err := CollectionUser.Find(ctx, bson.D{{Key: "name", Value: reqbody.Name}, {Key: "personal_info.dob", Value: dobStr}})
	if err != nil {
		return false, err
	}
	result, err = convertDbResultIntoUserStruct(data)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return true, err
	}
	return false, err
}

func validateByEmail(reqbody pojo.UserDetailsRequest) (bool, error) {
	dobStr := reqbody.PersonalInfo.DOB
	fmt.Println(dobStr)
	var result []*pojo.UserDetails
	data, err := CollectionUser.Find(ctx, bson.D{{Key: "email", Value: reqbody.Email}})
	if err != nil {
		return false, err
	}
	result, err = convertDbResultIntoUserStruct(data)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return true, err
	}
	return false, err
}

func UploadFile(path string) (string, error) {
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	fileURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	segments := strings.Split(fileURL.Path, "/")
	fileName := segments[len(segments)-1]
	fileName = uploadPath + fileName
	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	resp, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Close()
	size, err := io.Copy(file, resp)
	fmt.Println("size:", size)
	defer file.Close()
	return "File Downloaded successfully", nil
}

func SetValueInModel(req pojo.UserDetailsRequest) (pojo.UserDetails, error) {
	var data pojo.UserDetails
	data.Name = req.Name
	data.Role = req.Role
	data.PhoneNumber = req.PhoneNumber
	data.Email = req.Email
	data.Password = req.Password
	data.PersonalInfo.Name = req.PersonalInfo.Name
	data.PersonalInfo.FatherName = req.PersonalInfo.FatherName
	data.PersonalInfo.Age = req.PersonalInfo.Age
	data.PersonalInfo.DOB = req.PersonalInfo.DOB
	data.PersonalInfo.DocumentType = req.PersonalInfo.DocumentType
	data.PersonalInfo.Address.City = req.PersonalInfo.Address.City
	data.PersonalInfo.Address.Country = req.PersonalInfo.Address.Country
	data.PersonalInfo.Address.State = req.PersonalInfo.Address.State
	data.PersonalInfo.Address.ZipCode = req.PersonalInfo.Address.ZipCode

	data.PersonalInfo.Address.Street = req.PersonalInfo.Address.Street

	data.UploadedDocs.DocumentType = req.UploadedDocs.DocumentType
	data.UploadedDocs.DocumentIdentificationNo = req.UploadedDocs.DocumentIdentificationNo
	data.UploadedDocs.DocumentPath = req.UploadedDocs.DocumentPath

	return data, nil
}
func insertLoginData(email, password, userId string) {

	data, err := CollectionLogin.Find(ctx, bson.D{primitive.E{Key: "email", Value: email}})
	if err != nil {
		log.Println("Unable to fetch data from login details :", err)
	}
	fmt.Println(data)
	finalData, err := convertDbResultIntoLoginStruct(data)
	if err != nil {
		log.Println("Error while converting into login details struct :", err)
	}
	if finalData == nil {
		var request pojo.SignInInput
		request.Email = email
		request.Password = password
		// request.Active = true
		saveData, err := CollectionLogin.InsertOne(ctx, request)
		if err != nil {
			log.Println("Error while inserting into login details :", err)
		}
		fmt.Println("Saved Into Login Details :", saveData.InsertedID)
	} else {
		log.Println("User Already Exists!")
	}
}

func convertDbResultIntoLoginStruct(fetchDataCursor *mongo.Cursor) ([]*pojo.SignInInput, error) {
	var data []*pojo.SignInInput
	for fetchDataCursor.Next(ctx) {
		var db pojo.SignInInput
		err := fetchDataCursor.Decode(&db)
		if err != nil {
			return data, err
		}
		data = append(data, &db)
	}
	return data, nil
}

func convertDbResultIntoUserStruct(fetchDataCursor *mongo.Cursor) ([]*pojo.UserDetails, error) {
	var finaldata []*pojo.UserDetails
	for fetchDataCursor.Next(ctx) {
		var data pojo.UserDetails
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = append(finaldata, &data)
	}
	return finaldata, nil
}

func (e *Connection) UpdateUserDetailsById(reqData pojo.UserDetailsRequest, idStr string) (bson.M, string, error) {
	var updatedDocument bson.M
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return updatedDocument, "Error Occurred", err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	UpdateQuery := bson.D{}
	if reqData.Role != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "role", Value: reqData.Role})
	}
	if reqData.Name != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "name", Value: reqData.Name})
	}
	if reqData.Email != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "email", Value: reqData.Email})
	}
	if reqData.Password != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "password", Value: reqData.Password})
	}
	if reqData.PhoneNumber != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "phone_number", Value: reqData.PhoneNumber})
	}
	if reqData.PersonalInfo.Name != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.name", Value: reqData.PersonalInfo.Name})
	}
	if reqData.PersonalInfo.FatherName != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.father_name", Value: reqData.PersonalInfo.FatherName})
	}
	if reqData.PersonalInfo.Age != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.age", Value: reqData.PersonalInfo.Age})
	}
	if reqData.PersonalInfo.Age != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.dob", Value: reqData.PersonalInfo.DOB})
	}
	if reqData.PersonalInfo.DocumentType != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.document_type", Value: reqData.PersonalInfo.DocumentType})
	}
	if reqData.PersonalInfo.Address.City != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.city", Value: reqData.PersonalInfo.Address.City})
	}
	if reqData.PersonalInfo.Address.Street != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.street", Value: reqData.PersonalInfo.Address.Street})
	}
	if reqData.PersonalInfo.Address.ZipCode != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.zip_code", Value: reqData.PersonalInfo.Address.ZipCode})
	}
	if reqData.PersonalInfo.Address.State != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.state", Value: reqData.PersonalInfo.Address.State})
	}
	if reqData.PersonalInfo.Address.Country != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.country", Value: reqData.PersonalInfo.Address.Country})
	}

	update := bson.D{{Key: "$set", Value: UpdateQuery}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	r := CollectionUser.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if r != nil {
		return updatedDocument, "Error Occurred", r
	}
	fmt.Println(updatedDocument)
	if updatedDocument == nil {
		return updatedDocument, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}

	return updatedDocument, "Document Updated Successfully", nil
}

func convertDate(dateStr string) (time.Time, error) {

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Println(err)
		return date, err
	}
	return date, nil
}

func (e *Connection) VerifyUserDetails(req pojo.VerifyUser, adminMail string) ([]*pojo.UserDetails, string, error) {
	var finalData []*pojo.UserDetails
	var adminData []*pojo.UserDetails

	data, err := CollectionUser.Find(ctx, bson.D{primitive.E{Key: "email", Value: adminMail}})
	adminData, err = convertDbResultIntoUserStruct(data)
	if len(adminData) == 0 {
		return finalData, "", errors.New("Data not present in db acc. to given tokenId")
	}
	filter := bson.D{}
	flag := true
	if req.Id != "" {
		id, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return finalData, "", err
		}
		filter = append(filter, primitive.E{Key: "_id", Value: id})
		flag = false
	}
	if flag {
		if req.Email != "" {
			filter = append(filter, primitive.E{Key: "email", Value: bson.M{"$regex": req.Email}})
			flag = false
		}
	}
	UpdateQuery := bson.D{}
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "is_verified", Value: true})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "verified_by.id", Value: adminData[0].ID})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "verified_by.name", Value: adminData[0].Name})
	update := bson.D{{Key: "$set", Value: UpdateQuery}}

	CollectionUser.FindOneAndUpdate(ctx, filter, update)

	data, err = CollectionUser.Find(ctx, filter)
	if err != nil {
		return finalData, "", err
	}
	finalData, err = convertDbResultIntoUserStruct(data)
	if err != nil {
		return finalData, "", err
	}
	//Send mail method
	return finalData, "Voter verified successfully!", nil
}

func (e *Connection) SearchUsersDetailsById(idStr string) ([]*pojo.UserDetails, string, error) {
	var finalData []*pojo.UserDetails

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return finalData, "Error Occurred", err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	searchData, err := CollectionUser.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoUserStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	return finalData, "Data Fetch Successfully", nil
}

func (e *Connection) SearchFilterUserDetails(search pojo.SearchUser) ([]*pojo.UserDetails, string, error) {
	var finalData []*pojo.UserDetails
	query := bson.D{}

	if search.Id != "" {
		id, err := primitive.ObjectIDFromHex(search.Id)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		query = append(query, primitive.E{Key: "_id", Value: id})
	}
	if search.Role != "" {
		query = append(query, primitive.E{Key: "role", Value: search.Role})
	}
	if search.City != "" {
		query = append(query, primitive.E{Key: "personal_info.address.city", Value: search.City})
	}
	if search.IsVerified != false {
		query = append(query, primitive.E{Key: "is_verified", Value: search.IsVerified})
	}
	if search.State != "" {
		query = append(query, primitive.E{Key: "personal_info.address.state", Value: search.State})
	}
	if search.ZipCode != "" {
		query = append(query, primitive.E{Key: "personal_info.address.zip_code", Value: search.ZipCode})
	}

	searchData, err := CollectionUser.Find(ctx, query)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoUserStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	return finalData, "Data Fetch Successfully", nil
}

func (e *Connection) DeactivateUserDetails(idStr string) (bson.M, string, error) {

	var updatedDocument bson.M
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return updatedDocument, "Error Occurred", err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "is_verified", Value: false}}}}

	r := CollectionUser.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDocument)
	if r != nil {
		return updatedDocument, "Error Occurred", r
	}
	fmt.Println(updatedDocument)
	if updatedDocument == nil {
		return updatedDocument, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}

	return updatedDocument, "User details deactivate successfully!", nil
}

// ==========================candidate details=============================================

func (e *Connection) SaveCandidateDetails(reqBody pojo.CandidatesRequest) ([]*pojo.ElectionDetails, string, error) {
	var finalData []*pojo.ElectionDetails
	electionId, err := primitive.ObjectIDFromHex(reqBody.ElectionId)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	userId, err := primitive.ObjectIDFromHex(reqBody.UserId)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	userData, err := CollectionUser.Find(ctx, bson.D{primitive.E{Key: "_id", Value: userId}})
	if err != nil {
		return finalData, "Error Occurred", err
	}
	userRecord, _ := convertDbResultIntoUserStruct(userData)

	if len(userRecord) == 0 {
		return finalData, "Error Occurred", errors.New("Invalid UserId")
	}
	userName := userRecord[0].Name

	msg, err := UploadFile(reqBody.VoteSign)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", errors.New("Unable to upload file")
	}
	fmt.Println("Upload file:", msg)

	elecData, err := CollectionElection.Find(ctx, bson.D{primitive.E{Key: "_id", Value: electionId}})
	if err != nil {
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoElectionStruct(elecData)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	Candidates := finalData[0].Candidates

	for i := range Candidates {
		userid := Candidates[i].UserId
		canid, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			return finalData, "Error Occurred", err
		}

		if canid == userId {
			return finalData, "Error Occurred", errors.New("Candidate already registered for the given election")
		}
	}

	filter := bson.D{primitive.E{Key: "_id", Value: electionId}}
	UpdateQuery := bson.D{}
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "user_id", Value: userId})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "name", Value: userName})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "commitments", Value: reqBody.Commitments})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "vote_sign", Value: reqBody.VoteSign})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "nomination_status", Value: "not verified"})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "is_nomination_verified", Value: false})

	update := bson.D{{Key: "candidates", Value: UpdateQuery}}
	update = bson.D{{Key: "$push", Value: update}}

	CollectionElection.FindOneAndUpdate(ctx, filter, update)

	fetchData, err := CollectionElection.Find(ctx, filter)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoElectionStruct(fetchData)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	return finalData, "Candidates details saved successfully!", nil
}

func (e *Connection) VerifyCandidate(req pojo.VerifyCandidates, adminMail string) ([]*pojo.ElectionDetails, string, error) {
	var finalData []*pojo.ElectionDetails
	var adminData []*pojo.UserDetails

	data, err := CollectionUser.Find(ctx, bson.D{primitive.E{Key: "email", Value: adminMail}})
	adminData, err = convertDbResultIntoUserStruct(data)
	if len(adminData) == 0 {
		return finalData, "Error Occurred", errors.New("Data not present in db acc. to given tokenId")
	}
	filter := bson.D{}
	electionId, err := primitive.ObjectIDFromHex(req.ElectionId)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return finalData, "Error Occurred", err
	}

	cur, err := CollectionElection.Find(ctx, bson.D{primitive.E{Key: "_id", Value: electionId}})
	if err != nil {
		return finalData, "Error Occurred", errors.New("unable to query db")
	}
	finalData, err = convertDbResultIntoElctionStruct(cur)
	Candidates := finalData[0].Candidates

	for i := range Candidates {
		userid := Candidates[i].UserId
		canid, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		if canid == userId {
			Candidates[i].NominationStatus = "verified"
			Candidates[i].VerifyCandidates.Id = adminData[0].ID
			Candidates[i].VerifyCandidates.Name = adminData[0].Name
		}
	}

	filter = bson.D{primitive.E{Key: "_id", Value: electionId}}
	update := bson.D{primitive.E{Key: "$set", Value: finalData[0]}}
	CollectionElection.FindOneAndUpdate(ctx, filter, update)

	data, err = CollectionElection.Find(ctx, filter)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoElectionStruct(data)
	if err != nil {
		return finalData, "Error Occurred", err
	}

	return finalData, "Candidates verified successfully!", nil
}

//============================election details=======================================

func (e *Connection) SaveElectionDetails(reqBody pojo.ElectionDetailsRequest) ([]*pojo.ElectionDetails, string, error) {
	var finalData []*pojo.ElectionDetails
	setData, err := SetValueInElection(reqBody)
	if err != nil {
		return finalData, "Error occurred", err
	}
	insert, err := CollectionElection.InsertOne(ctx, setData)
	if err != nil {
		return finalData, "Error occurred", err
	}
	fetchData, err := CollectionElection.Find(ctx, bson.D{primitive.E{Key: "_id", Value: insert.InsertedID}})

	finalData, err = convertDbResultIntoElectionStruct(fetchData)
	if err != nil {
		return finalData, "Error occurred", err
	}
	return finalData, "Election details saved successfully!", nil
}

func SetValueInElection(req pojo.ElectionDetailsRequest) (pojo.ElectionDetails, error) {
	var data pojo.ElectionDetails
	data.ElectionDate = req.ElectionDate
	data.ElectionStatus = req.ElectionStatus
	data.Result = req.Result
	data.ResultDate = req.ResultDate
	data.Location = req.Location
	return data, nil
}

func convertDbResultIntoElectionStruct(fetchDataCursor *mongo.Cursor) ([]*pojo.ElectionDetails, error) {
	var finaldata []*pojo.ElectionDetails
	for fetchDataCursor.Next(ctx) {
		var data pojo.ElectionDetails
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = append(finaldata, &data)
	}
	return finaldata, nil
}

func (e *Connection) FindElectionById(idStr string) ([]*pojo.ElectionDetails, string, error) {
	var finalData []*pojo.ElectionDetails

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return finalData, "Error Occurred", err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	searchData, err := CollectionElection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoElectionStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	return finalData, "Data Fetched Successfully", nil
}

func (e *Connection) UpdateElectionDetailsById(reqData pojo.ElectionDetailsRequest, idStr string) (bson.M, string, error) {
	var updatedDocument bson.M
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return updatedDocument, "Error Occurred", err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	UpdateQuery := bson.D{}
	if reqData.ElectionDate != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "election_date", Value: reqData.ElectionDate})
	}
	if reqData.ElectionStatus != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "election_status", Value: reqData.ElectionStatus})
	}
	if reqData.ResultDate != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "result_date", Value: reqData.ResultDate})
	}
	if reqData.Location != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "location", Value: reqData.Location})
	}

	update := bson.D{{Key: "$set", Value: UpdateQuery}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	r := CollectionElection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if r != nil {
		return updatedDocument, "Error Occurred", r
	}
	fmt.Println(updatedDocument)
	if updatedDocument == nil {
		return updatedDocument, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}

	return updatedDocument, "Document Updated Successfully", nil
}

func (e *Connection) SearchFilterElectionDetails(search pojo.SearchFilterElectionReq) ([]*pojo.ElectionDetails, string, error) {
	var finalData []*pojo.ElectionDetails
	query := bson.D{}

	if search.Id != "" {
		id, err := primitive.ObjectIDFromHex(search.Id)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		query = append(query, primitive.E{Key: "_id", Value: id})
	}
	if search.ElectionDate != "" {
		query = append(query, primitive.E{Key: "election_date", Value: search.ElectionDate})
	}
	if search.ElectionStatus != "" {
		query = append(query, primitive.E{Key: "election_status", Value: search.ElectionStatus})
	}
	if search.CandidateName != "" {
		query = append(query, primitive.E{Key: "candidates.$[].name", Value: search.CandidateName})
	}
	if search.ResultDate != "" {
		query = append(query, primitive.E{Key: "result_date", Value: search.ResultDate})
	}
	if search.Location != "" {
		query = append(query, primitive.E{Key: "location", Value: search.Location})
	}

	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				query,
				bson.D{{Key: "election_status", Value: bson.M{"$ne": "Deactivated"}}},
			},
		},
	}

	searchData, err := CollectionElection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoElctionStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	return finalData, "Data Fetch Successfully", nil
}

func convertDbResultIntoElctionStruct(fetchDataCursor *mongo.Cursor) ([]*pojo.ElectionDetails, error) {
	var finaldata []*pojo.ElectionDetails
	for fetchDataCursor.Next(ctx) {
		var data pojo.ElectionDetails
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = append(finaldata, &data)
	}
	return finaldata, nil
}

func convertDbResultIntoElction(fetchDataCursor *mongo.Cursor) (*pojo.ElectionDetails, error) {
	var finaldata *pojo.ElectionDetails
	for fetchDataCursor.Next(ctx) {
		var data pojo.ElectionDetails
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = &data
	}
	return finaldata, nil
}

func (e *Connection) DeactivateElection(idStr string) (bson.M, string, error) {
	var updatedDocument bson.M
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return updatedDocument, "Error Occurred", err
	}
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "_id", Value: id}},
				bson.D{{Key: "election_status", Value: bson.M{"$ne": "Deactivated"}}},
			},
		},
	}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "election_status", Value: "Deactivated"}}}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	r := CollectionElection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if r != nil {
		return updatedDocument, "Error Occurred", r
	}

	if updatedDocument == nil {
		return updatedDocument, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}

	return updatedDocument, "Election details deactivated successfully!", nil
}

// ======================================Token=============================================
func (e *Connection) GenerateToken(request pojo.SignInInput) (string, string, error) {

	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "email", Value: request.Email}},
				bson.D{{Key: "password", Value: request.Password}},
			},
		},
	}

	// check if email exists and password is correct
	record, err := CollectionUser.Find(ctx, filter)
	if err != nil {
		return "", "Error Occurred", err
	}

	convertData, err := convertDbResultIntoLoginStruct(record)
	if err != nil {
		return "", "Error Occurred", err
	}

	if len(convertData) != 0 {
		tokenString, err := auth.GenerateJWT(request.Email)
		if err != nil {
			return "", "Error Occurred", err
		}
		return tokenString, "Login Successfully", err
	} else {
		return "", "Error Occurred", errors.New("Invalid Credentials")
	}
}

func (e *Connection) FetchRole(mailId string) string {
	data, err := CollectionUser.Find(ctx, bson.D{primitive.E{Key: "email", Value: mailId}})
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	convData, err := convertDbResultIntoUserStruct(data)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	if len(convData) == 0 {
		return ""
	}
	return convData[0].Role
}

// ======================================vote Result======================================
func (e *Connection) CastVote(reqData pojo.CastVote, mail string) ([]*pojo.ElectionDetails, string, error) {
	var finalData []*pojo.ElectionDetails
	var userData []*pojo.UserDetails
	var msg = ""
	data, err := CollectionUser.Find(ctx, bson.D{primitive.E{Key: "email", Value: mail}})
	if err != nil {
		return finalData, "Error Occurred", err
	}
	userData, err = convertDbResultIntoUserStruct(data)
	if len(userData) == 0 {
		return finalData, "", errors.New("Data not present in db acc. to given tokenId")
	}
	userId := userData[0].ID

	if (reqData.ElectionId != "") || (reqData.CandidateId != "") {
		electionid, err := primitive.ObjectIDFromHex(reqData.ElectionId)
		if err != nil {
			return finalData, "Error Occurred", err
		}

		candidateId, err := primitive.ObjectIDFromHex(reqData.CandidateId)
		if err != nil {
			return finalData, "Error Occurred", err
		}

		userFound, err := fetchElectionId(electionid, userId)
		fmt.Println("userFound:", userFound)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		if userFound {
			return finalData, "Error Occurred", errors.New("User already voted")
		}
		msg, err = addVoteInElectionDB(electionid, candidateId, userId)
	}

	updateUserDetails(userId, reqData)
	return finalData, msg, nil
}

func updateUserDetails(userId primitive.ObjectID, reqData pojo.CastVote) (*pojo.UserDetails, string, error) {
	var finalData *pojo.UserDetails
	var err error

	electionid, err := primitive.ObjectIDFromHex(reqData.ElectionId)
	if err != nil {
		return finalData, "Error Occurred", err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: userId}}
	UpdateQuery := bson.D{}

	UpdateQuery = append(UpdateQuery, primitive.E{Key: "election_id", Value: electionid})
	// update := bson.D{{Key: "$set", Value: UpdateQuery}}
	update := bson.D{{Key: "voted", Value: UpdateQuery}}

	update = bson.D{{Key: "$push", Value: update}}

	CollectionUser.FindOneAndUpdate(ctx, filter, update)
	data, err := CollectionUser.Find(ctx, filter)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoUser(data)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}

	return finalData, "Cast Vote by User", err
}

func fetchElectionId(electionid, userId primitive.ObjectID) (bool, error) {
	var data []*pojo.UserDetails
	query := bson.D{}
	// electionid, err := primitive.ObjectIDFromHex(req.ElectionId)
	// if err != nil {
	// 	return false, err
	// }
	// userid, err := primitive.ObjectIDFromHex(req.CandidateId)
	// if err != nil {
	// 	return false, err
	// }
	query = append(query, primitive.E{Key: "_id", Value: userId}, primitive.E{Key: "voted.election_id", Value: electionid})

	searchData, err := CollectionUser.Find(ctx, query)
	if err != nil {
		log.Println(err)
		return false, err
	}
	data, err = convertDbResultIntoUserStruct(searchData)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if len(data) != 0 {
		return true, err

	}
	return false, err
}
func convertDbResultIntoUser(fetchDataCursor *mongo.Cursor) (*pojo.UserDetails, error) {
	var finaldata *pojo.UserDetails
	for fetchDataCursor.Next(ctx) {
		var data pojo.UserDetails
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = &data
	}
	return finaldata, nil
}

func (e *Connection) SearchCandidateList(reqData pojo.SearchElectionResult) ([]pojo.CandidatesDetails, string, error) {
	var finalData []pojo.CandidatesDetails
	var err error

	query := bson.D{}

	if reqData.ElectionId != "" {
		id, err := primitive.ObjectIDFromHex(reqData.ElectionId)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		query = append(query, primitive.E{Key: "_id", Value: id})
	}
	if reqData.ElectionDate != "" {
		query = append(query, primitive.E{Key: "election_date", Value: reqData.ElectionDate})
	}
	if reqData.ElectionStatus != "" {
		query = append(query, primitive.E{Key: "election_status", Value: reqData.ElectionStatus})
	}
	if reqData.ResultDate != "" {
		query = append(query, primitive.E{Key: "result_date", Value: reqData.ResultDate})
	}
	if reqData.Location != "" {
		query = append(query, primitive.E{Key: "location", Value: reqData.Location})
	}

	searchData, err := CollectionElection.Find(ctx, query)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = convertDbResultIntoCandidateStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	return finalData, "Candidate Data Fetch Successfully", nil

}

func convertDbResultIntoCandidateStruct(fetchDataCursor *mongo.Cursor) ([]pojo.CandidatesDetails, error) {
	var finaldata []pojo.CandidatesDetails
	var eledata []*pojo.ElectionDetails
	for fetchDataCursor.Next(ctx) {
		var data pojo.ElectionDetails
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		eledata = append(eledata, &data)
	}
	for _, value := range eledata {
		fmt.Println(value)
		finaldata = value.Candidates
	}
	return finaldata, nil
}

func convertDbResultIntoCandidateResult(fetchDataCursor *mongo.Cursor) ([]pojo.CandidatesDetails, error) {
	var finaldata []pojo.CandidatesDetails
	var eledata []*pojo.ElectionDetails
	for fetchDataCursor.Next(ctx) {
		var data pojo.ElectionDetails
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		eledata = append(eledata, &data)
	}

	for _, value := range eledata {
		fmt.Println(value)
		finaldata = value.Candidates

	}
	for _, value := range finaldata {
		fmt.Println(value)
		voteCou := value.Votecount
		fmt.Println(voteCou)
	}
	return finaldata, nil
}

func addVoteInElectionDB(electionid, candidateId, userId primitive.ObjectID) (string, error) {
	var doc bson.M
	query := bson.D{}
	query = append(query, primitive.E{Key: "_id", Value: electionid})

	eleData, err := CollectionElection.Find(ctx, query)
	if err != nil {
		log.Println(err)
		return "Error Occurred", err
	}
	convData, err := convertDbResultIntoElctionStruct(eleData)
	if err != nil {
		log.Println(err)
		return "Error Occurred", err
	}
	candidateRecord := convData[0].Candidates
	datestr := convData[0].ElectionDate
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")
	fmt.Println("Current Time in String: ", currentTime.Format("2006-01-02"))
	if datestr == currentDate {
		for i := range candidateRecord {
			userid := candidateRecord[i].UserId
			canid, err := primitive.ObjectIDFromHex(userid)
			if err != nil {
				return "Error Occurred", err
			}
			if canid == candidateId {
				if candidateRecord[i].NominationStatus == "verified" {
					candidateRecord[i].Votecount = candidateRecord[i].Votecount + 1
				} else {
					return "Error Occurred", errors.New("Invalid Candidate")
				}
			}

		}
	}

	if datestr > currentDate {
		return "Election is not started", nil
	}
	if datestr < currentDate {
		return "Election is expired", nil
	}

	update := bson.D{primitive.E{Key: "$set", Value: convData[0]}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	CollectionElection.FindOneAndUpdate(ctx, query, update, opts).Decode(&doc)

	if doc == nil {
		return "Error Occurred while voting", errors.New("Unable to vote")
	}
	return "Voted Successfully", nil
}

func (e *Connection) ElectionResultById(idStr string) ([]pojo.CandidatesDetails, string, error) {
	var updatedDocument []*pojo.ElectionDetails
	var candidate []pojo.CandidatesDetails
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return candidate, "Error Occurred", err
	}
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "_id", Value: id}},
				bson.D{{Key: "election_status", Value: bson.M{"$eq": "Completed"}}},
			},
		},
	}
	data, err := CollectionElection.Find(ctx, filter)
	if err != nil {
		return candidate, "Error Occurred", err
	}
	fmt.Println("data:", data)
	updatedDocument, err = convertDbResultIntoElctionStruct(data)
	if updatedDocument == nil {
		return candidate, "Result not fetch successfully", errors.New("Given id is not valid or the result date is not today")
	}
	candidates := updatedDocument[0].Candidates

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Votecount > candidates[j].Votecount
	})
	counter := 0
	UpdateQuery := bson.D{}
	fmt.Println(candidates)
	for i := range candidates {
		if candidates[i].NominationStatus == "verified" {
			if counter == 0 {
				UpdateQuery = append(UpdateQuery, primitive.E{Key: "result.name", Value: candidates[i].Name})
				UpdateQuery = append(UpdateQuery, primitive.E{Key: "result.vote", Value: candidates[i].Votecount})
			}
			counter = counter + 1
			candidate = append(candidate, candidates[i])
		}
	}
	filter = bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: UpdateQuery}}
	CollectionElection.FindOneAndUpdate(ctx, filter, update)

	data, err = CollectionElection.Find(ctx, filter)
	if err != nil {
		return candidate, "Error Occurred", err
	}
	return candidate, "Result fetch successfully", nil
}
