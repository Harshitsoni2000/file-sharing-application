package cmd

import (
	"context"
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

// fileNameCmd represents the fileName command
var fileNameCmd = &cobra.Command{
	Use:   "get-file",
	Short: "Pass the File Name",
	Args:  cobra.ExactArgs(1),
	Long:  `Pass the name of file that you want to receive from the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		serverIPFlag, _ := cmd.Flags().GetString("server")
		if !strings.Contains(serverIPFlag, ":50051") {
			serverIPFlag = serverIPFlag + ":50051"
		}

		if len(filePath) <= 3 {
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

		pwd := strings.Split(filePath, "/")
		fileName := pwd[len(pwd)-1]

		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)

		if err != nil {
			log.Fatalf("Error while creating file on Local :: %s\n", err.Error())
			return
		}

		start := time.Now()

		client := pb.NewFileServiceClient(conn)
		dfr := &pb.DownloadFileRequest{FileName: fileName}
		stream, err := client.DownloadFile(context.Background(), dfr)

		if err != nil {
			log.Fatalf("Error while server streaming :: %s\n", err.Error())
		}

		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			if _, err := file.Write(chunk.ChunkData); err != nil {
				log.Fatalf("Error while writing chunk to file :: %s\n", err.Error())
			}
		}
		log.Printf("File Downloaded successfully in Time :: %f seconds\n", (float64(time.Now().UnixMilli())-float64(start.UnixMilli()))/1000.0)
	},
}

func init() {
	fileNameCmd.Flags().StringP("server", "s", "127.0.0.1", "Use -s and provide the server ip")
	rootCmd.AddCommand(fileNameCmd)
}
