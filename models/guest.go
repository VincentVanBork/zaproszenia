package models

import "cloud.google.com/go/datastore"

type Guest struct {
	K       *datastore.Key `datastore:"__key__"`
	Name    string
	Surname string
}
