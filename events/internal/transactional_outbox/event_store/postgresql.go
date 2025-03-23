package eventstore

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	transactionaloutbox "github.com/YankovskiyVS/eventProject/events/internal/transactional_outbox"
)

type Store struct {
	db *sql.DB
}

func NewStore() (*Store, error) {
	//Getting all required info from docker compose file
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST_2"), os.Getenv("PGPORT_2"), os.Getenv("PGUSER_2"),
		os.Getenv("PGPASSWORD_2"), os.Getenv("PGDATABASE_2"))
	db, err := sql.Open("postgres", connStr)
	if err != nil || db.Ping() != nil {
		log.Fatalf("failed to connect to database %v", err)
		return nil, err
	}
	return &Store{
		db: db,
	}, nil
}

// ClearLocksWithDurationBeforeDate clears all records which are too old
func (s Store) ClearLocksWithDurationBeforeDate(time time.Time) error {
	_, err := s.db.Exec(`
	UPDATE outbox 
		SET
			locked_by=NULL,
			locked_on=NULL
		WHERE locked_on < $1
		`,
		time,
	)
	if err != nil {
		log.Fatalf("failed to clear the old records %s", err)
		return err
	}
	return nil
}

// UpdateRecordLockByState updated the lock information based on the state
func (s Store) UpdateRecordLockByState(lockID string, lockedOn time.Time, state transactionaloutbox.RecordState) error {
	_, err := s.db.Exec(`
	UPDATE outbox 
		SET 
			locked_by=$1,
			locked_on=$2
		WHERE state = $3
		`,
		lockID,
		lockedOn,
		state,
	)
	if err != nil {
		log.Fatalf("failed to update the lock info %s", err)
		return err
	}
	return nil
}

// UpdateRecordByID updates the provided record based on its id
func (s Store) UpdateRecordByID(rec transactionaloutbox.Record) error {
	msgData := new(bytes.Buffer)
	enc := gob.NewEncoder(msgData)
	encErr := enc.Encode(rec.Message)
	if encErr != nil {
		return encErr
	}

	_, err := s.db.Exec(`
	UPDATE outbox 
		SET 
			data=$1,
			state=$2,
			created_on=$3,
			locked_by=$4,
			locked_on=$5,
			processed_on=$6,
		    number_of_attempts=$7,
		    last_attempted_on=$8,
		    error=$9
		WHERE id = $10
		`,
		msgData.Bytes(),
		rec.State,
		rec.CreatedOn,
		rec.LockID,
		rec.LockedOn,
		rec.ProcessedOn,
		rec.NumberOfAttempts,
		rec.LastAttemptOn,
		rec.Error,
		rec.ID,
	)
	if err != nil {
		log.Fatalf("failed to update the record %s", err)
		return err
	}
	return nil
}

// ClearLocksByLockID clears lock information of the records with the provided id
func (s Store) ClearLocksByLockID(lockID string) error {
	_, err := s.db.Exec(`
	UPDATE outbox 
		SET 
			locked_by=NULL,
			locked_on=NULL
		WHERE id = $1
		`,
		lockID,
	)
	if err != nil {
		log.Fatalf("clear lock info %s", err)
		return err
	}
	return nil
}

// GetRecordsByLockID returns the records of the provided id
func (s Store) GetRecordsByLockID(lockID string) ([]transactionaloutbox.Record, error) {
	rows, err := s.db.Query(`
	SELECT id, data, state, created_on,locked_by,locked_on,processed_on,number_of_attempts,last_attempted_on,error 
	FROM outbox 
	WHERE locked_by = $1`,
		lockID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var messages []transactionaloutbox.Record

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var rec transactionaloutbox.Record
		var data []byte
		scanErr := rows.Scan(
			&rec.ID,
			&data,
			&rec.State,
			&rec.CreatedOn,
			&rec.LockID,
			&rec.LockedOn,
			&rec.ProcessedOn,
			&rec.NumberOfAttempts,
			&rec.LastAttemptOn,
			&rec.Error,
		)
		if scanErr != nil {
			if scanErr == sql.ErrNoRows {
				return messages, nil
			}
			return messages, err
		}
		decErr := gob.NewDecoder(bytes.NewReader(data)).Decode(&rec.Message)
		if decErr != nil {
			return nil, decErr
		}

		messages = append(messages, rec)
	}
	if err = rows.Err(); err != nil {
		return messages, err
	}
	return messages, nil
}

// AddRecordTx stores the record in the db within the provided transaction tx
func (s Store) AddRecordTx(rec transactionaloutbox.Record, tx *sql.Tx) error {
	msgBuf := new(bytes.Buffer)
	msgEnc := gob.NewEncoder(msgBuf)
	encErr := msgEnc.Encode(rec.Message)

	if encErr != nil {
		return encErr
	}
	query := `
	INSERT INTO outbox 
	(id, 
	data, 
	state, 
	created_on,
	locked_by,
	locked_on,
	processed_on,
	number_of_attempts,
	last_attempted_on,
	error) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err := tx.Exec(query,
		rec.ID,
		msgBuf.Bytes(),
		rec.State,
		rec.CreatedOn,
		rec.LockID,
		rec.LockedOn,
		rec.ProcessedOn,
		rec.NumberOfAttempts,
		rec.LastAttemptOn,
		rec.Error,
	)
	if err != nil {
		log.Fatalf("failed to add record %s", err)
		return err
	}
	return nil
}

// RemoveRecordsBeforeDatetime removes records before the provided datetime
func (s Store) RemoveRecordsBeforeDatetime(expiryTime time.Time) error {
	_, err := s.db.Exec(`
	DELETE FROM outbox 
		WHERE created_on < $1
		`,
		expiryTime,
	)
	if err != nil {
		log.Fatalf("failed to delete old records %s", err)
		return err
	}
	return nil
}
