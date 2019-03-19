package core

import (
  //"context"
  "bytes"
  "strconv"
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  //"reflect"
  //"strconv"
)

type GitlabOwner struct {
  ID           int `json:"id"`
  Name         string `json:"name"`
  Username     string
  State        string
  AvatarUrl    string `json:"avatar_url"`
  WebUrl       string `json:"web_url"`
  //CreatedAt    string
  //Bio          string
  //Location     string
  //PublicEmail  string
  //Skype        string
  //Linkedin     string
  //Twitter      string
  //WebsiteUrl   string
  //Organization string
}


type GitlabRepository struct {
  Owner         GitlabOwner
  ID            int
  Name          string
  FullName      string
  CloneURL      string
  URL           string
  DefaultBranch string
  Description   string
  Homepage      string
}

type GitlabRepo struct {
  Owner         GitlabOwner
  ID            int
  Path          string
  Name          string
  Description   string `json:"description"`
  DefaultBranch string `json:"default_branch"`
  HttpUrlToRepo string `json:"http_url_to_repo"`
}

type Response struct {
  repos       []GitlabRepo
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []string) bool {
  if len(a) != len(b) {
    return false
  }
  for i, v := range a {
    if v != b[i] {
      return false
    }
  }
  return true

}

func HandleApiCall(apiUrl string, token string) ([]uint8, error, http.Header) {
  // pagination is stored in headers in the gitlab api
  body := make([]uint8, 0)
  url := apiUrl

  // while loop to catch all paginated keys
   resp, err := http.Get(url)
   if err != nil {
     return nil, err, nil
     }
   defer resp.Body.Close()
   header := resp.Header

   bodyNext, _ := ioutil.ReadAll(resp.Body)
   body = append(body, bodyNext...)
   return body, nil, header

}



func GetUserId(login string, apiUrl string, token string) (string, error) {
  // string are more efficient to concat as bytes
  var b bytes.Buffer
  users := make([]GitlabOwner, 0)

  b.WriteString(apiUrl)
  b.WriteString(fmt.Sprintf("/users?username=%s&private_token=%s", login, token))

  body, err, _ := HandleApiCall(b.String(), token)
  if err != nil {
    return "", err
    }

  _ = json.Unmarshal(body, &users)

  // assume there is only a single user with this username
  userId := strconv.Itoa(users[0].ID)

  return userId, nil
}


func GetUserOrOrganization(login string, sess *Session) (*GitlabOwner, error) {
  var b bytes.Buffer

  // make new api call to get the id from login
  userId, _ := GetUserId(login, *sess.Options.BaseUrl, sess.GithubAccessToken)
  owner := GitlabOwner{}

  b.WriteString(*sess.Options.BaseUrl)
  b.WriteString(fmt.Sprintf("/users/%s?private_token=%s", userId, sess.GithubAccessToken))

  body, err, _ := HandleApiCall(b.String(), sess.GithubAccessToken)
  if err != nil {
    return nil, err
    }
  errJson := json.Unmarshal(body, &owner)

  if errJson != nil {
    return nil, err
  }


  return &owner, nil
}

func GetRepositoriesFromOwner(login string, sess *Session) ([]GitlabRepo, error) {
  var b bytes.Buffer

  b.WriteString(*sess.Options.BaseUrl)
  b.WriteString(fmt.Sprintf("/projects?per_page=100&private_token=%s", sess.GithubAccessToken))
  apiUrl := b.String()
  url := apiUrl
  reps := make([]GitlabRepo, 0)

  // gitlab api pagination is in the header
  // while loop that breaks when the last page is reached
  for {
    body, err, header := HandleApiCall(url, sess.GithubAccessToken)
    if err != nil {
      return nil, err
      }
    repsNext := make([]GitlabRepo, 0)
    errJson := json.Unmarshal(body, &repsNext)
    reps = append(reps, repsNext...)
    if errJson != nil {
      return nil, errJson
    }
    if !Equal(header["X-Page"], header["X-Total-Pages"]) {
      url = apiUrl + fmt.Sprintf("&page=%s", header["X-Next-Page"][0])
      continue
      }
    break
   }


  return reps, nil
  //opt := &gitlab.ListProjectsOptions{
  //  //Simple: true,
  //}
  //forkParent:=gitlab.ForkParent{}

  //for {
  //  repos, resp, err := client.Projects.ListProjects(opt)
  //  if err != nil {
  //    return allRepos, err
  //  }
  //  for _, repo := range repos {

  //    if *repo.ForkedFromProject == forkParent {
  //      r := GitlabRepository{
  //        //Owner:         repo.Owner.Username,
  //        ID:            repo.ID,
  //        Name:          repo.Name,
  //        CloneURL:      repo.HTTPURLToRepo,
  //        URL:           repo.WebURL,
  //        DefaultBranch: repo.DefaultBranch,
  //        Description:   repo.Description,
  //      }
  //      allRepos = append(allRepos, &r)
  //    }
  //  }
  //  if resp.NextPage == 0 {
  //    break
  //  }
  //  opt.Page = resp.NextPage
  //}

  //return allRepos, nil
}

//func GetOrganizationMembers(login *string, client *gitlab.Client) ([]*GitlabOwner, error) {
//  var allMembers []*GitlabOwner
//  loginVal := *login
//  ctx := context.Background()
//  opt := &github.ListMembersOptions{}
//  for {
//    members, resp, err := client.Organizations.ListMembers(ctx, loginVal, opt)
//    if err != nil {
//      return allMembers, err
//    }
//    for _, member := range members {
//      allMembers = append(allMembers, &GithubOwner{Login: member.Login, ID: member.ID, Type: member.Type})
//    }
//    if resp.NextPage == 0 {
//      break
//    }
//    opt.Page = resp.NextPage
//  }
//  return allMembers, nil
//}
