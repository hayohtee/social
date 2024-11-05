package database

import (
	"context"
	"log"
	"math/rand/v2"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/hayohtee/social/internal/model"
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
			if err := repo.Users.Create(ctx, user); err != nil {
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

func generateUsers(faker *gofakeit.Faker, num int) []*model.User {
	users := make([]*model.User, num)

	for i := 0; i < num; i++ {
		users[i] = &model.User{
			Username: faker.Username(),
			Email:    faker.Email(),
			Password: faker.Password(true, true, true, true, false, 12),
		}
	}

	return users
}

func generatePosts(users []*model.User, faker *gofakeit.Faker, num int) []*model.Post {
	posts := make([]*model.Post, num)

	for i := 0; i < num; i++ {
		posts[i] = &model.Post{
			Content: faker.Sentence(100),
			Title:   faker.Sentence(10),
			UserID:  users[rand.IntN(len(users))].ID,
			Tags:    []string{},
		}
	}

	return posts
}

func generateComments(users []*model.User, posts []*model.Post, faker *gofakeit.Faker, num int) []*model.Comment {
	comments := make([]*model.Comment, num)

	for i := 0; i < num; i++ {
		comments[i] = &model.Comment{
			Content: faker.Comment(),
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
		}
	}

	return comments
}
