package main

import (
	pb "github.com/Harshitsoni2000/File_Sharing_Application/server/proto"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func sendEntireDir(filePath string, stream pb.FileService_DownloadDirServer) error {
	dirEntries, _ := os.ReadDir(filePath)

	for _, dirEntry := range dirEntries {
		nameWithPath := filePath + "/" + dirEntry.Name()
		if strings.Contains(nameWithPath, "DS_Store") {
			continue
		}
		log.Println(nameWithPath)
		if dirEntry.IsDir() {
			if err := stream.Send(&pb.DirChunk{DirName: &nameWithPath}); err != nil {
				log.Println("Error while sending dir chunk :: " + err.Error())
				return err
			}

			if err := sendEntireDir(nameWithPath, stream); err != nil {
				return err
			}
		} else {
			file, err := os.OpenFile(nameWithPath, os.O_RDONLY, 0644)

			if err != nil {
				log.Println("Error while opening file :: " + err.Error())
				return err
			}

			start := time.Now()
			buffer := make([]byte, chunkSize)
			for {
				n, err := file.Read(buffer)
				if n == 0 || err == io.EOF {
					break
				} else if err != nil {
					log.Println("Error while reading file :: " + err.Error())
					if err := file.Close(); err != nil {
						return err
					}
					return err
				}
				if err := stream.Send(&pb.DirChunk{FileName: &nameWithPath, ChunkData: buffer}); err != nil {
					log.Println("Error while sending file chunk :: " + err.Error())
					if err := file.Close(); err != nil {
						return err
					}
					return err
				}
			}
			if err := file.Close(); err != nil {
				return err
			}
			log.Printf("File %s Sent in Time :: %f seconds\n", dirEntry.Name(), (float64(time.Now().UnixMilli())-float64(start.UnixMilli()))/1000.0)
		}
	}
	return nil
}

func (fs *FileServer) DownloadDir(ddr *pb.DownloadDirRequest, stream pb.FileService_DownloadDirServer) error {
	dirName := ddr.DirName

	_, err := os.ReadDir(dirName)
	if err != nil {
		log.Println("Error while opening Dir :: " + err.Error())
		return err
	}

	pwd := strings.Split(dirName, "/")
	dir := pwd[len(pwd)-1]

	filePath := "./" + dir
	if err := stream.Send(&pb.DirChunk{DirName: &filePath}); err != nil {
		log.Println("Error while sending dir name :: " + err.Error())
		return err
	}

	if err := sendEntireDir(filePath, stream); err != nil {
		return err
	}
	return nil
}
