package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/bika-c/sqlite"
	"github.com/bika-c/sqlite/sqlitex"
)

type SQLite struct {
	db *sqlitex.Pool
}

type Conn sqlite.Conn

type Pragma struct {
	ForeignKeys       string
	RecursiveTriggers string
	JournalMode       string
	Synchronous       string
	TempStore         string
	MmapSize          string
	BusyTimeout       string
	Cache             string
}

func conj(s *strings.Builder, k, v string) {
	if v == "" {
		return
	}
	s.WriteString(fmt.Sprintf("PRAGMA %s=%v;", k, v))
}

func (p *Pragma) String() string {
	var s strings.Builder
	conj(&s, "foreign_keys", p.ForeignKeys)
	conj(&s, "recursive_triggers", p.RecursiveTriggers)
	conj(&s, "journal_mode", p.JournalMode)
	conj(&s, "synchronous", p.Synchronous)
	conj(&s, "temp_store", p.TempStore)
	conj(&s, "mmap_size", p.MmapSize)
	conj(&s, "busy_timeout", p.BusyTimeout)
	conj(&s, "cache", p.Cache)
	return s.String()
}

var DefaultPragma Pragma = Pragma{
	ForeignKeys:       "TRUE",
	RecursiveTriggers: "FALSE",
	JournalMode:       "WAL",
	Synchronous:       "normal",
	TempStore:         "memory",
	MmapSize:          "300000000",
	Cache:             "shared",
	BusyTimeout:       "1000",
}

func NewSQLite(path string, size int, opt Pragma) *SQLite {
	if db, err := sqlitex.Open(path, 0, size); err != nil {
		panic(err)
	} else {
		s := &SQLite{
			db: db,
		}
		conns := make([]*Conn, 0, size)
		for i := 0; i < size; i++ {
			c := db.Get(context.Background())
			sqlitex.ExecScript(c, opt.String())
			conns = append(conns, (*Conn)(c))
		}
		for i := 0; i < size; i++ {
			db.Put((*sqlite.Conn)(conns[i]))
		}
		Setup(s)

		return s
	}
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) Get() *Conn {
	return (*Conn)(s.db.Get(context.Background()))
}
func (s *SQLite) Put(c *Conn) {
	s.db.Put((*sqlite.Conn)(c))
}

func (c *Conn) Close() error {
	return c.Raw().Close()
}

func (c *Conn) Raw() *sqlite.Conn {
	return (*sqlite.Conn)(c)
}

func (c *Conn) Save() func(*error) {
	return sqlitex.Save((*sqlite.Conn)(c))
}

func (c *Conn) PrepareTransient(query string) *sqlite.Stmt {
	if s, _, err := c.Raw().PrepareTransient(query); err != nil {
		panic(err)
	} else {
		return s
	}
}

func (c *Conn) Prepare(query string) *sqlite.Stmt {
	return c.Raw().Prep(query)
}

func (c *Conn) Pragma(key, val string) error {
	return sqlitex.ExecScript(c.Raw(), fmt.Sprintf("PRAGMA %s = %s", key, val))
}

func Setup(db *SQLite) {
	conn := db.Get()
	defer db.Put(conn)
	if err := sqlitex.ExecScript(conn.Raw(), initDB); err != nil {
		panic(err.Error())
	}
}

const initDB = `
	CREATE TABLE IF NOT EXISTS Account (
		ID			INTEGER		PRIMARY KEY,
		UUID 		TEXT 		NOT NULL,
		Email 		TEXT 		NOT NULL,
		Password 	TEXT 		NOT NULL,
		CreatedAt	TEXT		DEFAULT (DATETIME('now', 'localtime')) NOT NULL,
		UpdatedAt	TEXT		DEFAULT (DATETIME('now', 'localtime')) NOT NULL,
		LastLogin 	TEXT,
		UNIQUE(UUID),
		UNIQUE(Email)
	) STRICT;
	CREATE TABLE IF NOT EXISTS Contact (
		ID			INTEGER		PRIMARY KEY,
		AccountID	INTEGER		NOT NULL,
		Street		TEXT,
		Unit		TEXT,
		Province	TEXT,
		Country		TEXT,
		PostCode	TEXT,
		Email		TEXT,
		Phone		TEXT,
		FOREIGN KEY (AccountID) REFERENCES AccountID ON DELETE CASCADE,
		UNIQUE(AccountID)
	) STRICT;
	CREATE TABLE IF NOT EXISTS User (
		ID			INTEGER		PRIMARY KEY,
		Name		TEXT		NOT NULL,
		AccountID	INTEGER		NOT NULL,
		FOREIGN KEY (AccountID) REFERENCES Account(ID) ON DELETE CASCADE,
		UNIQUE(AccountID)
	) STRICT;
	CREATE TABLE IF NOT EXISTS Institution (
		ID			INTEGER		PRIMARY KEY,
		Name		TEXT		NOT NULL,
		Description	TEXT,
		AccountID	INTEGER		NOT NULL,
		ContactID	INTEGER		NOT NULL,
		FOREIGN KEY (AccountID) REFERENCES Account(ID) ON DELETE CASCADE,
		FOREIGN KEY (ContactID) REFERENCES Contact(ID) ON DELETE CASCADE,
		UNIQUE(AccountID),
		UNIQUE(Name)
	) STRICT;
	CREATE TABLE IF NOT EXISTS Instructor (
		ID				INTEGER		PRIMARY KEY,
		Name			TEXT		NOT NULL,
		AccountID		INTEGER		NOT NULL,
		InstitutionID	INTEGER		NOT NULL,
		ContactID		INTEGER		NOT NULL,
		FOREIGN KEY (AccountID) 	REFERENCES Account(ID) ON DELETE CASCADE,
		FOREIGN KEY (ContactID) 	REFERENCES Contact(ID) ON DELETE CASCADE,
		FOREIGN KEY (InstitutionID) REFERENCES Institution(ID) ON DELETE CASCADE,
		UNIQUE(AccountID)
	) STRICT;
	CREATE TABLE IF NOT EXISTS UserAsset (
		ID				INTEGER		PRIMARY KEY,
		Balance			INTEGER		DEFAULT (0),
		UserID			INTEGER		NOT NULL,
		InstitutionID	INTEGER		NOT NULL,
		FOREIGN KEY (UserID) 		REFERENCES User(ID) ON DELETE CASCADE,
		FOREIGN KEY (InstitutionID) REFERENCES Institution(ID) ON DELETE CASCADE
	) STRICT;
	CREATE TABLE IF NOT EXISTS Course (
		ID				INTEGER		PRIMARY KEY,
		Title			TEXT		NOT NULL,
		InstitutionID	INTEGER		NOT NULL,
		Description		TEXT,
		FOREIGN KEY (InstitutionID) REFERENCES Institution(ID) ON DELETE CASCADE
	) STRICT;
	CREATE TABLE IF NOT EXISTS Section (
		ID				INTEGER		PRIMARY KEY,
		Title			TEXT		NOT NULL,
		Description		TEXT,
		Start			TEXT		NOT NULL,
		Duration		INTEGER		NOT NULL,
		At				INTEGER		NOT NULL,
		UnitPrice		INTEGER		NOT NULL,
		RepetitionID	INTEGER,
		InstructorID	INTEGER		NOT NULL,
		CourseID		INTEGER		NOT NULL,
		FOREIGN KEY (CourseID) 		REFERENCES Course(ID) ON DELETE CASCADE,
		FOREIGN KEY (InstructorID) 	REFERENCES Instructor(ID),
		FOREIGN KEY (RepetitionID) 	REFERENCES Repetition(ID),
		CHECK ((Duration / 1000000000 + At) < 86400)
	) STRICT;
	CREATE TABLE IF NOT EXISTS Class (
		ID				INTEGER		PRIMARY KEY,
		Title			TEXT		NOT NULL,
		'Index'			INTEGER		DEFAULT (1),
		Time			TEXT		NOT NULL,
		Duration		INTEGER		NOT NULL,
		UnitPrice		INTEGER		NOT NULL,
		Canceled		INTEGER		DEFAULT (FALSE),
		Completed		INTEGER		DEFAULT (FALSE),
		Remark			TEXT,
		InstructorID	INTEGER		NOT NULL,
		SectionID		INTEGER		NOT NULL,
		UpdatedAt		TEXT		NOT NULL,
		FOREIGN KEY (SectionID) 	REFERENCES Section(ID) ON DELETE CASCADE,
		FOREIGN KEY (InstructorID) 	REFERENCES Instructor(ID),
		CHECK (Canceled IN (TRUE, FALSE)),
		CHECK (Completed IN (TRUE, FALSE))
	) STRICT;
	CREATE TABLE IF NOT EXISTS Enrollment (
		ID				INTEGER		PRIMARY KEY,
		UserID			INTEGER,
		ClassID			INTEGER,
		FOREIGN KEY (UserID) 	REFERENCES User(ID) ON DELETE CASCADE,
		FOREIGN KEY (ClassID) 	REFERENCES Class(ID) ON DELETE CASCADE
	) STRICT;
	CREATE TABLE IF NOT EXISTS Repetition (
		ID				INTEGER		PRIMARY KEY,
		End				TEXT		NOT NULL,
		Type			TEXT		NOT NULL,
		Interval		INTEGER		DEFAULT (0),
		Record			INTEGER		DEFAULT (0),
		SectionID		INTEGER		NOT NULL,
		FOREIGN KEY (SectionID)		REFERENCES Section(ID)	ON DELETE CASCADE,
		CHECK (Type IN ('Daily', 'Weekly', 'Monthly'))
	) STRICT;

	/*
		TRIGGERS
	*/
	-- Create the triggers to update UpdatedAt and CreatedAt columns
	-- Account Table
	CREATE TRIGGER IF NOT EXISTS AccountUpdatedAt
	AFTER UPDATE ON Account
	FOR EACH ROW BEGIN
		UPDATE Account
			SET UpdatedAt = DATETIME(CURRENT_TIMESTAMP, 'localtime')
		WHERE ID = OLD.ID AND UUID = OLD.UUID;
	END;
	-- Class Table
	CREATE TRIGGER IF NOT EXISTS ClassUpdatedAt
	AFTER UPDATE ON Class
	FOR EACH ROW BEGIN
		UPDATE Class
			SET UpdatedAt = DATETIME(CURRENT_TIMESTAMP, 'localtime')
		WHERE ID = OLD.ID;
	END;

	-- Auto Delete Repetition Records
	CREATE TRIGGER IF NOT EXISTS DeleteRepetitionRecord
	AFTER UPDATE ON Section
	FOR EACH ROW WHEN (OLD.RepetitionID IS NOT NEW.RepetitionID)
	BEGIN
		DELETE FROM Repetition WHERE ID = OLD.RepetitionID;
	END;
	-- Auto Delete Contact Record for Institution
	CREATE TRIGGER IF NOT EXISTS DeleteContactInstitution
	AFTER UPDATE ON Institution
	FOR EACH ROW WHEN (OLD.ContactID IS NOT NEW.ContactID)
	BEGIN
		DELETE FROM Contact WHERE ID = OLD.ContactID;
	END;
	-- Auto Delete Contact Record for Instructor
	CREATE TRIGGER IF NOT EXISTS DeleteRepetitionInstructor
	AFTER UPDATE ON Instructor
	FOR EACH ROW WHEN (OLD.ContactID IS NOT NEW.ContactID)
	BEGIN
		DELETE FROM Contact WHERE ID = OLD.ContactID;
	END;
`

const mock = `
-- Account Table
INSERT INTO Account (UUID, Email, Password) VALUES
	('uuid1', 'email1@example.com', 'password1'),
	('uuid2', 'email2@example.com', 'password2'),
	('uuid3', 'email3@example.com', 'password3'),
	('uuid4', 'email4@example.com', 'password4'),
	('uuid5', 'email5@example.com', 'password5');

-- Contact Table
INSERT INTO Contact (AccountID, Street, Unit, Province, Country, PostCode, Email, Phone) VALUES
	(1, 'Street 1', 'Unit 1', 'Province 1', 'Country 1', 'PostCode 1', 'contact1@example.com', '1234567890'),
	(2, 'Street 2', 'Unit 2', 'Province 2', 'Country 2', 'PostCode 2', 'contact2@example.com', '2345678901'),
	(3, 'Street 3', 'Unit 3', 'Province 3', 'Country 3', 'PostCode 3', 'contact3@example.com', '3456789012'),
	(4, 'Street 4', 'Unit 4', 'Province 4', 'Country 4', 'PostCode 4', 'contact4@example.com', '4567890123'),
	(5, 'Street 5', 'Unit 5', 'Province 5', 'Country 5', 'PostCode 5', 'contact5@example.com', '5678901234');

-- User Table
INSERT INTO User (Name, AccountID, AssetID) VALUES
	('User 1', 1, NULL),
	('User 2', 2, NULL),
	('User 3', 3, NULL),
	('User 4', 4, NULL),
	('User 5', 5, NULL);

-- Institution Table
INSERT INTO Institution (Name, Description, AccountID, ContactID) VALUES
	('Institution 1', 'Description 1', 1, 1),
	('Institution 2', 'Description 2', 2, 2),
	('Institution 3', 'Description 3', 3, 3),
	('Institution 4', 'Description 4', 4, 4),
	('Institution 5', 'Description 5', 5, 5);

-- Instructor Table
INSERT INTO Instructor (Name, AccountID, InstitutionID, ContactID) VALUES
	('Instructor 1', 1, 1, 1),
	('Instructor 2', 2, 2, 2),
	('Instructor 3', 3, 3, 3),
	('Instructor 4', 4, 4, 4),
	('Instructor 5', 5, 5, 5);

-- UserAsset Table
INSERT INTO UserAsset (UserID, InstitutionID) VALUES
	(1, 1),
	(2, 2),
	(3, 3),
	(4, 4),
	(5, 5);

-- Course Table
INSERT INTO Course (Title, Description) VALUES
	('Course 1', 'Course 1 description'),
	('Course 2', 'Course 2 description'),
	('Course 3', 'Course 3 description'),
	('Course 4', 'Course 4 description'),
	('Course 5', 'Course 5 description');

-- Section Table
INSERT INTO Section (Title, Description, Start, Duration, At, UnitPrice, RepetitionID, InstructorID, CourseID) VALUES
	('Section 1', 'Section 1 description', '2023-07-05', 3600, 3600, 100, NULL, 1, 1),
	('Section 2', 'Section 2 description', '2023-07-06', 3600, 3600, 200, NULL, 2, 2),
	('Section 3', 'Section 3 description', '2023-07-07', 3600, 3600, 300, NULL, 3, 3),
	('Section 4', 'Section 4 description', '2023-07-08', 3600, 3600, 400, NULL, 4, 4),
	('Section 5', 'Section 5 description', '2023-07-09', 3600, 3600, 500, NULL, 5, 5);

-- Class Table
INSERT INTO Class (Title, 'Index', Time, Duration, UnitPrice, Canceled, Completed, Remark, InstructorID, SectionID, UpdatedAt) VALUES
	('Class 1', 1, '2023-07-05 10:00:00', 3600, 100, 0, 0, NULL, 1, 1, DATETIME('now', 'localtime')),
	('Class 2', 1, '2023-07-06 10:00:00', 3600, 200, 0, 0, NULL, 2, 2, DATETIME('now', 'localtime')),
	('Class 3', 1, '2023-07-07 10:00:00', 3600, 300, 0, 0, NULL, 3, 3, DATETIME('now', 'localtime')),
	('Class 4', 1, '2023-07-08 10:00:00', 3600, 400, 0, 0, NULL, 4, 4, DATETIME('now', 'localtime')),
	('Class 5', 1, '2023-07-09 10:00:00', 3600, 500, 0, 0, NULL, 5, 5, DATETIME('now', 'localtime'));
`
