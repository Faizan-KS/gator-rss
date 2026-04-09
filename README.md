This is an RSS Aggre**gator** made in Go with PostgreSQL as database for storing the rss feeds

Disclaaimer: Please do not ddos the rss feeds please be respectful

Prerequisites
1. To run the cli users need to have Go 1.24+ and PostgreSQL installed
2. Make a .gatorconfig.json file with {"db_url":"postgres://postgres:enter_your_local_postgres_password@localhost:5432/gator?sslmode=disable","current_user_name":""}
3. Install the cli using go install .

Usage
Assuming you have performed the pre-requisites
1. Open a terminal and use the "gator-rss" command
2. See the usage as the command exits
3. There are several commands but the initial command to use is register
4. Register the name for the rssfeed you want to follow

These are the following commands the user can interact with
1. register - register the users
2. login - login to the user
3. users - list all the users and see who the current user is
4. reset - reset the feeds
5. addfeed - add the rss_feed you want to have with name and url you 
6. follow - follow other users rss_feed
7. unfollow - unfollow other users rss_feed
8. allfeeds - show all feeds added by all the users
9. following - show all feeds you are following
10. agg - actively get the feeds of the rss_feeds you are following in intervals of X
11. myfeedposts - shows the recent posts of the feeds 
