package multipart

import (
    "os"
    //"fmt"
    "errors"
    "net/http"
    "io/ioutil"
    "mime/multipart"
    //"gopkg.in/mgo.v2/bson"
)

// UploadFile uploads a file to the server
func UploadFile(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    file, handle, err := r.FormFile("file")
    if err != nil {
        return nil, err
    }
    defer file.Close()
    //var path = nil
    mimeType := handle.Header.Get("Content-Type")
    switch mimeType {
        case "image/jpeg":
            return saveFile(w, file, handle)
        case "image/jpg":
            return saveFile(w, file, handle)
        case "image/png":
            return saveFile(w, file, handle)
        default:
            return nil, errors.New("The format file is not valid.")
    }
    return nil, nil
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) (interface{}, error) {
    data, err := ioutil.ReadAll(file);
    if err != nil {
        return nil, err
    }

    // ./upload does not exist
    if _, err := os.Stat("./upload"); os.IsNotExist(err) {
        _ = os.Mkdir("upload", os.ModePerm)
    }

    if err := ioutil.WriteFile(os.Getenv("UPLOAD_PATH") + handle.Filename, data, 0666); err != nil {
        return nil, err
    }
    return map[string]string{"filename": handle.Filename}, nil
    //return bson.M{"filename": handle.Filename}, nil
}
