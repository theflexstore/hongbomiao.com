import logging
from pathlib import Path
from typing import Any

import pythoncom
import win32com.client


def show_version_from_file(iads_config: Any, iads_config_path: Path) -> None:
    try:
        version = iads_config.VersionFromFile(iads_config_path)
        logging.info(f"{version = }")
    except Exception as e:
        logging.error(f"{iads_config_path = }, {e = }")


def execute_query(iads_config: Any, query: str) -> None:
    try:
        logging.info(f"{query = }")
        results: list[Any] | None = iads_config.Query(query)
        if results:
            for result in results:
                logging.info(f"{result = }")
    except Exception as e:
        logging.error(f"{e = }")


def process_config(iads_config_path: Path) -> None:
    iads_config: Any | None = None
    try:
        pythoncom.CoInitialize()
        iads_config = win32com.client.Dispatch("IadsConfigInterface.IadsConfig")

        show_version_from_file(iads_config, iads_config_path)
        iads_config.Open(iads_config_path, True)

        execute_query(iads_config, "select * from Desktops")
        execute_query(iads_config, "select System.RowNumber from Desktops")
        execute_query(iads_config, "select Parameter from ParameterDefaults")

        iads_config.Close(True)
    except Exception as e:
        logging.error(f"{e = }")
    finally:
        # Clean up COM resources
        if iads_config:
            iads_config = None
        pythoncom.CoUninitialize()


def main() -> None:
    iads_config_path = Path("pfConfig")
    process_config(iads_config_path)


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    main()
