package main

import (
    "FORUM-GO/databaseAPI"
    "FORUM-GO/webAPI"
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "net/http"
    "os"
)

type Post struct {
    Id         int
    Username   string
    Title      string
    Categories []string
    Content    string
    CreatedAt  string
    UpVotes    int
    DownVotes  int
    Comments   []Comment
}

type Comment struct {
    Id        int
    PostId    int
    Username  string
    Content   string
    CreatedAt string
}

var database *sql.DB

func main() {
    // Vérification et création de la base de données
    var _, err = os.Stat("database.db")
    if os.IsNotExist(err) {
        var file, err = os.Create("database.db")
        if err != nil {
            fmt.Println("Erreur de création de la base de données:", err)
            return
        }
        defer file.Close()
    }

    // Ouverture de la base de données
    database, err = sql.Open("sqlite3", "./database.db")
    if err != nil {
        fmt.Println("Erreur d'ouverture de la base de données:", err)
        return
    }
    defer database.Close()

    // Création des tables
    databaseAPI.CreateUsersTable(database)
    databaseAPI.CreatePostTable(database)
    databaseAPI.CreateCommentTable(database)
    databaseAPI.CreateVoteTable(database)
    databaseAPI.CreateCategoriesTable(database)
    databaseAPI.CreateCategories(database)
    databaseAPI.CreateCategoriesIcons(database)

    // Configuration de la base de données pour le webAPI
    webAPI.SetDatabase(database)

    // Serveur de fichiers statiques
    fs := http.FileServer(http.Dir("public"))

    // Configuration du routeur
    router := http.NewServeMux()

    // Routes principales
    router.HandleFunc("/", webAPI.Index)
    router.HandleFunc("/register", webAPI.Register)
    router.HandleFunc("/login", webAPI.Login)
    router.HandleFunc("/post", webAPI.DisplayPost)
    router.HandleFunc("/filter", webAPI.GetPostsByApi)
    router.HandleFunc("/newpost", webAPI.NewPost)

    // Routes API
    router.HandleFunc("/api/register", webAPI.RegisterApi)
    router.HandleFunc("/api/login", webAPI.LoginApi)
    router.HandleFunc("/api/logout", webAPI.LogoutAPI)
    router.HandleFunc("/api/createpost", webAPI.CreatePostApi)
    router.HandleFunc("/api/comments", webAPI.CommentsApi)
    router.HandleFunc("/api/vote", webAPI.VoteApi)
    router.HandleFunc("/api/deletepost", webAPI.DeletePostHandler)

    // Gestion des fichiers statiques
    router.Handle("/public/", http.StripPrefix("/public/", fs))

    // Démarrage du serveur
    fmt.Println("Démarrage du serveur sur http://localhost:8080/")
    err = http.ListenAndServe(":8080", router)
    if err != nil {
        fmt.Println("Erreur de démarrage du serveur:", err)
    }
}