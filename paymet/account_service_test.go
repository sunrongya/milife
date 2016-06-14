package paymet

import (
	es "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/utiltest"
	"testing"
)

func TestAccountServicePublishOpenAccountCommand(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		as := AccountService{Service: service}
		guid := as.OpenAccount("sry", "95588388383", 50)
		return &OpenAccountCommand{WithGuid: es.WithGuid{guid}, Name: "sry", Card: "95588388383", Balance: 50}
	})
}

func TestAccountServicePublishCreditAccountCommand(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		as := AccountService{Service: service}
		guid := es.NewGuid()
		as.CreditAccount(guid, 50)
		return &CreditAccountCommand{WithGuid: es.WithGuid{guid}, Amount: 50}
	})
}

func TestAccountServicePublishDebitAccountCommand(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		as := AccountService{Service: service}
		guid := es.NewGuid()
		as.DebitAccount(guid, 50)
		return &DebitAccountCommand{WithGuid: es.WithGuid{guid}, Amount: 50}
	})
}
