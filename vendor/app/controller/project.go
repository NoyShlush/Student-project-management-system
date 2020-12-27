package controller

import (
	"app/model"
	"app/shared/aws"
	"app/shared/session"
	"app/shared/view"
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type comments struct {
	Comment []uint32 `json:"comment"`
}

// NewIdeaGET displays the new idea form
func NewIdeaGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get user session
	if sess.Values["id"] != nil {

		//Check if the user id not student
		if sess.Values["is_student"] == nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		// Display the view
		v := view.New(r)
		v.Name = "project/new_idea"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		//Gett list of all supervisors
		supervisors, err := model.SelectListOfUsers(2)
		if err != nil {
			log.Println(err)
			supervisors = []model.User{}
		}

		managers, err := model.SelectListOfUsers(3)
		if err != nil {
			log.Println(err)
			managers = []model.User{}
		}

		v.Vars["supervisors"] = append(supervisors, managers...)

		// Get user ID form session
		studentId, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
		//Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		//Gett list of all students
		students, err := model.ListOfStudentsWithoutProject(uint32(studentId))
		if err != nil {
			log.Println(err)
			students = []model.User{}
		}

		v.Vars["students"] = students

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// NewIdeaPOST submit to the DB the new idea form
func NewIdeaPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"projectname", "short_description", "description", "supervisor_id", "student_id"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		NewIdeaGET(w, r)
		return
	}

	// Get form values
	supervisorId, err := strconv.ParseUint(r.FormValue("supervisor_id"), 10, 32)
	//Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		NewIdeaGET(w, r)
		return
	}

	// Get user ID form session
	studentId1, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
	//Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		NewIdeaGET(w, r)
		return
	}

	// Get form values
	studentId2, err := strconv.ParseInt(r.FormValue("student_id"), 10, 32)
	//Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		NewIdeaGET(w, r)
		return
	}

	students := model.Students{
		int(studentId1),
		int(studentId2),
	}

	studentJson, err := json.Marshal(students)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		NewIdeaGET(w, r)
		return
	}

	err = model.CreateStudentIdea(r.FormValue("projectname"), r.FormValue("description"), r.FormValue("short_description"), uint32(supervisorId), studentJson)

	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		NewIdeaGET(w, r)
		return
	}

	sess.AddFlash(view.Flash{"ההצעה נשלחה למנחה בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)
}

// NewProjectGET displays the new project form
func NewProjectGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get user session
	if sess.Values["id"] != nil {

		//Check if the user id not student
		if sess.Values["is_supervisor"] == nil && sess.Values["is_project_manager"] == nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		// Display the view
		v := view.New(r)
		v.Name = "project/page_project"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// NewProjectPOST submit to the DB the new project form
func NewProjectPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"projectname", "short_description", "description"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		NewProjectGET(w, r)
		return
	}

	// Get user ID form session
	projectManagerId, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
	//Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		NewProjectGET(w, r)
		return
	}

	students := model.Students{
		0,
		0,
	}

	studentJson, err := json.Marshal(students)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		NewProjectGET(w, r)
		return
	}

	err = model.CreateProject(r.FormValue("projectname"), r.FormValue("description"), r.FormValue("short_description"), uint32(projectManagerId), studentJson)

	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		NewProjectGET(w, r)
		return
	}

	sess.AddFlash(view.Flash{"הפרויקט עלה בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)
}

// ApprovalFormGET displays the approval form page
func ApprovalFormGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get list of all users
	if sess.Values["id"] != nil {

		// Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Get all information abut the project
		project, err := model.SelectProjectById(uint32(ID))
		if err != nil && err != model.ErrNoResult {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		approvalForm, err := model.SelectApprovalForm(project.FormId)
		if err != nil && err != model.ErrNoResult {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Display the view
		v := view.New(r)
		v.Name = "project/approval"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]
		v.Vars["approval_form"] = approvalForm

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ApprovalFormPOST submit the approval form to the DB
func ApprovalFormPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"synopses", "scopeoftheproject", "uniquefeatures"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		ApprovalFormGET(w, r)
		return
	}

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		ApprovalFormGET(w, r)
		return
	}

	formID, err := model.CreateApprovalForm(r.FormValue("synopses"), r.FormValue("scopeoftheproject"), r.FormValue("uniquefeatures"), uint32(ID))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		ApprovalFormGET(w, r)
		return
	}

	err = model.SetApprovalForm(uint32(ID), formID)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		ApprovalFormGET(w, r)
		return
	}

	sess.AddFlash(view.Flash{"השינוים עודכנו בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)
}

// ListProjectsGET displays the list of all projects
func ListProjectsGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get list of all users
	if sess.Values["id"] != nil {

		// Check if the student logged in
		if sess.Values["is_student"] != nil {
			// Get user ID form session
			studentId, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
			// Check if the number converted well
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/", http.StatusFound)
			}

			// Check if the student is has project
			project, err := model.SelectProjectByStudent(uint32(studentId))
			if err != nil && err != model.ErrNoResult {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/", http.StatusFound)
			}

			if err != model.ErrNoResult {
				// Display the view
				v := view.New(r)
				v.Name = "project/project"

				//View permissions
				v.Vars["token"] = csrfbanana.Token(w, r, sess)
				v.Vars["first_name"] = sess.Values["first_name"]
				v.Vars["is_auth"] = sess.Values["is_auth"]
				v.Vars["is_student"] = sess.Values["is_student"]
				v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
				v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

				//Extract the list of users
				var Students model.Students
				err = json.Unmarshal(project.StudentsId, &Students)
				if err != nil {
					log.Println(err)
					sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
					sess.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)
				}

				student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
				if err != nil {
					log.Println(err)
					sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
					sess.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)
				}

				student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
				if err != nil {
					log.Println(err)
					sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
					sess.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)
				}

				supervisors, err := model.SelectUserInfo(project.SupervisorId)
				if err != nil {
					log.Println(err)
					sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
					sess.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)
				}

				// Show comments
				var comments comments
				err = json.Unmarshal(project.CommentsId, &comments)
				if err != nil {
					log.Println(err)
					sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
					sess.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)
				}

				if project.FormId != 0 {
					approvalForm, err := model.SelectApprovalForm(project.FormId)
					if err != nil {
						log.Println(err)
						sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
						sess.Save(r, w)
						http.Redirect(w, r, "/", http.StatusFound)
					}
					v.Vars["approvalform"] = approvalForm
				}

				// Presentation of the comments in the page
				var message []string
				for _, id := range comments.Comment {
					comment, _ := model.SelectComment(id)
					by, _ := model.SelectUserInfo(comment.UserId)
					tempString := "נכתב ע\"י " + by.FirstName + " " + by.LastName + " ב-" + comment.CreatedAt.Format("15:04:05 01/02/06") + "\t" + comment.Message
					message = append(message, tempString)
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

				var projectFiles model.ProjectFiles
				err = json.Unmarshal(project.Files, &projectFiles)
				if err != nil {
					sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
					sess.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}

				// Set date format
				v.Vars["createdat"] = project.CreateAt.Format("15:04:05 01/02/06")
				v.Vars["updatedat"] = project.UpdateAt.Format("15:04:05 01/02/06")
				v.Vars["comments"] = message
				v.Vars["supervisors"] = supervisors
				v.Vars["students"] = student1.FirstName + " " + student1.LastName + " ו" + student2.FirstName + " " + student2.LastName
				v.Vars["project"] = project
				v.Vars["files"] = projectFiles

				v.Render(w)
				return
			}
		}

		// Display the view
		v := view.New(r)
		v.Name = "project/open_project"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		projects, err := model.ListOpenProjects()
		if err != nil {
			log.Println(err)
		}
		v.Vars["projects"] = projects
		var supervisors []string

		for i := range projects {
			supervisor, _ := model.SelectUserInfo(projects[i].SupervisorId)

			//HACK NOT USE IT AGAIN
			projects[i].Description = supervisor.FirstName + " " + supervisor.LastName
		}

		v.Vars["supervisors"] = supervisors

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SearchOpenProjectPOST for project with status open
func SearchOpenProjectPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {

		// Display the view
		v := view.New(r)
		v.Name = "project/open_project"

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
			ListProjectsGET(w, r)
			return
		}

		query := r.FormValue("query")

		//Get list of students by search
		projects, err := model.SearchOpenProject(query)
		if err != nil {
			log.Println(err)
			projects = []model.Project{}
		}

		v.Vars["projects"] = projects
		var supervisors []string

		for i := range projects {
			supervisor, _ := model.SelectUserInfo(projects[i].SupervisorId)

			//HACK NOT USE IT AGAIN
			projects[i].Description = supervisor.FirstName + " " + supervisor.LastName
		}

		v.Vars["supervisors"] = supervisors
		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// MoreGET displays an additional information about the project
func MoreGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get user session
	if sess.Values["id"] != nil {

		//Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Display the view
		v := view.New(r)
		v.Name = "project/more_project"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		//Get the project
		project, err := model.SelectProjectById(uint32(ID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		supervisors, err := model.SelectUserInfo(project.SupervisorId)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		var Students model.Students
		err = json.Unmarshal(project.StudentsId, &Students)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		//Check if the student logged in
		if sess.Values["is_student"] != nil && Students.StudentId1 == 0 && Students.StudentId2 == 0 {
			// Get user ID form session
			studentId, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
			//Check if the number converted well
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			//Gett list of all students
			students, err := model.ListOfStudentsWithoutProject(uint32(studentId))
			if err != nil {
				log.Println(err)
				students = []model.User{}
			}

			v.Vars["students"] = students

		} else if (sess.Values["is_supervisor"] != nil || sess.Values["is_project_manager"] != nil) && Students.StudentId1 != 0 && Students.StudentId2 != 0 {
			var students []model.User
			student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}
			students = append(students, student1)
			student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}
			students = append(students, student2)
			v.Vars["project_students"] = students
		}

		// Show comments
		if (project.Type == 1 && project.StatusId == 2) || project.FormId != 0 {
			var comments comments
			err = json.Unmarshal(project.CommentsId, &comments)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			// Presentation of the comments in the page
			var message []string
			for _, id := range comments.Comment {
				comment, _ := model.SelectComment(id)
				by, _ := model.SelectUserInfo(comment.UserId)
				tempString := "נכתב ע\"י " + by.FirstName + " " + by.LastName + " ב-" + comment.CreatedAt.Format("15:04:05 01/02/06") + "\t" + comment.Message
				message = append(message, tempString)
			}
			v.Vars["comments"] = message
		}

		// Present approval for if exist
		if project.FormId != 0 {
			approvalForm, err := model.SelectApprovalForm(project.FormId)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}
			v.Vars["approvalform"] = approvalForm
		}

		v.Vars["createdat"] = project.CreateAt.Format("15:04:05 01/02/06")
		v.Vars["updatedat"] = project.UpdateAt.Format("15:04:05 01/02/06")
		v.Vars["project"] = project
		v.Vars["supervisors"] = supervisors

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// MorePOST assign the project to the students
func MorePOST(w http.ResponseWriter, r *http.Request) {
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
		MoreGET(w, r)
		return
	}

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"student_id"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		MoreGET(w, r)
		return
	}

	// Get user ID form session
	studentId1, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
	//Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		MoreGET(w, r)
		return
	}

	// Get form values
	studentId2, err := strconv.ParseInt(r.FormValue("student_id"), 10, 32)
	//Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		MoreGET(w, r)
		return
	}

	students := model.Students{
		int(studentId1),
		int(studentId2),
	}

	studentJson, err := json.Marshal(students)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		MoreGET(w, r)
		return
	}

	err = model.AssignStudentsToProject(uint32(ID), studentJson)

	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		MoreGET(w, r)
		return
	}

	sess.AddFlash(view.Flash{"ההצעה נשלחה למנחה בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)
}

// ApprovalWaitingProjectsGET displays the list of project which waiting for approval
func ApprovalWaitingProjectsGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get user session
	if sess.Values["id"] != nil {

		// Display the view
		v := view.New(r)
		v.Name = "project/wait_project"

		// Get user ID form session
		ID, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
		//Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		projects, err := model.ListWaitProjects(uint32(ID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		v.Vars["projects"] = projects
		var supervisors []string

		for i, project := range projects {

			//HACK NOT USE IT AGAIN
			var Students model.Students
			err := json.Unmarshal(project.StudentsId, &Students)
			if err != nil {
				log.Println(err)
			}
			student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
			if err != nil {
				log.Println(err)
			}
			student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
			if err != nil {
				log.Println(err)
			}
			projects[i].Description = student1.FirstName + " " + student1.LastName + " ו" + student2.FirstName + " " + student2.LastName
		}

		v.Vars["supervisors"] = supervisors

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SearchApprovalWaitingProjectPOST search projects with status waiting for approval by supervisor
func SearchApprovalWaitingProjectPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {

		// Display the view
		v := view.New(r)
		v.Name = "project/wait_project"

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
			ApprovalWaitingProjectsGET(w, r)
			return
		}

		query := r.FormValue("query")

		// Get user ID form session
		ID, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)

		// Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ApprovalWaitingProjectsGET(w, r)
			return
		}

		projects, err := model.SearchWaitingProject(uint32(ID), query)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ApprovalWaitingProjectsGET(w, r)
			return
		}

		v.Vars["projects"] = projects
		var supervisors []string

		for i, project := range projects {

			//HACK NOT USE IT AGAIN
			var Students model.Students
			err := json.Unmarshal(project.StudentsId, &Students)
			if err != nil {
				log.Println(err)
			}
			student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
			if err != nil {
				log.Println(err)
			}
			student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
			if err != nil {
				log.Println(err)
			}
			projects[i].Description = student1.FirstName + " " + student1.LastName + " ו" + student2.FirstName + " " + student2.LastName
		}

		v.Vars["supervisors"] = supervisors

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ManagerApprovalWaitingProjectsGET displays the list of project which waiting for manager approval
func ManagerApprovalWaitingProjectsGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get user session
	if sess.Values["id"] != nil {

		// Display the view
		v := view.New(r)
		v.Name = "project/manager_wait_project"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		projects, err := model.ListWaitManagerApproval()
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		v.Vars["projects"] = projects
		var supervisors []string

		for i, project := range projects {

			//HACK NOT USE IT AGAIN
			var Students model.Students
			err := json.Unmarshal(project.StudentsId, &Students)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}
			student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}
			student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}
			projects[i].Description = student1.FirstName + " " + student1.LastName + " ו" + student2.FirstName + " " + student2.LastName
		}

		v.Vars["supervisors"] = supervisors

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SearchManagerApprovalWaitingProjectsPOST for project with status waiting for manager approval
func SearchManagerApprovalWaitingProjectsPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {

		// Display the view
		v := view.New(r)
		v.Name = "project/manager_wait_project"

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
			ManagerApprovalWaitingProjectsGET(w, r)
			return
		}

		query := r.FormValue("query")
		projects, err := model.SearchManagerApprovalProject(query)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ManagerApprovalWaitingProjectsGET(w, r)
			return
		}

		v.Vars["projects"] = projects
		var supervisors []string

		for i, project := range projects {

			//HACK NOT USE IT AGAIN
			var Students model.Students
			err := json.Unmarshal(project.StudentsId, &Students)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ManagerApprovalWaitingProjectsGET(w, r)
				return
			}

			student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ManagerApprovalWaitingProjectsGET(w, r)
				return
			}

			student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ManagerApprovalWaitingProjectsGET(w, r)
				return
			}

			projects[i].Description = student1.FirstName + " " + student1.LastName + " ו" + student2.FirstName + " " + student2.LastName
		}

		v.Vars["supervisors"] = supervisors

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// AddCommentsGET displays the add comment form to the supervisor
func AddCommentsGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get user session
	if sess.Values["id"] != nil {

		//Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Display the view
		v := view.New(r)
		v.Name = "project/add_comments_project"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		//Get the project
		project, err := model.SelectProjectById(uint32(ID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		//Extract the list of users
		var Students model.Students
		var students []model.User

		err = json.Unmarshal(project.StudentsId, &Students)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		students = append(students, student1)
		student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Get supervisors info
		supervisors, err := model.SelectUserInfo(project.SupervisorId)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		students = append(students, student2)
		v.Vars["project_students"] = students
		v.Vars["supervisors"] = supervisors
		v.Vars["project"] = project

		var comments comments
		err = json.Unmarshal(project.CommentsId, &comments)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Presentation of the comments in the page
		var message []string
		for _, id := range comments.Comment {
			comment, _ := model.SelectComment(id)
			by, _ := model.SelectUserInfo(comment.UserId)
			tempString := "נכתב ע\"י " + by.FirstName + " " + by.LastName + " ב-" + comment.CreatedAt.Format("15:04:05 01/02/06") + "\t" + comment.Message
			message = append(message, tempString)
		}

		v.Vars["comments"] = message
		v.Vars["createdat"] = project.CreateAt.Format("15:04:05 01/02/06")
		v.Vars["updatedat"] = project.UpdateAt.Format("15:04:05 01/02/06")

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// AddCommentsPOST Add comment to the project
func AddCommentsPOST(w http.ResponseWriter, r *http.Request) {
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
		ListProjectsGET(w, r)
		return
	}

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"comment"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		AddCommentsGET(w, r)
		return
	}

	// Get user ID form session
	supervisor, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
	// Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddCommentsGET(w, r)
		return
	}

	// Get comment from the form
	comment := r.FormValue("comment")

	// Add new comment
	commentID, err := model.CreateComment(comment, uint32(supervisor))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddCommentsGET(w, r)
		return
	}

	var comments comments

	//Get the project
	project, err := model.SelectProjectById(uint32(ID))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddCommentsGET(w, r)
		return
	}

	err = json.Unmarshal(project.CommentsId, &comments)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddCommentsGET(w, r)
		return
	}

	// Add to array of comments the new comment
	comments.Comment = append(comments.Comment, commentID)
	commentsJson, err := json.Marshal(comments)

	// Update the DB with the new comment
	err = model.UpdateComments(uint32(ID), commentsJson)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddCommentsGET(w, r)
		return
	}

	sess.AddFlash(view.Flash{"הערות נוספו בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/waiting", http.StatusFound)
}

// AddApprovalCommentsGET displays the add comment form to the project manager
func AddApprovalCommentsGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get user session
	if sess.Values["id"] != nil {

		//Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Display the view
		v := view.New(r)
		v.Name = "project/add_comments_project"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		//Get the project
		project, err := model.SelectProjectById(uint32(ID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Get supervisors info
		supervisors, err := model.SelectUserInfo(project.SupervisorId)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
		}

		//Extract the list of users
		var Students model.Students
		err = json.Unmarshal(project.StudentsId, &Students)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		var students []model.User
		student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		students = append(students, student1)
		student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		students = append(students, student2)

		v.Vars["project_students"] = students
		v.Vars["project"] = project
		v.Vars["supervisors"] = supervisors

		// Present approval for if exist
		if project.FormId != 0 {
			approvalForm, err := model.SelectApprovalForm(project.FormId)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}
			v.Vars["approvalform"] = approvalForm
		}

		var comments comments
		err = json.Unmarshal(project.CommentsId, &comments)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Presentation of the comments in the page
		var message []string
		for _, id := range comments.Comment {
			comment, _ := model.SelectComment(id)
			by, _ := model.SelectUserInfo(comment.UserId)
			tempString := "נכתב ע\"י " + by.FirstName + " " + by.LastName + " ב-" + comment.CreatedAt.Format("15:04:05 01/02/06") + "\t" + comment.Message
			message = append(message, tempString)
		}
		v.Vars["comments"] = message
		v.Vars["createdat"] = project.CreateAt.Format("15:04:05 01/02/06")
		v.Vars["updatedat"] = project.UpdateAt.Format("15:04:05 01/02/06")

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// AddCommentsPOST Add comment to the approval form
func AddApprovalCommentsPOST(w http.ResponseWriter, r *http.Request) {
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
		ListProjectsGET(w, r)
		return
	}

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"comment"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		AddApprovalCommentsGET(w, r)
		return
	}

	// Get user ID form session
	supervisor, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)

	// Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddApprovalCommentsGET(w, r)
		return
	}

	// Get comment from the form
	comment := r.FormValue("comment")

	// Add new comment
	commentID, err := model.CreateComment(comment, uint32(supervisor))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddApprovalCommentsGET(w, r)
		return
	}

	var comments comments

	//Get the project
	project, err := model.SelectProjectById(uint32(ID))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddApprovalCommentsGET(w, r)
		return
	}

	err = json.Unmarshal(project.CommentsId, &comments)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddApprovalCommentsGET(w, r)
		return
	}

	// Add to array of comments the new comment
	comments.Comment = append(comments.Comment, commentID)
	commentsJson, err := json.Marshal(comments)

	// Update the DB with the new comment
	err = model.UpdateComments(uint32(ID), commentsJson)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		AddApprovalCommentsGET(w, r)
		return
	}

	sess.AddFlash(view.Flash{"הערות נוספו בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/manager_waiting", http.StatusFound)
}

// EditProjectGET displays the edit project page
func EditProjectGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get list of all users
	if sess.Values["id"] != nil {

		// Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Get all information abut the project
		project, err := model.SelectProjectById(uint32(ID))
		if err != nil && err != model.ErrNoResult {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
		}

		// Display the view
		v := view.New(r)
		v.Name = "project/page_project"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]
		v.Vars["project"] = project

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// EditProjectPOST submit the updated project to the DB
func EditProjectPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"projectname", "short_description", "description"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		ListProjectsGET(w, r)
		return
	}

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		ListProjectsGET(w, r)
		return
	}

	err = model.UpdateProject(uint32(ID), r.FormValue("projectname"), r.FormValue("description"), r.FormValue("short_description"))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		ListProjectsGET(w, r)
		return
	}

	sess.AddFlash(view.Flash{"השינוים עודכנו בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)
}

// ApproveProjectGET approves the project or an idea
func ApproveProjectGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get list of all users
	if sess.Values["id"] != nil {

		// Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Approve the project
		err = model.ApproveProject(uint32(ID))
		if err != nil && err != model.ErrNoResult {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
		}

		sess.AddFlash(view.Flash{"הפרויקט מאושר", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/waiting", http.StatusFound)

	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// DeclineProjectGET decline the project or an idea
func DeclineProjectGET(w http.ResponseWriter, r *http.Request) { // Get session
	sess := session.Instance(r)

	// Get list of all users
	if sess.Values["id"] != nil {

		// Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Get all information abut the project
		project, err := model.SelectProjectById(uint32(ID))
		if err != nil && err != model.ErrNoResult {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
		}

		// If the project is idea it will be deleted if the project if project (By supervisor) it will remove the assigned students
		if project.Type == 1 {

			// Delete all comments before deleting of the project
			var comments comments
			err = json.Unmarshal(project.CommentsId, &comments)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			for _, id := range comments.Comment {
				go func() {
					model.DeleteComment(id)
				}()
			}

			err = model.DeleteProject(uint32(ID))
			if err != nil && err != model.ErrNoResult {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

		} else if project.Type == 2 {
			err = model.DeclineProject(uint32(ID))
			if err != nil && err != model.ErrNoResult {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}
		}

		sess.AddFlash(view.Flash{"הפרויקט נדחה", view.FlashWarning})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/waiting", http.StatusFound)

	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ApproveProjectGET approves the project or an idea
func ApproveManagerProjectGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get list of all users
	if sess.Values["id"] != nil {

		// Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Approve the project
		err = model.ApproveProjectByManager(uint32(ID))
		if err != nil && err != model.ErrNoResult {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Add default milestones
		err = model.CreateMilestone(uint32(ID), "כתיבת ספר חלק א'")
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		err = model.CreateMilestone(uint32(ID), "הכנת מצגת חלק א'")
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		err = model.CreateMilestone(uint32(ID), "כתיבת ספר חלק ב'")
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		err = model.CreateMilestone(uint32(ID), "הכנת מצגת חלק ב'")
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		err = model.CreateMilestone(uint32(ID), "כתיבת קוד")
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		sess.AddFlash(view.Flash{"הפרויקט מאושר", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/manager_waiting", http.StatusFound)

	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// DeclineManagerProjectGET delete the project or an idea
func DeclineManagerProjectGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get list of all users
	if sess.Values["id"] != nil {

		// Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		// Get all information abut the project
		project, err := model.SelectProjectById(uint32(ID))
		if err != nil && err != model.ErrNoResult {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}
		// Delete all comments before deleting of the project
		var comments comments
		err = json.Unmarshal(project.CommentsId, &comments)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		for _, id := range comments.Comment {
			go func() {
				model.DeleteComment(id)
			}()
		}

		err = model.DeleteProject(uint32(ID))
		if err != nil && err != model.ErrNoResult {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		sess.AddFlash(view.Flash{"הפרויקט נדחה", view.FlashWarning})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/manager_waiting", http.StatusFound)

	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ListRunningProjectsGET displays the list of running projects after all approval processes
func ListRunningProjectsGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get user session
	if sess.Values["id"] != nil {

		// Display the view
		v := view.New(r)
		v.Name = "project/running_project"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		// Get user ID form session
		ID, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)

		// Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		user, err := model.SelectUserInfo(uint32(ID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		var projects []model.Project
		if user.Role == 2 {
			projects, err = model.ListRunningBySupervisorProjects(user.ID)
		} else if user.Role == 3 {
			projects, err = model.ListRunningProjects()
		}

		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		v.Vars["projects"] = projects
		var supervisors []string

		for i, project := range projects {

			//HACK NOT USE IT AGAIN
			var Students model.Students
			err := json.Unmarshal(project.StudentsId, &Students)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			projects[i].Description = student1.FirstName + " " + student1.LastName + " ו" + student2.FirstName + " " + student2.LastName
		}

		v.Vars["supervisors"] = supervisors

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SearchRunningProjectPOST for project with status running
func SearchRunningProjectPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {

		// Display the view
		v := view.New(r)
		v.Name = "project/running_project"

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
			ListProjectsGET(w, r)
			return
		}

		query := r.FormValue("query")

		// Get user ID form session
		ID, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)

		// Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		user, err := model.SelectUserInfo(uint32(ID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		var projects []model.Project
		if user.Role == 2 {
			projects, err = model.SearchRunningProjectBySupervisor(user.ID, query)
		} else if user.Role == 3 {
			projects, err = model.SearchRunningProject(query)
		}

		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ListProjectsGET(w, r)
			return
		}

		v.Vars["projects"] = projects
		var supervisors []string

		for i, project := range projects {

			//HACK NOT USE IT AGAIN
			var Students model.Students
			err := json.Unmarshal(project.StudentsId, &Students)
			if err != nil {
				log.Println(err)
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			student1, err := model.SelectUserInfo(uint32(Students.StudentId1))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			student2, err := model.SelectUserInfo(uint32(Students.StudentId2))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				ListProjectsGET(w, r)
				return
			}

			projects[i].Description = student1.FirstName + " " + student1.LastName + " ו" + student2.FirstName + " " + student2.LastName
		}

		v.Vars["supervisors"] = supervisors

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// UploadPDFFilePOST upload pdf file of book to AWS
func UploadPDFFilePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

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

	var projectFiles model.ProjectFiles
	err = json.Unmarshal(project.Files, &projectFiles)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	// Upload file
	// upload of 10 MB files.
	r.ParseMultipartForm(10 * 1024 * 1024)

	file, handler, err := r.FormFile("File")
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	defer file.Close()
	extension := filepath.Ext(handler.Filename)
	if aws.Contains(aws.AllowTypes, handler.Header.Get("Content-Type")) && aws.Contains(aws.AllowExtensions, extension) {
		sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה בשנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	path, err := aws.FileToProject(file, handler)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	var newFile uint32
	if projectFiles.Book1PDF == 0 {
		newFile, err = model.CreateFile(path, 1)
		projectFiles.Book1PDF = int(newFile)
		if projectFiles.Book1PDF != 0 && projectFiles.Book1WORD != 0 && projectFiles.Presentation1 != 0 {
			go func() {
				err = model.SendToProjectManager(21, 20)
				if err != nil {
					log.Println(err)
				}
			}()
		}
	} else if projectFiles.Book2PDF == 0 && projectFiles.Book1PDF != 0 {
		newFile, err = model.CreateFile(path, 2)
		projectFiles.Book2PDF = int(newFile)
		if projectFiles.Book2PDF != 0 && projectFiles.Book2WORD != 0 && projectFiles.Presentation2 != 0 && projectFiles.SourceCode != 0 {
			go func() {
				err = model.SendToProjectManager(23, 22)
				if err != nil {
					log.Println(err)
				}
			}()
			err = model.MarkProjectAsDone(project.ID)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/project/projects", http.StatusFound)
				return
			}
			err = model.CreateArchiveProject(project)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/project/projects", http.StatusFound)
				return
			}
		}
	}
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	project.Files, err = json.Marshal(projectFiles)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	err = model.AddFileToListOfFiles(project.ID, project.Files)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	sess.AddFlash(view.Flash{"הקובץ עלה בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)

}

// UploadDOCFilePOST upload word file of book to AWS
func UploadDOCFilePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

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

	var projectFiles model.ProjectFiles
	err = json.Unmarshal(project.Files, &projectFiles)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	// Upload file
	// upload of 10 MB files.
	r.ParseMultipartForm(10 * 1024 * 1024)

	file, handler, err := r.FormFile("File")
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	defer file.Close()
	extension := filepath.Ext(handler.Filename)
	if aws.Contains(aws.AllowTypes, handler.Header.Get("Content-Type")) && aws.Contains(aws.AllowExtensions, extension) {
		sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה בשנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	path, err := aws.FileToProject(file, handler)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	var newFile uint32
	if projectFiles.Book1WORD == 0 {
		newFile, err = model.CreateFile(path, 1)
		projectFiles.Book1WORD = int(newFile)
		if projectFiles.Book1PDF != 0 && projectFiles.Book1WORD != 0 && projectFiles.Presentation1 != 0 {
			go func() {
				err = model.SendToProjectManager(21, 20)
				if err != nil {
					log.Println(err)
				}
			}()
		}
	} else if projectFiles.Book2WORD == 0 && projectFiles.Book1WORD != 0 {
		newFile, err = model.CreateFile(path, 2)
		projectFiles.Book2WORD = int(newFile)
		if projectFiles.Book2PDF != 0 && projectFiles.Book2WORD != 0 && projectFiles.Presentation2 != 0 && projectFiles.SourceCode != 0 {
			go func() {
				err = model.SendToProjectManager(23, 22)
				if err != nil {
					log.Println(err)
				}
			}()
			err = model.MarkProjectAsDone(project.ID)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/project/projects", http.StatusFound)
				return
			}
			err = model.CreateArchiveProject(project)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/project/projects", http.StatusFound)
				return
			}
		}
	}
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	project.Files, err = json.Marshal(projectFiles)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	err = model.AddFileToListOfFiles(project.ID, project.Files)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	sess.AddFlash(view.Flash{"הקובץ עלה בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)

}

// UploadPowerPointFilePOST upload presentation file to AWS
func UploadPowerPointFilePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

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

	var projectFiles model.ProjectFiles
	err = json.Unmarshal(project.Files, &projectFiles)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	// Upload file
	// upload of 10 MB files.
	r.ParseMultipartForm(10 * 1024 * 1024)

	file, handler, err := r.FormFile("File")
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	defer file.Close()
	extension := filepath.Ext(handler.Filename)
	if aws.Contains(aws.AllowTypes, handler.Header.Get("Content-Type")) && aws.Contains(aws.AllowExtensions, extension) {
		sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה בשנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	path, err := aws.FileToProject(file, handler)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	var newFile uint32
	if projectFiles.Presentation1 == 0 {
		newFile, err = model.CreateFile(path, 3)
		projectFiles.Presentation1 = int(newFile)
		if projectFiles.Book1PDF != 0 && projectFiles.Book1WORD != 0 && projectFiles.Presentation1 != 0 {
			go func() {
				err = model.SendToProjectManager(21, 20)
				if err != nil {
					log.Println(err)
				}
			}()
		}
	} else if projectFiles.Presentation2 == 0 && projectFiles.Presentation1 != 0 {
		newFile, err = model.CreateFile(path, 4)
		projectFiles.Presentation2 = int(newFile)
		if projectFiles.Book2PDF != 0 && projectFiles.Book2WORD != 0 && projectFiles.Presentation2 != 0 && projectFiles.SourceCode != 0 {
			go func() {
				err = model.SendToProjectManager(23, 22)
				if err != nil {
					log.Println(err)
				}
			}()
			err = model.MarkProjectAsDone(project.ID)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/project/projects", http.StatusFound)
				return
			}
			err = model.CreateArchiveProject(project)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/project/projects", http.StatusFound)
				return
			}
		}
	}
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	project.Files, err = json.Marshal(projectFiles)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	err = model.AddFileToListOfFiles(project.ID, project.Files)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	sess.AddFlash(view.Flash{"הקובץ עלה בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)

}

// UploadCodeFilePOST upload source code file to AWS
func UploadCodeFilePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

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

	var projectFiles model.ProjectFiles
	err = json.Unmarshal(project.Files, &projectFiles)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	// Upload file
	// upload of 10 MB files.
	r.ParseMultipartForm(10 * 1024 * 1024)

	file, handler, err := r.FormFile("File")
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	defer file.Close()
	extension := filepath.Ext(handler.Filename)
	if aws.Contains(aws.AllowTypes, handler.Header.Get("Content-Type")) && aws.Contains(aws.AllowExtensions, extension) {
		sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה בשנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	path, err := aws.FileToProject(file, handler)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	newFile, err := model.CreateFile(path, 5)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	projectFiles.SourceCode = int(newFile)
	project.Files, err = json.Marshal(projectFiles)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	if projectFiles.Book2PDF != 0 && projectFiles.Book2WORD != 0 && projectFiles.Presentation2 != 0 && projectFiles.SourceCode != 0 {
		go func() {
			err = model.SendToProjectManager(23, 22)
			if err != nil {
				log.Println(err)
			}
		}()
		err = model.MarkProjectAsDone(project.ID)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/project/projects", http.StatusFound)
			return
		}
		err = model.CreateArchiveProject(project)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/project/projects", http.StatusFound)
			return
		}
	}

	err = model.AddFileToListOfFiles(project.ID, project.Files)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/project/projects", http.StatusFound)
		return
	}

	sess.AddFlash(view.Flash{"הקובץ עלה בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/project/projects", http.StatusFound)

}
