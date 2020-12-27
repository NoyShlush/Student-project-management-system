grant all privileges on database postgres to postgres;

GRANT ALL ON schema public TO postgres;

CREATE TABLE user_table(
                           ID serial PRIMARY KEY,
                           IdNumber INT UNIQUE NOT NULL,
                           FirstName VARCHAR (20),
                           LastName VARCHAR (20),
                           Email VARCHAR (100) UNIQUE NOT NULL,
                           Role INT NOT NULL,
                           ContactInfomation json,
                           HashPassword VARCHAR (60) NOT NULL,
                           CreatedAt TIMESTAMP NOT NULL,
                           Block boolean NOT NULL,
                           ResetPasswordAt TIMESTAMP,
                           ResetToken uuid,
                           TokenValidUntil TIMESTAMP
);

INSERT INTO user_table (id, idnumber, firstname, lastname, email, role, contactinfomation, hashpassword, createdat, block, resetpasswordat, resettoken, tokenvaliduntil) VALUES (DEFAULT, 100000000, 'מנהל', 'מנהל', 'admin@braude.ac.il', 3, '{}', '$2a$10$xDhXPxyTp9UxNxJowi1YcOY.A8yg8XfFMmtZNKKKrQNvjdbFgk01u', '2019-11-09 15:34:21.649700', false, null, null, null);

CREATE TABLE role_table(
                           ID serial PRIMARY KEY,
                           Role VARCHAR (15) UNIQUE NOT NULL
);

INSERT INTO public.role_table ("id", "role") VALUES (DEFAULT, 'Student');
INSERT INTO public.role_table ("id", "role") VALUES (DEFAULT, 'Supervisor');
INSERT INTO public.role_table ("id", "role") VALUES (DEFAULT, 'Project Manager');

CREATE TABLE notifications_table(
                                    ID serial PRIMARY KEY,
                                    Content json NOT NULL,
                                    Type INT NOT NULL
);

INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"ברוך הבא למערכת","Body":"<p align=''right''>:לכניסה ראשונה לחץ על הקישור </p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"ישנה הודעה חדשה בצ''אט"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"ישנה הודעה חדשה בצ''אט","Body":"<p align=''right''>.ישנה הודעה חדשה בצ''אט</p><p align=''right''>כדי לראות את ההודעה, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"הוגש טופס הצעה חדש לפרויקט"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"הוגש טופס הצעה חדש לפרויקט","Body":"<p align=''right''>.הוגש טופס הצעה חדש לפרויקט</p><p align=''right''>כדי לראות את הטופס, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"המנחה אישר את הטופס"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"המנחה אישר את הטופס","Body":"<p align=''right''>.המנחה אישר את הטופס</p><p align=''right''>כדי לראות את הטופס, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"המנחה דחה את הטופס"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"המנחה דחה את הטופס","Body":"<p align=''right''>.המנחה אישר את הטופס</p><p align=''right''>כדי לראות את הטופס, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"המנחה הוסיף לטופס הערות"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"הוגש טופס אישור פרויקט חדש","Body":"<p align=''right''>.הוגש טופס אישור פרויקט חדש</p><p align=''right''>כדי לראות את הטופס, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"המנחה הוסיף לטופס הערות","Body":"<p align=''right''>.המנחה אישר את הטופס</p><p align=''right''>כדי לראות את ההערות, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"מנהל פרויקט אישר את הטופס"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"הוגש טופס אישור פרויקט חדש"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"מנהל פרויקט אישר את הטופס","Body":"<p align=''right''>.מנהל פרויקט אישר את הטופס</p><p align=''right''>כדי לראות את הטופס, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"מנהל פרויקט דחה את הטופס"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"מנהל פרויקט דחה את הטופס","Body":"<p align=''right''>.מנהל פרויקט דחה את הטופס</p><p align=''right''>כדי לראות את הטופס, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"מנהל פרויקט הוסיף הערות לטופס"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"מנהל פרויקט הוסיף הערות לטופס","Body":"<p align=''right''>.מנהל פרויקט הוסיף הערות לטופס</p><p align=''right''>כדי לראות את הטופס, אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"הוגש פרויקט חלק א'' חדש"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"הוגש פרויקט חלק א'' חדש","Body":"<p align=''right''>.הוגש פרויקט חלק א'' חדש</p><p align=''right''>כדי לראות את ההגשה הסופית של חלק א'', אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"הוגש פרויקט חלק ב'' חדש"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"הוגש פרויקט חלק ב'' חדש","Body":"<p align=''right''>.הוגש פרויקט חלק ב'' חדש</p><p align=''right''>כדי לראות את ההגשה הסופית של חלק ב'', אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"סטטוס הפרויקט עודכן, ניתן לבצע הגשה של חלק א''"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"סטטוס הפרויקט עודכן, ניתן לבצע הגשה של חלק א''","Body":"<p align=''right''>.סטטוס הפרויקט עודכן, ניתן לבצע הגשה של חלק א''</p><p align=''right''>כדי להגיש את חלק א'', אנא כנס למערכת ניהול פרויקטים</p>"}', 2);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"","Body":"סטטוס הפרויקט עודכן, ניתן לבצע הגשה של חלק ב''"}', 1);
INSERT INTO public.notifications_table (id, content, type) VALUES (DEFAULT, '{"Subject":"סטטוס הפרויקט עודכן, ניתן לבצע הגשה של חלק ב''","Body":"<p align=''right''>.סטטוס הפרויקט עודכן, ניתן לבצע הגשה של חלק ב''</p><p align=''right''>כדי להגיש את חלק ב'', אנא כנס למערכת ניהול פרויקטים</p>"}', 2);

CREATE TABLE notificationsType_table(
                                        ID serial PRIMARY KEY,
                                        Type VARCHAR (15) UNIQUE NOT NULL
);

INSERT INTO public.notificationstype_table ("id", "type") VALUES (DEFAULT, 'SMS');
INSERT INTO public.notificationstype_table ("id", "type") VALUES (DEFAULT, 'Email');

CREATE TABLE project_table(
                              ID serial PRIMARY KEY,
                              ProjectName VARCHAR (50) NOT NULL,
                              Description text NOT NULL,
                              ShortDescription VARCHAR (255) NOT NULL,
                              StatusId INT NOT NULL,
                              Files json,
                              Type INT NOT NULL,
                              FormId INT,
                              CreatedAt TIMESTAMP NOT NULL,
                              UpdateAt TIMESTAMP NOT NULL,
                              CommentsId json,
                              StudentsId json NOT NULL,
                              SupervisorId INT NOT NULL
);

CREATE TABLE projectType_table(
                                  ID   serial PRIMARY KEY,
                                  Type VARCHAR(20) UNIQUE NOT NULL
);

INSERT INTO public.projectType_table ("id", "type") VALUES (DEFAULT, 'StudentIdea');
INSERT INTO public.projectType_table ("id", "type") VALUES (DEFAULT, 'SupervisorProject');

CREATE TABLE status_table(
                             ID serial PRIMARY KEY,
                             Description VARCHAR(50) NOT NULL
);

INSERT INTO public.status_table ("id", "description") VALUES (DEFAULT, 'open');
INSERT INTO public.status_table ("id", "description") VALUES (DEFAULT, 'waiting_for_supervisor_approval');
INSERT INTO public.status_table ("id", "description") VALUES (DEFAULT, 'waiting_for_projectManager_approval');
INSERT INTO public.status_table ("id", "description") VALUES (DEFAULT, 'in_progress');
INSERT INTO public.status_table ("id", "description") VALUES (DEFAULT, 'done');

CREATE TABLE progressBar_table(
                                  ID serial PRIMARY KEY,
                                  ProgressBarId INT,
                                  Milestone VARCHAR(50) NOT NULL,
                                  Done boolean NOT NULL
);

CREATE TABLE file_table(
                           ID serial PRIMARY KEY,
                           Link VARCHAR(255) NOT NULL,
                           Type INT NOT NULL
);

CREATE TABLE fileType_table(
                               ID serial PRIMARY KEY,
                               Type VARCHAR(15) NOT NULL
);

INSERT INTO public.fileType_table ("id", "type") VALUES (DEFAULT, 'Book1');
INSERT INTO public.fileType_table ("id", "type") VALUES (DEFAULT, 'Book2');
INSERT INTO public.fileType_table ("id", "type") VALUES (DEFAULT, 'Presention1');
INSERT INTO public.fileType_table ("id", "type") VALUES (DEFAULT, 'Presention2');
INSERT INTO public.fileType_table ("id", "type") VALUES (DEFAULT, 'SourceCode');
INSERT INTO public.fileType_table ("id", "type") VALUES (DEFAULT, 'Guidelines');

CREATE TABLE comment_table(
                              ID serial PRIMARY KEY,
                              Message text NOT NULL,
                              UserId INT NOT NULL,
                              CreatedAt TIMESTAMP NOT NULL
);

CREATE TABLE approvalForm_table(
                                   ID serial PRIMARY KEY,
                                   Synopses text NOT NULL,
                                   ScopeOfTheProject text NOT NULL,
                                   UniqueFeatures text NOT NULL
);

CREATE TABLE projectArchive_table(
                                     ID serial PRIMARY KEY,
                                     CreatedAt TIMESTAMP NOT NULL,
                                     ApprovedToPresent boolean NOT NULL,
                                     ApprovedFiles json NOT NULL,
                                     ProjectId INT NOT NULL
);

CREATE TABLE chat_table(
                           ID serial PRIMARY KEY,
                           CreatedAt TIMESTAMP NOT NULL,
                           Message text NOT NULL,
                           ChatId INT NOT NULL,
                           ReadBy json NOT NULL,
                           SendBy INT NOT NULL
);