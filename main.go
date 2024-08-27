package main

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func upload(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// TODO need to check the file size at this point or if possible earlier
	// we also need to check users plan and total available storage
	// to prevent user from uploading something beyond their plan
	// alternatively to save seconds on performance
	// we can let the file go through usual flow,
	// and somewhere later check if total files size is exceeding plan,
	// if yes flag it and simple check flag here, should faster in theory

	// TODO here we need to start uploading file to the s3 storage async, should be non-blocking, so user should receive response without waiting full upload
	// but i wonder if upload fails how do we inform user, the goal is to shred some seconds since uploading to s3 would take time

	// TODO won't be required when uploading to s3
	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// TODO won't be required when uploading to s3
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	// TODO zero idea on what is centralized HTTPErrorHandler
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO use prod and debug env flags to determine domain
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "im alive!")
	})

	e.POST("/upload", upload)

	// TODO curios how this logger works
	e.Logger.Fatal(e.Start(":1323"))
}
