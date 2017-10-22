#!/bin/bash

MYDB=my://xodb:xodb@localhost/xodb

DEST=$1

if [ -z "$DEST" ]; then
  DEST=models
fi

EXTRA=$2

XOBIN=$(which xo)
if [ -e ./xo ]; then
  XOBIN=./xo
fi

set -ex

rm -rf $DEST/*.xo.go

# mysql enum list query
$XOBIN $MYDB -a -N -M -B -T Enum -F MyEnums -o $DEST $EXTRA << ENDSQL
SELECT
  DISTINCT column_name AS enum_name
FROM information_schema.columns
WHERE data_type = 'enum' AND table_schema = %%schema string%%
ENDSQL

# mysql enum value list query
$XOBIN $MYDB -N -M -B -1 -T MyEnumValue -F MyEnumValues -o $DEST $EXTRA << ENDSQL
SELECT
  SUBSTRING(column_type, 6, CHAR_LENGTH(column_type) - 6) AS enum_values
FROM information_schema.columns
WHERE data_type = 'enum' AND table_schema = %%schema string%% AND column_name = %%enum string%%
ENDSQL

# mysql autoincrement list query
$XOBIN $MYDB -N -M -B -T MyAutoIncrement -F MyAutoIncrements -o $DEST $EXTRA << ENDSQL
SELECT
  table_name
FROM information_schema.tables
WHERE auto_increment IS NOT null AND table_schema = %%schema string%%
ENDSQL

# mysql proc list query
$XOBIN $MYDB -a -N -M -B -T Proc -F MyProcs -o $DEST $EXTRA << ENDSQL
SELECT
  r.routine_name AS proc_name,
  p.dtd_identifier AS return_type
FROM information_schema.routines r
INNER JOIN information_schema.parameters p
  ON p.specific_schema = r.routine_schema AND p.specific_name = r.routine_name AND p.ordinal_position = 0
WHERE r.routine_schema = %%schema string%%
ENDSQL

# mysql proc parameter list query
$XOBIN $MYDB -a -N -M -B -T ProcParam -F MyProcParams -o $DEST $EXTRA << ENDSQL
SELECT
  dtd_identifier AS param_type
FROM information_schema.parameters
WHERE ordinal_position > 0 AND specific_schema = %%schema string%% AND specific_name = %%proc string%%
ORDER BY ordinal_position
ENDSQL

# mysql table list query
$XOBIN $MYDB -a -N -M -B -T Table -F MyTables -o $DEST $EXTRA << ENDSQL
SELECT
  table_name,
  table_comment
FROM information_schema.tables
WHERE table_schema = %%schema string%% AND table_type = %%relkind string%%
ENDSQL

# mysql table column list query
$XOBIN $MYDB -a -N -M -B -T Column -F MyTableColumns -o $DEST $EXTRA << ENDSQL
SELECT
  ordinal_position AS field_ordinal,
  column_name,
  IF(data_type = 'enum', column_name, column_type) AS data_type,
  IF(is_nullable = 'YES', false, true) AS not_null,
  column_default AS default_value,
  IF(column_key = 'PRI', true, false) AS is_primary_key,
  column_comment
FROM information_schema.columns
WHERE table_schema = %%schema string%% AND table_name = %%table string%%
ORDER BY ordinal_position
ENDSQL

# mysql table foreign key list query
$XOBIN $MYDB -a -N -M -B -T ForeignKey -F MyTableForeignKeys -o $DEST $EXTRA << ENDSQL
SELECT
  constraint_name AS foreign_key_name,
  column_name AS column_name,
  referenced_table_name AS ref_table_name,
  referenced_column_name AS ref_column_name
FROM information_schema.key_column_usage
WHERE referenced_table_name IS NOT NULL AND table_schema = %%schema string%% AND table_name = %%table string%%
ENDSQL

# mysql table index list query
$XOBIN $MYDB -a -N -M -B -T Index -F MyTableIndexes -o $DEST $EXTRA << ENDSQL
SELECT
  DISTINCT index_name,
  NOT non_unique AS is_unique
FROM information_schema.statistics
WHERE index_name <> 'PRIMARY' AND index_schema = %%schema string%% AND table_name = %%table string%%
ENDSQL

# mysql index column list query
$XOBIN $MYDB -a -N -M -B -T IndexColumn -F MyIndexColumns -o $DEST $EXTRA << ENDSQL
SELECT
  seq_in_index AS seq_no,
  column_name
FROM information_schema.statistics
WHERE index_schema = %%schema string%% AND table_name = %%table string%% AND index_name = %%index string%%
ORDER BY seq_in_index
ENDSQL
