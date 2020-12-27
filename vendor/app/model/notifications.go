package model

import (
	"app/shared/aws"
	"app/shared/database"
	sendemail "app/shared/email"
	"encoding/json"
	"strconv"
	"time"
)

// *****************************************************************************
// Notifications
// *****************************************************************************

// Notifications table contains the information for each notification
type Notifications struct {
	ID      uint32          `db:"ID"`
	Content json.RawMessage `db:"Content"`
	Type    uint            `db:"Type"`
}

type Content struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

const SEND_AFTER = 1

// SelectNotificationsById gets the notification by ID
func SelectNotificationsById(id uint32) (Notifications, error) {
	var err error

	result := Notifications{}
	err = database.SQL.QueryRow("SELECT * FROM notifications_table WHERE id = $1 LIMIT 1", id).Scan(
		&result.ID,
		&result.Content,
		&result.Type)

	return result, StandardizeError(err)
}

// SendChatNotificationToStudents send email and sms notifications after X minutes
func SendChatNotificationToStudents(projectID, emailNotificationID, smsNotificationID, messageID uint32) error {

	// Waiting time
	time.Sleep(time.Minute * SEND_AFTER)

	var readBy ReadBy
	var studentsID Students
	var userInfo1, userInfo2 ContactInfo
	var email, sms Content

	jsonRead, err := WhoRead(messageID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(jsonRead, &readBy)
	if err != nil {
		return StandardizeError(err)
	}
	project, err := SelectProjectById(projectID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(project.StudentsId, &studentsID)
	if err != nil {
		return StandardizeError(err)
	}
	user1, err := SelectUserInfo(uint32(studentsID.StudentId1))
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(user1.ContactInfomation, &userInfo1)
	if err != nil {
		return StandardizeError(err)
	}
	user2, err := SelectUserInfo(uint32(studentsID.StudentId2))
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(user2.ContactInfomation, &userInfo2)
	if err != nil {
		return StandardizeError(err)
	}

	template, err := SelectNotificationsById(emailNotificationID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(template.Content, &email)
	if err != nil {
		return StandardizeError(err)
	}
	template, err = SelectNotificationsById(smsNotificationID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(template.Content, &sms)
	if err != nil {
		return StandardizeError(err)
	}

	if readBy.UserId1 == 0 && readBy.UserId2 == 0 {
		err = sendemail.SendEmail(user1.Email, email.Subject, email.Body)
		if err != nil {
			return StandardizeError(err)
		}
		if userInfo1.Telephone != 0 {
			aws.SendSMS("+972"+strconv.Itoa(userInfo1.Telephone), sms.Body)
		}

		err = sendemail.SendEmail(user2.Email, email.Subject, email.Body)
		if err != nil {
			return StandardizeError(err)
		}
		if userInfo2.Telephone != 0 {
			aws.SendSMS("+972"+strconv.Itoa(userInfo2.Telephone), sms.Body)
		}

	} else if readBy.UserId1 != 0 && readBy.UserId2 == 0 {
		err = sendemail.SendEmail(user2.Email, email.Subject, email.Body)
		if err != nil {
			return StandardizeError(err)
		}
		if userInfo2.Telephone != 0 {
			aws.SendSMS("+972"+strconv.Itoa(userInfo2.Telephone), sms.Body)
		}
	} else if readBy.UserId1 == 0 && readBy.UserId2 != 0 {
		err = sendemail.SendEmail(user1.Email, email.Subject, email.Body)
		if err != nil {
			return StandardizeError(err)
		}
		if userInfo1.Telephone != 0 {
			aws.SendSMS("+972"+strconv.Itoa(userInfo1.Telephone), sms.Body)
		}
	}
	return nil
}

// SendChatNotificationToStudentSupervisor send email and sms notifications after X minutes
func SendChatNotificationToStudentSupervisor(projectID, sendBy, emailNotificationID, smsNotificationID, messageID uint32) error {

	// Waiting time
	time.Sleep(time.Minute * SEND_AFTER)

	var readBy ReadBy
	var studentsID Students
	var userInfo1, userInfo2 ContactInfo
	var user1, user2 User
	var email, sms Content

	jsonRead, err := WhoRead(messageID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(jsonRead, &readBy)
	if err != nil {
		return StandardizeError(err)
	}
	project, err := SelectProjectById(projectID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(project.StudentsId, &studentsID)
	if err != nil {
		return StandardizeError(err)
	}

	if sendBy == uint32(studentsID.StudentId1) {
		user1, err = SelectUserInfo(project.SupervisorId)
		if err != nil {
			return StandardizeError(err)
		}
		err = json.Unmarshal(user1.ContactInfomation, &userInfo1)
		if err != nil {
			return StandardizeError(err)
		}
		user2, err = SelectUserInfo(uint32(studentsID.StudentId2))
		if err != nil {
			return StandardizeError(err)
		}
		err = json.Unmarshal(user2.ContactInfomation, &userInfo2)
		if err != nil {
			return StandardizeError(err)
		}
	} else {
		user1, err = SelectUserInfo(uint32(studentsID.StudentId1))
		if err != nil {
			return StandardizeError(err)
		}
		err = json.Unmarshal(user1.ContactInfomation, &userInfo1)
		if err != nil {
			return StandardizeError(err)
		}
		user2, err = SelectUserInfo(project.SupervisorId)
		if err != nil {
			return StandardizeError(err)
		}
		err = json.Unmarshal(user2.ContactInfomation, &userInfo2)
		if err != nil {
			return StandardizeError(err)
		}
	}

	template, err := SelectNotificationsById(emailNotificationID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(template.Content, &email)
	if err != nil {
		return StandardizeError(err)
	}
	template, err = SelectNotificationsById(smsNotificationID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(template.Content, &sms)
	if err != nil {
		return StandardizeError(err)
	}

	if readBy.UserId1 == 0 && readBy.UserId2 == 0 {
		err = sendemail.SendEmail(user1.Email, email.Subject, email.Body)
		if err != nil {
			return StandardizeError(err)
		}
		if userInfo1.Telephone != 0 {
			aws.SendSMS("+972"+strconv.Itoa(userInfo1.Telephone), sms.Body)
		}

		err = sendemail.SendEmail(user2.Email, email.Subject, email.Body)
		if err != nil {
			return StandardizeError(err)
		}
		if userInfo2.Telephone != 0 {
			aws.SendSMS("+972"+strconv.Itoa(userInfo2.Telephone), sms.Body)
		}

	} else if readBy.UserId1 != 0 && readBy.UserId2 == 0 {
		err = sendemail.SendEmail(user2.Email, email.Subject, email.Body)
		if err != nil {
			return StandardizeError(err)
		}
		if userInfo2.Telephone != 0 {
			aws.SendSMS("+972"+strconv.Itoa(userInfo2.Telephone), sms.Body)
		}
	} else if readBy.UserId1 == 0 && readBy.UserId2 != 0 {
		err = sendemail.SendEmail(user1.Email, email.Subject, email.Body)
		if err != nil {
			return StandardizeError(err)
		}
		if userInfo1.Telephone != 0 {
			aws.SendSMS("+972"+strconv.Itoa(userInfo1.Telephone), sms.Body)
		}
	}
	return nil
}

// SendChatNotificationToStudents send email and sms notifications after X minutes
func SendToProjectManager(emailNotificationID, smsNotificationID uint32) error {

	var email, sms Content
	var Info ContactInfo

	projectManager, err := SelectProjectManager()
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(projectManager.ContactInfomation, &Info)
	if err != nil {
		return StandardizeError(err)
	}
	template, err := SelectNotificationsById(emailNotificationID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(template.Content, &email)
	if err != nil {
		return StandardizeError(err)
	}
	template, err = SelectNotificationsById(smsNotificationID)
	if err != nil {
		return StandardizeError(err)
	}
	err = json.Unmarshal(template.Content, &sms)
	if err != nil {
		return StandardizeError(err)
	}
	err = sendemail.SendEmail(projectManager.Email, email.Subject, email.Body)
	if err != nil {
		return StandardizeError(err)
	}
	if Info.Telephone != 0 {
		aws.SendSMS("+972"+strconv.Itoa(Info.Telephone), sms.Body)
	}

	return nil
}
