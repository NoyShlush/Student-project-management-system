package route

import (
	"net/http"

	"app/controller"
	"app/route/middleware/acl"
	hr "app/route/middleware/httprouterwrapper"
	"app/route/middleware/logrequest"
	"app/route/middleware/pprofhandler"
	"app/shared/session"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Load returns the routes and middleware
func Load() http.Handler {
	return middleware(routes())
}

// LoadHTTPS returns the HTTP routes and middleware
func LoadHTTPS() http.Handler {
	return middleware(routes())
}

// LoadHTTP returns the HTTPS routes and middleware
func LoadHTTP() http.Handler {
	return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *httprouter.Router {
	r := httprouter.New()

	// Set 404 handler
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.
		New().
		ThenFunc(controller.Static)))

	// Home page
	r.GET("/", hr.Handler(alice.
		New().
		ThenFunc(controller.IndexGET)))

	// Login
	r.GET("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginGET)))
	r.POST("/login", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.LoginPOST)))
	r.GET("/logout", hr.Handler(alice.
		New().
		ThenFunc(controller.LogoutGET)))

	// Student management
	r.GET("/admin/student_management", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.StudentManagementGET)))
	r.POST("/admin/student_management", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchStudentPOST)))
	r.GET("/admin/student_management/create_student", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CreateUserGET)))
	r.POST("/admin/student_management/create_student", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CreateStudentPOST)))
	r.GET("/admin/student_management/waiting_student", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.WaitingStudentGET)))
	r.POST("/admin/student_management/waiting_student", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchWaitingStudentPOST)))
	r.POST("/admin/student_management/upload_file", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadStudentFilePOST)))
	r.GET("/admin/student_management/edit/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UpdateUserGET)))
	r.POST("/admin/student_management/edit/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UpdateUserPOST)))
	r.GET("/admin/student_management/block/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BlockUserGET)))

	// Supervisor management
	r.GET("/admin/supervisor_management", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SupervisorManagementGET)))
	r.POST("/admin/supervisor_management", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchSupervisorPOST)))
	r.GET("/admin/supervisor_management/create_supervisor", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CreateUserGET)))
	r.POST("/admin/supervisor_management/create_supervisor", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.CreateSupervisorPOST)))
	r.GET("/admin/supervisor_management/waiting_supervisor", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.WaitingSupervisorGET)))
	r.POST("/admin/supervisor_management/waiting_supervisor", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchWaitingSupervisorPOST)))
	r.POST("/admin/supervisor_management/upload_file", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadSupervisorFilePOST)))
	r.GET("/admin/supervisor_management/edit/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UpdateUserGET)))
	r.POST("/admin/supervisor_management/edit/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UpdateUserPOST)))
	r.GET("/admin/supervisor_management/block/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.BlockUserGET)))

	//Password
	r.GET("/forgot_password", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.ForgetPasswordGET)))
	r.POST("/forgot_password", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.ForgetPasswordPOST)))
	r.GET("/chnage_password/:token", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.ChangePasswordGET)))
	r.POST("/chnage_password/:token", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.ChangePasswordPOST)))

	// Enable Pprof
	r.GET("/debug/pprof/*pprof", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(pprofhandler.Handler)))

	//User
	r.GET("/user/update", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UpdateInfoGET)))
	r.POST("/user/update", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UpdateInfoPOST)))
	r.GET("/user/update/:token", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.FirstInfoUpdateGET)))
	r.POST("/user/update/:token", hr.Handler(alice.
		New(acl.DisallowAuth).
		ThenFunc(controller.FirstInfoUpdatePOST)))

	//Guidelines
	r.GET("/guidline", hr.Handler(alice.
		New().
		ThenFunc(controller.GuidLineGET)))

	//Project
	r.GET("/project/projects", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ListProjectsGET)))
	r.POST("/project/projects", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchOpenProjectPOST)))
	r.GET("/project/new_idea", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NewIdeaGET)))
	r.POST("/project/new_idea", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NewIdeaPOST)))
	r.GET("/project/new_project", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NewProjectGET)))
	r.POST("/project/new_project", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.NewProjectPOST)))
	r.GET("/project/more/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MoreGET)))
	r.POST("/project/more/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MorePOST)))
	r.GET("/project/waiting", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ApprovalWaitingProjectsGET)))
	r.POST("/project/waiting", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchApprovalWaitingProjectPOST)))
	r.GET("/project/manager_waiting", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ManagerApprovalWaitingProjectsGET)))
	r.POST("/project/manager_waiting", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchManagerApprovalWaitingProjectsPOST)))
	r.GET("/project/running", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ListRunningProjectsGET)))
	r.POST("/project/running", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchRunningProjectPOST)))
	r.GET("/project/add_comments/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AddCommentsGET)))
	r.POST("/project/add_comments/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AddCommentsPOST)))
	r.GET("/project/approval_add_comments/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AddApprovalCommentsGET)))
	r.POST("/project/approval_add_comments/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.AddApprovalCommentsPOST)))
	r.GET("/project/edit/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditProjectGET)))
	r.POST("/project/edit/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.EditProjectPOST)))
	r.GET("/project/approved/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ApproveProjectGET)))
	r.GET("/project/declined/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.DeclineProjectGET)))
	r.GET("/project/manager_approved/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ApproveManagerProjectGET)))
	r.GET("/project/manager_declined/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.DeclineManagerProjectGET)))
	r.GET("/project/approval/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ApprovalFormGET)))
	r.POST("/project/approval/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ApprovalFormPOST)))
	r.POST("/project/upload_file/part1_book_pdf", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadPDFFilePOST)))
	r.POST("/project/upload_file/part1_book_doc", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadDOCFilePOST)))
	r.POST("/project/upload_file/part1_presentation", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadPowerPointFilePOST)))
	r.POST("/project/upload_file/part2_book_pdf", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadPDFFilePOST)))
	r.POST("/project/upload_file/part2_book_doc", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadDOCFilePOST)))
	r.POST("/project/upload_file/part2_code", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadCodeFilePOST)))
	r.POST("/project/upload_file/part2_presentation", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadPowerPointFilePOST)))

	//Chat
	r.GET("/chat", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ChatGET)))
	r.POST("/chat", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ChatPOST)))
	r.GET("/chat/chat_list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ChatListGET)))
	r.GET("/chat/chat/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ChatGET)))
	r.POST("/chat/chat/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ChatPOST)))
	r.POST("/chat/upload_file/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.UploadChatFilePOST)))
	r.POST("/chat/chat_list", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.SearchChatListPOST)))

	//Progress Bar
	r.GET("/progressbar", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ProgressBarStudentGET)))
	r.GET("/progressbar/progressbar/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.ProgressBarSupervisorGET)))
	r.GET("/progressbar/done/:id1/:id2", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MilestoneDoneGET)))
	r.GET("/progressbar/remove/:id1/:id2", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MilestoneRemoveGET)))
	r.POST("/progressbar/progressbar/:id", hr.Handler(alice.
		New(acl.DisallowAnon).
		ThenFunc(controller.MilestoneAddPOST)))

	//Project Archive
	r.GET("/archive", hr.Handler(alice.
		New().
		ThenFunc(controller.ListArchiveProjectsGET)))
	r.POST("/archive", hr.Handler(alice.
		New().
		ThenFunc(controller.SearchArchiveProjectsPOST)))
	r.GET("/archive/:id", hr.Handler(alice.
		New().
		ThenFunc(controller.ArchiveProjectGET)))

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, session.Store, session.Name)
	cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "token"
	csrfbanana.SingleToken = false
	h = cs

	// Log every request
	h = logrequest.Handler(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
