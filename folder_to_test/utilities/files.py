from typing import List

from fastapi import UploadFile


async def get_names_with_files(files: List[UploadFile]) -> dict:
    if files:
        return {file.filename: file for file in files if file.filename}
    return {}
