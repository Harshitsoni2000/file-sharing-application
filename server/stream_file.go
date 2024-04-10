package main

import (
	"errors"
	pb "github.com/Harshitsoni2000/File_Sharing_Application/server/proto"
	"io"
	"log"
	"os"
	"time"
)

func (fs *FileServer) DownloadFile(dfr *pb.DownloadFileRequest, stream pb.FileService_DownloadFileServer) error {
	fileName := dfr.FileName
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)

	if errors.Is(err, os.ErrNotExist) {
		log.Println("File does not exist at the given path!")
		return err
	} else if err != nil {
		log.Println("Error while opening file :: " + err.Error())
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error while closing file :: ", err.Error())
			return
		}
	}(file)

	start := time.Now()
	buffer := make([]byte, chunkSize)
	for {
		n, err := file.Read(buffer)

		if n == 0 || err == io.EOF {
			log.Println("Reached EOF")
			break
		} else if err != nil {
			log.Println("Error while reading file :: " + err.Error())
			return err
		}
		if err := stream.Send(&pb.FileChunk{ChunkData: buffer}); err != nil {
			log.Println("Error while sending file chunk :: " + err.Error())
			return err
		}
	}
	log.Printf("File %s Sent in Time :: %f seconds\n", dfr.FileName, (float64(time.Now().UnixMilli())-float64(start.UnixMilli()))/1000.0)
	return nil
}
