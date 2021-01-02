# Student project management system

Students in their final year are required to characterize and develop their final project. The final project serves as a reflection of the knowledge they have accumulated during their studies and demonstrate their skills.
This project is the most significant challenge of the degree and requires a lot of effort, work, self-learning, thinking, sleepless nights, and sometimes quite a few nerves. Also, sometimes the students are required to present their final project in job interviews, and it can affect their acceptance into a desirable work that they have dreamed so long.

Nowadays, when the students starting work on their final project, they are facing difficulties such as no source of all previous projects, a list of all available projects and supervisors, a list of potential partners for the final project, progress tracking and more. 
After a short discussion with students, I found that a lot of students complain that there is no single platform that can answer to the real students' needs and use it at all stages of the project. I found that not only the difficulty of the project itself but also the lack of a dedicated system for the students makes the project much more complex and difficult.

The system I was developed (SPMS) over the last two semesters in our final project meets the needs described above. This system was developed using modern technologies such as a cloud environment and a new stack of programming languages. The developed system is web-based and allows any user to access it using his computer. The system is convenient and intuitive to the end-user, it can be easily integrated into any educational institution that requires a final project management system. At the end of the day the system developed to help students, supervisors and project managers to succeed.

## Technologies and languages 

I have obtained the following knowledge, technologies, and tools:
### Languages
* Go 1.13
* HTML5
* CSS
* Bootstrap
* JavaScript
### Technologies
* PostgresSQL
* Google SMTP
* AWS
	* S3
	* RDS 
	* EC2 
	* Route 53
	* SNS

## During the project development I focused on the following items:

* **Projects archive** – The SPMS stores all completed projects on the AWS S3. Any students can access the project archive and see projects. 
* **Improve the existing approval process** – The SPMS provides a fully managed process through the system by using web forms. The students can choose a project which has been proposed by the supervisor or offer their ideas. 
Students and supervisors who saw our project told us that it makes the whole approval flow much more intuitive and straightforward.
* **Fast and simple communication channel** – The SPMS should become a single system to communicate between students and supervisors on different projects. Students and supervisors can send messages to each other and attach files to the message. The system sent a notification by email or SMS when a new message added by the students or supervisor or if the user did not react more than 5 minutes. 
* **Project progress bar** – Once the project is approved, the supervisor can define and add new milestones on the progress bar of the project. 
When the milestones are presented, the supervisor can mark the milestone as done or remove them. It makes progress more transparent for both sides and helps them to meet their expectations during the project.
* **Find a partner** – The SPMS is providing a list of students who still not attached to the project. Each student can decide which contact information he wants to share with others. The most common communication and social media channels are existing in our system. For example, phone, email, Facebook Twiter, and telegram.
* **Single management system** – The SPMS is a single system to manage all the projects. The system includes the list of projects which have been proposed by the supervisor. Students can propose their idea and the relevant supervisor need to approve it. One of the most significant advantages is that the supervisor and the project manager will be able to get more visibility and control over the running projects. 
* **High performance** – The SPMS developed by using a server-side language with high performance, and the location of the server on the cloud environment increases it much more.
AWS services are costly, so we choose the most chipper EC2 server and RDS database. Initially, I was afraid that it will affect server performance. However, when I did the load tests, I was surprised to find out that our system can handle a very large scale and reply within 16 milliseconds.
* **Cloud-based environment** – The SPMS hosted on AWS (Amazon Web Services) and will use different AWS services such as EC2, S3, RDS, SNS, Route 53, and more.
* **Supplying a system with a high usability level** – During the development, we presented our UI to the potential system users. Their feedback helps us to create a system with high usability, which includes a modern UI and intuitive usage.


## Installation

1. Clone this repository ``` git clone https://github.com/NoySlush/Student-project-management-system.git ```
2. Execute the migration file from ```/config/megretion.sql``` over your DB 
3. Fill the missing keys in the configuration file ```/config/config.json```
4. Run the project.

## Structure

Recently, the folder structure changed. After looking at all the forks 
and reusing my project in different places, I decided to move the Go code to the 
**app** folder inside the **vendor** folder so the github path is not littered 
throughout the many imports. I did not want to use relative paths so the vendor
folder seemed like the best option.

The project is organized into the following folders:

~~~
config		- application settings, database schema and migration file 
static		- location of statically served files like CSS and JS
template	- HTML templates

vendor/app/controller		- page logic organized by HTTP methods (GET, POST)
vendor/app/shared		- packages for templates, Postgres, cryptography, sessions, and json
vendor/app/model		- database queries
vendor/app/route		- route information and middleware
~~~

## Configuration

I removed all the keys from the configuration file due to security reasons. 
If you want to add any of your own settings, you can add them to config.json. 
This is config.json:

~~~ json
{
	"Database": {
		"Username": "postgres",
		"Password": "postgres",
		"DBName": "postgres",
		"Hostname": "127.0.0.1"
	},
	"AWS": {
		"Region": "us-east-2",
		"Bucket": "fs.spms-project.info",
		"AccessKey": "",
		"SecretKey": ""
	},
	"Email": {
		"Username": "",
		"Password": "",
		"Hostname": "",
		"Port": 587,
		"From": ""
	},
	"Server": {
		"Hostname": "",
		"UseHTTP": true,
		"UseHTTPS": true,
		"HTTPPort": 3080,
		"HTTPSPort": 8443,
		"CertFile": "tls/cert1.pem",
		"KeyFile": "tls/privkey1.pem"
	},
	"Session": {
		"SecretKey": "0yhE@#q4*$%uFugE@c5@&67WwM6*un8d",
		"Name": "spms-session",
		"Options": {
			"Path": "/",
			"Domain": "",
			"MaxAge": 28800,
			"Secure": true,
			"HttpOnly": true
		}
	},
	"Template": {
		"Root": "base",
		"Children": [
			"partial/menu",
			"partial/footer"
		]
	},
	"View": {
		"BaseURI": "/",
		"Extension": "tmpl",
		"Folder": "template",
		"Name": "blank",
		"Caching": true
	}
}
~~~

## Maintainers

This project is mantained by:
* [Noy Shlush](https://github.com/NoyShlush)
* [Nofar Shmilovich](https://github.com/NofarShmil)


## Video

[![SPMS Video](https://yt-embed.herokuapp.com/embed?v=N5_xTTRaE3g)](https://www.youtube.com/watch?v=N5_xTTRaE3g "SPMS")
