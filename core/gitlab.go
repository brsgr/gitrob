package core

import (
  //"context"
  "bytes"
  "strconv"
  "fmt"
  "net/http"
	"io/ioutil"
	"encoding/json"
  //"strconv"
)

type GitlabOwner struct {
  ID           int `json:"id"`
  Name         string `json:"name"`
  Username     string
  State        string
  AvatarUrl    string `json:"avatar_url"`
  WebUrl			 string `json:"web_url"`
  //CreatedAt    string
  //Bio				   string
  //Location     string
  //PublicEmail  string
	//Skype				 string
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
	Owner				  GitlabOwner
	ID						int
	Path					string
	Name					string
	Description		string `json:"description"`
	DefaultBranch string `json:"default_branch"`
	HttpUrlToRepo string `json:"http_url_to_repo"`
}

type Response struct {
	repos				[]GitlabRepo
}

func GetUserId(login string, apiUrl string, token string) (string, error) {
  // string are more efficient to concat as bytes
  var b bytes.Buffer
	users := make([]GitlabOwner, 0)

  b.WriteString(apiUrl)
  b.WriteString(fmt.Sprintf("/users?username=%s&private_token=%s", login, token))

  resp, err := http.Get(b.String())

  if err != nil {
    return "", err
    }

	defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  _ = json.Unmarshal(body, &users)

	userId := strconv.Itoa(users[0].ID)

  return userId, nil
}


func GetUserOrOrganization(login string, sess *Session) (*GitlabOwner, error) {
  var b bytes.Buffer

  // make new api call to get the id from login
  userId, _ := GetUserId(login, *sess.Options.BaseUrl, sess.GithubAccessToken)

  b.WriteString(*sess.Options.BaseUrl)
  b.WriteString(fmt.Sprintf("/users/%s?private_token=%s", userId, sess.GithubAccessToken))

  resp, err := http.Get(b.String())

  if err != nil {
    return nil, err
  }

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	owner := GitlabOwner{}
	errJson := json.Unmarshal(body, &owner)

	if errJson != nil {
		return nil, err
	}


  return &owner, nil
}

func GetRepositoriesFromOwner(login string, sess *Session) ([]GitlabRepo, error) {
  //var b bytes.Buffer

  //b.WriteString(*sess.Options.BaseUrl)
  //b.WriteString(fmt.Sprintf("/projects?primate_token=%s", sess.GithubAccessToken))

	urlString := "https://git.ckmnet.co/api/v4/projects?private_token=7qNxJyvYUP4shCGMP2-Y"
  resp, err := http.Get(urlString)
  if err != nil {
    return nil, err
  }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	reps := make([]GitlabRepo, 0)
	//responseObject := Response{}

	errJson := json.Unmarshal(body, &reps)

	if errJson != nil {
		return nil, errJson
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
