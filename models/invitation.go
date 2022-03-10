package models

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/lithammer/shortuuid/v4"
)

type Invitation struct {
	K     *datastore.Key `datastore:"__key__"`
	Token string

	IsWedding    bool
	IsReception  bool
	HasCompanion bool

	Email string

	Hotel     bool
	Transport bool
}

type Requirements struct {
	K            *datastore.Key `datastore:"__key__"`
	Hotel        bool
	Transport    bool
	InvitationId int
}

type InvitData struct {
	IsWedding    bool
	IsReception  bool
	HasCompanion bool
}

func CreateInvitation(
	ctx context.Context,
	client *datastore.Client,
	isWedding bool,
	isReception bool,
	hasCompanion bool,
) (*datastore.Key, error) {

	invit := &Invitation{
		Token: shortuuid.New(),

		IsWedding:    isWedding,
		IsReception:  isReception,
		HasCompanion: hasCompanion,

		Email:     "",
		Hotel:     false,
		Transport: false,
	}

	k := datastore.IncompleteKey("Invitation", nil)

	key, err := client.Put(ctx, k, invit)
	return key, err
}

func GetAllInvitations(ctx context.Context, client *datastore.Client) ([]*Invitation, error) {
	// Print out previous visits.
	q := datastore.NewQuery("Invitation")

	invitations := make([]*Invitation, 0)
	_, err := client.GetAll(ctx, q, &invitations)
	return invitations, err
}

func GetInvitation(ctx context.Context, client *datastore.Client, key *datastore.Key) (*Invitation, error) {
	invit := &Invitation{}
	err := client.Get(ctx, key, invit)
	return invit, err
}
