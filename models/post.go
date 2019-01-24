package models

// Post is a data object for the post table
type Post struct {
	ID        int
	Version   int
	Parent    int
	Title     string
	Body      string
	Created   string
	Published string
}

// PostDataReader defines methods to retrieve post data
type PostDataReader interface {
	GetPosts() (posts []Post, err error)
}

// GetPosts returns a list of posts
func (db *DB) GetPosts() (posts []Post, err error) {
	rows, err := db.Query(`
		SELECT id, title, body, published
		FROM post
		WHERE parent = 0
		LIMIT 10, 0
		SORT BY created DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Body, &post.Published)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}
