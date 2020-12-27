package controller

import (
	"app/shared/session"
	"app/shared/view"
	"net/http"
)

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	session := session.Instance(r)

	if session.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "index/auth"
		v.Vars["first_name"] = session.Values["first_name"]
		v.Vars["is_auth"] = session.Values["is_auth"]
		v.Vars["is_student"] = session.Values["is_student"]
		v.Vars["is_supervisor"] = session.Values["is_supervisor"]
		v.Vars["is_project_manager"] = session.Values["is_project_manager"]

		v.Render(w)
	} else {
		// Display the view
		v := view.New(r)
		v.Name = "index/anon"
		v.Render(w)
		return
	}
}
