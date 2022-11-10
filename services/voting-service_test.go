package services

import (
	"onlineVoting/pojo"
	"testing"
)

func TestCases(t *testing.T) {
	uploadDocumentsDetails := pojo.UploadDocumentsDetails{DocumentType: "", DocumentIdentificationNo: "", DocumentPath: "D:/SanDiskMemoryZone_QuickStartGuide.pdf"}
	got, _ := UploadFile(uploadDocumentsDetails.DocumentPath)
	want := "File Downloaded successfully"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
