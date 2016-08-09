package paymet

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
)

func TestAccountRestore(t *testing.T) {
	acc := &Account{}
	acc.HandleAccountOpenedEvent(&AccountOpenedEvent{Name: "sry", Card: "955884334444", Balance: 100})
	acc.HandleAccountCreditedEvent(&AccountCreditedEvent{Amount: 100})
	acc.HandleAccountDebitedEvent(&AccountDebitedEvent{Amount: 50})
	acc.HandleAccountDebitFailedEvent(&AccountDebitFailedEvent{})

	assert.Equal(t, "sry", acc._name)
	assert.Equal(t, BankCard("955884334444"), acc._card)
	assert.Equal(t, Money(150), acc._balance)
}

func TestOpenAccountCommand(t *testing.T) {
	command := &OpenAccountCommand{Name: "sry", Card: "955884334444", Balance: 100}
	events := []es.Event{&AccountOpenedEvent{Name: "sry", Card: "955884334444", Balance: 100}}

	assert.Equal(t, events, new(Account).ProcessOpenAccountCommand(command), "")
}

func TestCreditAccountCommand(t *testing.T) {
	command := &CreditAccountCommand{Amount: 100}
	events := []es.Event{&AccountCreditedEvent{Amount: 100}}

	assert.Equal(t, events, new(Account).ProcessCreditAccountCommand(command), "")
}

func TestDebitAccountCommand(t *testing.T) {
	tests := []struct {
		account *Account
		command *DebitAccountCommand
		events  []es.Event
	}{
		{
			&Account{_balance: 50},
			&DebitAccountCommand{Amount: 50},
			[]es.Event{&AccountDebitedEvent{Amount: 50}},
		},
		{
			&Account{_balance: 150},
			&DebitAccountCommand{Amount: 50},
			[]es.Event{&AccountDebitedEvent{Amount: 50}},
		},
		{
			&Account{_balance: 100},
			&DebitAccountCommand{Amount: 101},
			[]es.Event{&AccountDebitFailedEvent{}},
		},
	}
	for _, v := range tests {
		assert.Equal(t, v.events, v.account.ProcessDebitAccountCommand(v.command), "")
	}
}

func TestDebitAccountBecauseOfTransferCommand(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	command := &DebitAccountBecauseOfTransferCommand{mTDetails: details}
	events := []es.Event{&AccountDebitedBecauseOfTransferEvent{mTDetails: details}}
	account := &Account{_balance: 45}

	assert.Equal(t, events, account.ProcessDebitAccountBecauseOfTransferCommand(command), "")
}

func TestDebitAccountBecauseOfTransferCommand2Failed(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      50,
		Transaction: es.NewGuid(),
	}
	command := &DebitAccountBecauseOfTransferCommand{mTDetails: details}
	events := []es.Event{&AccountDebitedBecauseOfTransferFailedEvent{mTDetails: details}}
	account := &Account{_balance: 45}

	assert.Equal(t, events, account.ProcessDebitAccountBecauseOfTransferCommand(command), "")
}

func TestCreditAccountBecauseOfTransferCommand(t *testing.T) {
	details := mTDetails{
		From:        es.NewGuid(),
		To:          es.NewGuid(),
		Amount:      40,
		Transaction: es.NewGuid(),
	}
	command := &CreditAccountBecauseOfTransferCommand{mTDetails: details}
	events := []es.Event{&AccountCreditedBecauseOfTransferEvent{mTDetails: details}}
	account := &Account{_balance: 100}

	assert.Equal(t, events, account.ProcessCreditAccountBecauseOfTransferCommand(command), "")
}
