#!/bin/sh

psql2pgcopy() {
	echo "
		COPY (

			SELECT
				'fuji'::TEXT AS name,
				3776::INTEGER AS height_in_meter,
				3.776::DOUBLE PRECISION AS height_in_km,
				42::BIGINT AS id,
				false::BOOL AS active,
				NULL::BYTEA AS data,
				CLOCK_TIMESTAMP()::TIMESTAMP WITH TIME ZONE AS updated,
				'cafef00d-dead-beaf-face-864299792458'::UUID AS related

			UNION ALL

			SELECT
				'takao'::TEXT AS name,
				599::INTEGER AS height_in_meter,
				0.599::DOUBLE PRECISION AS height_in_km,
				43::BIGINT AS id,
				true::BOOL AS active,
				NULL::BYTEA AS data,
				NULL::TIMESTAMP WITH TIME ZONE AS updated,
				NULL::UUID AS related

		)
		TO STDOUT
		WITH (
			FORMAT BINARY
		)
	" |
		env PGUSER=postgres env LC_ALL=C psql \
			>./sample.d/input.pgcopy
}

#psql2pgcopy

export ENV_COL_INFO_JSONL_NAME=./sample.d/types.jsonl

cat sample.d/input.pgcopy |
	./pgcopy2sqlvalues
