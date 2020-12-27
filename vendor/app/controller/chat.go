package controller

import (
	"app/model"
	"app/shared/aws"
	"app/shared/session"
	"app/shared/view"
	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// ChatGET handles the chat page
func ChatGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "chat/chat"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		if v.Vars["is_student"] != nil {
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
				if err == model.ErrNoResult {
					log.Println(err)
					sess.AddFlash(view.Flash{"לא נפתח עבורך צ'אט, עלייך תחילה להירשם לפרויקט", view.FlashError})
					sess.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)
					return
				} else {
					log.Println(err)
					sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
					sess.Save(r, w)
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
			}

			if project.StatusId == 4 {
				//Get list of all message
				chat, err := model.SelectChat(project.ID, uint32(userId))
				if err != nil {
					log.Println(err)
					chat = []model.Chat{}
				}

				c := map[string]interface{}{
					"key": "value",
				}

				var chats []interface{}
				for _, m := range chat {
					author, _ := model.SelectUserInfo(m.SendBy)
					a := author.FirstName + " " + author.LastName
					var f bool
					if strings.HasPrefix(m.Message, "http://fs.spms-project.info/chat/") {
						f = true
					} else {
						f = false
					}
					c["chat"] = map[string]interface{}{
						"chat":    m,
						"author":  a,
						"is_file": f,
					}
					chats = append(chats, c["chat"])
				}
				v.Vars["chat"] = chats
				v.Vars["current_user"] = userId
				//v.Vars["is_ready"] = true
			} else {
				log.Println(err)
				sess.AddFlash(view.Flash{"לא נפתח עבורך צ'אט, עלייך תחילה להירשם לפרויקט", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

		if v.Vars["is_supervisor"] != nil || v.Vars["is_project_manager"] != nil {
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

			//Get list of all message
			chat, err := model.SelectChat(uint32(projectID), uint32(userId))
			if err != nil {
				log.Println(err)
				chat = []model.Chat{}
			}

			c := map[string]interface{}{
				"key": "value",
			}

			var chats []interface{}
			for _, m := range chat {
				author, _ := model.SelectUserInfo(m.SendBy)
				a := author.FirstName + " " + author.LastName
				var f bool
				if strings.HasPrefix(m.Message, "http://fs.spms-project.info/chat/") {
					f = true
				} else {
					f = false
				}
				c["chat"] = map[string]interface{}{
					"chat":    m,
					"author":  a,
					"is_file": f,
				}
				chats = append(chats, c["chat"])
			}

			v.Vars["chat"] = chats
			v.Vars["current_user"] = userId
		}

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ChatListGET handles the chat list page
func ChatListGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "chat/chat_list"

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

		projects, err := model.ListRunningBySupervisorProjects(uint32(userId))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		v.Vars["userId"] = userId
		v.Vars["projects"] = projects

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ChatPOST submit to the DB the new message
func ChatPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	if sess.Values["is_student"] != nil {
		// Validate with required fields
		if validate, _ := view.Validate(r, []string{"message"}); !validate {
			sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

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

		message := r.FormValue("message")
		messageID, err := model.CreateNewMessage(message, uint32(userId), project.ID)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		go func() {
			err = model.SendChatNotificationToStudentSupervisor(project.ID, uint32(userId), 3, 2, uint32(messageID))
			if err != nil {
				log.Println(err)
			}
		}()

		sess.Save(r, w)
		http.Redirect(w, r, "/chat/", http.StatusFound)

	} else if sess.Values["is_project_manager"] != nil || sess.Values["is_supervisor"] != nil {
		// Validate with required fields
		if validate, _ := view.Validate(r, []string{"message"}); !validate {
			sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		userId, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
		//Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

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

		message := r.FormValue("message")
		messageID, err := model.CreateNewMessage(message, uint32(userId), uint32(projectID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		go func() {
			err = model.SendChatNotificationToStudents(uint32(projectID), 3, 2, uint32(messageID))
			if err != nil {
				log.Println(err)
			}
		}()

		sess.Save(r, w)
		http.Redirect(w, r, "/chat/chat/"+strconv.Itoa(int(projectID)), http.StatusFound)
	}
}

// SearchChatListPOST handles the search for supervisor list chat
func SearchChatListPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "chat/chat_list"

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
			ChatListGET(w, r)
			return
		}

		query := r.FormValue("query")

		userId, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
		//Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		//Get list of projects by search
		projects, err := model.SearchRunningProjectBySupervisor(uint32(userId), query)
		if err != nil {
			log.Println(err)
			projects = []model.Project{}
		}
		v.Vars["projects"] = projects

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// UploadChatFilePOST check the file and upload it to the AWS S3
func UploadChatFilePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	if sess.Values["is_student"] != nil {

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

		// Upload file
		// upload of 5 MB files.
		r.ParseMultipartForm(5 * 1024 * 1024)

		file, handler, err := r.FormFile("File")
		if err != nil {
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		defer file.Close()
		extension := filepath.Ext(handler.Filename)
		if aws.Contains(aws.AllowTypes, handler.Header.Get("Content-Type")) && aws.Contains(aws.AllowExtensions, extension) {
			sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה בשנית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		path, err := aws.FileToChat(file, handler)
		if err != nil {
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		messageID, err := model.CreateNewMessage(path, uint32(userId), project.ID)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		go func() {
			err = model.SendChatNotificationToStudentSupervisor(project.ID, uint32(userId), 3, 2, uint32(messageID))
			if err != nil {
				log.Println(err)
			}
		}()

		sess.Save(r, w)
		http.Redirect(w, r, "/chat/", http.StatusFound)

	} else if sess.Values["is_project_manager"] != nil || sess.Values["is_supervisor"] != nil {

		userId, err := strconv.ParseInt(sess.Values["id"].(string), 10, 32)
		//Check if the number converted well
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

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

		// Upload file
		// upload of 5 MB files.
		r.ParseMultipartForm(5 * 1024 * 1024)

		file, handler, err := r.FormFile("File")
		if err != nil {
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		defer file.Close()
		extension := filepath.Ext(handler.Filename)
		if aws.Contains(aws.AllowTypes, handler.Header.Get("Content-Type")) && aws.Contains(aws.AllowExtensions, extension) {
			sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה בשנית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		path, err := aws.FileToChat(file, handler)
		if err != nil {
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		messageID, err := model.CreateNewMessage(path, uint32(userId), uint32(projectID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			ChatGET(w, r)
			return
		}

		go func() {
			err = model.SendChatNotificationToStudents(uint32(projectID), 3, 2, uint32(messageID))
			if err != nil {
				log.Println(err)
			}
		}()

		sess.Save(r, w)
		http.Redirect(w, r, "/chat/chat/"+strconv.Itoa(int(projectID)), http.StatusFound)
	}
}
