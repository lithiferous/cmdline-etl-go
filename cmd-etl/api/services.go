package api

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jszwec/csvutil"
	db "github.com/lithiferous/cmd-etl/db/sqlc"
	"github.com/lithiferous/cmd-etl/util"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func RunDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Err(err).Msgf("cannot create new migrate instance: %s", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Err(err).Msgf("failed to run migrate up: %s", err)
	}

	log.Info().Msg("db migrated successfully")
}

func UploadFilesToDB(ctx context.Context, files []string) {
	state := ctx.Value(AppState{}).(*State)

	// iterate over filenames in `files`
	// #todo: make a goroutine instead
	for _, f := range files {
		fpath := fmt.Sprintf("%s/%s", state.Config.DirSource, f)
		fi, err := os.Stat(fpath)
		if err != nil {
			log.Err(err).Msgf("%s", err)
			return
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			log.Err(err).Msg("directory is not supported by `add` method")
		case mode.IsRegular():
			file, err := os.Open(fpath)
			if err != nil {
				log.Err(err).Msgf("cannot open file - %s: %s", fpath, err)
			}

			r := csv.NewReader(file)
			dec, err := csvutil.NewDecoder(r)
			if err != nil {
				log.Err(err).Msgf("cannot decode - %s", err)
			}

			numRecords := 0
			numErrors := 0

			for {
				numRecords++
				var row SnapshotRow
				err = dec.Decode(&row)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Err(err).Msgf("error on row (%d): %s", numRecords, err)
					numErrors++
					continue
				}
				// use filename of type `snapshot_YYMMDD.csv` to record date of snapshot
				replacer := strings.NewReplacer("snapshot_", "", ".csv", "")
				snapshot_date, err := time.Parse("20060102", replacer.Replace(fi.Name()))

				if err != nil {
					log.Err(err).Msgf("failed to parse date, incorrect format should be %%Y%%m%%d for file (%s): %s", file.Name(), err)
				}

				arg := db.CreateSnapshotParams{
					UserName:    row.Name,
					StoreName:   row.StoreName,
					CreditLimit: decimal.NewFromFloat(row.CreditLimit),
					SnapshotAt: pgtype.Timestamp{
						Time:  snapshot_date,
						Valid: !snapshot_date.IsZero(),
					},
				}
				// access only store in context's state to create a snapshot in db
				state := ctx.Value(AppState{}).(*State)
				_, err = (*state).Store.CreateSnapshot(ctx, arg)

				if err != nil {
					log.Err(err).Msgf("failed to create a snapshot in db: %s", err)
				}
			}
			log.Info().Msgf("total records processed - %d, rows discarded due to errors - %d", numRecords, numErrors)

			defer file.Close()

		}
	}

}

func FindFilesToUpload(ctx context.Context) []string {
	// access context's state to query database and infer environment
	// variable to check for raw files on folder
	state := ctx.Value(AppState{}).(*State)
	src_dir := (*state).Config.DirSource

	// search source directory for latest files
	files, err := util.ReadDir(src_dir)
	if err != nil {
		log.Err(err).Msgf("missing files in raw directory: %s", err)
	}

	// keep track of dates in file names
	var file_dates []time.Time
	// record file names to upload
	var new_files []string

	for _, f := range files {
		// use filename of type `snapshot_YYMMDD.csv` to record date of snapshot
		replacer := strings.NewReplacer("snapshot_", "", ".csv", "")
		sd, err := time.Parse("20060102", replacer.Replace(f))
		if err != nil {
			log.Err(err).Msgf("failed to parse date, incorrect format should be %%Y%%m%%d for file (%s): %s", f, err)
		}
		file_dates = append(file_dates, sd)
	}

	// mapping snapshotAt -> createdAt
	// to persist which files were recorded with keys on snapshotAt
	// and what time they were loaded to database with createdAt
	db_map := make(map[time.Time]time.Time)
	empty_date := time.Unix(0, 0)
	for _, d := range file_dates {
		// unix start time as a placeholder
		db_map[d] = empty_date
	}

	// get snapshot list from db
	snapshots, err := (*state).Store.ListSnapshots(ctx)
	if err != nil {
		log.Err(err).Msgf("%s", err)
	}

	// iterate over snapshots from db and record value when it was inserted
	// i.e.: snapshot_date -> inserted_at
	for _, s := range snapshots {
		db_map[s.SnapshotAt.Time] = s.CreatedAt.Time
	}

	for sd, cd := range db_map {
		sf := fmt.Sprintf("snapshot_%d%02d%02d.csv",
			sd.Year(), sd.Month(), sd.Day())
		if cd == empty_date {
			new_files = append(new_files, sf)
			log.Info().Msgf("file `%s` not loaded to db", sf)
		} else {
			log.Info().Msgf("file `%s` was uploaded to db at `%s`", sf, cd)
		}
	}
	return new_files

}
