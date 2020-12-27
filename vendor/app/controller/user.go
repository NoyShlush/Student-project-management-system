package controller

import (
	"app/model"
	sendemail "app/shared/email"
	"app/shared/passhash"
	"app/shared/session"
	"app/shared/view"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO: Next remove
type contactInfo struct {
	Telephone int    `json:"Telephone"`
	Facebook  string `json:"Facebook"`
	Twitter   string `json:"Twitter"`
	Telegram  string `json:"Telegram"`
}

const (
	// Name of the session variable that tracks login attempts
	sessLoginAttempt      = "login_attempt"
	changePasswordAttempt = "change_password_attempt"
)

// loginAttempt increments the number of login attempts in sessions variable
func loginAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values[sessLoginAttempt] == nil {
		sess.Values[sessLoginAttempt] = 1
	} else {
		sess.Values[sessLoginAttempt] = sess.Values[sessLoginAttempt].(int) + 1
	}
}

// ChangePasswordAttempt increments the number of change password attempts in sessions variable
func ChangePasswordAttempt(sess *sessions.Session) {
	// Log the attempt
	if sess.Values[changePasswordAttempt] == nil {
		sess.Values[changePasswordAttempt] = 1
	} else {
		sess.Values[changePasswordAttempt] = sess.Values[changePasswordAttempt].(int) + 1
	}
}

// StudentManagementGET displays list of all student
func StudentManagementGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/admin_student"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		//Gett list of all students
		users, err := model.SelectListOfUsers(1)
		if err != nil {
			log.Println(err)
			users = []model.User{}
		}
		v.Vars["users"] = users

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// CreateUserGET displays the register page
func CreateUserGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/create"
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		// Refill any form fields
		view.Repopulate([]string{"id_number", "email"}, r.Form, v.Vars)
		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// CreateStudentPOST handles the registration form submission
func CreateStudentPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"id_number", "email"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		CreateUserGET(w, r)
		return
	}

	// Get form values
	idNumber, err := strconv.ParseUint(r.FormValue("id_number"), 10, 32)
	//Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		CreateUserGET(w, r)
		return
	}

	email := r.FormValue("email")
	tempPassword, _ := uuid.NewRandom()
	password, _ := passhash.HashString(tempPassword.String())

	// If password hashing failed
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		CreateUserGET(w, r)
		return
	}

	// Get database result
	_, erre := model.CheckIfTheUserExists(uint32(idNumber), email)

	if erre == model.ErrNoResult { // If success (no user exists with that email)
		//Create user with student role
		userId, ex := model.CreateNewUser(uint32(idNumber), email, password, 1)
		// Will only error if there is a problem with the query
		if ex != nil {
			log.Println(ex)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
		} else {
			token, err := model.CreateChangePasswordToken(userId)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאה בשרת", view.FlashError})
				sess.Save(r, w)
				ForgetPasswordGET(w, r)
				return
			}
			link := "http://" + r.Host + "/chnage_password/" + token
			go func() {
				err = sendemail.SendEmail(email, "", "<p align='right'>:לכניסה ראשונה לחץ על הקישור </p>"+link)
			}()
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאה בשרת", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/admin/student_management", http.StatusFound)
				return
			}
			sess.AddFlash(view.Flash{"החשבון עבור " + email + " נוצר בהצלחה ", view.FlashSuccess})
			sess.Save(r, w)
			http.Redirect(w, r, "/admin/student_management", http.StatusFound)
			return
		}
	} else if erre != nil && erre != model.ErrNoResult { // Catch all other errors
		log.Println(erre)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
	} else { // Else the user already exists
		if erre != model.ErrNoResult {
			sess.AddFlash(view.Flash{"החשבון קיים במערכת אנא בדוק את מספר תעודת הזהות או מייל ", view.FlashError})
			sess.Save(r, w)
		}
	}

	// Display the page
	CreateUserGET(w, r)
}

// UpdateUserGET displays the update user page
func UpdateUserGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	userID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tempUser, err := model.SelectUser(uint32(userID))
	if err != nil && err != model.ErrNoResult {
		log.Println(err)
		sess.AddFlash(view.Flash{"משתמש זה אינו קיים במערכת", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "admin/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["is_auth"] = sess.Values["is_auth"]
	v.Vars["is_student"] = sess.Values["is_student"]
	v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
	v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

	// Fills the existing fields on the view
	v.Vars["idnumber"] = &tempUser.IdNumber
	v.Vars["email"] = &tempUser.Email
	v.Render(w)
}

// UpdateUserPOST handles the update form to update user data
func UpdateUserPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"idnumber", "email"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		UpdateUserGET(w, r)
		return
	}

	//Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	ID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		UpdateUserGET(w, r)
	}

	// Get form values
	email := r.FormValue("email")
	idNumber, err := strconv.ParseUint(r.FormValue("idnumber"), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		UpdateUserGET(w, r)
		return
	}

	// Get the real user data before the update
	realUser, err := model.SelectUser(uint32(ID))
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		UpdateUserGET(w, r)
		return
	}

	if realUser.Email != email && realUser.IdNumber == uint32(idNumber) {
		_, err := model.CheckIfTheUserExists(0, email)
		if err == model.ErrNoResult {
			err = model.UpdateUser(uint32(ID), uint32(idNumber), email)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				UpdateUserGET(w, r)
				return
			}
		} else {
			sess.AddFlash(view.Flash{"החשבון קיים במערכת אנא בדוק את מספר תעודת הזהות או מייל", view.FlashError})
			sess.Save(r, w)
			UpdateUserGET(w, r)
			return
		}
	} else if realUser.IdNumber != uint32(idNumber) && realUser.Email == email {
		_, err := model.CheckIfTheUserExists(uint32(idNumber), "")
		if err == model.ErrNoResult {
			err = model.UpdateUser(uint32(ID), uint32(idNumber), email)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				UpdateUserGET(w, r)
				return
			}
		} else {
			sess.AddFlash(view.Flash{"החשבון קיים במערכת אנא בדוק את מספר תעודת הזהות או מייל", view.FlashError})
			sess.Save(r, w)
			UpdateUserGET(w, r)
			return
		}
	} else if realUser.IdNumber != uint32(idNumber) && realUser.Email != email {
		_, err := model.CheckIfTheUserExists(uint32(idNumber), email)
		if err == model.ErrNoResult {
			err = model.UpdateUser(uint32(ID), uint32(idNumber), email)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
				sess.Save(r, w)
				UpdateUserGET(w, r)
				return
			}
		} else {
			sess.AddFlash(view.Flash{"החשבון קיים במערכת אנא בדוק את מספר תעודת הזהות או מייל", view.FlashError})
			sess.Save(r, w)
			UpdateUserGET(w, r)
			return
		}
	} else {
		sess.AddFlash(view.Flash{"לא עודכנו פרטים חדשים", view.FlashWarning})
		sess.Save(r, w)
		rexPath := regexp.MustCompile(`/[^/]*/[^/]*$`)
		http.Redirect(w, r, rexPath.ReplaceAllString(r.RequestURI, ""), http.StatusFound)
	}

	sess.AddFlash(view.Flash{"החשבון עבור " + email + " עודכן בהצלחה ", view.FlashSuccess})
	sess.Save(r, w)
	rexPath := regexp.MustCompile(`/[^/]*/[^/]*$`)
	http.Redirect(w, r, rexPath.ReplaceAllString(r.RequestURI, ""), http.StatusFound)
	return

	// Display the page
	UpdateUserGET(w, r)
}

// SearchStudentPOST handles the search for student
func SearchStudentPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/admin_student"

		//View permissions
		//TODO: ANDREY CHECK CSRF TOKEN
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
			StudentManagementGET(w, r)
			return
		}

		query := r.FormValue("query")

		//Get list of students by search
		users, err := model.SearchUser(query, 1)
		if err != nil {
			log.Println(err)
			users = []model.User{}
		}
		v.Vars["users"] = users

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// BlockUserGET handles the block request of the user
func BlockUserGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	if sess.Values["id"] != nil {
		// Get values from URI
		var params httprouter.Params
		params = context.Get(r, "params").(httprouter.Params)
		userID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			rexPath := regexp.MustCompile(`/[^/]*/[^/]*$`)
			http.Redirect(w, r, rexPath.ReplaceAllString(r.RequestURI, ""), http.StatusFound)
			return
		}

		user, err := model.SelectUser(uint32(userID))
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"לא קיים משתמש כזה במערכת", view.FlashError})
			sess.Save(r, w)
			rexPath := regexp.MustCompile(`/[^/]*/[^/]*$`)
			http.Redirect(w, r, rexPath.ReplaceAllString(r.RequestURI, ""), http.StatusFound)
			return
		}

		//Check if the user is blocked
		if user.Block != true {
			err = model.BlockUser(uint32(userID))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"לא ניתן לחסום משתמש זה", view.FlashError})
				sess.Save(r, w)
				rexPath := regexp.MustCompile(`/[^/]*/[^/]*$`)
				http.Redirect(w, r, rexPath.ReplaceAllString(r.RequestURI, ""), http.StatusFound)
				return
			}

			sess.AddFlash(view.Flash{"המשתמש נחסם", view.FlashSuccess})
			sess.Save(r, w)
			rexPath := regexp.MustCompile(`/[^/]*/[^/]*$`)
			http.Redirect(w, r, rexPath.ReplaceAllString(r.RequestURI, ""), http.StatusFound)
		} else {
			err = model.UnBlockUser(uint32(userID))
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"לא ניתן לחסום משתמש זה", view.FlashError})
				sess.Save(r, w)
				rexPath := regexp.MustCompile(`/[^/]*/[^/]*$`)
				http.Redirect(w, r, rexPath.ReplaceAllString(r.RequestURI, ""), http.StatusFound)
				return
			}

			sess.AddFlash(view.Flash{"המשתמש שוחרר", view.FlashSuccess})
			sess.Save(r, w)
			rexPath := regexp.MustCompile(`/[^/]*/[^/]*$`)
			http.Redirect(w, r, rexPath.ReplaceAllString(r.RequestURI, ""), http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// UploadStudentFilePOST handles the file upload
func UploadStudentFilePOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// upload of 1 MB files.
	r.ParseMultipartForm(1 << 2)

	file, handler, err := r.FormFile("File")
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאה בשרת אנא נסה לעלות שנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/student_management", http.StatusFound)
		return
	}

	defer file.Close()
	extension := filepath.Ext(handler.Filename)
	if handler.Header.Get("Content-Type") == "application/vnd.ms-excel" && extension != ".csv" {
		sess.AddFlash(view.Flash{"קובץ אינו תקיו, אנא נסה שנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/student_management", http.StatusFound)
		return
	}

	// Create a temporary file within our temp directory that follows
	//TODO: ADD SYSTEM DETECTION
	tempFile, err := ioutil.TempFile("C:\\temp-file", "upload-*.csv")
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאה בשרת, אנא נסה להעלות את הקובץ שנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/student_management", http.StatusFound)
		return
	}

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאה בשרת, אנא נסה להעלות את הקובץ שנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/student_management", http.StatusFound)
		return
	}

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	csvFile, err := os.Open(tempFile.Name())
	if err != nil {
		sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה שנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/student_management", http.StatusFound)
		return
	}

	table := csv.NewReader(csvFile)
	lines, err := table.ReadAll()
	if err != nil {
		log.Println("error reading all lines: %v", err)
	}

	idNumber := make([]string, len(lines)-1)
	email := make([]string, len(lines)-1)
	for i, line := range lines {
		if i == 0 {
			// skip header line
			continue
		}
		idNumber[i-1] = line[0]
		email[i-1] = line[1]
	}

	var wg sync.WaitGroup
	var listOfExisting []string

	for z, _ := range idNumber {
		// Increment the wait group counter
		wg.Add(1)

		go func() {
			// Decrement the counter when the go routine completes
			defer wg.Done()
			// Call the function check
			idTemp, _ := strconv.ParseUint(idNumber[z], 10, 32)
			_, erre := model.CheckIfTheUserExists(uint32(idTemp), email[z])
			if erre == nil {
				listOfExisting = append(listOfExisting, strconv.FormatInt(int64(z+2), 10))
			}
		}()
		// Wait for all the checkWebsite calls to finish
		wg.Wait()
	}

	if len(listOfExisting) != 0 {
		list := strings.Join(listOfExisting, ", ")
		sess.AddFlash(view.Flash{"החשבונות עבור שורות " + list + " קיימים כבר ", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/student_management", http.StatusFound)
		return
	}

	for j, _ := range idNumber {
		// Increment the wait group counter
		wg.Add(1)
		go func() {
			// Decrement the counter when the go routine completes
			defer wg.Done()
			// Call the function check
			idNumber, _ := strconv.ParseUint(idNumber[j], 10, 32)
			tempPassword, _ := uuid.NewRandom()
			password, _ := passhash.HashString(tempPassword.String())
			model.CreateNewUser(uint32(idNumber), email[j], password, 1)
		}()
		// Wait for all the checkWebsite calls to finish
		wg.Wait()
	}

	sess.AddFlash(view.Flash{"כל המשתמשים נוצרו בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/admin/student_management", http.StatusFound)
}

// WaitingStudentGET handles the waiting student list
func WaitingStudentGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/wait"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		//Gett list of all students
		users, err := model.SelectListOfWaiting(1)
		if err != nil {
			log.Println(err)
			users = []model.User{}
		}
		v.Vars["users"] = users

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SearchWaitingStudentPOST handles the search for waiting student
func SearchWaitingStudentPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/wait"

		//View permissions
		//TODO: ANDREY CHECK CSRF TOKEN
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
			WaitingStudentGET(w, r)
			return
		}

		query := r.FormValue("query")

		//Get list of students by search
		users, err := model.SearchWaitingUser(query, 1)
		if err != nil {
			log.Println(err)
			users = []model.User{}
		}
		v.Vars["users"] = users

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// LoginGET displays the login page
func LoginGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "login/login"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)

	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

// LoginPOST handles the login form submission
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values[sessLoginAttempt] != nil && sess.Values[sessLoginAttempt].(int) >= 3 {
		log.Println("Brute force login prevented")
		sess.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"email", "password"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get database result
	result, err := model.SelectUserByEmail(email)

	// Determine if user exists
	if err == model.ErrNoResult {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"סיסמא או מייל שגויים- נסיון " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"התרחשה שגיאה, אנא נסה שנית מאוחר יותר", view.FlashError})
		sess.Save(r, w)
	} else if passhash.MatchString(result.HashPassword, password) {
		if result.Block == true {
			// User inactive and display inactive message
			sess.AddFlash(view.Flash{"המשתמש חסום, אנא פנה למנהל מערכת", view.FlashNotice})
			sess.Save(r, w)
		} else {
			// Login successfully
			session.Empty(sess)
			sess.AddFlash(view.Flash{"התחברת בהצלחה!", view.FlashSuccess})
			sess.Values["id"] = result.UserID()
			sess.Values["email"] = email
			sess.Values["first_name"] = result.FirstName
			sess.Values["is_auth"] = true
			switch result.Role {
			case 1:
				sess.Values["is_student"] = result.Role
			case 2:
				sess.Values["is_supervisor"] = result.Role
			case 3:
				sess.Values["is_project_manager"] = result.Role
			default:
				sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashWarning})
			}
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		loginAttempt(sess)
		sess.AddFlash(view.Flash{"סיסמא או מייל שגויים- נסיון " + fmt.Sprintf("%v", sess.Values[sessLoginAttempt]), view.FlashWarning})
		sess.Save(r, w)
	}

	// Show the login page again
	LoginGET(w, r)
}

// LogoutGET clears the session and logs the user out
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// If user is authenticated
	if sess.Values["id"] != nil {
		session.Empty(sess)
		sess.AddFlash(view.Flash{"להתראות!", view.FlashNotice})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// ForgetPasswordGET displays the forget password page
func ForgetPasswordGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "password/forget"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)

	// Refill any form fields
	view.Repopulate([]string{"email"}, r.Form, v.Vars)
	v.Render(w)
}

// ForgetPasswordPOST handles the chnage password form submission
func ForgetPasswordPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	if sess.Values[changePasswordAttempt] != nil && sess.Values[changePasswordAttempt].(int) >= 3 {
		log.Println("Brute force login prevented")
		sess.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
		sess.Save(r, w)
		ForgetPasswordGET(w, r)
		return
	}

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"email"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		ForgetPasswordGET(w, r)
		return
	}

	// Form values
	email := r.FormValue("email")

	// Get database result
	result, err := model.SelectUserByEmail(email)

	// Determine if user exists
	if err == model.ErrNoResult {
		ChangePasswordAttempt(sess)
		sess.AddFlash(view.Flash{"לא קיים משתמש כזה- נסיון " + fmt.Sprintf("%v", sess.Values[changePasswordAttempt]), view.FlashWarning})
		sess.Save(r, w)
	} else if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"התרחשה שגיאה, אנא נסה שנית מאוחר יותר", view.FlashError})
		sess.Save(r, w)
	} else if result.Block == true {
		log.Println("The user is blocked " + result.Email)
		sess.AddFlash(view.Flash{"המשתמש שלך חסום, אנא פנה למנהל מערכת", view.FlashError})
		sess.Save(r, w)
	} else {
		//Chnage password
		token, err := model.CreateChangePasswordToken(result.ID)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאה בשרת", view.FlashError})
			sess.Save(r, w)
			ForgetPasswordGET(w, r)
			return
		}
		link := "http://" + r.Host + "/chnage_password/" + token
		go func() {
			err = sendemail.SendEmail(result.Email, "שינוי סיסמא למערכת ניהול פרויקטי גמר", "<p align='right'>:הקישור לשינוי סיסמא הוא</p>"+link)
		}()
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאה בשרת", view.FlashError})
			sess.Save(r, w)
			ForgetPasswordGET(w, r)
			return
		}
		session.Empty(sess)
		sess.AddFlash(view.Flash{"הקישור לשינוי סיסמא נשלח לכתובת המייל", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Show the page
	ForgetPasswordGET(w, r)
}

// ChangePasswordGET displays the chnage password page
func ChangePasswordGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	token := params.ByName("token")
	_, validUntil, err := model.CheckUserToken(token)
	if err != nil {
		sess.AddFlash(view.Flash{"הקישור אינו תקין, אנא התחל תהליך מחדש", view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	} else if time.Now().After(validUntil) {
		sess.AddFlash(view.Flash{"פג תוקף, אנא התחל תהליך מחדש", view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	} else {
		// Display the view
		v := view.New(r)
		v.Name = "password/change"
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Render(w)
	}

}

// ChangePasswordPOST handles the chnage password form submission
func ChangePasswordPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	token := params.ByName("token")
	id, _, err := model.CheckUserToken(token)
	if err != nil {
		sess.AddFlash(view.Flash{"הקישור אינו תקין, אנא התחל תהליך מחדש", view.FlashError})
		sess.Save(r, w)
		LoginGET(w, r)
		return
	}

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"password", "password2"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		ChangePasswordGET(w, r)
		return
	}

	// Form values
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	if password != password2 {
		sess.AddFlash(view.Flash{"הסיסמא אינה תואמת", view.FlashError})
		sess.Save(r, w)
		ChangePasswordGET(w, r)
		return
	}

	password, err = passhash.HashString(password)

	// If password hashing failed
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		ChangePasswordGET(w, r)
		return
	}

	// Get database result
	err = model.UpdateChangePassword(id, password)

	if err != nil {
		// Display error message
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
	} else {
		// Login successfully
		session.Empty(sess)
		sess.AddFlash(view.Flash{"הסיסמא שונתה בהצלחה", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/user/update/"+token, http.StatusFound) //update info page with token
		return
	}

	// Show the update info page again
	http.Redirect(w, r, "/user/update/"+token, http.StatusFound) //update info page with token
}

// UpdateInfoGET displays the update contact information page
func UpdateInfoGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// Get values from session
	id := sess.Values["id"]
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", id), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tempUser, err := model.SelectUserInfo(uint32(userID))
	if err != nil && err != model.ErrNoResult {
		log.Println(err)
		sess.AddFlash(view.Flash{"משתמש זה אינו קיים במערכת", view.FlashError})
		sess.Save(r, w)
		return
	} else if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "user/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["is_auth"] = sess.Values["is_auth"]
	v.Vars["is_student"] = sess.Values["is_student"]
	v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
	v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

	var contactInformation contactInfo
	err = json.Unmarshal(tempUser.ContactInfomation, &contactInformation)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound) //MainPage
		return
	}

	// Fills the existing fields on the view
	v.Vars["firstname"] = &tempUser.FirstName
	v.Vars["lastname"] = &tempUser.LastName
	if contactInformation.Telephone == 0 {
		v.Vars["telephone"] = ""
	} else {
		v.Vars["telephone"] = &contactInformation.Telephone
	}
	v.Vars["facebook"] = &contactInformation.Facebook
	v.Vars["twitter"] = &contactInformation.Twitter
	v.Vars["telegram"] = &contactInformation.Telegram

	v.Render(w)

}

// UpdateInfoPOST handles the update contact information data
func UpdateInfoPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"firstname", "lastname"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		UpdateInfoGET(w, r)
		return
	}

	// Get values from session
	id := sess.Values["id"]
	ID, err := strconv.ParseUint(fmt.Sprintf("%v", id), 10, 32)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		UpdateInfoGET(w, r)
		return
	}

	// Get form values
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	facebook := r.FormValue("facebook")
	twitter := r.FormValue("twitter")
	telegram := r.FormValue("telegram")

	telephone, err := strconv.ParseInt(r.FormValue("telephone"), 10, 32)

	contactInformation := contactInfo{
		int(telephone),
		facebook,
		twitter,
		telegram,
	}

	contactInfoJson, err := json.Marshal(contactInformation)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		UpdateInfoGET(w, r)
		return
	}

	err = model.UpdateInfo(uint32(ID), firstName, lastName, contactInfoJson)

	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
	}

	sess.AddFlash(view.Flash{"פרטייך עודכנו בהצלחה", view.FlashSuccess})
	sess.Save(r, w)

	// Display the page
	http.Redirect(w, r, "/", http.StatusFound) //MainPage
}

// FirstInfoUpdateGET displays the update contact information page
func FirstInfoUpdateGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	token := params.ByName("token")
	id, _, err := model.CheckUserToken(token)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tempUser, err := model.SelectUserInfo(id)
	if err != nil && err == model.ErrNoResult {
		log.Println(err)
		sess.AddFlash(view.Flash{"משתמש זה אינו קיים במערכת", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		tempUser, err = model.SelectUser(id)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	// Display the view
	v := view.New(r)
	v.Name = "user/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["first_name"] = sess.Values["first_name"]
	v.Vars["is_auth"] = sess.Values["is_auth"]
	v.Vars["is_student"] = sess.Values["is_student"]
	v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
	v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

	if tempUser.ContactInfomation != nil {
		var contactInformation contactInfo
		err = json.Unmarshal(tempUser.ContactInfomation, &contactInformation)
		if err != nil {
			log.Println(err)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound) //MainPage
			return
		}

		// Fills the existing fields on the view
		v.Vars["firstname"] = &tempUser.FirstName
		v.Vars["lastname"] = &tempUser.LastName
		if contactInformation.Telephone == 0 {
			v.Vars["telephone"] = ""
		} else {
			v.Vars["telephone"] = &contactInformation.Telephone
		}
		v.Vars["facebook"] = &contactInformation.Facebook
		v.Vars["twitter"] = &contactInformation.Twitter
		v.Vars["telegram"] = &contactInformation.Telegram
	}

	v.Render(w)

}

// FirstInfoUpdatePOST handles the update contact information data
func FirstInfoUpdatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"firstname", "lastname"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound) //MainPage
		return
	}

	// Get values from URI
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	token := params.ByName("token")
	id, _, err := model.CheckUserToken(token)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound) //MainPage
		return
	}

	// Get form values
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	facebook := r.FormValue("facebook")
	twitter := r.FormValue("twitter")
	telegram := r.FormValue("telegram")

	telephone, err := strconv.ParseInt(r.FormValue("telephone"), 10, 32)

	contactInformation := contactInfo{
		int(telephone),
		facebook,
		twitter,
		telegram,
	}

	contactInfoJson, err := json.Marshal(contactInformation)
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound) //MainPage
		return
	}

	err = model.UpdateInfo(id, firstName, lastName, contactInfoJson)

	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
	}

	sess.AddFlash(view.Flash{"פרטייך עודכנו בהצלחה", view.FlashSuccess})
	sess.Save(r, w)

	// Display the page
	http.Redirect(w, r, "/", http.StatusFound) //MainPage
}

// SupervisorManagementGET displays list of all supervisors
func SupervisorManagementGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/admin_supervisor"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		//Gett list of all students
		users, err := model.SelectListOfUsers(2)
		if err != nil {
			log.Println(err)
			users = []model.User{}
		}
		v.Vars["users"] = users

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SearchSupervisorPOST handles the search for supervisor
func SearchSupervisorPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/admin_supervisor"

		//View permissions
		//TODO: ANDREY CHECK CSRF TOKEN
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
			SupervisorManagementGET(w, r)
			return
		}

		query := r.FormValue("query")

		//Get list of students by search
		users, err := model.SearchUser(query, 2)
		if err != nil {
			log.Println(err)
			users = []model.User{}
		}
		v.Vars["users"] = users

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// CreateSupervisorPOST handles the registration form submission
func CreateSupervisorPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, _ := view.Validate(r, []string{"id_number", "email"}); !validate {
		sess.AddFlash(view.Flash{"אנא מלא שדות חסרים", view.FlashError})
		sess.Save(r, w)
		CreateUserGET(w, r)
		return
	}

	// Get form values
	idNumber, err := strconv.ParseUint(r.FormValue("id_number"), 10, 32)
	//Check if the number converted well
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		CreateUserGET(w, r)
		return
	}

	email := r.FormValue("email")
	tempPassword, _ := uuid.NewRandom()
	password, _ := passhash.HashString(tempPassword.String())

	// If password hashing failed
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
		CreateUserGET(w, r)
		return
	}

	// Get database result
	_, erre := model.CheckIfTheUserExists(uint32(idNumber), email)

	if erre == model.ErrNoResult { // If success (no user exists with that email)
		//Create user with student role
		userId, ex := model.CreateNewUser(uint32(idNumber), email, password, 2)
		// Will only error if there is a problem with the query
		if ex != nil {
			log.Println(ex)
			sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
			sess.Save(r, w)
		} else {
			token, err := model.CreateChangePasswordToken(userId)
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאה בשרת", view.FlashError})
				sess.Save(r, w)
				ForgetPasswordGET(w, r)
				return
			}
			link := "http://" + r.Host + "/chnage_password/" + token
			go func() {
				err = sendemail.SendEmail(email, "ברוך הבא למערכת", "<p align='right'>:לכניסה ראשונה לחץ על הקישור </p>"+link)
			}()
			if err != nil {
				log.Println(err)
				sess.AddFlash(view.Flash{"שגיאה בשרת", view.FlashError})
				sess.Save(r, w)
				http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
				return
			}
			sess.AddFlash(view.Flash{"החשבון עבור " + email + " נוצר בהצלחה ", view.FlashSuccess})
			sess.Save(r, w)
			http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
			return
		}
	} else if erre != nil && erre != model.ErrNoResult { // Catch all other errors
		log.Println(erre)
		sess.AddFlash(view.Flash{"שגיאת שרת פנימית", view.FlashError})
		sess.Save(r, w)
	} else { // Else the user already exists
		if erre != model.ErrNoResult {
			sess.AddFlash(view.Flash{"החשבון קיים במערכת אנא בדוק את מספר תעודת הזהות או מייל ", view.FlashError})
			sess.Save(r, w)
		}
	}

	// Display the page
	CreateUserGET(w, r)
}

// WaitingSupervisorGET handles the waiting student list
func WaitingSupervisorGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/wait"

		//View permissions
		v.Vars["token"] = csrfbanana.Token(w, r, sess)
		v.Vars["first_name"] = sess.Values["first_name"]
		v.Vars["is_auth"] = sess.Values["is_auth"]
		v.Vars["is_student"] = sess.Values["is_student"]
		v.Vars["is_supervisor"] = sess.Values["is_supervisor"]
		v.Vars["is_project_manager"] = sess.Values["is_project_manager"]

		//Get list of all students
		users, err := model.SelectListOfWaiting(2)
		if err != nil {
			log.Println(err)
			users = []model.User{}
		}
		v.Vars["users"] = users

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SearchWaitingSupervisorPOST handles the search for waiting student
func SearchWaitingSupervisorPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//Get list of all users
	if sess.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "admin/wait"

		//View permissions
		//TODO: ANDREY CHECK CSRF TOKEN
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
			WaitingSupervisorGET(w, r)
			return
		}

		query := r.FormValue("query")

		//Get list of students by search
		users, err := model.SearchWaitingUser(query, 2)
		if err != nil {
			log.Println(err)
			users = []model.User{}
		}
		v.Vars["users"] = users

		v.Render(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// UploadSupervisorFilePOST handles the file upload
func UploadSupervisorFilePOST(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// upload of 1 MB files.
	r.ParseMultipartForm(1 << 2)

	file, handler, err := r.FormFile("File")
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאה בשרת, אנא נסה בשנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
		return
	}

	defer file.Close()
	extension := filepath.Ext(handler.Filename)
	if handler.Header.Get("Content-Type") == "application/vnd.ms-excel" && extension != ".csv" {
		sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה בשנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
		return
	}

	// Create a temporary file within our temp directory that follows
	//TODO: ADD SYSTEM DETECTION
	tempFile, err := ioutil.TempFile("C:\\temp-file", "upload-*.csv")
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאה בשרת, אנא נסה בשנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
		return
	}

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		sess.AddFlash(view.Flash{"שגיאה בשרת, אנא נסה להעלות את הקובץ שנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
		return
	}

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	csvFile, err := os.Open(tempFile.Name())
	if err != nil {
		sess.AddFlash(view.Flash{"קובץ אינו תקין, אנא נסה שנית", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
		return
	}

	table := csv.NewReader(csvFile)
	lines, err := table.ReadAll()
	if err != nil {
		log.Println("error reading all lines: %v", err)
	}

	idNumber := make([]string, len(lines)-1)
	email := make([]string, len(lines)-1)
	for i, line := range lines {
		if i == 0 {
			// skip header line
			continue
		}
		idNumber[i-1] = line[0]
		email[i-1] = line[1]
	}

	var wg sync.WaitGroup
	var listOfExisting []string

	for z, _ := range idNumber {
		// Increment the wait group counter
		wg.Add(1)

		go func() {
			// Decrement the counter when the go routine completes
			defer wg.Done()
			// Call the function check
			idTemp, _ := strconv.ParseUint(idNumber[z], 10, 32)
			_, erre := model.CheckIfTheUserExists(uint32(idTemp), email[z])
			if erre == nil {
				listOfExisting = append(listOfExisting, strconv.FormatInt(int64(z+2), 10))
			}
		}()
		// Wait for all the checkWebsite calls to finish
		wg.Wait()
	}

	if len(listOfExisting) != 0 {
		list := strings.Join(listOfExisting, ", ")
		sess.AddFlash(view.Flash{"החשבונות עבור שורות " + list + " קיימים כבר ", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
		return
	}

	for j, _ := range idNumber {
		// Increment the wait group counter
		wg.Add(1)
		go func() {
			// Decrement the counter when the go routine completes
			defer wg.Done()
			// Call the function check
			idNumber, _ := strconv.ParseUint(idNumber[j], 10, 32)
			tempPassword, _ := uuid.NewRandom()
			password, _ := passhash.HashString(tempPassword.String())
			model.CreateNewUser(uint32(idNumber), email[j], password, 2)
		}()
		// Wait for all the checkWebsite calls to finish
		wg.Wait()
	}

	sess.AddFlash(view.Flash{"כל המשתמשים נוצרו בהצלחה", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/admin/supervisor_management", http.StatusFound)
}
