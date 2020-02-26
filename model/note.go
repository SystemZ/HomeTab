package model

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/golang-commonmark/markdown"
	"html/template"
	"log"
	"strings"
	"time"
)

type Note struct {
	Id        uint          `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	Title     string        `gorm:"column:title" json:"title"`
	Body      string        `gorm:"column:body" json:"body"`
	BodyShort string        `gorm:"-"`
	BodyMd    template.HTML `gorm:"-"`
	ProjectId uint          `gorm:"column:project_id" json:"projectId"`
	Tags      string        `gorm:"-"`
	CreatedAt *time.Time    `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt *time.Time    `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time    `gorm:"column:deleted_at" json:"deletedAt"`
}

type NoteTag struct {
	Id        uint       `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	NoteId    uint       `gorm:"column:note_id" json:"note_id"`
	Name      string     `gorm:"column:name" json:"name"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (note *Note) BeforeSave(scope *gorm.Scope) (err error) {
	//FIXME this is wrong method
	// delete all tags
	DB.Unscoped().Where("note_id = ?", note.Id).Delete(NoteTag{})
	// add all tags from the scratch
	tagsList := strings.Split(note.Tags, ",")
	for _, tag := range tagsList {
		newTag := NoteTag{
			NoteId: note.Id,
			Name:   strings.TrimSpace(tag),
		}
		DB.Save(&newTag)
	}
	return nil
}

func CreateNote(title string, body string, tags string, projectId uint) Note {
	// prepare data
	var note Note
	note.Title = title
	note.Body = body
	note.ProjectId = projectId
	now := time.Now()
	note.CreatedAt = &now
	note.UpdatedAt = &now
	// save note to DB
	err := DB.Save(&note).Error
	if err != nil {
		log.Printf("%v", err)
	}
	// add all tags
	tagsList := strings.Split(tags, ",")
	for _, tag := range tagsList {
		newTag := NoteTag{
			NoteId: note.Id,
			Name:   strings.TrimSpace(tag),
		}
		DB.Save(&newTag)
	}
	// all done, return note ID created
	return note
}

func NoteLongList() (result []Note) {
	query := `
SELECT
  notes.id,
  notes.project_id,
  notes.title,
  notes.body,
  LEFT(notes.body , 16),
  (SELECT GROUP_CONCAT(note_tags.name SEPARATOR ',') FROM note_tags WHERE note_tags.note_id = notes.id) AS tags,
  notes.created_at,
  notes.updated_at
FROM notes
GROUP BY notes.id
ORDER BY notes.id DESC
`
	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list Note
		err := rows.Scan(&list.Id, &list.ProjectId, &list.Title, &list.Body, &list.BodyShort, &list.Tags, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return
		}
		result = append(result, list)
	}
	return result
}

func OneNote(id uint) (result []Note) {
	query := `
SELECT
  notes.id,
  notes.project_id,
  notes.title,
  notes.body,
  LEFT(notes.body , 16),
  (SELECT GROUP_CONCAT(note_tags.name SEPARATOR ',') FROM note_tags WHERE note_tags.note_id = notes.id) AS tags,
  notes.created_at,
  notes.updated_at
FROM notes
WHERE notes.id = ?
GROUP BY notes.id
ORDER BY notes.id DESC
`
	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list Note
		err := rows.Scan(&list.Id, &list.ProjectId, &list.Title, &list.Body, &list.BodyShort, &list.Tags, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return
		}

		md := markdown.New(markdown.XHTMLOutput(true))
		list.BodyMd = template.HTML(md.RenderToString([]byte(list.Body)))
		result = append(result, list)
	}
	return result
}
