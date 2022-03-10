package models

import (
	"cloud.google.com/go/datastore"
	"context"
)

type Guest struct {
	K          *datastore.Key `datastore:"__key__"`
	Name       string
	Surname    string
	Invitation *datastore.Key
}

type GuestData struct {
	Name       string
	Surname    string
	Invitation string
}

func AssignGuest(ctx context.Context, client *datastore.Client, data GuestData) (*datastore.Key, error) {
	UpdatedInvitation := Invitation{}
	InvitationKey, keyerr := datastore.DecodeKey(data.Invitation)
	if keyerr != nil {
		return nil, keyerr
	}
	getErr := client.Get(ctx, InvitationKey, &UpdatedInvitation)

	if getErr != nil {
		return nil, getErr
	}

	newGuest := &Guest{
		Name:       data.Name,
		Surname:    data.Surname,
		Invitation: InvitationKey,
	}
	k := datastore.IncompleteKey("Guest", nil)
	NewGuestKey, err := client.Put(ctx, k, newGuest)
	return NewGuestKey, err
}

func GetGuest(ctx context.Context, client *datastore.Client, key *datastore.Key) (*Guest, error) {
	guest := &Guest{}
	err := client.Get(ctx, key, guest)
	return guest, err
}

func GetAllGuests(ctx context.Context, client *datastore.Client) ([]*Guest, error) {
	// Print out previous visits.
	q := datastore.NewQuery("Guest")

	guests := make([]*Guest, 0)
	_, err := client.GetAll(ctx, q, &guests)
	return guests, err
}
