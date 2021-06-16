package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	mainpb "github.com/LibenHailu/grpc_file_stream/file_stream/filepb"
	"github.com/LibenHailu/grpc_file_stream/file_stream/service/client"
	pb "github.com/LibenHailu/peer_to_peer_file_share/peer-copy/filepb"
)

var (
	serverAddr *string
	server     *string
)

func InitFileClient(srvAddr *string, svr *string) {
	serverAddr = srvAddr
	server = svr
}

func connect(srvAddr *string) *grpc.ClientConn {
	conn, err := grpc.Dial(*srvAddr, grpc.WithInsecure())

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	grpclog.Println("client started...")
	return conn

}

func DownloadFile(fileName string) {

	//get connection
	conn := connect(serverAddr)
	fmt.Println(*serverAddr)
	defer conn.Close()

	if *server == "server" {

		c := mainpb.NewFileServiceClient(conn)
		req := &mainpb.ServeFileRequest{
			FileName: fileName,
		}
		resStream, err := c.DownloadFile(context.Background(), req)
		if err != nil {
			log.Fatalf("error downloading file: %v", err)
		}
		fileData := bytes.Buffer{}
		fileSize := 0
		for {
			msg, err := resStream.Recv()

			if err == io.EOF {
				// we've reached the end of stream
				log.Println("recived all chunks")
				break
			}
			if err != nil {
				log.Fatalf("error while reciving chunk %v", err)
			}
			// log.Printf("Response from GreetManyTimes: %v ", msg.ChunkData)
			chunk := msg.GetChunkData()
			size := len(chunk)

			fileSize += size

			// if fileSize > maxFileSize {
			// 	return logError(status.Errorf(codes.InvalidArgument, "file is too large: %d > %d", fileSize, maxFileSize))
			// }

			_, err = fileData.Write(chunk)
			if err != nil {
				log.Fatal("couldn't write chunk data: %v", err)
			}

		}

		clientSave := client.NewDiskFileStore("C:/Users/Liben/go/src/github.com/LibenHailu/peer_to_peer_file_share/peer/file")
		clientSave.Save(fileData, fileName)

	} else {

		c := pb.NewFileServiceClient(conn)
		req := &pb.ServeFileRequest{
			FileName: fileName,
		}
		resStream, err := c.DownloadFile(context.Background(), req)
		if err != nil {
			log.Fatalf("error downloading file: %v", err)
		}
		fileData := bytes.Buffer{}
		fileSize := 0
		for {
			msg, err := resStream.Recv()

			if err == io.EOF {
				// we've reached the end of stream
				log.Println("recived all chunks")
				break
			}
			if err != nil {
				log.Fatalf("error while reciving chunk %v", err)
			}
			log.Printf("Response from GreetManyTimes: %v ", msg.ChunkData)
			chunk := msg.GetChunkData()
			size := len(chunk)

			fileSize += size

			// if fileSize > maxFileSize {
			// 	return logError(status.Errorf(codes.InvalidArgument, "file is too large: %d > %d", fileSize, maxFileSize))
			// }

			_, err = fileData.Write(chunk)
			if err != nil {
				log.Fatal("couldn't write chunk data: %v", err)
			}

		}

		clientSave := client.NewDiskFileStore("C:/Users/Liben/go/src/github.com/LibenHailu/peer_to_peer_file_share/peer/file")
		clientSave.Save(fileData, fileName)

	}
}
