package database

import (
	"context"
	"log"
	"math/rand/v2"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/hayohtee/social/internal/data"
	"github.com/hayohtee/social/internal/repository"
)

func Seed(repo repository.Repository) {
	faker := gofakeit.NewFaker(source.NewCrypto(), true)
	ctx := context.Background()
	var wg sync.WaitGroup

	users := generateUsers(faker, 200)
	for _, user := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := repo.Users.Create(ctx, user, nil); err != nil {
				log.Println("error creating user:", err)
			}
		}()
	}
	wg.Wait()

	posts := generatePosts(users, faker, 400)
	for _, post := range posts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := repo.Posts.Create(ctx, post); err != nil {
				log.Println("error creating post:", err)
			}
		}()
	}
	wg.Wait()

	comments := generateComments(users, posts, faker, 600)
	for _, comment := range comments {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := repo.Comments.Create(ctx, comment); err != nil {
				log.Println("error creating comment:", err)
			}
		}()
	}
	wg.Wait()
	log.Println("seeding complete")
}

func generateUsers(faker *gofakeit.Faker, num int) []*data.User {
	users := make([]*data.User, num)

	for i := 0; i < num; i++ {
		user := data.User{
			Username: faker.Username(),
			Email:    faker.Email(),
		}

		password := faker.Password(true, true, true, true, false, 12)
		err := user.Password.Set(password)
		if err != nil {
			log.Println(err)
			continue
		}

		users[i] = &user
	}

	return users
}

func generatePosts(users []*data.User, faker *gofakeit.Faker, num int) []*data.Post {
	posts := make([]*data.Post, num)

	for i := 0; i < num; i++ {
		post := data.Post{
			Content: faker.Sentence(100),
			Title:   faker.Sentence(10),
			UserID:  users[rand.IntN(len(users))].ID,
		}
		numTags := rand.IntN(5)
		tags := make([]string, numTags)
		for j := 0; j < numTags; j++ {
			tags[j] = faker.Word()
		}
		post.Tags = tags
		posts[i] = &post
	}

	return posts
}

func generateComments(users []*data.User, posts []*data.Post, faker *gofakeit.Faker, num int) []*data.Comment {
	comments := make([]*data.Comment, num)

	for i := 0; i < num; i++ {
		comments[i] = &data.Comment{
			Content: faker.Comment(),
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
		}
	}

	return comments
}
