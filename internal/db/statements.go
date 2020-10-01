package db

var statements = map[string]string{
	// inserts a new row into the callhandlings table
	"create-call": `
  INSERT INTO calls (call_id, sid, conversation_id, ANI, DNIS, status)
    values(?, ?, ?, ?, ?, ?)
  `,
	// soft deletes a callhandling by id
	// "delete-call": `
	// UPDATE
	//   calls
	// SET
	//   deleted_at = NOW()
	// WHERE
	//   callhandling_id = UUID_TO_BIN(?)
	//   AND deleted_at IS NULL
	// `,
	// gets a single callhandling row by id
	"get-call": `
  SELECT
    *
  FROM
    calls
  WHERE
    call_id = ?
  `,
	// update a single callhandling row by ID
	// "update-callhandling": `
	// UPDATE
	//   callhandlings
	// SET
	//   name = ?
	// WHERE
	//   callhandling_id = UUID_TO_BIN(?)
	//   AND deleted_at IS NULL
	// `,
}
