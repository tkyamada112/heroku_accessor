package main

import (
        "os"
        "context"
        "flag"
        "fmt"
        "log"
        "strings"
        "github.com/heroku/heroku-go/v3"
)

type HerokuAPI struct {
        cli *heroku.Service
}

type Apps struct {
        name    string
}

type Collaborators struct {
        name    string
}

type AddUsers struct {
        email   string
}

func GetAPI(username string, password string) *HerokuAPI {
        heroku.DefaultTransport.Username = username
        heroku.DefaultTransport.Password = password

        h := heroku.NewService(heroku.DefaultClient)
        return &HerokuAPI {
                cli: h,
        }
}

func (h *HerokuAPI) ListApps() []Apps {
        herokuapps, err := h.cli.AppList(context.TODO(), &heroku.ListRange{Field: "name"})
        if err != nil {
		log.Fatal(err)
        }
        apps := []Apps{}
        for _, h := range herokuapps {
                apps = append(apps, Apps{
                        name: h.Name,
                })
        }
        return apps
}

func (h *HerokuAPI) ListCollaborators(app string) []Collaborators {
        cbs, err := h.cli.CollaboratorList(context.TODO(), app, &heroku.ListRange{Field: "id"})
        if err != nil {
                log.Fatal(err)
        }

        collaborators := []Collaborators{}

        for _, c := range cbs {
                collaborators = append(collaborators, Collaborators{
                        name: c.User.Email,
                })
        }
        return collaborators
}

func ReadUsers() []AddUsers {
        f, err := os.Open("Users")
        if err != nil{
                fmt.Println("error")
        }
        defer f.Close()

        lines := []AddUsers{}

        buf := make([]byte, 1024)
        for {
                n, err := f.Read(buf)
                if n == 0{
                    break
                }
                if err != nil{
                    break
                }
                for _, line := range strings.Split(string(buf[:n]), "\n") {
                        lines = append(lines, AddUsers{
                                email: string(line),
                        })
                }
        }
        return lines
}

func main() {
        var (
                userName      string
                userPass      string
                typeOperation string
                appName       string
        )

        flagSilent := true

        flag.StringVar(&userName, "username", "", "heroku api username")
        flag.StringVar(&userPass, "password", "", "heroku api password")
        flag.StringVar(&typeOperation, "type", "", "operation type 'showall', 'showuser', 'updateuser'")
        flag.StringVar(&appName, "name", "", "heroku application name")
        flag.Parse()

        h := GetAPI(userName, userPass)
        switch typeOperation{
        case "showall":
                if appName == "" {
                        apps := h.ListApps()
                        for _, a := range apps {
                                fmt.Println(a.name)
                        }
                } else {
                        fmt.Println("Invalid Arguments.")
                }
        case "showuser":
                if appName != "" {
                        users := h.ListCollaborators(appName)
                        for _, u := range users {
                                fmt.Println(u.name)
                        }
                } else {
                        fmt.Println("Set Application name. option with '-name xxxxxxx'")
                }
        case "updateuser":
                if appName != "" {
                        AddUsers := ReadUsers()

                        ExistingUsers := h.ListCollaborators(appName)

                        for  _, eu := range ExistingUsers {
                                h.cli.CollaboratorDelete(context.TODO(), appName, eu.name)
                        }
                        for _, au := range AddUsers {
                                opts := heroku.CollaboratorCreateOpts{Silent: &flagSilent, User: au.email}
                                h.cli.CollaboratorCreate(context.TODO(), appName, opts)
                                
                        }
                        fmt.Println("User Added.")
                } else {
                        fmt.Println("Set Application name. option with '-name xxxxxxx'")
                }
        default:
                fmt.Println("Invalid option.... \n\nThe first argument '-usename xxx'\nSecond argument '-password xxx'\nand last '-type showall or showuser or updateuser.'\n")
        }
}