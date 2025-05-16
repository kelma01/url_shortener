package entities

//dbde tutulacak olan url tuplelarinin schemasi
type URL struct {
	ID 	  		int    		`db:"id"`
	OriginalURL	string		`db:"original_url"`
	ShortURL 	string		`db:"short_url"`
	CreatedAt	time.Time	`db:"created_at"`
	DeletedAt	*time.Time	`db:"deleted_at"`	//* işareti burada nullable anlamina gelmekte
	ExpiresAt	*time.Time	`db:"expires_at"`
	UsageCount	int			`db:"usage_count"`
}