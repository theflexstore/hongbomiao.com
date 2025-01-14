import logging
import urllib.request
from pathlib import Path

import nvidia.dali.fn as fn
import nvidia.dali.types as types
import torch
from nvidia.dali import pipeline_def
from nvidia.dali.plugin.pytorch import DALIGenericIterator


def download_sample_images(data_path: Path) -> None:
    # Create main directory if it doesn not exist
    data_path.mkdir(parents=True, exist_ok=True)

    # Create a class subdirectory (e.g., "class0")
    class_dir = data_path / "class0"
    class_dir.mkdir(parents=True, exist_ok=True)

    # Sample image URLs
    image_urls: list[str] = [
        "https://raw.githubusercontent.com/pytorch/hub/master/images/dog.jpg"
    ]

    # Download images into the class subdirectory
    for i, url in enumerate(image_urls):
        try:
            filename = f"image_{i}.jpg"
            filepath = class_dir / filename
            if not filepath.exists():
                logging.info(f"Downloading {url} to {filepath}")
                urllib.request.urlretrieve(url, str(filepath))
        except Exception as e:
            logging.error(f"Error downloading {url}: {e}")


@pipeline_def(batch_size=2, num_threads=2, device_id=None)
def image_pipeline(data_path: Path):
    jpegs, labels = fn.readers.file(
        file_root=data_path, random_shuffle=True, initial_fill=2
    )

    images = fn.decoders.image(jpegs, device="cpu")
    images = fn.resize(images, resize_x=224, resize_y=224)
    images = fn.crop_mirror_normalize(
        images,
        mean=[0.485 * 255, 0.456 * 255, 0.406 * 255],
        std=[0.229 * 255, 0.224 * 255, 0.225 * 255],
        dtype=types.FLOAT,
    )
    images = fn.transpose(images, perm=[2, 0, 1])

    return images, labels


def get_num_samples(data_path: Path) -> int:
    image_files = list(data_path.rglob("*.jpg"))
    return len(image_files)


def main() -> None:
    BATCH_SIZE: int = 2
    NUM_THREADS: int = 2

    # Create data directory and download sample images
    data_path = Path("data")
    download_sample_images(data_path)

    # Get total number of samples
    num_samples = get_num_samples(data_path)
    if num_samples == 0:
        logging.error("No images available in the directory.")
        return

    pipe = image_pipeline(
        data_path=data_path, batch_size=BATCH_SIZE, num_threads=NUM_THREADS
    )
    pipe.build()

    dali_iter = DALIGenericIterator(
        pipelines=[pipe],
        output_map=["data", "label"],
        reader_name="Reader",
        auto_reset=True,
    )

    logging.info("Pipeline created successfully!")
    logging.info(f"Ready to process images from {data_path}")

    try:
        for i, data in enumerate(dali_iter):
            images: torch.Tensor = data[0]["data"]
            labels: torch.Tensor = data[0]["label"]
            logging.info(f"Batch {i}: Image shape: {images.shape}, Labels: {labels}")
    except StopIteration:
        logging.info("Finished processing all images.")


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    main()
