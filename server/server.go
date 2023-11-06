package server

import (
	"homework/jsonutil"
	pb "homework/pb"
	"io"
	"log"
	"os"
	"path/filepath"
)

const uploadDirectory = "./uploads/"

type FileUploadServer struct {
	pb.UnimplementedFileUploadServer
}

func (s *FileUploadServer) UploadFile(stream pb.FileUpload_UploadFileServer) error {

	if _, err := os.Stat(uploadDirectory); os.IsNotExist(err) {
		os.MkdirAll(uploadDirectory, os.ModePerm)
	}

	var fileName string
	var fullFilePath string
	var file *os.File

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileName = chunk.FileName
		fullFilePath = filepath.Join(uploadDirectory, fileName)

		if file == nil {
			var err error
			file, err = getFileForOutput(fullFilePath)
			if err != nil {
				return err
			}
			defer file.Close()
		}

		data := chunk.Data
		_, writeErr := file.Write(data)
		if writeErr != nil {
			return writeErr
		}
	}

	newFileName := generateModifiedFileName(fileName)
	fullNewFilePath := filepath.Join(uploadDirectory, newFileName)

	if err := applyJSONTransformations(fullFilePath, fullNewFilePath); err != nil {
		log.Printf("Error processing '%s' as JSON file: %v\n", fileName, err)
	}

	response := &pb.UploadResponse{
		Success: true,
	}
	return stream.SendAndClose(response)
}

func getFileForOutput(fullFilePath string) (*os.File, error) {
	var openMode int
	var file *os.File

	if _, fileErr := os.Stat(fullFilePath); fileErr == nil {
		if removeErr := os.Remove(fullFilePath); removeErr != nil {
			return nil, removeErr
		}
		openMode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	} else {
		openMode = os.O_CREATE | os.O_WRONLY
	}

	var openErr error
	file, openErr = os.OpenFile(fullFilePath, openMode, 0644)
	if openErr != nil {
		return nil, openErr
	}
	return file, nil
}

func applyJSONTransformations(inputFilePath, outputFilePath string) error {
	inputData, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	updatedJSON, err := jsonutil.UpdateJSON(inputData)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFilePath, updatedJSON, 0644); err != nil {
		return err
	}

	return nil
}

func generateModifiedFileName(fileName string) string {
	baseName := filepath.Base(fileName)
	ext := filepath.Ext(fileName)
	modifiedFileName := baseName[:len(baseName)-len(ext)] + "_modified" + ext
	return modifiedFileName
}
