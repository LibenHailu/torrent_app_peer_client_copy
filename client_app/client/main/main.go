package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	fileclient "github.com/LibenHailu/grpc_file_stream/file_stream/file_client"
	filesource "github.com/LibenHailu/grpc_file_stream/file_stream/file_source"
	"github.com/LibenHailu/peer_to_peer_file_share/peer-copy/client_app/client"
)

func main() {
	g := gin.Default()
	g.GET("/file/:fileName", func(ctx *gin.Context) {
		name := ctx.Param("fileName")
		fmt.Println(name)
		serverAddr, server := filesource.SearchAddressForThefile(name)
		fmt.Println(*serverAddr)
		fmt.Println(*server)
		client.InitFileClient(serverAddr, server)
		client.DownloadFile(name)
		ctx.JSON(http.StatusOK, gin.H{
			"server_address": fmt.Sprintf("server address %s", *serverAddr),
			"server":         fmt.Sprintf("served by %s", *server),
			"result":         fmt.Sprintf("file %s saved successfuly", name),
		})
	})
	g.POST("file/upload", func(ctx *gin.Context) {

		fileclient.UploadFileClient(ctx)

		// ctx.JSON(http.StatusOK, gin.H{
		// 	"data": fmt.Sprintf("file %s uploaded successfully", file.Filename),
		// })
		// return

		// Upload the file to specific dst.
		// c.SaveUploadedFile(file, dst)

		// c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	g.Run()
}
