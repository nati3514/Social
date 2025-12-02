package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/nati3514/Social/internal/store"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve",
	"frank", "grace", "heidi", "ivan", "judy",
	"ken", "laura", "mallory", "nancy", "oscar",
	"peggy", "quinn", "ruth", "sam", "trent",
	"ursula", "victor", "wendy", "xavier", "yvonne",
	"zack", "brian", "claire", "dennis", "elena",
	"felix", "george", "hannah", "ian", "jess",
	"karen", "leo", "mike", "nina", "oliver",
	"paula", "quentin", "rachel", "steve", "tina",
	"ulysses", "vanessa", "will", "xena", "yuri",
}

var titles = []string{
	"Mastering the Art of Focus",
	"Why Simplicity Wins Every Time",
	"The Power of Small Habits",
	"Building Resilient Systems",
	"How to Stay Creative Under Pressure",
	"Scaling Without Losing Control",
	"Lessons from Failed Projects",
	"Designing for Humans, Not Machines",
	"The Future of Remote Work",
	"How Open Source Changed My Career",
	"A Beginner's Guide to APIs",
	"Think Fast, Code Faster",
	"What I Learned from Debugging All Night",
	"Turning Ideas into Action",
	"Automation That Actually Works",
	"Why Documentation Matters",
	"Balancing Innovation and Stability",
	"From Chaos to Clean Code",
	"Breaking Down Complex Problems",
	"The Hidden Cost of Technical Debt",
}

var contents = []string{
	"In today's fast-paced tech world, focusing on small, consistent improvements often leads to the biggest results.",
	"Building systems that scale starts with understanding the trade-offs between performance and maintainability.",
	"Every developer faces burnout at some point — what matters is how you recover and grow from it.",
	"Design isn't just about visuals; it's about how people experience your product every day.",
	"Automation is powerful, but only when you automate the right things.",
	"The best code isn't the one that works — it's the one that's easy to understand six months later.",
	"Innovation starts with curiosity, not just a big idea.",
	"Good documentation saves more time than any single line of code you'll ever write.",
	"Failing fast is great, but learning fast is better.",
	"Security isn't a feature; it's a mindset that should be built into every stage of development.",
	"Remote work has redefined collaboration, making trust and transparency more important than ever.",
	"Debugging is like detective work — every clue brings you closer to the truth.",
	"Small teams can build big things when they focus on clarity and communication.",
	"Consistency beats talent when it comes to long-term success in software projects.",
	"Technical debt is inevitable — but unmanaged debt can quietly kill your velocity.",
	"A solid architecture grows with your product; a rigid one eventually breaks it.",
	"Open-source contributions are one of the best ways to learn from real-world code.",
	"Sometimes the hardest part of building software is knowing when to stop adding features.",
	"Clean code isn't about rules — it's about empathy for the next developer.",
	"Great engineers don't just write code; they solve problems that matter.",
}

var tags = []string{
	"golang",
	"webdev",
	"productivity",
	"startups",
	"design",
	"ai",
	"machine-learning",
	"cloud",
	"opensource",
	"devops",
	"security",
	"career",
	"frontend",
	"backend",
	"databases",
	"testing",
	"remote-work",
	"innovation",
	"clean-code",
	"automation",
}

var comments = []string{
	"Great insights! Thanks for sharing.",
	"I never thought about it that way before.",
	"This was super helpful, appreciate the breakdown!",
	"Can you share more about your approach to this?",
	"Loved this post — really made me think.",
	"Simple but powerful explanation!",
	"I've been struggling with this concept, and this helped a lot.",
	"Totally agree with your point about consistency.",
	"This is gold. Bookmarking for later!",
	"I tried this and it actually worked perfectly.",
	"Could you expand on the scaling part?",
	"Your writing style makes complex topics so easy to follow.",
	"Awesome example — makes the concept crystal clear.",
	"Interesting perspective! I might try this in my next project.",
	"I've had similar issues — glad to see I'm not alone.",
	"This deserves more attention, well said!",
	"Great content as always — keep it coming!",
	"Love how you explained this in plain English.",
	"This comment section is full of great discussions!",
	"I learned something new today, thank you!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	// Begin a transaction for the entire seeding operation
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error beginning transaction:", err)
		return
	}
	defer tx.Rollback() // Rollback if we don't commit

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		return
	}

	log.Println("Database seeded successfully")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123456",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(contents))],
			Tags: []string{
				titles[rand.Intn(len(titles))],
				titles[rand.Intn(len(titles))],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}

	}

	return cms
}
