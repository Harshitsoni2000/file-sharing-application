package cmd

import (
	"context"
	"errors"
	pb "github.com/Harshitsoni2000/File_Sharing_Application/client/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var dirNameCmd = &cobra.Command{
	Use:   "get-dir",
	Short: "Pass the Directory Name",
	Args:  cobra.ExactArgs(1),
	Long: `Pass the name of directory that you want to receive from the server.
You will also receive all subdirectories and files present in them.`,
	Run: func(cmd *cobra.Command, args []string) {
		dirName := args[0]
		serverIPFlag, _ := cmd.Flags().GetString("server")
		if !strings.Contains(serverIPFlag, ":50051") {
			serverIPFlag = serverIPFlag + ":50051"
		}

		if len(dirName) <= 3 {
			log.Fatalf("Provide either fileName or dirName \n")
		}
		conn, err := grpc.Dial(serverIPFlag, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Fatalf("Error while connecting to GRPC Server : %s\n", err.Error())
		}

		defer func(conn *grpc.ClientConn) {
			if err := conn.Close(); err != nil {
				log.Fatal("Error while closing connection : ", err.Error())
			}
		}(conn)
		start := time.Now()

		client := pb.NewFileServiceClient(conn)
		ddr := &pb.DownloadDirRequest{DirName: dirName}
		stream, err := client.DownloadDir(context.Background(), ddr)
		if err != nil {
			log.Fatalf("Error while connecting to stream :: %s\n", err.Error())
		}

		for {
			dirChunk, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while reading from Stream :: %v\n", err)
			}

			if dirChunk.DirName != nil {
				dirName := *dirChunk.DirName
				if err := os.Mkdir(dirName, 0744); !errors.Is(err, os.ErrExist) && err != nil {
					log.Fatalf("Error while creating directory on Local :: %s\n", err.Error())
					return
				}
			} else {
				fileName := *dirChunk.FileName
				file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
				if err != nil {
					log.Fatalf("Error while opening file :: %s\n", err.Error())
				}

				if _, err := file.Write(dirChunk.ChunkData); err != nil {
					log.Fatalf("Error while writing chunk to file :: %s\n", err.Error())
				}
			}
		}
		log.Printf("Directory Downloaded successfully in Time :: %f seconds\n", (float64(time.Now().UnixMilli())-float64(start.UnixMilli()))/1000.0)
	},
}

func init() {
	dirNameCmd.Flags().StringP("server", "s", "127.0.0.1", "Use -s and provide the server ip")
	rootCmd.AddCommand(dirNameCmd)
}
