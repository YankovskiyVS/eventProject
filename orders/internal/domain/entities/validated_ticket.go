package entities

type ValidatedTickets struct {
	Tickets
	isValidated bool
}

func (vt *ValidatedTickets) isValid() bool {
	return vt.isValidated
}

func NewValidatedTickets(tickets *Tickets) (*ValidatedTickets, error) {
	if err := tickets.validate(); err != nil {
		return nil, err
	}
	return &ValidatedTickets{
		Tickets:     *tickets,
		isValidated: true,
	}, nil
}
