package server

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	mclient "github.com/LibenHailu/grpc_file_stream/file_stream/file_client"
	mainpb "github.com/LibenHailu/grpc_file_stream/file_stream/filepb"
	pb "github.com/LibenHailu/peer_to_peer_file_share/peer-copy/filepb"
)

var (
	port       = flag.Int("port", 10000, "The server port")
	ip         = flag.String("ip", "127.0.0.1", "The client ip")
	grpcServer *grpc.Server
)

type file_server struct {
	pb.UnimplementedFileServiceServer
}

func (*file_server) DownloadFile(req *pb.ServeFileRequest, res pb.FileService_DownloadFileServer) error {
	bufferSize := 64 * 1024 //64KiB, tweak this as desired
	file, err := os.Open("./file/" + req.GetFileName())

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	buff := make([]byte, bufferSize)
	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		resp := &pb.ServeFileResponse{
			ChunkData: buff[:bytesRead],
		}
		err = res.Send(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}

	}
	return nil
}
func InitFileServer() {
	grpclog.Println("start server...")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		grpclog.Fatal("failed to listen: %v", err)
	}

	grpcServer = grpc.NewServer()
	pb.RegisterFileServiceServer(grpcServer, new(file_server))

	file, err := os.Open("../file")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not found connect %v ", err)
	}

	defer cc.Close()

	c := mainpb.NewFileServiceClient(cc)
	// fmt.Println(*ip)
	// fmt.Println(int32(*port))
	// fmt.Println(list)
	mclient.RegisterPeers(c, *ip, int32(*port), list)

	grpcServer.Serve(lis)

	grpclog.Println("server shutdown...")

}
