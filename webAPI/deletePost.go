package webAPI

import (
	"FORUM-GO/databaseAPI"
	"fmt"
	"net/http"
	"strconv"
)

// DeletePostHandler gère la suppression d'un post
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Vérification de l'authentification
	cookie, err := r.Cookie("SESSION")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Vérifier si le cookie est valide
	cookieExists := databaseAPI.CheckCookie(database, cookie.Value)
	if !cookieExists {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Récupération des paramètres
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Erreur de traitement du formulaire: %v", err), http.StatusBadRequest)
		return
	}

	username := databaseAPI.GetUser(database, cookie.Value)
	postIdStr := r.FormValue("postId")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Identifiant de post invalide", http.StatusBadRequest)
		return
	}

	// Vérification des droits
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM posts WHERE id = ? AND username = ?", postId, username).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "Vous n'avez pas l'autorisation de supprimer ce post", http.StatusUnauthorized)
		return
	}

	// Suppression du post avec transaction
	tx, err := database.Begin()
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	
	// Suppression des votes
	_, err = tx.Exec("DELETE FROM votes WHERE post_id = ?", postId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	
	// Suppression des commentaires
	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", postId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	
	// Suppression du post
	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", postId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	
	if err = tx.Commit(); err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}