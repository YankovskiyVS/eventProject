package entities

type ValidatedTickets struct {
	Tickets
	isValidated bool
}

func (vt *ValidatedTickets) isValid() bool {
	return vt.isValidated
}

func NewValidateTickets(tickets *Tickets) (*ValidatedTickets, error) {
	if err := tickets.validate(); err != nil {
		return nil, err
	}
	return &ValidatedTickets{
		Tickets:     *tickets,
		isValidated: true,
	}, nil
}
