package entities

import "time"

//dbde tanimlanmis olan columnlarin go'da nasil kullanilacaklarini tanimlayan struct semasi
type URL struct {
	ID 	  		int    		`json:"id"`
	OriginalURL	string		`json:"original_url"`
	ShortURL 	string		`json:"short_url"`
	CreatedAt	time.Time	`json:"created_at"`
	DeletedAt	*time.Time	`json:"deleted_at"`	//* işareti burada nullable anlamina gelmekte
	ExpiresAt	*time.Time	`json:"expires_at"`
	UsageCount	int			`json:"usage_count"`
}
func (URL) TableName() string {
    return "url_table"
}