package tm_poznan

// https://app.quicktype.io/

type Welcome struct {
	Data  map[string]DatumValue `json:"data"`
	Links WelcomeLinks          `json:"links"`
	Meta  Meta                  `json:"meta"`
}

type DatumValue struct {
	Data []DatumDatum `json:"data"`
}

type DatumDatum struct {
	Type          FluffyType       `json:"type"`
	ID            int64            `json:"id"`
	Attributes    PurpleAttributes `json:"attributes"`
	Relationships Relationships    `json:"relationships"`
	Links         DatumLinks       `json:"links"`
}

type PurpleAttributes struct {
	TypesID         int64    `json:"types_id"`
	GenresID        int64    `json:"genres_id"`
	Archive         int64    `json:"archive"`
	Position        int64    `json:"position"`
	Bilety24URL     string   `json:"bilety24_url"`
	ImgURL          string   `json:"img_url"`
	Img2URL         string   `json:"img2_url"`
	PosterURL       string   `json:"poster_url"`
	Title           string   `json:"title"`
	Excerpt         string   `json:"excerpt"`
	ReleaseDate     *string  `json:"release_date"`
	AgeRestrictions *string  `json:"age_restrictions"`
	Duration        *string  `json:"duration"`
	Location        Location `json:"location"`
}

type DatumLinks struct {
	Slug string `json:"slug"`
}

type Relationships struct {
	Callendar Callendar `json:"callendar"`
}

type Callendar struct {
	Data []CallendarDatum `json:"data"`
}

type CallendarDatum struct {
	Type       PurpleType       `json:"type"`
	ID         int64            `json:"id"`
	Attributes FluffyAttributes `json:"attributes"`
}

type FluffyAttributes struct {
	ShowsID     int64   `json:"shows_id"`
	Publication string  `json:"publication"`
	Bilety24URL string  `json:"bilety24_url"`
	TagsTitle   *string `json:"tags_title"`
}

type WelcomeLinks struct {
	First string      `json:"first"`
	Last  string      `json:"last"`
	Prev  interface{} `json:"prev"`
	Next  interface{} `json:"next"`
}

type Meta struct {
	CurrentPage int64  `json:"current_page"`
	From        int64  `json:"from"`
	LastPage    int64  `json:"last_page"`
	Path        string `json:"path"`
	PerPage     string `json:"per_page"`
	To          int64  `json:"to"`
	Total       int64  `json:"total"`
}

type Location string

const (
	Away  Location = "away"
	Guest Location = "guest"
	Own   Location = "own"
)

type PurpleType string

const (
	ShowsCallendar PurpleType = "shows_callendar"
)

type FluffyType string

const (
	Shows FluffyType = "shows"
)
