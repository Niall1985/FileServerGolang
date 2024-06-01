# Go File Server

This project is a simple file server implemented in Go that allows users to upload, download, list, and delete files. The server saves uploaded files in an `uploads` directory and provides endpoints to manage the files.

## Prerequisites

- Go (1.16 or higher)
- PowerShell (for testing and interacting with the server)

## Installation

1. **Create a New Go Module**

   If you're starting a new project, you'll want to initialize it as a Go module. Navigate to your project directory and run:

   ```bash
   go mod init github.com/yourusername/your-repository-name
   ```

   Replace `github.com/yourusername/your-repository-name` with the desired module path.

2. **Create the Project Structure**

   Inside your project directory, create the necessary directories and files:

   ```bash
   mkdir uploads
   ```

   This will create the `uploads` directory where uploaded files will be stored.

3. **Create the `main.go` File**

   Create a `main.go` file in your project directory and import the `FileServerGolang` package to use it. Here's an example `main.go` file:

   ```go
   package main

   import (
       "github.com/Niall1985/FileServerGolang"
   )

   func main() {
       fileserver.StartServer(":8080")
   }
   ```
4. **Install the `FileServerGolang` Package**

   Install the `FileServerGolang` package using the `go get` command:

   ```bash
   go get github.com/Niall1985/FileServerGolang
   ```

## Running the Server

### Using `go run`

You can run the server directly with `go run`. First, create a new `main.go` file in your project directory and import the `fileserver` package:

- Run the server using:

```bash
go run main.go
```

### Using `go build`

Alternatively, you can build and run the server executable:

1. **Build the Project**

   ```bash
   go build -o file-server
   ```

2. **Run the Executable**

   ```bash
   ./file-server
   ```

The server will start listening on port 8080.

## PowerShell Commands

### Uploading a File

To upload a file from your local system to the server, use the following PowerShell script:

```powershell
$FilePath = "C:\Path\To\Your\File.ext"
$Uri = "http://localhost:8080/upload"
$fileName = [System.IO.Path]::GetFileName($FilePath)
$fileContent = [System.IO.File]::ReadAllBytes($FilePath)
$boundary = [System.Guid]::NewGuid().ToString()

$bodyLines = @(
    "--$boundary",
    "Content-Disposition: form-data; name=`"file`"; filename=`"$fileName`"",
    "Content-Type: application/octet-stream",
    "",
    [System.Text.Encoding]::Default.GetString($fileContent),
    "--$boundary--"
)

$body = $bodyLines -join "`r`n"
$bytes = [System.Text.Encoding]::Default.GetBytes($body)

$webRequest = [System.Net.HttpWebRequest]::Create($Uri)
$webRequest.Method = "POST"
$webRequest.ContentType = "multipart/form-data; boundary=$boundary"
$webRequest.ContentLength = $bytes.Length

$stream = $webRequest.GetRequestStream()
$stream.Write($bytes, 0, $bytes.Length)
$stream.Close()

$response = $webRequest.GetResponse()
$reader = New-Object System.IO.StreamReader($response.GetResponseStream())
$result = $reader.ReadToEnd()

Write-Host $result
```

### Viewing the List of Files

To view the list of files in the `uploads` directory:

```powershell
$Uri = "http://localhost:8080/list"
Invoke-RestMethod -Uri $Uri -Method Get
```

### Downloading a File

To download a specific file from the server:

```powershell
$FileName = "example.pdf"  # Replace with the actual filename
$Uri = "http://localhost:8080/download/$FileName"
$OutputPath = "C:\path\to\save\downloaded_example.pdf"  # Replace with the desired output path
Invoke-WebRequest -Uri $Uri -OutFile $OutputPath
```

### Deleting a File

To delete a specific file from the server:

```powershell
$FileName = "example.pdf"  # Replace with the actual filename
$Uri = "http://localhost:8080/delete/$FileName"
Invoke-RestMethod -Uri $Uri -Method Delete
```

## Server Endpoints

- **Upload a File**: `POST /upload`
- **Download a File**: `GET /download/{filename}`
- **List Files**: `GET /list`
- **Delete a File**: `DELETE /delete/{filename}`

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any changes or suggestions.

## License

This project is licensed under the MIT License.

