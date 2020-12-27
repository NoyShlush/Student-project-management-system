package controller

import (
	"app/model"
	"app/shared/session"
	"app/shared/view"
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

// ListProjectsGET displays the list of all projects
func ListArchiveProjectsGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "archive/archive"

	//View permissions
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["is_auth"] = sess.Values["is_auth"]
	v.Vars["is_student"] = sess.Values["is_student"]
	v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
	v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

	archiveProjects, err := model.ListArchiveProject()
	if err != nil {
		log.Println(err)
	}

	// Get all archive projects
	var projects []interface{}
	var students model.Students
	i := map[string]interface{}{
		"key": "value",
	}

	for _, a := range archiveProjects {
		t, _ := model.SelectProjectById(a.ProjectId)
		s, _ := model.SelectUserInfo(t.SupervisorId)
		json.Unmarshal(t.StudentsId, &students)
		u1, _ := model.SelectUserInfo(uint32(students.StudentId1))
		u2, _ := model.SelectUserInfo(uint32(students.StudentId2))
		i["project"] = map[string]interface{}{
			"project":    t,
			"supervisor": s.FirstName + " " + s.LastName,
			"students":   u1.FirstName + " " + u1.LastName + " ו" + u2.FirstName + " " + u2.LastName,
		}
		projects = append(projects, i)
	}

	v.Vars["archiveProjects"] = projects
	v.Render(w)
}

// ArchiveProjectGET displays an additional information about the archive project
func ArchiveProjectGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "archive/more_archive"

	//View permissions
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["is_auth"] = sess.Values["is_auth"]
	v.Vars["is_student"] = sess.Values["is_student"]
	v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
	v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

	archiveProject, err := model.SelectArchiveProjectById(uint32(ID))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var filesJson model.ProjectFiles
	err = json.Unmarshal(archiveProject.ApprovedFiles, &filesJson)

	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	// Get all files
	Book1PDF, err := model.SelectFile(uint32(filesJson.Book1PDF))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	v.Vars["Book1PDF"] = Book1PDF
	Book1WORD, err := model.SelectFile(uint32(filesJson.Book1WORD))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	v.Vars["Book1WORD"] = Book1WORD
	Book2PDF, err := model.SelectFile(uint32(filesJson.Book2PDF))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	v.Vars["Book2PDF"] = Book2PDF
	Book2WORD, err := model.SelectFile(uint32(filesJson.Book2WORD))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	v.Vars["Book2WORD"] = Book2WORD
	Presentation1, err := model.SelectFile(uint32(filesJson.Presentation1))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	v.Vars["Presentation1"] = Presentation1
	Presentation2, err := model.SelectFile(uint32(filesJson.Presentation2))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	v.Vars["Presentation2"] = Presentation2
	SourceCode, err := model.SelectFile(uint32(filesJson.SourceCode))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	v.Vars["SourceCode"] = SourceCode

	//Get the project
	project, err := model.SelectProjectById(archiveProject.ProjectId)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	supervisors, err := model.SelectUserInfo(project.SupervisorId)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var Students model.Students
	err = json.Unmarshal(project.StudentsId, &Students)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var students []model.User
	student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	students = append(students, student1)

	student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	students = append(students, student2)

	v.Vars["project_students"] = students
	v.Vars["createdat"] = project.CreateAt.Format("15:04:05 01/02/06")
	v.Vars["updatedat"] = project.UpdateAt.Format("15:04:05 01/02/06")
	v.Vars["project"] = project
	v.Vars["supervisors"] = supervisors

	v.Render(w)

}

// SearchArchiveProjectsPOST displays the list of all projects
func SearchArchiveProjectsPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	// Display the view
	v := view.New(r)
	v.Name = "archive/archive"

	//View permissions
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["is_auth"] = sess.Values["is_auth"]
	v.Vars["is_student"] = sess.Values["is_student"]
	v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
	v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"query"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	query := r.FormValue("query")

	archiveProjects, err := model.SearchArchiveProject(query)
	if err != nil {
		log.Println(err)
	}

	// Get all archive projects
	var projects []interface{}
	var students model.Students
	i := map[string]interface{}{
		"key": "value",
	}

	for _, a := range archiveProjects {
		t, _ := model.SelectProjectById(a.ProjectId)
		s, _ := model.SelectUserInfo(t.SupervisorId)
		json.Unmarshal(t.StudentsId, &students)
		u1, _ := model.SelectUserInfo(uint32(students.StudentId1))
		u2, _ := model.SelectUserInfo(uint32(students.StudentId2))
		i["project"] = map[string]interface{}{
			"project":    t,
			"supervisor": s.FirstName + " " + s.LastName,
			"students":   u1.FirstName + " " + u1.LastName + " ו" + u2.FirstName + " " + u2.LastName,
		}
		projects = append(projects, i)
	}

	v.Vars["archiveProjects"] = projects
	v.Render(w)
}
