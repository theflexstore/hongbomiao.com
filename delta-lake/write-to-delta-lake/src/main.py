import config
import pandas as pd
import pyarrow as pa
from deltalake.writer import write_deltalake


def main():
    df = pd.read_parquet(config.parquet_path, engine="pyarrow")
    storage_options = {
        "AWS_DEFAULT_REGION": config.aws_default_region,
        "AWS_ACCESS_KEY_ID": config.aws_access_key_id,
        "AWS_SECRET_ACCESS_KEY": config.aws_secret_access_key,
        "AWS_S3_ALLOW_UNSAFE_RENAME": "true",
    }
    schema = pa.schema(
        [
            ("timestamp", pa.float64()),
            ("current", pa.float64()),
            ("voltage", pa.float64()),
            ("temperature", pa.float64()),
        ]
    )
    write_deltalake(
        config.s3_path,
        df,
        mode="append",
        schema=schema,
        storage_options=storage_options,
    )


if __name__ == "__main__":
    main()
