package repository

type User struct {
	ID             int64  `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	HashedPassword string `db:"hashed_password" json:"hashed_password"`
}

type Url struct {
	ID    int64  `db:"id" json:"id"`
	Alias string `db:"alias" json:"alias"`
	URL   string `db:"url" json:"url"`
	//UserId int64  `db:"user_id" json:"user_id"`
}
