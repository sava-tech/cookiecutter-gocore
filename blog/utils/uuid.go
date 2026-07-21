package utils

import (
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgtype"
)

func ConvertToPgUUID(id uuid.UUID) pgtype.UUID {
    var pgUUID pgtype.UUID
	copy(pgUUID.Bytes[:], id[:]) // copy bytes into fixed array
    pgUUID.Valid = true
    return pgUUID
}

func ConvertToPgUUIDFromString(id string) (pgtype.UUID, error) {
    parsedUUID, err := uuid.Parse(id)
    if err != nil {
        return pgtype.UUID{}, err
    }
    return ConvertToPgUUID(parsedUUID), nil
}   

func ConvertFromPgUUID(pgUUID pgtype.UUID) (uuid.UUID, error) {
    if !pgUUID.Valid {
        return uuid.Nil, nil // or return an error if you prefer
    }
    var id uuid.UUID
    copy(id[:], pgUUID.Bytes[:]) // copy bytes from fixed array to slice
    return id, nil
}