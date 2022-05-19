package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pr struct {
	ID         primitive.ObjectID `bson:"_id"`
	Owner      string             `bson:"owner"`
	Repo       string             `bson:"repo"`
	Number     int                `bson:"number"`
	Branch     string             `bson:"branch"`
	State      string             `bson:"state"`
	Image      string             `bson:"image"`
	LastCommit string             `bson:"last_commit"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func (p *Pr) DomainName() string {
	return "sky-commands.coflnet.com"
}

func (p *Pr) DomainPath() string {
	return fmt.Sprintf("/pr/%s/%s/%d", strings.ToLower(p.Owner), strings.ToLower(p.Repo), p.Number)
}

func (p *Pr) CompleteDomain() string {
	return fmt.Sprintf("%s%s", p.DomainName(), p.DomainPath())
}

func (p *Pr) FullImageWithTag() string {
	return fmt.Sprintf("harbor.flou.dev/%s:%s", p.Image, p.Tag())
}

func (p *Pr) Tag() string {
	return fmt.Sprintf("%s-%s-%s-%s", strings.ToLower(p.Owner), strings.ToLower(p.Repo), strings.ToLower(p.Branch), strings.ToLower(p.LastCommit))
}

func (p *Pr) GitUrl() string {
	return fmt.Sprintf("git://github.com/%s/%s.git", p.Owner, p.Repo)
}

func (p *Pr) KubernetesResourceName() string {

	b := strings.ToLower(p.Branch)
	curLength := len("pr-env-") + len(strings.ToLower(p.Repo)) + len(strconv.Itoa(p.Number))
	maxLen := 63 - curLength - 3

	if len(b) > maxLen {
		b = b[:maxLen]
	}

	return fmt.Sprintf("pr-env-%s-%s-%d", strings.ToLower(p.Repo), b, p.Number)
}
