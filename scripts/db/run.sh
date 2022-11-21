#/bin/sh
psql root < cpu_usage.sql &&
echo "====== BEGING inserting data =======" &&
psql  -d homework -c "\COPY cpu_usage FROM cpu_usage.csv CSV HEADER" &&
echo "====== DONE inserting data ======="
