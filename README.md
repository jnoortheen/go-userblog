# AppStory
This is a port of [AppStory](https://github.com/jnoortheen/appstory) to Go. A multi-user blog using 
[go-buffalo](https://github.com/gobuffalo/buffalo).

## Tech Stack
- Go
- Buffalo
- Bootstrap3

## Database Setup

using sqlite3 as data backend

### Create Your Databases

Ok, so you've edited the "database.yml" file and started sqlite3, now Buffalo can create the databases in that file for you:

	$ buffalo db create -a

To run all pending migrations

	$ buffalo db migrate up

## Starting the Application

Run the binary found here at [dist](./dist)
	$ ./muserblog

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) blog home page.

## Functionality
-[ ] User Accounts:
  Users account activity is implemented using secured cookies. Usernames are maintained to be unique. The password is stored as Hash values with salted compound.
  User activities  
  - Sign-up
  - Sign-in 
  - Sign-out 
-[ ] Blog Posts:
  Registered users can post to the blog. That post can later be edited or deleted by the author.
-[ ] Commenting:
  Posts can be commented by any of registered user. Unregistered users can view the comments. Registered users can edit or delete their own comments.
-[ ] Like/Dislike:
  Each of the post can be like/disliked by other registered users (Not the Author). User can't dislike a post if they haven't already liked it.
