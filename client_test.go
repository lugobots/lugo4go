package client_test

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNewClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

}
