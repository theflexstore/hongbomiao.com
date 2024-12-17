import logging
import shutil
import subprocess
import tempfile
from datetime import datetime
from decimal import Decimal
from pathlib import Path
from typing import Any
from zoneinfo import ZoneInfo

import pandas as pd
import pythoncom
import win32com.client


class IadsUtil:
    IADS_CONFIG_FILE_NAME = "pfConfig"
    IADS_METADATA_FILE_NAME = "IadsArchiveInfo.txt"
    IRIG_TIME_COLUMN_NAME = "TIME"
    UNIX_TIME_COLUMN_NAME = "_timestamp_ns"
    PARQUET_ZSTD_COMPRESSION_LEVEL = 19

    @staticmethod
    def get_iads_signals(iads_config: Any, query: str) -> set[str]:
        logging.info(f"Executing: {query}")
        results = iads_config.Query(query)
        if not results:
            return set()

        # Check for duplicates
        signal_set = {s.rstrip("\x00") for s in results}
        if len(signal_set) < len(results):
            duplicated_signal_set = {s for s in signal_set if results.count(s) > 1}
            logging.warning(f"Duplicate signals found: {duplicated_signal_set}")
        return signal_set

    @staticmethod
    def create_iads_signal_group(
        iads_config: Any, signal_set: set[str], group_name: str
    ) -> None:
        signals = ",".join(signal_set)
        iads_config.Query(f"update DataGroups set * = |||{group_name}|{signals}|")
        iads_config.Save()

    @staticmethod
    def get_irig_times(iads_metadata_path: Path) -> tuple[str, str, int]:
        irig_start_time = ""
        irig_end_time = ""
        year = None
        with open(iads_metadata_path, "r") as file:
            for line in file:
                if line.startswith("ArchiveStartTime"):
                    irig_start_time = line.split("=")[1].strip()
                elif line.startswith("ArchiveEndTime"):
                    irig_end_time = line.split("=")[1].strip()
                elif line.startswith("FlightDate"):
                    # Parse date using datetime (format: MM/DD/YY)
                    date_str = line.split("=")[1].strip()
                    date = datetime.strptime(date_str, "%m/%d/%y")
                    year = date.year

        if not irig_start_time or not irig_end_time or year is None:
            raise ValueError(
                "Could not find start time, end time, or year in archive info file"
            )

        return irig_start_time, irig_end_time, year

    @staticmethod
    def copy_iads_config(iads_config_path: Path) -> Path:
        # Create a copy with '_temp' suffix
        temp_iads_config_path = (
            iads_config_path.parent / f"{iads_config_path.name}_temp"
        )
        shutil.copy2(iads_config_path, temp_iads_config_path)
        return temp_iads_config_path

    @staticmethod
    def export_to_parquet(
        iads_config_path: Path,
        irig_start_time: str,
        irig_end_time: str,
        group_name: str,
        parquet_output_dir_path: Path,
        iads_manager_exe_path: Path,
    ) -> Path:
        parquet_file_path = parquet_output_dir_path / Path(
            f"{iads_config_path.parent.name}.parquet"
        )
        cmd = [
            str(iads_manager_exe_path),
            "/hide",
            "/DataExport",
            str(iads_config_path),
            "Parquet",
            f"{group_name}|\\ParquetCompressionType=ZSTD\\ParquetZSTDCompressionLevel={IadsUtil.PARQUET_ZSTD_COMPRESSION_LEVEL}",
            irig_start_time,
            irig_end_time,
            str(parquet_file_path),
        ]
        logging.info(f"Executing command: {' '.join(cmd)}")
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode == 0:
            return parquet_file_path
        else:
            raise ValueError(f"Error output: {result}")

    @staticmethod
    def convert_irig_to_unix_time_ns(
        irig_time_ns: int, year: int, timezone: str
    ) -> int:
        # Get Unix timestamp for Jan 1 of the specified year in the given timezone
        local_time = datetime(year, 1, 1, 0, 0, 0, tzinfo=ZoneInfo(timezone))
        year_start_ns = int(
            Decimal(str(local_time.timestamp())) * Decimal("1000000000")
        )
        return year_start_ns + irig_time_ns

    @staticmethod
    def get_iads_dataframe(
        iads_manager_exe_path: Path, iads_data_path: Path, timezone: str
    ) -> pd.DataFrame:
        iads_config_path = iads_data_path / Path(IadsUtil.IADS_CONFIG_FILE_NAME)
        iads_metadata_path = iads_data_path / Path(IadsUtil.IADS_METADATA_FILE_NAME)
        iads_config: Any | None = None
        temp_iads_config_path: Path | None = None

        try:
            with tempfile.TemporaryDirectory() as tmp_dir:
                tmp_dir_path = Path(tmp_dir)
                logging.info(f"Created temporary directory: {tmp_dir_path}")

                # Copy IADS config file
                temp_iads_config_path = IadsUtil.copy_iads_config(iads_config_path)
                logging.info(
                    f"Created a copy of IADS config file: {temp_iads_config_path}"
                )

                # Get IADS config
                pythoncom.CoInitialize()
                iads_config = win32com.client.Dispatch("IadsConfigInterface.IadsConfig")
                iads_config.Open(str(temp_iads_config_path), False)

                # Get signals
                query = "select Parameter from ParameterDefaults"
                signal_set = IadsUtil.get_iads_signals(iads_config, query)

                # Create signal group
                group_name = "AllSignals"
                logging.info(
                    f"Creating signal group with query. group_name = '{group_name}'"
                )
                IadsUtil.create_iads_signal_group(iads_config, signal_set, group_name)
                logging.info(f"Signal group '{group_name}' created successfully")

                # Get IRIG times and year
                irig_start_time, irig_end_time, year = IadsUtil.get_irig_times(
                    iads_metadata_path
                )

                # Export to parquet
                parquet_file_path = IadsUtil.export_to_parquet(
                    temp_iads_config_path,
                    irig_start_time,
                    irig_end_time,
                    group_name,
                    tmp_dir_path,
                    iads_manager_exe_path,
                )
                logging.info("Export parquet completed successfully")

                # Read the exported parquet file
                df = pd.read_parquet(parquet_file_path)

                # Convert IRIG time to Unix time
                df[IadsUtil.UNIX_TIME_COLUMN_NAME] = df[
                    IadsUtil.IRIG_TIME_COLUMN_NAME
                ].apply(
                    lambda irig_time_ns: IadsUtil.convert_irig_to_unix_time_ns(
                        irig_time_ns, year, timezone
                    )
                )
                logging.info(
                    f"Added {IadsUtil.UNIX_TIME_COLUMN_NAME} column with Unix time in nanoseconds ({year = })"
                )
                return df

        except Exception as e:
            logging.error(f"{e = }")
            return None

        finally:
            if iads_config is not None:
                pythoncom.CoUninitialize()

            if temp_iads_config_path is not None:
                temp_iads_config_path.unlink()
                logging.info(f"Delete the temporary pfConfig: {temp_iads_config_path}")


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    iads_data_path = Path(r"C:\iads_data")
    iads_manager_exe_path = Path(
        r"C:\Program Files\IADS\DataManager\IadsDataManager.exe"
    )
    timezone = "America/Los_Angeles"
    df = IadsUtil.get_iads_dataframe(iads_manager_exe_path, iads_data_path, timezone)
    logging.info(f"{df = }")
