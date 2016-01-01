// Copyright 2014 The Too Authors. All rights reserved.

package too

import "github.com/garyburd/redigo/redis"

type Engine struct {
	c     redis.Conn
	class string

	AutoUpdate bool

	Likes    Rater
	Dislikes Rater

	Similars    Similars
	Suggestions Suggestions
}

// New returns a new engine for given class connected to Redis server at addr.
func New(url, class string) (*Engine, error) {
	c, err := redis.DialURL(url)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		c:          c,
		class:      class,
		AutoUpdate: true,
	}
	e.Likes = Rater{e, "likes"}
	e.Dislikes = Rater{e, "dislikes"}
	e.Similars = Similars{e}
	e.Suggestions = Suggestions{e}
	return e, nil
}

func (e Engine) DisableAutoUpdate() {
	e.AutoUpdate = false
}

func (e Engine) EnableAutoUpdate() {
	e.AutoUpdate = true
}

func (e Engine) Update(user User) error {
	err := e.Similars.update(user)
	if err != nil {
		return err
	}

	err = e.Suggestions.update(user)
	if err != nil {
		return err
	}
	return nil
}
