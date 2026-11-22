package main

import "github.com/corentings/chess/v2"

// The geometric subset only selects candidates. Every returned certificate
// makes a real legal checking move and verifies that no legal reply exists.
func contactMateProof(p *chess.Position, child *chess.SearchPosition) (chess.Move, bool) {
	if !p.SearchHasProtectedQueenContact(p.Turn()) {
		return chess.Move{}, false
	}
	for _, m := range p.AppendSearchProtectedQueenContacts(nil) {
		if child.Update(p, &m).SearchStatus() == chess.Checkmate {
			return m, true
		}
	}
	return chess.Move{}, false
}

// A hypothetical pass diagnoses a threat only; it never evaluates a forced loss.
func passedContactThreat(p *chess.Position, pass, child *chess.SearchPosition) (chess.Move, bool) {
	if !p.SearchHasProtectedQueenContact(p.Turn().Other()) {
		return chess.Move{}, false
	}
	return contactMateProof(pass.Update(p, nil), child)
}
