package main

import (
	"bufio"
	"fmt"
	"os"
    "strings"
    "net/url"
	"github.com/jzelinskie/geddit"
)

var session *geddit.LoginSession
// Please don't handle errors this way.
func main() {
	//Read reddit username and password
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username: ")
	usr, _ := reader.ReadString('\n')
    usr = strings.TrimSpace(usr)
	fmt.Print("Enter Password: ")
	pwd, _ := reader.ReadString('\n')
    pwd = strings.TrimSpace(pwd)

	// Login to reddit
	session, _ = geddit.NewLoginSession(
		usr,
		pwd,
		"gedditAgent v1",
	)

	// Set listing options
	subOpts := geddit.ListingOptions{
		Limit: 100,
	}
    
	// Get specific subreddit submissions, sorted by new
	submissions, _ := session.SubredditSubmissions("Dota2", geddit.TopSubmissions, subOpts)
    fmt.Printf("Len: %d\n\n", len(submissions))
    subOpts = geddit.ListingOptions{
		Limit: 35,
	}
    submissionshot, _ := session.SubredditSubmissions("Dota2", geddit.HotSubmissions, subOpts)
    submissions = append(submissions,submissionshot...)

    fmt.Printf("Len: %d Len: %d\n\n", len(submissionshot),len(submissions))
    var comments []*geddit.Comment
    for _, s := range submissions {
        comment, err := session.Comments(s)
        if err != nil {
            continue
        }
        comments = append(comments, commentDetect("Except that's what they do. Riki has had so many reworks with the intention of making the hero competitively viable instead of just buffing him because it WOULD break 1-2k. They are actively trying to find a position for him everywhere, not decent in competitive and literally unbeatable for lower skill.", comment)...)
    }
    fmt.Printf("Len commets: %d\n\n", len(comments))
    fmt.Print("Enter Meme: ")
    urlquery, _ := reader.ReadString('\n')
    urlquery = url.QueryEscape(strings.TrimSpace(urlquery))
    fmt.Printf("5/7 %s\n\n",urlquery)
    
    knowyourmemes, _ := getMemes(urlquery)
    for _, s := range knowyourmemes {
        fmt.Printf("Meme: %d %s\n\n", len(knowyourmemes), s.Body)
    }
    //http://rkgk.api.searchify.com/v1/indexes/kym_production/instantlinks?query=jotain&fetch=*
    
	// Print title and author of each submission
    // comments, _ := session.Comments(submissions[0])
    // for _, c := range comments {
	// 	fmt.Printf("Comment: %s\n\n",c.String())
	// }
	// for _, s := range submissions {
        
	// 	fmt.Printf("Title: %s\nAuthor: %s Comments: %s\n\n", s.Title, s.Author,comments)
	// }

	// Upvote the first post
	//session.Vote(submissions[0], geddit.UpVote)
}

func commentDetect(detect string, comments []*geddit.Comment)([]*geddit.Comment)  {
    var finalcomments []*geddit.Comment
    for _, c := range comments {
        if strings.Contains(c.Body, detect) {
            finalcomments = append(finalcomments, c)
        }
        if len(c.Replies) > 0 {
            finalcomments = append(finalcomments,commentDetect(detect, c.Replies)...)
        }
        
    }
    return finalcomments
}