from pyspark.sql import DataFrame, SparkSession
from utils.trip import load_trips, preprocess_trips
from utils.zone import load_zones, preprocess_zones


def get_top_routes(
    spark: SparkSession, trips: DataFrame, zones: DataFrame, sql_query: str
) -> DataFrame:
    trips.createOrReplaceTempView("trips")
    zones.createOrReplaceTempView("zones")
    return spark.sql(sql_query)


def main(
    data_dirname: str,
    query_dirname: str,
    trip_filenames: list[str],
    zone_filename: str,
    query_filename: str,
) -> None:
    trip_data_paths = [f"{data_dirname}/{f}" for f in trip_filenames]
    zone_data_path = f"{data_dirname}/{zone_filename}"

    with open(f"{query_dirname}/{query_filename}", "r") as f:
        sql_query = f.read()

    spark = (
        SparkSession.builder.master("local[*]")
        .appName("find-taxi-top-routes-sql")
        .config("spark.ui.port", "4040")
        .getOrCreate()
    )

    trips = load_trips(spark, trip_data_paths)
    zones = load_zones(spark, zone_data_path)

    trips = preprocess_trips(trips)
    print((trips.count(), len(trips.columns)))
    trips.show()

    zones = preprocess_zones(zones)
    print((zones.count(), len(zones.columns)))
    zones.show()

    top_routes = get_top_routes(spark, trips, zones, sql_query)
    print((top_routes.count(), len(top_routes.columns)))
    top_routes.show(truncate=False)

    spark.stop()


if __name__ == "__main__":
    # https://www.nyc.gov/site/tlc/about/tlc-trip-record-data.page
    data_dirname = "data"
    query_dirname = "src/queries"
    trip_filenames = [
        "yellow_tripdata_2019-01.parquet",
        "yellow_tripdata_2019-02.parquet",
        "yellow_tripdata_2019-03.parquet",
        "yellow_tripdata_2019-04.parquet",
        "yellow_tripdata_2019-05.parquet",
        "yellow_tripdata_2019-06.parquet",
        "yellow_tripdata_2019-07.parquet",
        "yellow_tripdata_2019-08.parquet",
        "yellow_tripdata_2019-09.parquet",
        "yellow_tripdata_2019-10.parquet",
        "yellow_tripdata_2019-11.parquet",
        "yellow_tripdata_2019-12.parquet",
        "yellow_tripdata_2020-01.parquet",
        "yellow_tripdata_2020-02.parquet",
        "yellow_tripdata_2020-03.parquet",
        "yellow_tripdata_2020-04.parquet",
        "yellow_tripdata_2020-05.parquet",
        "yellow_tripdata_2020-06.parquet",
        "yellow_tripdata_2020-07.parquet",
        "yellow_tripdata_2020-08.parquet",
        "yellow_tripdata_2020-09.parquet",
        "yellow_tripdata_2020-10.parquet",
        "yellow_tripdata_2020-11.parquet",
        "yellow_tripdata_2020-12.parquet",
        "yellow_tripdata_2021-01.parquet",
        "yellow_tripdata_2021-02.parquet",
        "yellow_tripdata_2021-03.parquet",
        "yellow_tripdata_2021-04.parquet",
        "yellow_tripdata_2021-05.parquet",
        "yellow_tripdata_2021-06.parquet",
        "yellow_tripdata_2021-07.parquet",
        "yellow_tripdata_2021-08.parquet",
        "yellow_tripdata_2021-09.parquet",
        "yellow_tripdata_2021-10.parquet",
        "yellow_tripdata_2021-11.parquet",
        "yellow_tripdata_2021-12.parquet",
        "yellow_tripdata_2022-01.parquet",
        "yellow_tripdata_2022-02.parquet",
        "yellow_tripdata_2022-03.parquet",
        "yellow_tripdata_2022-04.parquet",
        "yellow_tripdata_2022-05.parquet",
        "yellow_tripdata_2022-06.parquet",
    ]
    zone_filename = "taxi_zones.csv"
    query_filename = "query.sql"
    main(data_dirname, query_dirname, trip_filenames, zone_filename, query_filename)
