// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: track.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const associateArtistToTrack = `-- name: AssociateArtistToTrack :exec
INSERT INTO artist_tracks (artist_id, track_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`

type AssociateArtistToTrackParams struct {
	ArtistID int32
	TrackID  int32
}

func (q *Queries) AssociateArtistToTrack(ctx context.Context, arg AssociateArtistToTrackParams) error {
	_, err := q.db.Exec(ctx, associateArtistToTrack, arg.ArtistID, arg.TrackID)
	return err
}

const countTopTracks = `-- name: CountTopTracks :one
SELECT COUNT(DISTINCT l.track_id) AS total_count
FROM listens l
WHERE l.listened_at BETWEEN $1 AND $2
`

type CountTopTracksParams struct {
	ListenedAt   time.Time
	ListenedAt_2 time.Time
}

func (q *Queries) CountTopTracks(ctx context.Context, arg CountTopTracksParams) (int64, error) {
	row := q.db.QueryRow(ctx, countTopTracks, arg.ListenedAt, arg.ListenedAt_2)
	var total_count int64
	err := row.Scan(&total_count)
	return total_count, err
}

const countTopTracksByArtist = `-- name: CountTopTracksByArtist :one
SELECT COUNT(DISTINCT l.track_id) AS total_count
FROM listens l
JOIN artist_tracks at ON l.track_id = at.track_id
WHERE l.listened_at BETWEEN $1 AND $2
AND at.artist_id = $3
`

type CountTopTracksByArtistParams struct {
	ListenedAt   time.Time
	ListenedAt_2 time.Time
	ArtistID     int32
}

func (q *Queries) CountTopTracksByArtist(ctx context.Context, arg CountTopTracksByArtistParams) (int64, error) {
	row := q.db.QueryRow(ctx, countTopTracksByArtist, arg.ListenedAt, arg.ListenedAt_2, arg.ArtistID)
	var total_count int64
	err := row.Scan(&total_count)
	return total_count, err
}

const countTopTracksByRelease = `-- name: CountTopTracksByRelease :one
SELECT COUNT(DISTINCT l.track_id) AS total_count
FROM listens l
JOIN tracks t ON l.track_id = t.id
WHERE l.listened_at BETWEEN $1 AND $2
AND t.release_id = $3
`

type CountTopTracksByReleaseParams struct {
	ListenedAt   time.Time
	ListenedAt_2 time.Time
	ReleaseID    int32
}

func (q *Queries) CountTopTracksByRelease(ctx context.Context, arg CountTopTracksByReleaseParams) (int64, error) {
	row := q.db.QueryRow(ctx, countTopTracksByRelease, arg.ListenedAt, arg.ListenedAt_2, arg.ReleaseID)
	var total_count int64
	err := row.Scan(&total_count)
	return total_count, err
}

const deleteTrack = `-- name: DeleteTrack :exec
DELETE FROM tracks WHERE id = $1
`

func (q *Queries) DeleteTrack(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteTrack, id)
	return err
}

const getAllTracksFromArtist = `-- name: GetAllTracksFromArtist :many
SELECT t.id, t.musicbrainz_id, t.duration, t.release_id, t.title
FROM tracks_with_title t
JOIN artist_tracks at ON t.id = at.track_id
WHERE at.artist_id = $1
`

func (q *Queries) GetAllTracksFromArtist(ctx context.Context, artistID int32) ([]TracksWithTitle, error) {
	rows, err := q.db.Query(ctx, getAllTracksFromArtist, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TracksWithTitle
	for rows.Next() {
		var i TracksWithTitle
		if err := rows.Scan(
			&i.ID,
			&i.MusicBrainzID,
			&i.Duration,
			&i.ReleaseID,
			&i.Title,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopTracksByArtistPaginated = `-- name: GetTopTracksByArtistPaginated :many
SELECT
    t.id,
    t.title,
    t.musicbrainz_id,
    t.release_id,
    r.image,
    COUNT(*) AS listen_count,
    get_artists_for_track(t.id) AS artists
FROM listens l
JOIN tracks_with_title t ON l.track_id = t.id
JOIN releases r ON t.release_id = r.id
JOIN artist_tracks at ON at.track_id = t.id
WHERE l.listened_at BETWEEN $1 AND $2
  AND at.artist_id = $5
GROUP BY t.id, t.title, t.musicbrainz_id, t.release_id, r.image
ORDER BY listen_count DESC, t.id
LIMIT $3 OFFSET $4
`

type GetTopTracksByArtistPaginatedParams struct {
	ListenedAt   time.Time
	ListenedAt_2 time.Time
	Limit        int32
	Offset       int32
	ArtistID     int32
}

type GetTopTracksByArtistPaginatedRow struct {
	ID            int32
	Title         string
	MusicBrainzID *uuid.UUID
	ReleaseID     int32
	Image         *uuid.UUID
	ListenCount   int64
	Artists       []byte
}

func (q *Queries) GetTopTracksByArtistPaginated(ctx context.Context, arg GetTopTracksByArtistPaginatedParams) ([]GetTopTracksByArtistPaginatedRow, error) {
	rows, err := q.db.Query(ctx, getTopTracksByArtistPaginated,
		arg.ListenedAt,
		arg.ListenedAt_2,
		arg.Limit,
		arg.Offset,
		arg.ArtistID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopTracksByArtistPaginatedRow
	for rows.Next() {
		var i GetTopTracksByArtistPaginatedRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.MusicBrainzID,
			&i.ReleaseID,
			&i.Image,
			&i.ListenCount,
			&i.Artists,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopTracksInReleasePaginated = `-- name: GetTopTracksInReleasePaginated :many
SELECT
    t.id,
    t.title,
    t.musicbrainz_id,
    t.release_id,
    r.image,
    COUNT(*) AS listen_count,
    get_artists_for_track(t.id) AS artists
FROM listens l
JOIN tracks_with_title t ON l.track_id = t.id
JOIN releases r ON t.release_id = r.id
WHERE l.listened_at BETWEEN $1 AND $2
  AND t.release_id = $5
GROUP BY t.id, t.title, t.musicbrainz_id, t.release_id, r.image
ORDER BY listen_count DESC, t.id
LIMIT $3 OFFSET $4
`

type GetTopTracksInReleasePaginatedParams struct {
	ListenedAt   time.Time
	ListenedAt_2 time.Time
	Limit        int32
	Offset       int32
	ReleaseID    int32
}

type GetTopTracksInReleasePaginatedRow struct {
	ID            int32
	Title         string
	MusicBrainzID *uuid.UUID
	ReleaseID     int32
	Image         *uuid.UUID
	ListenCount   int64
	Artists       []byte
}

func (q *Queries) GetTopTracksInReleasePaginated(ctx context.Context, arg GetTopTracksInReleasePaginatedParams) ([]GetTopTracksInReleasePaginatedRow, error) {
	rows, err := q.db.Query(ctx, getTopTracksInReleasePaginated,
		arg.ListenedAt,
		arg.ListenedAt_2,
		arg.Limit,
		arg.Offset,
		arg.ReleaseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopTracksInReleasePaginatedRow
	for rows.Next() {
		var i GetTopTracksInReleasePaginatedRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.MusicBrainzID,
			&i.ReleaseID,
			&i.Image,
			&i.ListenCount,
			&i.Artists,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopTracksPaginated = `-- name: GetTopTracksPaginated :many
SELECT
    t.id,
    t.title,
    t.musicbrainz_id,
    t.release_id,
    r.image,
    COUNT(*) AS listen_count,
    get_artists_for_track(t.id) AS artists
FROM listens l
JOIN tracks_with_title t ON l.track_id = t.id
JOIN releases r ON t.release_id = r.id
WHERE l.listened_at BETWEEN $1 AND $2
GROUP BY t.id, t.title, t.musicbrainz_id, t.release_id, r.image
ORDER BY listen_count DESC, t.id
LIMIT $3 OFFSET $4
`

type GetTopTracksPaginatedParams struct {
	ListenedAt   time.Time
	ListenedAt_2 time.Time
	Limit        int32
	Offset       int32
}

type GetTopTracksPaginatedRow struct {
	ID            int32
	Title         string
	MusicBrainzID *uuid.UUID
	ReleaseID     int32
	Image         *uuid.UUID
	ListenCount   int64
	Artists       []byte
}

func (q *Queries) GetTopTracksPaginated(ctx context.Context, arg GetTopTracksPaginatedParams) ([]GetTopTracksPaginatedRow, error) {
	rows, err := q.db.Query(ctx, getTopTracksPaginated,
		arg.ListenedAt,
		arg.ListenedAt_2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopTracksPaginatedRow
	for rows.Next() {
		var i GetTopTracksPaginatedRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.MusicBrainzID,
			&i.ReleaseID,
			&i.Image,
			&i.ListenCount,
			&i.Artists,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTrack = `-- name: GetTrack :one
SELECT 
  t.id, t.musicbrainz_id, t.duration, t.release_id, t.title,
  get_artists_for_track(t.id) AS artists,
  r.image
FROM tracks_with_title t
JOIN releases r ON t.release_id = r.id
WHERE t.id = $1 LIMIT 1
`

type GetTrackRow struct {
	ID            int32
	MusicBrainzID *uuid.UUID
	Duration      int32
	ReleaseID     int32
	Title         string
	Artists       []byte
	Image         *uuid.UUID
}

func (q *Queries) GetTrack(ctx context.Context, id int32) (GetTrackRow, error) {
	row := q.db.QueryRow(ctx, getTrack, id)
	var i GetTrackRow
	err := row.Scan(
		&i.ID,
		&i.MusicBrainzID,
		&i.Duration,
		&i.ReleaseID,
		&i.Title,
		&i.Artists,
		&i.Image,
	)
	return i, err
}

const getTrackByMbzID = `-- name: GetTrackByMbzID :one
SELECT id, musicbrainz_id, duration, release_id, title FROM tracks_with_title
WHERE musicbrainz_id = $1 LIMIT 1
`

func (q *Queries) GetTrackByMbzID(ctx context.Context, musicbrainzID *uuid.UUID) (TracksWithTitle, error) {
	row := q.db.QueryRow(ctx, getTrackByMbzID, musicbrainzID)
	var i TracksWithTitle
	err := row.Scan(
		&i.ID,
		&i.MusicBrainzID,
		&i.Duration,
		&i.ReleaseID,
		&i.Title,
	)
	return i, err
}

const getTrackByTitleAndArtists = `-- name: GetTrackByTitleAndArtists :one
SELECT t.id, t.musicbrainz_id, t.duration, t.release_id, t.title
FROM tracks_with_title t
JOIN artist_tracks at ON at.track_id = t.id
WHERE t.title = $1
  AND at.artist_id = ANY($2::int[])
GROUP BY t.id, t.title, t.musicbrainz_id, t.duration, t.release_id
HAVING COUNT(DISTINCT at.artist_id) = cardinality($2::int[])
`

type GetTrackByTitleAndArtistsParams struct {
	Title   string
	Column2 []int32
}

func (q *Queries) GetTrackByTitleAndArtists(ctx context.Context, arg GetTrackByTitleAndArtistsParams) (TracksWithTitle, error) {
	row := q.db.QueryRow(ctx, getTrackByTitleAndArtists, arg.Title, arg.Column2)
	var i TracksWithTitle
	err := row.Scan(
		&i.ID,
		&i.MusicBrainzID,
		&i.Duration,
		&i.ReleaseID,
		&i.Title,
	)
	return i, err
}

const insertTrack = `-- name: InsertTrack :one
INSERT INTO tracks (musicbrainz_id, release_id, duration)
VALUES ($1, $2, $3)
RETURNING id, musicbrainz_id, duration, release_id
`

type InsertTrackParams struct {
	MusicBrainzID *uuid.UUID
	ReleaseID     int32
	Duration      int32
}

func (q *Queries) InsertTrack(ctx context.Context, arg InsertTrackParams) (Track, error) {
	row := q.db.QueryRow(ctx, insertTrack, arg.MusicBrainzID, arg.ReleaseID, arg.Duration)
	var i Track
	err := row.Scan(
		&i.ID,
		&i.MusicBrainzID,
		&i.Duration,
		&i.ReleaseID,
	)
	return i, err
}

const updateReleaseForAll = `-- name: UpdateReleaseForAll :exec
UPDATE tracks SET release_id = $2
WHERE release_id = $1
`

type UpdateReleaseForAllParams struct {
	ReleaseID   int32
	ReleaseID_2 int32
}

func (q *Queries) UpdateReleaseForAll(ctx context.Context, arg UpdateReleaseForAllParams) error {
	_, err := q.db.Exec(ctx, updateReleaseForAll, arg.ReleaseID, arg.ReleaseID_2)
	return err
}

const updateTrackDuration = `-- name: UpdateTrackDuration :exec
UPDATE tracks SET duration = $2
WHERE id = $1
`

type UpdateTrackDurationParams struct {
	ID       int32
	Duration int32
}

func (q *Queries) UpdateTrackDuration(ctx context.Context, arg UpdateTrackDurationParams) error {
	_, err := q.db.Exec(ctx, updateTrackDuration, arg.ID, arg.Duration)
	return err
}

const updateTrackMbzID = `-- name: UpdateTrackMbzID :exec
UPDATE tracks SET musicbrainz_id = $2
WHERE id = $1
`

type UpdateTrackMbzIDParams struct {
	ID            int32
	MusicBrainzID *uuid.UUID
}

func (q *Queries) UpdateTrackMbzID(ctx context.Context, arg UpdateTrackMbzIDParams) error {
	_, err := q.db.Exec(ctx, updateTrackMbzID, arg.ID, arg.MusicBrainzID)
	return err
}

const updateTrackPrimaryArtist = `-- name: UpdateTrackPrimaryArtist :exec
UPDATE artist_tracks SET is_primary = $3
WHERE artist_id = $1 AND track_id = $2
`

type UpdateTrackPrimaryArtistParams struct {
	ArtistID  int32
	TrackID   int32
	IsPrimary bool
}

func (q *Queries) UpdateTrackPrimaryArtist(ctx context.Context, arg UpdateTrackPrimaryArtistParams) error {
	_, err := q.db.Exec(ctx, updateTrackPrimaryArtist, arg.ArtistID, arg.TrackID, arg.IsPrimary)
	return err
}
