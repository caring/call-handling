package db

var statements = map[string]string{
  // inserts a new row into the callhandlings table
  "create-callhandling": `
  INSERT INTO callhandlings (callhandling_id, name)
    values(UUID_TO_BIN(?), ?)
  `,
  // soft deletes a callhandling by id
  "delete-callhandling": `
  UPDATE
    callhandlings
  SET
    deleted_at = NOW()
  WHERE
    callhandling_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
  // gets a single callhandling row by id
  "get-callhandling": `
  SELECT
    callhandling_id, name
  FROM
    callhandlings
  WHERE
    callhandling_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
  // update a single callhandling row by ID
  "update-callhandling": `
  UPDATE
    callhandlings
  SET
    name = ?
  WHERE
    callhandling_id = UUID_TO_BIN(?)
    AND deleted_at IS NULL
  `,
}
