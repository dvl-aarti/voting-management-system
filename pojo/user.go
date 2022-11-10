package pojo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDetails struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string                 `bson:"name,omitempty" json:"name" validate:"required,min=2,max=100"`
	Role         string                 `bson:"role,omitempty" json:"role,omitempty"`
	Email        string                 `bson:"email,omitempty" json:"email" validate:"email,required"`
	Password     string                 `bson:"password,omitempty" json:"password" validate:"required"`
	PhoneNumber  string                 `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	PersonalInfo PersonalInfoDetails    `bson:"personal_info,omitempty" json:"personal_info,omitempty"`
	IsVerified   bool                   `bson:"is_verified,omitempty" json:"is_verified"`
	VerifiedBy   VerifyDetails          `bson:"verified_by,omitempty" json:"verified_by,omitempty"`
	UploadedDocs UploadDocumentsDetails `bson:"uploaded_docs,omitempty" json:"uploaded_docs,omitempty"`
	Voted        []Voted                `bson:"voted,omitempty" json:"voted,omitempty"`
}

type UserDetailsRequest struct {
	Name         string                 `bson:"name,omitempty" json:"name" validate:"required,min=2,max=100"`
	Role         string                 `bson:"role,omitempty" json:"role,omitempty"`
	Email        string                 `bson:"email,omitempty" json:"email" validate:"email,required"`
	Password     string                 `bson:"password,omitempty" json:"password" validate:"required"`
	PhoneNumber  string                 `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	PersonalInfo PersonalInfoDetails    `bson:"personal_info,omitempty" json:"personal_info,omitempty"`
	IsVerified   bool                   `bson:"is_verified,omitempty" json:"is_verified"`
	VerifiedBy   VerifyDetails          `bson:"verified_by,omitempty" json:"verified_by,omitempty"`
	UploadedDocs UploadDocumentsDetails `bson:"uploaded_docs,omitempty" json:"uploaded_docs,omitempty"`
}

type PersonalInfoDetails struct {
	Name         string         `bson:"name,omitempty" json:"name,omitempty"`
	FatherName   string         `bson:"father_name,omitempty" json:"father_name,omitempty"`
	DOB          string         `bson:"dob,omitempty" json:"dob,omitempty"`
	Age          string         `bson:"age,omitempty" json:"age,omitempty"`
	VoterId      string         `bson:"voter_id,omitempty" json:"voter_id,omitempty"`
	DocumentType string         `bson:"document_type,omitempty" json:"document_type,omitempty"`
	Address      AddressDetails `bson:"address,omitempty" json:"address,omitempty"`
}

type AddressDetails struct {
	Street  string `bson:"street,omitempty" json:"street,omitempty"`
	City    string `bson:"city,omitempty" json:"city,omitempty"`
	State   string `bson:"state,omitempty" json:"state,omitempty"`
	ZipCode string `bson:"zip_code,omitempty" json:"zip_code,omitempty"`
	Country string `bson:"country,omitempty" json:"country,omitempty"`
}

type VerifyDetails struct {
	Id   primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name,omitempty" json:"name,omitempty"`
}

type resultDetails struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	Vote int64  `bson:"vote,omitempty" json:"vote,omitempty"`
}

type Voted struct {
	ElectionId string `bson:"election_id,omitempty" json:"election_id,omitempty"`
}

type UploadDocumentsDetails struct {
	DocumentType             string `bson:"document_type,omitempty" json:"document_type,omitempty"`
	DocumentIdentificationNo string `bson:"document_identification_no,omitempty" json:"document_identification_no,omitempty"`
	DocumentPath             string `bson:"document_path,omitempty" json:"document_path,omitempty"`
}

type SignInInput struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `bson:"email,omitempty" json:"email,omitempty"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
	UserId   string             `bson:"user_id,omitempty" json:"user_id"`
}

type Response struct {
	Success    string      `json:"success,omitempty"`
	SucessCode string      `json:"successCode,omitempty"`
	SucessMsg  string      `json:"successMsg,omitempty"`
	Response   interface{} `json:"response,omitempty"`
}

type VerifyUser struct {
	Id         string `bson:"id,omitempty" json:"id,omitempty"`
	Role       string `bson:"role,omitempty" json:"role,omitempty"`
	Name       string `bson:"name,omitempty" json:"name,omitempty"`
	Email      string `bson:"email,omitempty" json:"email"`
	IsVerified bool   `bson:"is_verified,omitempty" json:"is_verified"`
}

type SearchUser struct {
	Id         string `bson:"_id,omitempty" json:"id,omitempty"`
	Role       string `bson:"role,omitempty" json:"role,omitempty"`
	City       string `bson:"city,omitempty" json:"city,omitempty"`
	State      string `bson:"state,omitempty" json:"state,omitempty"`
	ZipCode    string `bson:"zip_code,omitempty" json:"zip_code,omitempty"`
	Country    string `bson:"country,omitempty" json:"country,omitempty"`
	IsVerified bool   `bson:"is_verified,omitempty" json:"is_verified,omitempty"`
}

type ElectionDetails struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	ElectionDate   string              `bson:"election_date,omitempty" json:"election_date,omitempty"`
	ResultDate     string              `bson:"result_date,omitempty" json:"result_date,omitempty"`
	ElectionStatus string              `bson:"election_status,omitempty" json:"election_status,omitempty"`
	Result         resultDetails       `bson:"result,omitempty" json:"result" validate:"required"`
	Location       string              `bson:"location,omitempty" json:"location,omitempty"`
	Candidates     []CandidatesDetails `bson:"candidates,omitempty" json:"candidates,omitempty"`
}

type ElectionDetailsRequest struct {
	ElectionDate   string        `bson:"election_date,omitempty" json:"election_date,omitempty"`
	ResultDate     string        `bson:"result_date,omitempty" json:"result_date,omitempty"`
	ElectionStatus string        `bson:"election_status,omitempty" json:"election_status,omitempty"`
	Result         resultDetails `bson:"result,omitempty" json:"result" validate:"required"`
	Location       string        `bson:"location,omitempty" json:"location,omitempty"`
}

type CandidatesDetails struct {
	UserId           string        `bson:"user_id" json:"user_id"`
	Name             string        `bson:"name" json:"name"`
	Commitments      []string      `bson:"commitments" json:"commitments"`
	Votecount        int64         `bson:"vote_count" json:"vote_count"`
	VoteSign         string        `bson:"vote_sign" json:"vote_sign"`
	NominationStatus string        `bson:"nomination_status" json:"nomination_status"`
	VerifyCandidates VerifyDetails `bson:"nomaination_verified_by,omitempty" json:"nomaination_verified_by,omitempty"`
}
type CandidatesRequest struct {
	ElectionId       string        `bson:"election_id,omitempty" json:"election_id,omitempty"`
	UserId           string        `bson:"user_id,omitempty" json:"user_id,omitempty"`
	CandidateName    string        `bson:"candidate_name,omitempty" json:"candidate_name,omitempty"`
	Commitments      []string      `bson:"commitments,omitempty" json:"commitments,omitempty"`
	VoteSign         string        `bson:"vote_sign,omitempty" json:"vote_sign,omitempty"`
	NominationStatus string        `bson:"nomination_status" json:"nomination_status"`
	VerifyCandidates VerifyDetails `bson:"nomaination_verified_by,omitempty" json:"nomaination_verified_by,omitempty"`
}

type VerifyCandidates struct {
	ElectionId string `bson:"election_id,omitempty" json:"election_id,omitempty"`
	UserId     string `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

type CastVote struct {
	ElectionId  string `bson:"election_id,omitempty" json:"election_id,omitempty"`
	CandidateId string `bson:"candidate_id,omitempty" json:"candidate_id,omitempty"`
}

type SearchFilterElectionReq struct {
	Id             string `bson:"_id,omitempty" json:"_id,omitempty"`
	Location       string `bson:"location,omitempty" json:"location,omitempty"`
	ElectionDate   string `bson:"election_date,omitempty" json:"election_date,omitempty"`
	ResultDate     string `bson:"result_date,omitempty" json:"result_date,omitempty"`
	ElectionStatus string `bson:"election_status,omitempty" json:"election_status,omitempty"`
	CandidateName  string `bson:"candidate_name,omitempty" json:"candidate_name,omitempty"`
}

type SearchElectionResult struct {
	ElectionId     string `bson:"election_id,omitempty" json:"election_id,omitempty"`
	Location       string `bson:"location,omitempty" json:"location,omitempty"`
	ElectionDate   string `bson:"election_date,omitempty" json:"election_date,omitempty"`
	ResultDate     string `bson:"result_date,omitempty" json:"result_date,omitempty"`
	ElectionStatus string `bson:"election_status,omitempty" json:"election_status,omitempty"`
}
