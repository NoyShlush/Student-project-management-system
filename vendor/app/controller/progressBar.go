package controller

import (
	"app/model"
	"app/shared/session"
	"app/shared/view"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

	//	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	//	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

// ProgressBarStudentGET handles the progress bar page for student
func ProgressBarStudentGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "progressbar/progressbar"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		// Get user ID form session
		userId, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
		//Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		project, err := model.SelectProjectByStudent(uint32(userId))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		//Get list of all milestones
		progressBar, err := model.SelectListOfMilestones(project.ID)
		if err != nil {
			log.Println(err)
			progressBar = []model.ProgressBar{}
		}

		done := 0
		for _, key := range progressBar {
			if key.Done {
				done++
			}
		}
		v.Vars["percentage"] = int((float32(done) / float32(len(progressBar))) * 100)
		v.Vars["progressBar"] = progressBar
		v.Render(w)

	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ProgressBarSupervisorGET handles the progress bar page for supervisor
func ProgressBarSupervisorGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "progressbar/progressbar"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		// Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		projectID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		//Get list of all milestones
		progressBar, err := model.SelectListOfMilestones(uint32(projectID))
		if err != nil {
			log.Println(err)
			progressBar = []model.ProgressBar{}
		}
		done := 0
		for _, key := range progressBar {
			if key.Done {
				done++
			}
		}

		v.Vars["progressBar"] = progressBar
		v.Vars["percentage"] = int((float32(done) / float32(len(progressBar))) * 100)
		v.Render(w)

	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// MilestoneDoneGET mark as done milestone on the DB
func MilestoneDoneGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	progressBarID, err := strconv.ParseUint(params.ByName("id1"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	milestoneID, err := strconv.ParseUint(params.ByName("id2"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = model.MarkAsDone(uint32(milestoneID))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
	}

	sess.Save(r, w)
	http.Redirect(w, r, "/progressbar/progressbar/"+strconv.Itoa(int(progressBarID)), http.StatusFound)

}

// MilestoneRemoveGET delete from the DB the milestone
func MilestoneRemoveGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)

	progressBarID, err := strconv.ParseUint(params.ByName("id1"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	milestoneID, err := strconv.ParseUint(params.ByName("id2"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = model.DeleteMilestone(uint32(milestoneID))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
	}

	sess.Save(r, w)
	http.Redirect(w, r, "/progressbar/progressbar/"+strconv.Itoa(int(progressBarID)), http.StatusFound)
}

// MilestoneAddPOST add to the DB the new milestone
func MilestoneAddPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"milestone"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		ProgressBarSupervisorGET(w, r)
		return
	}

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)

	progressBarID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	milestone := r.FormValue("milestone")
	err = model.CreateMilestone(uint32(progressBarID), milestone)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
	}

	sess.Save(r, w)
	ProgressBarSupervisorGET(w, r)
}
