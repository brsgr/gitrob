package core

import (
  //"context"
  "github.com/xanzy/go-gitlab"
)

type GitlabOwner struct {
  Login     string
  ID        int
  Type      string
  Name      string
  AvatarURL string
  URL       string
  Company   string
  Blog      string
  Location  string
  Email     string
  Bio       string
}

type GitlabRepository struct {
  Owner         string
  ID            int
  Name          string
  FullName      string
  CloneURL      string
  URL           string
  DefaultBranch string
  Description   string
  Homepage      string
}

func GetUserOrOrganization(login int, client *gitlab.Client) (*GitlabOwner, error) {
  //ctx := context.Background()
	//tmp.Println(ctx)
  user, _, err := client.Users.GetUser(login)
  if err != nil {
    return nil, err
  }
  return &GitlabOwner{
    Login:     user.Username,
    ID:        user.ID,
    Name:      user.Name,
    AvatarURL: user.AvatarURL,
    Location:  user.Location,
    Email:     user.Email,
    Bio:       user.Bio,
  }, nil
}

func GetRepositoriesFromOwner(login *string, client *gitlab.Client) ([]*GitlabRepository, error) {
  var allRepos []*GitlabRepository
  //loginVal := *login
	//tmp.Println(loginVal)
  //ctx := context.Background()
	//tmp.Println(ctx)
  opt := &gitlab.ListProjectsOptions{
    //Simple: true,
  }
	forkParent:=gitlab.ForkParent{}

  for {
    repos, resp, err := client.Projects.ListProjects(opt)
    if err != nil {
      return allRepos, err
    }
    for _, repo := range repos {

      if *repo.ForkedFromProject == forkParent {
        r := GitlabRepository{
          Owner:         repo.Owner.Username,
          ID:            repo.ID,
          Name:          repo.Name,
          CloneURL:      repo.HTTPURLToRepo,
          URL:           repo.WebURL,
          DefaultBranch: repo.DefaultBranch,
          Description:   repo.Description,
        }
        allRepos = append(allRepos, &r)
      }
    }
    if resp.NextPage == 0 {
      break
    }
    opt.Page = resp.NextPage
  }

  return allRepos, nil
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
