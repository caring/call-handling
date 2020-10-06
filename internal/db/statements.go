package db

var statements = map[string]string{
	// inserts a new row into the calls table
	"create-call": `
  INSERT INTO calls (call_id, sid, conversation_id, ANI, DNIS, status)
    values(?, ?, ?, ?, ?, ?)
  `,
	// gets a single call row by id
	"get-call": `
  SELECT
    *
  FROM
    calls
  WHERE
    call_id = ?
  `,
	// inserts a new row into the events table
	"create-event": `
  INSERT INTO events (call_id, type, identity_id, timestamp, meta)
	values(?, ?, ?, ?, ?)
	`,
	// gets a single event row by id
	"get_event": `
  SELECT
    *
  FROM
    events
  WHERE
    call_id = ?
	`,
}
