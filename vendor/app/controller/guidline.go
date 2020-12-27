package controller

import (
	"app/shared/session"
	"app/shared/view"
	"github.com/josephspurrier/csrfbanana"
	"net/http"
)

func GuidLineGET(w http.ResponseWriter, r *http.Request) {

	// Get session
	sess := session.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "guidline/guidline"

	//Get list of all users
	if sess.Values["id"] != nil {
		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]
		v.Render(w)
	} else {
		v.Render(w)
	}

}
