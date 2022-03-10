package controllers

import (
	"cloud.google.com/go/datastore"
)

type Controller struct {
	Objects *datastore.Client
}
