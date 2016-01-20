package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
	"github.com/jzelinskie/geddit"
    "sort"
)
//Match ...
type Match struct {
    CommentID string
    Answered bool
}
var session *geddit.LoginSession
var matches []*Match
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
        comments = append(
            comments, 
            commentDetect("It should hit all non creep units IMO, much like you don't pick PL into Ember you wouldn't pick PL into Riki. Ember in a way has a 5 second CD Riki ult", comment)...)
    }
    fmt.Printf("Len commets: %d\n\n", len(comments))
    
    // fmt.Print("Enter Meme: ")
    // urlquery, _ := reader.ReadString('\n')
    // urlquery = url.QueryEscape(strings.TrimSpace(urlquery))
    // fmt.Printf("5/7 %s\n\n",urlquery)
    
    // knowyourmemes, _ := getMemes(urlquery)
    // for _, s := range knowyourmemes {
    //     fmt.Printf("Meme: %d %s\n\n", len(knowyourmemes), s.Body)
    // }
    //http://rkgk.api.searchify.com/v1/indexes/kym_production/instantlinks?query=jotain&fetch=*
    
    // Print title and author of each submission
    for _, c := range comments {
         fmt.Printf("id:%s Comment: %s\n\n",c.FullID, c.String())
    }
    // for _, s := range submissions {
        // 	fmt.Printf("Title: %s\nAuthor: %s Comments: %s\n\n", s.Title, s.Author,comments)
    // }
    
    // Upvote the first post
    //session.Vote(submissions[0], geddit.UpVote)
}

func commentDetect(detect string, comments []*geddit.Comment)([]*geddit.Comment)  {
    var finalcomments []*geddit.Comment
    var temp *Match
    
    for _, c := range comments {
        if strings.Contains(c.Body, detect) {
            i := sort.Search(len(matches),func(i int) bool { return matches[i].CommentID >= c.FullID })
            if i < len(matches) && matches[i].CommentID == c.FullID {
                fmt.Printf("Match found. Do nothing.\n\n")
                // Match found. Do nothing.
            } else {
                finalcomments = append(finalcomments, c)
                matches = append(matches, temp)
                copy(matches[i+1:], matches[i:])
                new := Match{
                    Answered: false,
                    CommentID: c.FullID,
                }
                matches[i] = &new
            }
        }
        if len(c.Replies) > 0 {
            finalcomments = append(finalcomments,commentDetect(detect, c.Replies)...)
        }
        
    }
    return finalcomments
}