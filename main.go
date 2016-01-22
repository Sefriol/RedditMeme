package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
	"github.com/jzelinskie/geddit"
    "sort"
    "sync"
)
//Match ...
type Match struct {
    CommentID string
    Answered bool
}
var session *geddit.LoginSession
var matches []*Match
var memeTrigger = "Would you kindly explain "
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
        Limit: 50,
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
    var comments = make(chan *geddit.Comment, 1000)
    var wg sync.WaitGroup
    for _, s := range submissions {
        wg.Add(1)
        comment, err := session.Comments(s)
        if err != nil {
            continue
        }
        go func() {    
            // fmt.Printf("CommentDetect. Title: %s\n\n", s.Title)
            // fmt.Printf("CommentDetected. ID: %d\n\n", wg)
            CommentDetect("Dota", comment, comments, wg)
            defer wg.Done()
        }()
    }
    wg.Wait()
    close(comments)
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
    for elem := range comments {
        wg.Add(1)
        go func() {
            MemeCheck(elem, wg)
        }()
    }
    wg.Wait()
    // for _, s := range submissions {
        // 	fmt.Printf("Title: %s\nAuthor: %s Comments: %s\n\n", s.Title, s.Author,comments)
    // }
    
    // Upvote the first post
    //session.Vote(submissions[0], geddit.UpVote)
}

func CommentDetect(detect string, comment []*geddit.Comment,  comments chan *geddit.Comment, wg sync.WaitGroup)  {
    var temp *Match
    for _, c := range comment {
        if strings.Contains(c.Body, detect) {
            i := sort.Search(len(matches),func(i int) bool { return matches[i].CommentID >= c.FullID })
            if i < len(matches) && matches[i].CommentID == c.FullID {
                fmt.Printf("Match found. Do nothing.\n\n")
                // Match found. Do nothing.
            } else {
                comments <- c
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
            go func() {
                wg.Add(1)
                commentDetect(detect, c.Replies,comments,wg)
            }()
        }
    }
    defer wg.Done()
}

func MemeCheck(comment *geddit.Comment,wg sync.WaitGroup)  {
    //Trim text for meme search
    commenText :=  comment.Body[strings.Index(comment.Body, memeTrigger) + len(memeTrigger):]    
    
    replyKYM := "p&gt;p&gt;Read the full meme in [knowyourmeme](%s)! ---- ^This ^message ^was ^created ^by ^a ^bot [^[Contact ^creator]](http://np.reddit.com/message/compose/?to=&amp;amp;subject=TweetsInCommentsBot)[^[Github]](https://github.com/)"
    replyDotaMeme := `---- 
    p&gt;p&gt; ^Well ^memed? ^No? ^Improve ^this ^meme [^here](http://np.reddit.com/message/compose/?to=&amp;amp;subject=ALTER)
    p&gt;^This ^message ^was ^created ^by ^a ^bot [^[Contact ^creator]](http://np.reddit.com/message/compose/?to=&amp;amp;subject=)[^[Github]](https://github.com/)`
    fmt.Printf("id: %s Comment: %s\n\n",comment.FullID, commenText)
    urlquery := url.QueryEscape(strings.TrimSpace(commenText))
    fmt.Printf("urlquery %s\n\n",urlquery)
    
    knowyourmemes, err := getMemes(urlquery)
    if err != nill {
        meme := knowyourmemes[0][:strings.Index(knowyourmemes[0], "h2. Origin")]
    }
    knowyourmemes
    defer wg.Done()
}